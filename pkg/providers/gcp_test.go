package providers

import (
	"strings"
	"testing"
)

func TestGetInstallConfigGCP(t *testing.T) {
	const config = `
apiVersion: v1
metadata:
  name: name
baseDomain: basednsdomain
controlPlane:
  hyperthreading: Enabled
  name: master
  replicas: 3
  platform:
    gcp:
      type: n1-standard-4
compute:
- hyperthreading: Enabled
  name: worker
  replicas: 3
  platform:
      gcp:
        type: n1-standard-4
networking:
  clusterNetwork:
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  machineCIDR: 10.0.0.0/16
  networkType: OpenShiftSDN
  serviceNetwork:
  - 172.30.0.0/16
platform:
  gcp:
    projectID: projectid
    region: region
pullSecret: ""
sshKey: sshkey
`
	type args struct {
		instConfig InstallerConfigGCP
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "succeed",
			args: args{
				instConfig: InstallerConfigGCP{
					Name:          "name",
					BaseDnsDomain: "basednsdomain",
					Region:        "region",
					SSHKey:        "sshkey",
					ProjectID:     "projectid",
				},
			},
			want:    config,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInstallConfigGCP(tt.args.instConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInstallConfigGCP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if strings.ReplaceAll("\t", got, got) != strings.ReplaceAll("\t", tt.want, tt.want) {
				t.Errorf("GetInstallConfigGCP() = %v, want %v", got, tt.want)
			}
		})
	}
}
