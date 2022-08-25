package util

import (
	"reflect"
	"testing"
)

func TestReplaceApiServerAddr(t *testing.T) {

	testConfig := `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: xxxx
    server: https://127.0.0.1:6443
  name: cluster.local
contexts:
- context:
    cluster: cluster.local
    user: kubernetes-admin
  name: kubernetes-admin@cluster.local
current-context: kubernetes-admin@cluster.local
kind: Config
preferences: {}
users:
- name: kubernetes-admin
  user:
    client-certificate-data: xxxx 
    client-key-data: xxxx`

	wantConfig := `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: xxxx
    server: https://1.1.1.1:6443
  name: cluster.local
contexts:
- context:
    cluster: cluster.local
    user: kubernetes-admin
  name: kubernetes-admin@cluster.local
current-context: kubernetes-admin@cluster.local
kind: Config
preferences: {}
users:
- name: kubernetes-admin
  user:
    client-certificate-data: xxxx
    client-key-data: xxxx
`

	type args struct {
		data []byte
		addr string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "case 1",
			args: args{
				data: []byte(testConfig),
				addr: "https://1.1.1.1:6443",
			},
			want:    []byte(wantConfig),
			wantErr: false,
		},
		{
			name: "case 2",
			args: args{
				data: []byte(testConfig),
				addr: "1.1.1.1",
			},
			want:    []byte(wantConfig),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReplaceApiServerAddr(tt.args.data, tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceApiServerAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReplaceApiServerAddr() got = %v, want %v", got, tt.want)
			}
		})
	}
}
