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
kubecm import -n dev_129_cluster -l /tmp/configs/config_dev_182_cluster
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