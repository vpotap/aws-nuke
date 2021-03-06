package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

type RDSDBCluster struct {
	svc                *rds.RDS
	id                 string
	deletionProtection bool
}

func init() {
	register("RDSDBCluster", ListRDSClusters)
}

func ListRDSClusters(sess *session.Session) ([]Resource, error) {
	svc := rds.New(sess)

	params := &rds.DescribeDBClustersInput{}
	resp, err := svc.DescribeDBClusters(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, instance := range resp.DBClusters {
		resources = append(resources, &RDSDBCluster{
			svc:                svc,
			id:                 *instance.DBClusterIdentifier,
			deletionProtection: *instance.DeletionProtection,
		})
	}

	return resources, nil
}

func (i *RDSDBCluster) Remove() error {
	if (i.deletionProtection) {
		modifyParams := &rds.ModifyDBClusterInput{
			DBClusterIdentifier: &i.id,
			DeletionProtection:  aws.Bool(false),
		}
		_, err := i.svc.ModifyDBCluster(modifyParams)
		if err != nil {
			return err
		}
	}

	params := &rds.DeleteDBClusterInput{
		DBClusterIdentifier: &i.id,
		SkipFinalSnapshot:   aws.Bool(true),
	}

	_, err := i.svc.DeleteDBCluster(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *RDSDBCluster) String() string {
	return i.id
}
