package main

import (
	"errors"
	"io/ioutil"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

func deviceDirectories() ([]string, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	devicesRoot := filepath.Join(user.HomeDir, "Library/Developer/CoreSimulator/Devices")
	dirInfos, err := ioutil.ReadDir(devicesRoot)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, ff := range dirInfos {
		if ff.IsDir() {
			path := filepath.Join(devicesRoot, ff.Name())
			result = append(result, path)
		}
	}
	return result, nil
}

func appDirectories(deviceDir string) ([]string, error) {
	appsRoot := filepath.Join(deviceDir, "data/Containers/Bundle/Application")
	dirInfos, err := ioutil.ReadDir(appsRoot)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, ff := range dirInfos {
		if ff.IsDir() {
			path := filepath.Join(appsRoot, ff.Name())
			result = append(result, path)
		}
	}
	return result, nil
}

func infoPlistPath(appDir string) (string, error) {
	infos, err := ioutil.ReadDir(appDir)
	if err != nil {
		return "", err
	}

	appBundleDir := ""
	for _, ff := range infos {
		if ff.IsDir() && filepath.Ext(ff.Name()) == ".app" {
			appBundleDir = filepath.Join(appDir, ff.Name())
			break
		}
	}

	if appBundleDir == "" {
		return "", errors.New("app bundle not found")
	}

	return filepath.Join(appBundleDir, "Info.plist"), nil
}

func dataDirectory(deviceDir string, appBundleId string) (string, error) {
	deviceUUID := filepath.Base(deviceDir)
	cmd := exec.Command("xcrun", "simctl", "get_app_container", deviceUUID, appBundleId, "data")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
