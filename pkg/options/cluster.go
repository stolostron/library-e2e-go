package options

import (
	"fmt"
)

// ClusterName ...
type ClusterName struct {
	//Owner string
	Cloud string
	uid   string
}

//NewClusterName returns a new ClusterName with a uid.
func NewClusterName(cloud string) (*ClusterName, error) {
	switch cloud {
	case "aws", "gcp", "azure":
		uid, err := GetUID()
		if err != nil {
			return nil, err
		}
		return &ClusterName{
			//Owner: GetOwner(),
			Cloud: cloud,
			uid:   uid,
		}, nil
	default:
		return nil, fmt.Errorf("Can not generate clusterName as the cloud %s is unsupported", cloud)
	}
}

//String format the clusterName as string
func (c *ClusterName) String() string {
	return fmt.Sprintf("%s-%s", c.Cloud, c.uid)
}

//GetUID return the uid for that ClusterName
func (c *ClusterName) GetUID() string {
	return c.uid
}
