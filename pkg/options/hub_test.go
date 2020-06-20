package options

import "testing"

func TestGetHubKubeConfig(t *testing.T) {
	type args struct {
		configDir string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: args{
				configDir: "../../test/unit/resources/hub",
			},
			want: "../../test/unit/resources/hub/kubeconfig.yaml",
		},
		{
			name: "success no config Dir",
			args: args{
				configDir: "",
			},
			want: "",
		},
		{
			name: "success wrongPath",
			args: args{
				configDir: "wrongPath",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHubKubeConfig(tt.args.configDir); got != tt.want {
				t.Errorf("GetHubKubeConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
