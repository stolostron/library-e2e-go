package options

import "testing"

func TestGetHubKubeConfig(t *testing.T) {
	type args struct {
		configDir string
		scenario  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: args{
				configDir: "../../test/unit/resources/hubs",
				scenario:  "fake-scenario1",
			},
			want: "../../test/unit/resources/hubs/fake-scenario1/kubeconfig.yaml",
		},
		{
			name: "success no config Dir",
			args: args{
				configDir: "",
				scenario:  "fake-scenario1",
			},
			want: "",
		},
		{
			name: "success no scenario",
			args: args{
				configDir: "../../test/unit/resources/hubs",
				scenario:  "",
			},
			want: "",
		},
		{
			name: "success wrongPath",
			args: args{
				configDir: "wrongPath",
				scenario:  "fake-scenario1",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHubKubeConfig(tt.args.configDir, tt.args.scenario); got != tt.want {
				t.Errorf("GetHubKubeConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
