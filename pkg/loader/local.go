package loader

import (
	"io/ioutil"
	"strings"
)

func isLocal(path string) bool {
	if strings.HasPrefix(path, "file://") {
		return true
	} else if !strings.Contains(path, "://") {
		return true
	}

	return false
}

func localGet(path string) ([]byte, error) {
	path = strings.TrimLeft(path, "file://")
	return ioutil.ReadFile(path)
}
