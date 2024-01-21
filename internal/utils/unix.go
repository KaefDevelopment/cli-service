//go:build !windows

package utils

func MakeHiddenConfigFolder(path string) error {
	return nil // plug for Unix systems
}
