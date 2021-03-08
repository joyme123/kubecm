package manager

import "fmt"

var NameConflictError = fmt.Errorf("error conflict")
var NameNotExistError = fmt.Errorf("name doesn't exist")
var ConfigFileIsUsing = fmt.Errorf("config file is using now, can't be removed")
