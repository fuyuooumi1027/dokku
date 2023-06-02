package registry

import (
	"strings"

	"github.com/dokku/dokku/plugins/common"
)

// ReportSingleApp is an internal function that displays the registry report for one or more apps
func ReportSingleApp(appName string, format string, infoFlag string) error {
	if err := common.VerifyAppName(appName); err != nil {
		return err
	}

	flags := map[string]common.ReportFunc{
		"--registry-computed-image-repo":      reportComputedImageRepo,
		"--registry-image-repo":               reportImageRepo,
		"--registry-computed-push-on-release": reportComputedPushOnRelease,
		"--registry-global-push-on-release":   reportGlobalPushOnRelease,
		"--registry-push-on-release":          reportPushOnRelease,
		"--registry-computed-server":          reportComputedServer,
		"--registry-global-server":            reportGlobalServer,
		"--registry-server":                   reportServer,
		"--registry-tag-version":              reportTagVersion,
	}

	flagKeys := []string{}
	for flagKey := range flags {
		flagKeys = append(flagKeys, flagKey)
	}

	trimPrefix := false
	uppercaseFirstCharacter := true
	infoFlags := common.CollectReport(appName, infoFlag, flags)
	return common.ReportSingleApp("registry", appName, infoFlag, infoFlags, flagKeys, format, trimPrefix, uppercaseFirstCharacter)
}

func reportComputedImageRepo(appName string) string {
	imageRepo := reportImageRepo(appName)
	imageRepo = strings.TrimSpace(imageRepo)
	if imageRepo == "" {
		imageRepo = common.GetAppImageRepo(appName)
	}

	return imageRepo
}

func reportImageRepo(appName string) string {
	return common.PropertyGet("registry", appName, "image-repo")
}

func reportComputedPushOnRelease(appName string) string {
	value := reportPushOnRelease(appName)
	value = strings.TrimSpace(value)
	if value == "" {
		value = reportGlobalPushOnRelease(appName)
	}

	if value == "" {
		value = DefaultProperties["push-on-release"]
	}

	return value
}

func reportGlobalPushOnRelease(appName string) string {
	return common.PropertyGet("registry", "--global", "push-on-release")
}

func reportPushOnRelease(appName string) string {
	return common.PropertyGet("registry", appName, "push-on-release")
}

func reportComputedServer(appName string) string {
	server := getRegistryServerForApp(appName)
	return strings.TrimSpace(server)
}

func reportGlobalServer(appName string) string {
	return common.PropertyGet("registry", "--global", "server")
}

func reportServer(appName string) string {
	return common.PropertyGet("registry", appName, "server")
}

func reportTagVersion(appName string) string {
	tagVersion := common.PropertyGet("registry", appName, "tag-version")
	return strings.TrimSpace(tagVersion)
}
