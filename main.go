package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"howett.net/plist"
)

func main() {
	devices, err := deviceDirectories()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	for _, device := range devices {
		d, err := deviceInfo(device)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		apps, err := appDirectories(device)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		for _, app := range apps {
			printApp(app, d, os.Stdout)
		}
	}
}

type App struct {
	Name    string `plist:"CFBundleDisplayName"`
	ID      string `plist:"CFBundleIdentifier"`
	Version string `plist:"CFBundleShortVersionString"`
}

func (a App) String() string {
	return fmt.Sprintf("%v\t\t%v\t%v", a.Name, a.ID, a.Version)
}

func printApp(appDir string, device *Device, out *os.File) error {
	app, err := appInfo(appDir)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(out, "%v\t%v\t\t%v\n", app, device, appDir)
	return err
}

func appInfo(appDir string) (*App, error) {
	dirs, err := ioutil.ReadDir(appDir)
	if err != nil {
		return nil, err
	}

	appBundleDir := ""
	for _, dir := range dirs {
		if dir.IsDir() && strings.HasSuffix(dir.Name(), ".app") {
			appBundleDir = filepath.Join(appDir, dir.Name())
			break
		}
	}
	if appBundleDir == "" {
		return nil, errors.New("app bundle not found")
	}

	plistPath := filepath.Join(appBundleDir, "Info.plist")
	file, err := os.Open(plistPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var app App
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	_, err = plist.Unmarshal(data, &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func appDirectories(deviceDir string) ([]string, error) {
	appsRootDir := filepath.Join(deviceDir, "data/Containers/Bundle/Application")
	appDirs, err := ioutil.ReadDir(appsRootDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	result := []string{}
	for _, ff := range appDirs {
		if ff.IsDir() {
			path := filepath.Join(appsRootDir, ff.Name())
			result = append(result, path)
		}
	}

	return result, nil
}

type Device struct {
	Name    string `plist:"name"`
	Runtime string `plist:"runtime"`
}

func (d Device) String() string {
	r := strings.Split(d.Runtime, ".")
	return fmt.Sprintf("%v\t%v", d.Name, r[len(r)-1])
}

func deviceInfo(deviceDir string) (*Device, error) {
	plistPath := filepath.Join(deviceDir, "device.plist")
	file, err := os.Open(plistPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var device Device
	_, err = plist.Unmarshal(data, &device)
	if err != nil {
		return nil, err
	}
	return &device, nil
}

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
