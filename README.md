# kubecm

![resource/term.gif](resource/term.gif)

kubecm(kube config manager) is used for manager config file of kubernetes.

## Install

```bash
go get github.com/joyme123/kubecm
```

## Usage


list config files
```bash
kubecm list
```

import config file
```bash
# import fro local filesystem
kubecm import -n dev_129_cluster -l /tmp/configs/config_dev_182_cluster

# import via ssh with password 
kubecm import dev_0_101_cluster --from=ssh://root@192.168.0.101:/etc/kubernetes/kubectl.kubeconfig  -p mypassword

# import via ssh with key, default from $HOME/.ssh/id_rsa
kubecm import dev_0_102_cluster --from=ssh://root@192.168.0.102:/etc/kubernetes/kubectl.kubeconfig 
```

use config file
```bash
kubecm use -n dev_129_cluster
```

rename config file
```bash
kubecm rename -n dev_129_cluster -t dev_cluster
```

remove config file
```bash
kubecm remove -n dev_129_cluster
```

## ZSH 配置

zsh 命令行提示:

在对应的主题里添加下面的配置
```bash
kubecm_prompt() {
  echo " %{$fg[green]%}k8s:$(kubecm list -c)%{$reset_color%}"
}
local kubecm='$(kubecm_prompt)'
```