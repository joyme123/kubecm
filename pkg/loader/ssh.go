package loader

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func IsSSH(path string) bool {
	if strings.HasPrefix(path, "ssh://") {
		return true
	}
	return false
}

func DefaultSSHKeyPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir error: %v", err)
	}
	privateKeyPath := path.Join(dir, ".ssh", "id_rsa")

	return privateKeyPath, nil
}

func SSHGetWithPassword(path string, port int, password string) ([]byte, error) {
	params, err := getParams(path)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: params.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		BannerCallback:  ssh.BannerDisplayStderr(),
	}

	return sshGet(config, params.ip, port, params.path)
}

func SSHGetWithPrivateKey(filepath string, port int, privateKeyPath string) ([]byte, error) {
	params, err := getParams(filepath)
	if err != nil {
		return nil, err
	}

	privateKeyData, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("get ssh private key error: %v", err)
	}
	signer, err := ssh.ParsePrivateKey(privateKeyData)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: params.user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return sshGet(config, params.ip, port, params.path)
}

type sshParams struct {
	user string
	ip   string
	path string
}

func getParams(path string) (*sshParams, error) {
	path = strings.TrimLeft(path, "ssh://")
	strs := strings.SplitN(path, "@", 2)
	if len(strs) != 2 {
		return nil, fmt.Errorf("ssh path format error")
	}
	user := strs[0]
	strs = strings.SplitN(strs[1], ":", 2)
	if len(strs) != 2 {
		return nil, fmt.Errorf("ssh path format error")
	}
	ip := strs[0]
	path = strs[1]

	return &sshParams{user, ip, path}, nil
}

func sshGet(config *ssh.ClientConfig, ip string, port int, path string) ([]byte, error) {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), config)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("cat " + path); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
