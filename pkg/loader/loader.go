package loader

import (
	"fmt"
	"github.com/joyme123/kubecm/pkg/types"
)

func Load(syncInfo *types.Sync) ([]byte, error) {
	if isLocal(syncInfo.From) {
		return localGet(syncInfo.From)
	} else if isSSH(syncInfo.From) {
		if len(syncInfo.Password) > 0 {
			return sshGetWithPassword(syncInfo.From, syncInfo.SSHPort, syncInfo.Password)
		} else {
			return sshGetWithPrivateKey(syncInfo.From, syncInfo.SSHPort, syncInfo.PrivateKey)
		}
	} else {
		return nil, fmt.Errorf("unsupport path: %s", syncInfo.From)
	}
}
