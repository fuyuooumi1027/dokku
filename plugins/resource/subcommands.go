package resource

import (
	"errors"
	"strings"

	"github.com/dokku/dokku/plugins/common"
)

// CommandLimit implements resource:limit
func CommandLimit(appName string, processType string, r Resource) error {
	if appName == "" {
		return errors.New("Please specify an app to run the command on")
	}

	return setResourceType(appName, processType, r, "limit")
}

// CommandLimitClear implements resource:limit-clear
func CommandLimitClear(appName string, processType string) error {
	if appName == "" {
		return errors.New("Please specify an app to run the command on")
	}

	clearByResourceType(appName, processType, "limit")
	return nil
}

// CommandReport displays a resource report for one or more apps
func CommandReport(appName string, infoFlag string) error {
	if strings.HasPrefix(appName, "--") {
		infoFlag = appName
		appName = ""
	}

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

// CommandReserve implements resource:reserve
func CommandReserve(appName string, processType string, r Resource) error {
	if appName == "" {
		return errors.New("Please specify an app to run the command on")
	}

	return setResourceType(appName, processType, r, "reserve")
}

// CommandReserveClear implements resource:reserve-clear
func CommandReserveClear(appName string, processType string) error {
	if appName == "" {
		return errors.New("Please specify an app to run the command on")
	}

	clearByResourceType(appName, processType, "reserve")
	return nil
}
