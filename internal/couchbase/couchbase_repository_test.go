package couchbase_test

import (
	"context"
	"fmt"
	"github.com/couchbase/gocb/v2"
	"github.com/docker/go-connections/nat"
	"github.com/harundurmus/go-to-do-app/internal/config"
	"github.com/harundurmus/go-to-do-app/internal/todo"
	"github.com/harundurmus/go-to-do-app/pkg/couchbase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
	"log"
	"os/exec"
	"testing"
)

type CouchbaseTestSuite struct {
	suite.Suite
	container      testcontainers.Container
	conf           config.Couchbase
	cb             *couchbase.Couchbase
	bucket         *gocb.Bucket
	todoRepository todo.Repository
}

const (
	DBUsername                   = "admin"
	DBPassword                   = "password"
	todoAppBucket                = "todo"
	_deliveryPointCollection     = "deliverypoint"
	_incorrectShipmentCollection = "incorrectshipment"
	_shipmentCollection          = "shipment"
	_todoCollection              = "todo"
)

func (s *CouchbaseTestSuite) SetupSuite() {
	SetupTestCouchbaseInstance(s.T())

	s.conf = prepareConfig()
	s.cb = createCouchbaseConnection(s.T(), s.conf)

	var err error
	s.bucket, err = s.cb.Bucket(todoAppBucket)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.bucket)

	s.todoRepository = todo.NewRepository(s.cb.Cluster(), s.bucket)
}

func SetupTestCouchbaseInstance(t *testing.T) {
	startTestContainer(t)
	prepareDB()
}

func prepareConfig() config.Couchbase {
	return config.Couchbase{
		URL:      "localhost",
		Username: DBUsername,
		Password: DBPassword,
		Buckets: []config.BucketConfig{{
			Name: todoAppBucket,
			Scopes: []config.ScopeConfig{
				{
					Name: "",
					Collections: []config.CollectionConfig{
						{Name: _shipmentCollection},
						{Name: _todoCollection},
						{Name: _deliveryPointCollection},
						{Name: _incorrectShipmentCollection},
					},
				},
			},
		}},
	}
}

func prepareDB() {
	const baseURL = "curl -s -u Administrator:password -X POST http://localhost:8091%s"
	const initializeNodeCommand = "/nodes/self/controller/settings " +
		"-d path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fdata " +
		"-d index_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fadata " +
		"-d cbas_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fedata " +
		"-d eventing_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fidata"
	runCurlCommand(baseURL, initializeNodeCommand)

	const startServicesCommand = "/node/controller/setupServices -d services=kv%2Cindex%2Cn1ql%2Cfts"
	runCurlCommand(baseURL, startServicesCommand)

	const setBucketMemoryQuotaCommand = "/pools/default -d memoryQuota=256 -d indexMemoryQuota=256 -d ftsMemoryQuota=256"
	runCurlCommand(baseURL, setBucketMemoryQuotaCommand)

	const createAdminUserCommand = "/settings/web -d port=8091 -d username=admin -d password=password"
	runCurlCommand(baseURL, createAdminUserCommand)

	const newBaseURL = "curl -s -u admin:password -X POST http://localhost:8091%s"
	setIndexStorageSetting := "/settings/indexes " +
		"-d indexerThreads=0 " +
		"-d logLevel=info " +
		"-d maxRollbackPoints=5 " +
		"-d memorySnapshotInterval=200 " +
		"-d stableSnapshotInterval=5000 " +
		"-d storageMode=forestdb"
	runCurlCommand(newBaseURL, setIndexStorageSetting)
}

func createCouchbaseConnection(t *testing.T, conf config.Couchbase) *couchbase.Couchbase {
	err := createBuckets(t, conf)
	assert.Nil(t, err)

	cb, err := couchbase.New(conf)
	assert.Nil(t, err)
	assert.NotNil(t, cb)

	return cb
}

func runCurlCommand(baseURL, path string) {
	logger := zap.NewNop()
	command := fmt.Sprintf(baseURL, path)
	output, err := exec.Command("/bin/sh", "-c", command).CombinedOutput()
	if err != nil {
		logger.Debug(string(output))
		log.Fatal(err)
	}
}

func startTestContainer(t *testing.T) testcontainers.Container {
	port, err := nat.NewPort("tcp", "8091")
	assert.Nil(t, err)

	req := testcontainers.ContainerRequest{
		Image: "couchbase:community-7.0.0",
		ExposedPorts: []string{
			"8091:8091/tcp",
			"8092:8092/tcp",
			"8093:8093/tcp",
			"8094:8094/tcp",
			"11207:11207/tcp",
			"11210:11210/tcp",
			"11211:11211/tcp",
		},
		WaitingFor: wait.ForListeningPort(port),
	}

	ctx := context.Background()
	container, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})

	assert.Nil(t, err)
	return container
}

func (s *CouchbaseTestSuite) TeardownSuite() {
	err := s.container.Terminate(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func TestCouchbase(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, new(CouchbaseTestSuite))
}

func (s *CouchbaseTestSuite) CreateBuckets() error {
	return createBuckets(s.T(), s.conf)
}

func createBuckets(t *testing.T, conf config.Couchbase) error {
	cluster, err := gocb.Connect(
		conf.URL,
		gocb.ClusterOptions{
			Username: conf.Username,
			Password: conf.Password,
		},
	)
	assert.Nil(t, err)

	for _, bucketConf := range conf.Buckets {
		err := cluster.Buckets().CreateBucket(
			gocb.CreateBucketSettings{
				BucketSettings: gocb.BucketSettings{
					Name:           bucketConf.Name,
					FlushEnabled:   true,
					RAMQuotaMB:     100,
					NumReplicas:    1,
					BucketType:     gocb.CouchbaseBucketType,
					EvictionPolicy: gocb.EvictionPolicyTypeFull,
				},
				ConflictResolutionType: "",
			},
			nil,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
