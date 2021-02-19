module github.com/open-cluster-management/library-e2e-go

go 1.13

replace k8s.io/client-go => k8s.io/client-go v0.18.2

require (
	github.com/ghodss/yaml v1.0.0
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/klog v1.0.0
)
