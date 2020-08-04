package options

import (
	"strings"
	"testing"
)

func TestNewClusterName(t *testing.T) {
	err := LoadOptions("../../test/unit/resources/options.yaml")
	if err != nil {
		t.Error(err)
	}
	type args struct {
		owner string
		cloud string
	}
	tests := []struct {
		name    string
		args    args
		want    *ClusterName
		wantErr bool
	}{
		{
			name: "succeed",
			args: args{
				owner: "owner",
				cloud: "azure",
			},
			want: &ClusterName{
				Owner: "owner",
				Cloud: "azure",
			},
			wantErr: false,
		},
		{
			name: "failed unsupported",
			args: args{
				owner: "owner",
				cloud: "mycloud",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClusterName(tt.args.owner, tt.args.cloud)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClusterName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !strings.HasPrefix(got.String(), tt.want.String()) {
					t.Errorf("NewClusterName() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
