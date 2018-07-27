package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	flag.Parse()
	appIndex := -1
	if flag.NArg() > 0 {
		n, err := strconv.Atoi(flag.Arg(0))
		if err == nil {
			appIndex = n
		}
	}

	apps := sortedApps()

	if appIndex == -1 {
		for i, app := range apps {
			fmt.Println(i, "-", app)
		}
	} else {
		printDataDir(apps[appIndex])
	}
}

func sortedApps() Apps {
	apps := Apps([]*App{})

	devices, err := deviceDirectories()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	for _, deviceDir := range devices {
		device, err := readDevice(deviceDir)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		appDirs, _ := appDirectories(deviceDir)
		for _, appDir := range appDirs {
			app, err := readApp(device, appDir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			apps = append(apps, app)
		}
	}

	sort.Sort(apps)

	return apps
}

func printDataDir(app *App) {
	dir, err := dataDirectory(app.device.path, app.BundleIdentifier)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	fmt.Println(dir)
}
