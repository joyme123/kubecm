package util

import (
	"fmt"
	"io/ioutil"
	"os"
)

func Copy(src, dst string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, data, 0755)
}

func Move(src, dst string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(dst, data, 0755); err != nil {
		return err
	}
	return os.Remove(src)
}

func EnsureDir(dir string) error {
	f, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	if f.IsDir() {
		return nil
	}
	return fmt.Errorf("%s is not a dir", dir)
}

func EnsureFile(file string) error {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		_, err := os.Create(file)
		return err
	}
	return nil
}

func GetHomeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "/"
	}

	return dir
}
