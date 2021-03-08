package manager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"

	"github.com/joyme123/kubecm/pkg/util"

	"github.com/google/uuid"
)

type Interface interface {
	List() (*Configuration, error)
	Import(name string, location string) error
	Remove(name string) error
	Rename(src string, dst string) error
	Use(name string) error
}

type impl struct {
	m sync.Mutex

	configPath string
	configDir  string
	kubePath   string
	conf       *Configuration
}

func NewInterface(configDir string, configPath string, kubePath string) (Interface, error) {
	i := &impl{
		configDir:  configDir,
		configPath: configPath,
		kubePath:   kubePath,
		conf:       &Configuration{},
	}

	err := i.init()
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (i *impl) init() error {
	kubeDir := path.Dir(i.kubePath)
	dirs := []string{i.configDir, kubeDir}
	for _, dir := range dirs {
		if err := util.EnsureDir(dir); err != nil {
			return err
		}
	}
	err := util.EnsureFile(i.configPath)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(i.configPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, i.conf)
	if err != nil {
		return nil
	}

	return nil
}

func (i *impl) List() (*Configuration, error) {
	return i.conf, nil
}

func (i *impl) Import(name string, location string) error {
	i.m.Lock()
	defer i.m.Unlock()

	_, err := os.Open(location)
	if err != nil {
		return fmt.Errorf("open location %s error: %v", location, err)
	}

	index := i.search(name)
	if index >= 0 {
		return NameConflictError
	}

	id := uuid.New().String()
	newLocation := path.Join(i.configDir, id)
	if err := util.Copy(location, newLocation); err != nil {
		return err
	}

	item := ConfigItem{
		Name:      name,
		Location:  newLocation,
		TimeStamp: time.Now(),
	}
	i.conf.Items = append(i.conf.Items, item)
	return i.write()
}

func (i *impl) Remove(name string) error {
	i.m.Lock()
	defer i.m.Unlock()

	index := i.search(name)
	if index == -1 {
		return NameNotExistError
	}

	removeItem := i.conf.Items[index]
	if removeItem.Name == i.conf.Current {
		return ConfigFileIsUsing
	}
	removeFile := removeItem.Location
	if err := os.Remove(removeFile); err != nil && !os.IsNotExist(err) {
		return err
	}

	i.conf.Items = append(i.conf.Items[0:index], i.conf.Items[index+1:]...)

	return i.write()
}

func (i *impl) Rename(srcName string, dstName string) error {
	i.m.Lock()
	defer i.m.Unlock()

	index := i.search(srcName)

	if index == -1 {
		return NameNotExistError
	}
	i.conf.Items[index].Name = dstName

	if i.conf.Current == srcName {
		i.conf.Current = dstName
	}
	return i.write()
}

func (i *impl) Use(name string) error {
	i.m.Lock()
	defer i.m.Unlock()

	index := i.search(name)
	if index == -1 {
		return NameNotExistError
	}

	i.conf.Current = name
	item := i.conf.Items[index]

	// create symbolic link from config file to kube config file
	kubefile, err := os.Lstat(i.kubePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err == nil {
		if kubefile.Mode()&os.ModeSymlink != 0 {
			// symlink, remove it
			if err := os.Remove(i.kubePath); err != nil {
				return err
			}
		} else {
			// kube config file, backup it
			if err := util.Move(i.kubePath, i.kubePath+"_kubecm_backup"); err != nil {
				return err
			}
		}
	}

	// create symbolic file
	err = os.Symlink(item.Location, i.kubePath)
	if err != nil {
		return fmt.Errorf("create symlink from %s to %s error: %v", item.Location, i.kubePath, err)
	}

	return i.write()
}

func (i *impl) search(name string) int {
	for k, item := range i.conf.Items {
		if item.Name == name {
			return k
		}
	}

	return -1
}

func (i *impl) write() error {
	data, err := json.Marshal(i.conf)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", err)
	}

	err = ioutil.WriteFile(i.configPath, data, 0755)
	if err != nil {
		return fmt.Errorf("write file to %s error: %v", i.configPath, err)
	}
	return nil
}
