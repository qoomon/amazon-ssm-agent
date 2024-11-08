//go:build windows
// +build windows

package osdetect

import (
	"fmt"

	"github.com/aws/amazon-ssm-agent/agent/log"
	"github.com/aws/amazon-ssm-agent/agent/platform"

	c "github.com/aws/amazon-ssm-agent/agent/plugins/configurepackage/envdetect/constants"
)

// https://msdn.microsoft.com/en-us/library/aa394239%28v=vs.85%29.aspx

var getOSInfo = func(osData platform.Win32_OperatingSystem) (platform.Win32_OperatingSystem, error) {
	return platform.GetSingleWMIObject(osData)
}

func DetectPkgManager(platform string, version string, family string) (string, error) {
	return c.PackageManagerWindows, nil
}

func DetectInitSystem() (string, error) {
	return c.InitWindows, nil
}

func DetectPlatform(log log.T) (string, string, string, error) {
	if wmiData, err := getOSInfo(platform.Win32_OperatingSystem{}); err != nil {
		log.Errorf("Failed to fetch OS details from WMI, proceeding without 'Version': %v", err)
		return c.PlatformWindows, "", c.PlatformFamilyWindows, nil
	} else {
		version := wmiData.Version
		if isWindowsNano(wmiData.OperatingSystemSKU) {
			version = fmt.Sprint(version, "nano")
		}
		return c.PlatformWindows, version, c.PlatformFamilyWindows, nil
	}
}

func isWindowsNano(operatingSystemSKU uint32) bool {
	return operatingSystemSKU == c.SKUProductStandardNanoServer ||
		operatingSystemSKU == c.SKUProductDatacenterNanoServer
}
