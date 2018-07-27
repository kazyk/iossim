package main

import (
	"os"

	"howett.net/plist"
)

func readPlist(path string, v interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := plist.NewDecoder(file)
	err = decoder.Decode(v)
	return err
}
