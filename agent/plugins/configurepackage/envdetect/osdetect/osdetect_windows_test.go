//go:build windows
// +build windows

package osdetect

import (
	"fmt"
	"testing"

	"github.com/aws/amazon-ssm-agent/agent/mocks/log"
	"github.com/aws/amazon-ssm-agent/agent/platform"
	c "github.com/aws/amazon-ssm-agent/agent/plugins/configurepackage/envdetect/constants"
	"github.com/stretchr/testify/assert"
)

func TestDetectPkgManager(t *testing.T) {
	result, err := DetectPkgManager("", "", "") // parameters only matter for linux

	assert.NoError(t, err)
	assert.Equal(t, c.PackageManagerWindows, result)
}

func TestDetectInitSystem(t *testing.T) {
	result, err := DetectInitSystem()

	assert.NoError(t, err)
	assert.Equal(t, c.InitWindows, result)
}

func TestDetectPlatform(t *testing.T) {
	tests := []struct {
		name            string
		getOSInfo       func(platform.Win32_OperatingSystem) (platform.Win32_OperatingSystem, error)
		expectedVersion string
	}{
		{
			name: "WMI data is empty",
			getOSInfo: func(platform.Win32_OperatingSystem) (platform.Win32_OperatingSystem, error) {
				return platform.Win32_OperatingSystem{}, nil
			},
			expectedVersion: "",
		},
		{
			name: "WMI throws an error",
			getOSInfo: func(platform.Win32_OperatingSystem) (platform.Win32_OperatingSystem, error) {
				return platform.Win32_OperatingSystem{Version: "10"}, fmt.Errorf("Error while fetching WMI data")
			},
			expectedVersion: "",
		},
		{
			name: "Windows Nano SKU",
			getOSInfo: func(platform.Win32_OperatingSystem) (platform.Win32_OperatingSystem, error) {
				return platform.Win32_OperatingSystem{Version: "10", OperatingSystemSKU: 144}, nil
			},
			expectedVersion: "10nano",
		},
		{
			name: "Regular Windows SKU",
			getOSInfo: func(platform.Win32_OperatingSystem) (platform.Win32_OperatingSystem, error) {
				return platform.Win32_OperatingSystem{Version: "10", OperatingSystemSKU: 100}, nil
			},
			expectedVersion: "10",
		},
		{
			name: "No Windows SKU data",
			getOSInfo: func(platform.Win32_OperatingSystem) (platform.Win32_OperatingSystem, error) {
				return platform.Win32_OperatingSystem{Version: "10"}, nil
			},
			expectedVersion: "10",
		},
	}

	for _, d := range tests {
		t.Run(d.name, func(t *testing.T) {
			getOSInfo = d.getOSInfo
			platform, version, platformFamily, err := DetectPlatform(log.NewMockLog())
			assert.Equal(t, c.PlatformWindows, platform)
			assert.Equal(t, c.PlatformFamilyWindows, platformFamily)
			assert.Equal(t, d.expectedVersion, version)
			assert.NoError(t, err)
		})
	}
}
