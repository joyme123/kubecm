package types

import (
	"fmt"
	"os"
	"time"

	"github.com/gookit/color"
)

type Configuration struct {
	Current string       `json:"current" yaml:"current"`
	Items   []ConfigItem `json:"items" yaml:"items"`
}

type ConfigItem struct {
	Name      string    `json:"name" yaml:"name"`
	Location  string    `json:"location" yaml:"location"`
	TimeStamp time.Time `json:"timestamp" yaml:"timestamp"`
	Sync      *Sync     `json:"sync,omitempty" yaml:"sync,omitempty"`
}

type Sync struct {
	From       string `json:"from,omitempty" yaml:"from,omitempty"`
	SSHPort    int    `json:"sshPort,omitempty" yaml:"sshPort,omitempty"`
	Password   string `json:"password,omitempty" yaml:"password,omitempty"`
	PrivateKey string `json:"privateKey,omitempty" yaml:"privateKey,omitempty"`
}

func (c *Configuration) Print(onlyCurrent bool) {
	if c == nil {
		fmt.Println("config is empty")
		return
	}

	kubeConfigPath := os.Getenv("KUBECONFIG")
	current := c.Current

	if kubeConfigPath != "" {
		for _, item := range c.Items {
			if item.Location == kubeConfigPath {
				current = item.Name
				break
			}
		}
	}

	for _, item := range c.Items {
		if item.Name == current {
			if onlyCurrent {
				color.Println(item.Name)
				break
			} else {
				color.Success.Println(item.Name + "*")
			}
		} else if !onlyCurrent {
			fmt.Println(item.Name)
		}
	}
}
