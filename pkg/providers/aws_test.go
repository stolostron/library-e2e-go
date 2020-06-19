package providers

import (
	"strings"
	"testing"
)

func TestGetInstallConfigAWS(t *testing.T) {
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
    aws:
      rootVolume:
        iops: 4000
        size: 500
        type: io1
      type: m4.xlarge
compute:
- hyperthreading: Enabled
  name: worker
  replicas: 3
  platform:
    aws:
    rootVolume:
      iops: 2000
      size: 500
      type: io1
    type: m4.large
networking:
  clusterNetwork:
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  machineCIDR: 10.0.0.0/16
  networkType: OpenShiftSDN
  serviceNetwork:
  - 172.30.0.0/16
platform:
  aws:
    region: region
pullSecret: ""
sshKey: sshkey
`

	type args struct {
		instConfig InstallerConfigAWS
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
				instConfig: InstallerConfigAWS{
					Name:          "name",
					BaseDnsDomain: "basednsdomain",
					Region:        "region",
					SSHKey:        "sshkey",
				},
			},
			want:    config,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInstallConfigAWS(tt.args.instConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInstallConfigGCP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if strings.ReplaceAll("\t", got, got) != strings.ReplaceAll("\t", tt.want, tt.want) {
				t.Errorf("GetInstallConfigAWS() = %v, want %v", got, tt.want)
			}
		})
	}
}
