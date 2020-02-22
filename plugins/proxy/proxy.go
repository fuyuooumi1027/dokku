package proxy

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dokku/dokku/plugins/common"
	"github.com/dokku/dokku/plugins/config"
	columnize "github.com/ryanuber/columnize"
)

type PortMap struct {
	ContainerPort int
	HostPort      int
	Scheme        string
}

func (p PortMap) String() string {
	return fmt.Sprintf("%s:%d:%d", p.Scheme, p.HostPort, p.ContainerPort)
}

// IsAppProxyEnabled returns true if proxy is enabled; otherwise return false
func IsAppProxyEnabled(appName string) bool {
	err := common.VerifyAppName(appName)
	if err != nil {
		common.LogFail(err.Error())
	}

	proxyEnabled := true
	disableProxy := config.GetWithDefault(appName, "DOKKU_DISABLE_PROXY", "")
	if disableProxy != "" {
		proxyEnabled = false
	}
	return proxyEnabled
}

// ReportSingleApp is an internal function that displays the app report for one or more apps
func ReportSingleApp(appName string, infoFlag string) error {
	if err := common.VerifyAppName(appName); err != nil {
		return err
	}

	proxyEnabled := "false"
	if IsAppProxyEnabled(appName) {
		proxyEnabled = "true"
	}

	var proxyPortMap []string
	for _, portMap := range getProxyPortMap(appName) {
		proxyPortMap = append(proxyPortMap, portMap.String())
	}

	infoFlags := map[string]string{
		"--proxy-enabled":  proxyEnabled,
		"--proxy-type":     getAppProxyType(appName),
		"--proxy-port-map": strings.Join(proxyPortMap, " "),
	}

	trimPrefix := false
	uppercaseFirstCharacter := true
	return common.ReportSingleApp("proxy", appName, infoFlag, infoFlags, trimPrefix, uppercaseFirstCharacter)
}

func addProxyPorts(appName string, proxyPortMap []PortMap) error {
	allPortMaps := getProxyPortMap(appName)
	allPortMaps = append(allPortMaps, proxyPortMap...)

	return setProxyPorts(appName, allPortMaps)
}

func filterAppProxyPorts(appName string, scheme string, hostPort int) []PortMap {
	var filteredProxyMaps []PortMap
	proxyPortMap := getProxyPortMap(appName)
	for _, portMap := range proxyPortMap {
		if portMap.Scheme == scheme && portMap.HostPort == hostPort {
			filteredProxyMaps = append(filteredProxyMaps, portMap)
		}
	}

	return filteredProxyMaps
}

func getAppName(args []string) (appName string, err error) {
	if len(args) >= 1 {
		appName = args[0]
	} else {
		err = errors.New("Please specify an app to run the command on")
	}

	return
}

func getAppProxyType(appName string) string {
	return config.GetWithDefault(appName, "DOKKU_APP_PROXY_TYPE", "nginx")
}

func getProxyPortMap(appName string) []PortMap {
	value := config.GetWithDefault(appName, "DOKKU_PROXY_PORT_MAP", "")
	return parseProxyPortMapString(value)
}

func listAppProxyPorts(appName string) error {
	proxyPortMap := getProxyPortMap(appName)

	if len(proxyPortMap) == 0 {
		return errors.New("No port mappings configured for app")
	}

	var lines []string
	if os.Getenv("DOKKU_QUIET_OUTPUT") == "" {
		lines = append(lines, "-----> scheme:host port:container port")
	}

	for _, portMap := range proxyPortMap {
		lines = append(lines, portMap.String())
	}

	common.LogInfo1Quiet(fmt.Sprintf("Port mappings for %s", appName))
	config := columnize.DefaultConfig()
	config.Delim = ":"
	config.Prefix = "    "
	config.Empty = ""
	fmt.Println(columnize.Format(lines, config))
	return nil
}

func setProxyPorts(appName string, proxyPortMap []PortMap) error {
	var value []string
	for _, portMap := range uniqueProxyPortMap(proxyPortMap) {
		value = append(value, portMap.String())
	}

	entries := map[string]string{
		"DOKKU_PROXY_PORT_MAP": strings.Join(value, " "),
	}
	return config.SetMany(appName, entries, false)
}

func removeProxyPorts(appName string, proxyPortMap []PortMap) error {
	var toRemove map[string]bool

	for _, portMap := range proxyPortMap {
		toRemove[portMap.String()] = true
	}

	var toSet []PortMap
	existingPortMaps := getProxyPortMap(appName)
	for _, portMap := range existingPortMaps {
		if toRemove[portMap.String()] {
			continue
		}

		toSet = append(toSet, portMap)
	}

	if len(toSet) == 0 {
		keys := []string{"DOKKU_PROXY_PORT_MAP"}
		return config.UnsetMany(appName, keys, false)
	}

	return setProxyPorts(appName, toSet)
}

func parseProxyPortMapString(stringPortMap string) []PortMap {
	var proxyPortMap []PortMap

	for _, v := range strings.Split(strings.TrimSpace(stringPortMap), "") {
		parts := strings.SplitN(v, ":", 3)
		if len(parts) != 3 {
			common.LogWarn(fmt.Sprintf("Invalid port map %s", v))
			continue
		}

		hostPort, err := strconv.Atoi(parts[1])
		if err != nil {
			common.LogWarn(fmt.Sprintf("Invalid port map %s", v))
		}

		containerPort, err := strconv.Atoi(parts[2])
		if err != nil {
			common.LogWarn(fmt.Sprintf("Invalid port map %s", v))
		}

		proxyPortMap = append(proxyPortMap, PortMap{
			ContainerPort: containerPort,
			HostPort:      hostPort,
			Scheme:        parts[0],
		})
	}

	return uniqueProxyPortMap(proxyPortMap)
}

func uniqueProxyPortMap(proxyPortMap []PortMap) []PortMap {
	var unique []PortMap
	var existingPortMaps map[string]bool

	for _, portMap := range proxyPortMap {
		if existingPortMaps[portMap.String()] {
			continue
		}

		existingPortMaps[portMap.String()] = true
		unique = append(unique, portMap)
	}

	return unique
}
