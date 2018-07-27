package main

import (
	"fmt"
)

type App struct {
	path             string
	device           *Device
	BundleIdentifier string `plist:"CFBundleIdentifier"`
	Name             string `plist:"CFBundleName"`
	Version          string `plist:"CFBundleVersion"`
}

func (a App) String() string {
	return fmt.Sprintf("%v %v %v", a.Name, a.Version, a.device)
}

func readApp(device *Device, appDir string) (*App, error) {
	plistPath, err := infoPlistPath(appDir)
	if err != nil {
		return nil, err
	}
	app := &App{path: appDir, device: device}
	err = readPlist(plistPath, app)
	return app, err
}

type Apps []*App

func (a Apps) Len() int {
	return len(a)
}

func (a Apps) Less(i, j int) bool {
	l := a[i]
	r := a[j]
	if l.Name != r.Name {
		return l.Name < r.Name
	}
	if l.BundleIdentifier != r.BundleIdentifier {
		return l.BundleIdentifier < r.BundleIdentifier
	}
	if l.Version != r.Version {
		return l.Version < r.Version
	}
	if l.device != r.device {
		return l.device.String() < r.device.String()
	}
	return false
}

func (a Apps) Swap(i, j int) {
	t := a[i]
	a[i] = a[j]
	a[j] = t
}
