package appjson

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dokku/dokku/plugins/common"
)

// TriggerAppJSONProcessDeployParallelism returns the max number of processes to deploy in parallel
func TriggerAppJSONProcessDeployParallelism(appName string, processType string) error {
	appJSON, err := getAppJSON(appName)
	if err != nil {
		return err
	}

	parallelism := 1
	for procType, formation := range appJSON.Formation {
		if procType != processType {
			continue
		}

		if formation.MaxParallel == nil {
			continue
		}

		if *formation.MaxParallel > 0 {
			parallelism = *formation.MaxParallel
		}
	}

	fmt.Println(parallelism)
	return nil
}

// TriggerInstall initializes app restart policies
func TriggerInstall() error {
	if err := common.PropertySetup("app-json"); err != nil {
		return fmt.Errorf("Unable to install the app-json plugin: %s", err.Error())
	}

	directory := filepath.Join(common.MustGetEnv("DOKKU_LIB_ROOT"), "data", "app-json")
	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}

	return common.SetPermissions(directory, 0755)
}

// TriggerPostAppCloneSetup creates new app-json files
func TriggerPostAppCloneSetup(oldAppName string, newAppName string) error {
	err := common.PropertyClone("app-json", oldAppName, newAppName)
	if err != nil {
		return err
	}

	return nil
}

// TriggerPostAppRenameSetup renames app-json files
func TriggerPostAppRenameSetup(oldAppName string, newAppName string) error {
	if err := common.PropertyClone("app-json", oldAppName, newAppName); err != nil {
		return err
	}

	if err := common.PropertyDestroy("app-json", oldAppName); err != nil {
		return err
	}

	return nil
}

// TriggerPostDelete destroys the app-json data for a given app container
func TriggerPostDelete(appName string) error {
	directory := filepath.Join(common.MustGetEnv("DOKKU_LIB_ROOT"), "data", "app-json", appName)
	dataErr := os.RemoveAll(directory)
	propertyErr := common.PropertyDestroy("app-json", appName)

	if dataErr != nil {
		return dataErr
	}

	return propertyErr
}

// TriggerPostDeploy is a trigger to execute the postdeploy deployment task
func TriggerPostDeploy(appName string, imageTag string) error {
	image, err := common.GetDeployingAppImageName(appName, imageTag, "")
	if err != nil {
		return err
	}

	return executeScript(appName, image, imageTag, "postdeploy")
}

// TriggerPreDeploy is a trigger to execute predeploy and release deployment tasks
func TriggerPreDeploy(appName string, imageTag string) error {
	image := common.GetAppImageName(appName, imageTag, "")
	if err := refreshAppJSON(appName, image); err != nil {
		return err
	}

	if err := executeScript(appName, image, imageTag, "predeploy"); err != nil {
		return err
	}

	if err := executeScript(appName, image, imageTag, "release"); err != nil {
		return err
	}

	if err := setScale(appName, image); err != nil {
		return err
	}

	if common.PropertyGet("common", appName, "deployed") == "true" {
		return nil
	}

	// Ensure that a failed postdeploy does not trigger twice
	if common.PropertyGet("app-json", appName, "heroku.postdeploy") == "executed" {
		return nil
	}

	if err := common.PropertyWrite("app-json", appName, "heroku.postdeploy", "executed"); err != nil {
		return err
	}

	return executeScript(appName, image, imageTag, "heroku.postdeploy")
}
