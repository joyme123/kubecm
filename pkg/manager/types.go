package manager

import (
	"fmt"
	"time"

	"github.com/gookit/color"
)

type Configuration struct {
	Current string       `json:"current"`
	Items   []ConfigItem `json:"items"`
}

type ConfigItem struct {
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	TimeStamp time.Time `json:"timestamp"`
}

func (c *Configuration) Print() {
	if c == nil {
		fmt.Println("config is empty")
		return
	}
	for _, item := range c.Items {
		if item.Name == c.Current {
			color.Success.Println(item.Name + "*")
		} else {
			fmt.Println(item.Name)
		}
	}
}
