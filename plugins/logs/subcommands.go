package logs

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/dokku/dokku/plugins/common"
)

// CommandDefault displays recent log output
func CommandDefault(appName string, num int64, process string, tail, quiet bool) error {
	if err := common.VerifyAppName(appName); err != nil {
		return err
	}

	if !common.IsDeployed(appName) {
		return fmt.Errorf("App %s has not been deployed", appName)
	}

	s := common.GetAppScheduler(appName)
	t := strconv.FormatBool(tail)
	q := strconv.FormatBool(quiet)
	n := strconv.FormatInt(num, 10)

	if err := common.PlugnTrigger("scheduler-logs", []string{s, appName, process, t, q, n}...); err != nil {
		return err
	}
	return nil
}

// CommandFailed shows the last failed deploy logs
func CommandFailed(appName string, allApps bool) error {
	if allApps {
		return common.RunCommandAgainstAllAppsSerially(GetFailedLogs, "logs:failed")
	}

	if appName == "" {
		common.LogWarn("Deprecated: logs:failed specified without app, assuming --all")
		return common.RunCommandAgainstAllAppsSerially(GetFailedLogs, "logs:failed")
	}

	if err := common.VerifyAppName(appName); err != nil {
		return err
	}

	return GetFailedLogs(appName)
}

// CommandReport displays a logs report for one or more apps
func CommandReport(appName string, infoFlag string) error {
	if len(appName) == 0 {
		apps, err := common.DokkuApps()
		if err != nil {
			return err
		}
		for _, appName := range apps {
			if err := ReportSingleApp(appName, infoFlag); err != nil {
				return err
			}
		}
		return nil
	}

	return ReportSingleApp(appName, infoFlag)
}

// CommandSet sets or clears a logs property for an app
func CommandSet(appName string, property string, value string) error {
	if property == "vector-sink" && value != "" {
		_, err := valueToConfig(appName, value)
		if err != nil {
			return err
		}
	}

	common.CommandPropertySet("logs", appName, property, value, DefaultProperties, GlobalProperties)
	if property == "vector-sink" {
		common.LogVerboseQuiet(fmt.Sprintf("Writing updated vector config to %s", filepath.Join(common.MustGetEnv("DOKKU_LIB_ROOT"), "data", "logs", "vector.json")))
		return writeVectorConfig()
	}
	return nil
}

// CommandVectorLogs tails the log output for the vector container
func CommandVectorLogs() error {
	if !common.ContainerExists(vectorContainerName) {
		return errors.New("Vector container does not exist")
	}

	if !common.ContainerIsRunning(vectorContainerName) {
		return errors.New("Vector container is not running")
	}

	common.LogInfo1Quiet("Tailing vector container logs")
	common.LogVerboseQuietContainerLogsTail(vectorContainerName)

	return nil
}

// CommandVectorStart starts a new vector container
// or starts an existing one if it already exists
func CommandVectorStart(vectorImage string) error {
	common.LogInfo2("Starting vector container")
	if !common.ContainerExists(vectorContainerName) {
		common.LogVerbose("Ensuring vector configuration exists")
		if err := writeVectorConfig(); err != nil {
			return err
		}

		return startVectorContainer(vectorImage)
	}

	if common.ContainerIsRunning(vectorContainerName) {
		common.LogVerbose("Container already running")
		return nil
	}

	if !common.ContainerStart(vectorContainerName) {
		return errors.New("Unable to start vector container")
	}

	return nil
}

// CommandVectorStop stops and removes an existing vector container
func CommandVectorStop() error {
	common.LogInfo2Quiet("Stopping and removing vector container")
	return killVectorContainer()
}
