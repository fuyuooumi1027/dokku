package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dokku/dokku/plugins/common"
	"github.com/dokku/dokku/plugins/ps"
)

// main entrypoint to all triggers
func main() {
	parts := strings.Split(os.Args[0], "/")
	trigger := parts[len(parts)-1]
	flag.Parse()

	var err error
	switch trigger {
	case "app-restart":
		appName := flag.Arg(0)
		err = ps.TriggerAppRestart(appName)
	case "core-post-deploy":
		appName := flag.Arg(0)
		err = ps.TriggerCorePostDeploy(appName)
	case "core-post-extract":
		appName := flag.Arg(0)
		sourceWorkDir := flag.Arg(1)
		err = ps.TriggerCorePostExtract(appName, sourceWorkDir)
	case "install":
		err = ps.TriggerInstall()
	case "post-app-clone":
		oldAppName := flag.Arg(0)
		newAppName := flag.Arg(1)
		err = ps.TriggerPostAppClone(oldAppName, newAppName)
	case "post-app-clone-setup":
		oldAppName := flag.Arg(0)
		newAppName := flag.Arg(1)
		err = ps.TriggerPostAppCloneSetup(oldAppName, newAppName)
	case "post-app-rename":
		oldAppName := flag.Arg(0)
		newAppName := flag.Arg(1)
		err = ps.TriggerPostAppRename(oldAppName, newAppName)
	case "post-app-rename-setup":
		oldAppName := flag.Arg(0)
		newAppName := flag.Arg(1)
		err = ps.TriggerPostAppRenameSetup(oldAppName, newAppName)
	case "post-create":
		appName := flag.Arg(0)
		err = ps.TriggerPostCreate(appName)
	case "post-delete":
		appName := flag.Arg(0)
		err = ps.TriggerPostDelete(appName)
	case "post-stop":
		appName := flag.Arg(0)
		err = ps.TriggerPostStop(appName)
	case "pre-deploy":
		appName := flag.Arg(0)
		imageTag := flag.Arg(1)
		err = ps.TriggerPreDeploy(appName, imageTag)
	case "procfile-get-command":
		appName := flag.Arg(0)
		processType := flag.Arg(1)
		port := common.ToInt(flag.Arg(2), 5000)
		err = ps.TriggerProcfileGetCommand(appName, processType, port)
	case "ps-can-scale":
		appName := flag.Arg(0)
		canScale := common.ToBool(flag.Arg(1))
		err = ps.TriggerPsCanScale(appName, canScale)
	case "ps-current-scale":
		appName := flag.Arg(0)
		err = ps.TriggerPsCurrentScale(appName)
	case "ps-set-scale":
		appName := flag.Arg(0)
		appName, args := common.ShiftString(flag.Args())
		skipDeploy, args := common.ShiftString(args)
		clearExisting, processTuples := common.ShiftString(args)
		err = ps.TriggerPsSetScale(appName, common.ToBool(skipDeploy), common.ToBool(clearExisting), processTuples)
	case "report":
		appName := flag.Arg(0)
		err = ps.ReportSingleApp(appName, "", "")
	default:
		err = fmt.Errorf("Invalid plugin trigger call: %s", trigger)
	}

	if err != nil {
		common.LogFailWithError(err)
	}
}
