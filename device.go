package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Device struct {
	path    string
	Name    string `plist:"name"`
	Runtime string `plist:"runtime"`
}

func (d Device) String() string {
	r := strings.Split(d.Runtime, ".")
	return fmt.Sprintf("%v(%v)", d.Name, r[len(r)-1])
}

func readDevice(deviceDir string) (*Device, error) {
	plistPath := filepath.Join(deviceDir, "device.plist")
	device := &Device{path: deviceDir}
	err := readPlist(plistPath, device)
	return device, err
}
