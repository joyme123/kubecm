package loader

import (
    "io/ioutil"
    "strings"
)

func IsLocal(path string) bool {
    if strings.HasPrefix(path, "file://") {
        return true
    } else if !strings.Contains(path, "://") {
        return true
    }

    return false
}

func LocalGet(path string) ([]byte, error) {
    path = strings.TrimLeft(path, "file://")
    return ioutil.ReadFile(path)
}
