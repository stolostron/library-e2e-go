package providers

import (
	"strings"
	"testing"
)

func TestGetInstallConfigAzure(t *testing.T) {
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
    azure:
      osDisk:
      diskSizeGB: 128
    type:  Standard_D4s_v3
compute:
- hyperthreading: Enabled
  name: worker
  replicas: 3
  platform:
    azure:
      type:  Standard_D2s_v3
      osDisk:
      diskSizeGB: 128
      zones:
      - "1"
      - "2"
      - "3"
networking:
  clusterNetwork:
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  machineCIDR: 10.0.0.0/16
  networkType: OpenShiftSDN
  serviceNetwork:
  - 172.30.0.0/16
platform:
  azure:
    baseDomainResourceGroupName: basedomainrgn
    region: {region
pullSecret: "" # skip, hive will inject based on it's secrets
sshKey: sshkey 
`
	type args struct {
		instConfig InstallerConfigAzure
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
				instConfig: InstallerConfigAzure{
					Name:          "name",
					BaseDnsDomain: "basednsdomain",
					Region:        "region",
					SSHKey:        "sshkey",
					BaseDomainRGN: "basedomainrgn",
				},
			},
			want:    config,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInstallConfigAzure(tt.args.instConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInstallConfigGCP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if strings.ReplaceAll("\t", got, got) != strings.ReplaceAll("\t", tt.want, tt.want) {
				t.Errorf("GetInstallConfigAzure() = %v, want %v", got, tt.want)
			}
		})
	}
}
