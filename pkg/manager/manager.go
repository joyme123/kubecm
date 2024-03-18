package manager

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"

	"github.com/joyme123/kubecm/pkg/loader"
	"github.com/joyme123/kubecm/pkg/types"
	"github.com/joyme123/kubecm/pkg/util"

	"github.com/ghodss/yaml"
	"github.com/google/uuid"
	"k8s.io/klog/v2"
)

type Interface interface {
	List() (*types.Configuration, error)
	Import(name string, configData []byte, override bool, apiServerAddr string) error
	Remove(name string) error
	Rename(src string, dst string) error
	Use(name string, currentSession bool) error
	SaveSyncInfo(name string, syncInfo *types.Sync) error
	Sync(name string) map[string]error
	Save(filePath string) error
	Load(filePath string) error
}

type impl struct {
	m sync.Mutex

	configPath string
	configDir  string
	kubePath   string
	conf       *types.Configuration
}

func NewInterface(configDir string, configPath string, kubePath string) (Interface, error) {
	i := &impl{
		configDir:  configDir,
		configPath: configPath,
		kubePath:   kubePath,
		conf:       &types.Configuration{},
	}

	startTime := time.Now()
	err := i.init()
	if err != nil {
		return nil, err
	}
	cost := time.Since(startTime)
	klog.V(4).Infof("init cost: %d ms", cost.Milliseconds())
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
	err = yaml.Unmarshal(data, i.conf)
	if err != nil {
		return nil
	}

	return nil
}

func (i *impl) List() (*types.Configuration, error) {
	return i.conf, nil
}

func (i *impl) Import(name string, configData []byte, override bool, apiServerAddr string) error {
	i.m.Lock()
	defer i.m.Unlock()

	index := i.search(name)
	if !override && index >= 0 {
		return NameConflictError
	}
	return i.importByIndex(index, name, configData, apiServerAddr)
}

// if index < 0, config will append
func (i *impl) importByIndex(index int, name string, configData []byte, apiServerAddr string) error {
	var newLocation string
	if index >= 0 {
		newLocation = i.conf.Items[index].Location
	} else {
		id := uuid.New().String()
		newLocation = path.Join(i.configDir, id)
	}

	if apiServerAddr != "" {
		var err error
		configData, err = util.ReplaceApiServerAddr(configData, apiServerAddr)
		if err != nil {
			return err
		}
	}

	if err := util.Copy(configData, newLocation); err != nil {
		return err
	}

	if index >= 0 {
		i.conf.Items[index].TimeStamp = time.Now()
	} else {
		item := types.ConfigItem{
			Name:      name,
			Location:  newLocation,
			TimeStamp: time.Now(),
		}
		i.conf.Items = append(i.conf.Items, item)
	}
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

func (i *impl) Use(name string, currentSession bool) error {
	i.m.Lock()
	defer i.m.Unlock()

	index := i.search(name)
	if index == -1 {
		return NameNotExistError
	}

	i.conf.Current = name
	item := i.conf.Items[index]
	if currentSession {
		_ = os.Setenv("KUBECONFIG", item.Location)
		_ = os.Setenv("KUBECM_SHELL", "true")
		shell := os.Getenv("SHELL")
		cmd := exec.Command(shell)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	if os.Getenv("KUBECM_SHELL") == "true" && os.Getenv("KUBECONFIG") != "" {
		fmt.Printf("You can't change system kubeconfig after use `kubecm use xxx -c`. \nPlease execute cmd `exit && kubecm use %s` \n", name)
		return nil
	}

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

func (i *impl) SaveSyncInfo(name string, syncInfo *types.Sync) error {
	k := i.search(name)
	if k < 0 {
		return nil
	}
	i.conf.Items[k].Sync = syncInfo
	return i.write()
}

// Sync sync target config. if name is empty, sync all
func (i *impl) Sync(name string) map[string]error {
	syncFunc := func(index int, name string, info *types.Sync) error {
		if info == nil {
			return nil
		}
		data, err := loader.Load(info)
		if err != nil {
			return err
		}
		return i.importByIndex(index, name, data, "")
	}

	res := map[string]error{}

	if len(name) == 0 {
		for index, item := range i.conf.Items {
			if item.Sync == nil {
				continue
			}
			res[item.Name] = syncFunc(index, item.Name, item.Sync)
		}
		return res
	}

	index := i.search(name)
	if index < 0 {
		res[name] = fmt.Errorf("target name not exist")
		return res
	}
	item := i.conf.Items[index]
	res[name] = syncFunc(index, item.Name, item.Sync)
	return res
}

func (i *impl) Save(filePath string) error {
	_, err := os.Stat(filePath)
	if err == nil {
		return fmt.Errorf("file %s is already exist", filePath)
	}

	if !os.IsNotExist(err) {
		return err
	}

	var savedConfig types.SavedConfig
	savedConfig.Current = i.conf.Current

	for _, item := range i.conf.Items {
		data, err := ioutil.ReadFile(item.Location)
		if err != nil {
			fmt.Printf("kubeconfig %s is missing on location: %s", item.Name, item.Location)
			continue
		}
		savedConfig.Items = append(savedConfig.Items, types.SavedConfigItem{
			Name:    item.Name,
			Content: string(data),
		})
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	confData, err := yaml.Marshal(savedConfig)
	if err != nil {
		return err
	}
	_, err = f.Write(confData)
	if err != nil {
		return err
	}

	return nil
}

func (i *impl) Load(filePath string) error {
	confData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var savedConf types.SavedConfig
	if err := yaml.Unmarshal(confData, &savedConf); err != nil {
		return err
	}

	// merge with exist configurations
	exists := make(map[string]struct{})
	for _, item := range i.conf.Items {
		exists[item.Name] = struct{}{}
	}

	for _, item := range savedConf.Items {
		if _, ok := exists[item.Name]; ok {
			fmt.Printf("%s exists, skip\n", item.Name)
			continue
		}
		if err := i.Import(item.Name, []byte(item.Content), false, ""); err != nil {
			fmt.Printf("import %s failed, err: %v\n", item.Name, err)
			continue
		} else {
			fmt.Printf("import %s completed\n", item.Name)
		}
	}

	if i.conf.Current == "" {
		i.conf.Current = savedConf.Current
	}

	return nil
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
	data, err := yaml.Marshal(i.conf)
	if err != nil {
		return fmt.Errorf("yaml marshal error: %v", err)
	}

	err = ioutil.WriteFile(i.configPath, data, 0755)
	if err != nil {
		return fmt.Errorf("write file to %s error: %v", i.configPath, err)
	}
	return nil
}
