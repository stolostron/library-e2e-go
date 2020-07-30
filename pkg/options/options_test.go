package options

import (
	"strings"
	"testing"
)

func TestLoadOptions(t *testing.T) {
	type args struct {
		optionsFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				optionsFile: "../../test/unit/resources/options.yaml",
			},
			wantErr: false,
		},
		{
			name: "failed not found",
			args: args{
				optionsFile: "wrongPath",
			},
			wantErr: true,
		},
		{
			name: "failed malformed",
			args: args{
				optionsFile: "../../test/unit/resources/options-malformed.yaml",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadOptions(tt.args.optionsFile); (err != nil) != tt.wantErr {
				t.Errorf("LoadOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetClusterName(t *testing.T) {
	err := LoadOptions("../../test/unit/resources/options.yaml")
	if err != nil {
		t.Error(err)
	}
	type args struct {
		cloud string
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
				cloud: "azure",
			},
			want:    "my-azure-",
			wantErr: false,
		},
		{
			name: "failed unsupported",
			args: args{
				cloud: "mycloud",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetClusterName(tt.args.cloud)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClusterName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.HasPrefix(got, tt.want) {
				t.Errorf("GetClusterName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRegion(t *testing.T) {
	err := LoadOptions("../../test/unit/resources/options.yaml")
	if err != nil {
		t.Error(err)
	}
	type args struct {
		cloud string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "succeed azure",
			args: args{
				cloud: "azure",
			},
			want:    "my_azure_region",
			wantErr: false,
		},
		{
			name: "succeed aws",
			args: args{
				cloud: "aws",
			},
			want:    "my_aws_region",
			wantErr: false,
		},
		{
			name: "succeed gcp",
			args: args{
				cloud: "gcp",
			},
			want:    "my_gcp_region",
			wantErr: false,
		},
		{
			name: "failed not supported",
			args: args{
				cloud: "mycloud",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRegion(tt.args.cloud)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRegion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetRegion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBaseDomain(t *testing.T) {
	err := LoadOptions("../../test/unit/resources/options.yaml")
	if err != nil {
		t.Error(err)
	}
	type args struct {
		cloud string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "succeed azure",
			args: args{
				cloud: "azure",
			},
			want:    "my_azure_baseDnsDomain",
			wantErr: false,
		},
		{
			name: "succeed aws",
			args: args{
				cloud: "aws",
			},
			want:    "my_aws_baseDnsDomain",
			wantErr: false,
		},
		{
			name: "succeed gcp",
			args: args{
				cloud: "gcp",
			},
			want:    "my_gcp_baseDnsDomain",
			wantErr: false,
		},
		{
			name: "failed not supported",
			args: args{
				cloud: "mycloud",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBaseDomain(tt.args.cloud)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBaseDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBaseDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringWithCharset(t *testing.T) {
	type args struct {
		length  int
		charset string
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
				length:  4,
				charset: "1234567",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := StringWithCharset(tt.args.length, tt.args.charset)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringWithCharset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
