package util

import (
	"net/url"
	"strings"

	"k8s.io/client-go/tools/clientcmd"
)

// ReplaceApiServerAddr ...
func ReplaceApiServerAddr(data []byte, addr string) ([]byte, error) {
	config, err := clientcmd.NewClientConfigFromBytes(data)
	if err != nil {
		return nil, err
	}
	rawConfig, err := config.RawConfig()
	if err != nil {
		return nil, err
	}

	for k := range rawConfig.Clusters {
		if strings.HasPrefix(addr, "http") {
			rawConfig.Clusters[k].Server = addr
		} else {
			httpAddr := rawConfig.Clusters[k].Server
			u, err := url.Parse(httpAddr)
			if err != nil {
				return nil, err
			}
			u.Host = addr + ":" + u.Port()
			rawConfig.Clusters[k].Server = u.String()
		}
	}
	return clientcmd.Write(rawConfig)
}
