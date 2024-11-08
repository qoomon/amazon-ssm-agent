//go:build darwin
// +build darwin

package osdetect

import (
	"os/exec"
	"strings"

	"github.com/aws/amazon-ssm-agent/agent/log"

	c "github.com/aws/amazon-ssm-agent/agent/plugins/configurepackage/envdetect/constants"
)

func DetectPkgManager(platform string, version string, family string) (string, error) {
	return c.PackageManagerMac, nil
}

func DetectInitSystem() (string, error) {
	return c.InitLaunchd, nil
}

func DetectPlatform(_ log.T) (string, string, string, error) {
	cmdOut, err := exec.Command("/usr/bin/sw_vers", "-productVersion").Output()
	if err != nil {
		return "", "", "", err
	}

	return c.PlatformDarwin, extractDarwinVersion(cmdOut), c.PlatformFamilyDarwin, nil
}

func extractDarwinVersion(data []byte) string {
	return strings.TrimSpace(string(data))
}
