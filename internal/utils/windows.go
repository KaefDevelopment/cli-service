//go:build windows

package utils

import (
	"log"
	"syscall"
)

func MakeHiddenConfigFolder(path string) error {
	fileNameW, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := syscall.SetFileAttributes(fileNameW, syscall.FILE_ATTRIBUTE_HIDDEN); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
