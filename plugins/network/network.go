package network

import (
	"fmt"
	"strconv"
	"strings"

	common "github.com/dokku/dokku/plugins/common"
	config "github.com/dokku/dokku/plugins/config"
	proxy "github.com/dokku/dokku/plugins/proxy"

	sh "github.com/codeskyblue/go-sh"
)

// return the ipaddr for a given app container
func GetContainerIpaddress(appName string, procType string, isHerokuishContainer bool, containerId string) string {
	if procType != "web" {
		return ""
	}

	ipAddress := "127.0.0.1"
	if !proxy.IsAppProxyEnabled(appName) {
		return ipAddress
	}

	b, err := common.DockerInspect(containerId, "'{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}'")
	if err != nil || len(b) == 0 {
		// docker < 1.9 compatibility
		b, err = common.DockerInspect(containerId, "'{{ .NetworkSettings.IPAddress }}'")
	}

	if err == nil {
		return string(b[:])
	}

	return ""
}

// return the port for a given app container
func GetContainerPort(appName string, procType string, isHerokuishContainer bool, containerId string) string {
	if procType != "web" {
		return ""
	}

	dockerfilePorts := make([]string, 0)
	port := ""

	if isHerokuishContainer {
		configValue := config.GetWithDefault(appName, "DOKKU_DOCKERFILE_PORTS", "")
		if configValue != "" {
			dockerfilePorts = strings.Split(configValue, " ")
		}
	}

	if len(dockerfilePorts) > 0 {
		for _, p := range dockerfilePorts {
			if strings.HasSuffix(p, "/udp") {
				continue
			}
			port = strings.TrimSuffix(p, "/tcp")
			if port != "" {
				break
			}
		}
	} else {
		port = "5000"
	}

	if !proxy.IsAppProxyEnabled(appName) {
		b, err := sh.Command("docker", "port", containerId, port).Output()
		if err == nil {
			port = strings.Split(string(b[:]), ":")[1]
		}
	}

	return port
}

// builds network config files
func BuildConfig(appName string) {
	err := common.VerifyAppName(appName)
	if err != nil {
		common.LogFail(err.Error())
	}

	if !common.IsDeployed(appName) {
		return
	}

	appRoot := strings.Join([]string{common.MustGetEnv("DOKKU_ROOT"), appName}, "/")
	scaleFile := strings.Join([]string{appRoot, "DOKKU_SCALE"}, "/")
	if !common.FileExists(scaleFile) {
		return
	}

	image := common.GetAppImageName(appName, "", "")
	isHerokuishContainer := common.IsImageHerokuishBased(image)

	common.LogInfo1(fmt.Sprintf("Ensuring network configuration is in sync for %s", appName))
	lines, err := common.FileToSlice(scaleFile)
	if err != nil {
		return
	}
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		procParts := strings.SplitN(line, "=", 2)
		if len(procParts) != 2 {
			continue
		}
		procType := procParts[0]
		procCount, err := strconv.Atoi(procParts[1])
		if err != nil {
			continue
		}

		containerIndex := 0
		for containerIndex < procCount {
			containerIndex += 1
			containerIdFile := fmt.Sprintf("%v/CONTAINER.%v.%v", appRoot, procType, containerIndex)

			containerId := common.ReadFirstLine(containerIdFile)
			if containerId == "" || !common.ContainerIsRunning(containerId) {
				continue
			}

			ipAddress := GetContainerIpaddress(appName, procType, isHerokuishContainer, containerId)
			port := GetContainerPort(appName, procType, isHerokuishContainer, containerId)

			if ipAddress != "" {
				_, err := sh.Command("plugn", "trigger", "network-write-ipaddr", appName, procType, containerIndex, ipAddress).Output()
				if err != nil {
					common.LogWarn(err.Error())
				}
			}

			if port != "" {
				_, err := sh.Command("plugn", "trigger", "network-write-port", appName, procType, containerIndex, port).Output()
				if err != nil {
					common.LogWarn(err.Error())
				}
			}
		}
	}
}
