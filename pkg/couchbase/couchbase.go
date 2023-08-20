package couchbase

import (
	"errors"
	"fmt"
	"github.com/couchbase/gocb/v2"
	"github.com/harundurmus/go-to-do-app/internal/config"
	"github.com/labstack/gommon/log"
	"strings"
	"time"
)

const _timeout = 10 * time.Second
const defaultNamespace = "default"

type Couchbase struct {
	cluster *gocb.Cluster
	buckets map[string]*gocb.Bucket
}

func New(conf config.Couchbase) (*Couchbase, error) {
	cluster, err := gocb.Connect(conf.URL, gocb.ClusterOptions{
		Username: conf.Username,
		Password: conf.Password,
	})
	if err != nil {
		return nil, err
	}
	exists, err := verifyBucketsExists(cluster, conf.Buckets)
	if err != nil {
		log.Debug("Error verifying buckets exists", err)
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("all required buckets does not exists on the cluster")
	}
	buckets, err := connectToBuckets(cluster, conf.Buckets)
	if err != nil {
		log.Debug("Error connecting to buckets", err)
		return nil, err
	}
	return &Couchbase{
		cluster: cluster,
		buckets: buckets,
	}, nil
}

func (c *Couchbase) Bucket(name string) (*gocb.Bucket, error) {
	bucket, ok := c.buckets[name]
	if ok {
		return bucket, nil
	}
	return nil, fmt.Errorf("bucket not found")
}

func (c *Couchbase) Cluster() *gocb.Cluster {
	return c.cluster
}

func connectToBuckets(cluster *gocb.Cluster, confs []config.BucketConfig) (map[string]*gocb.Bucket, error) {
	buckets := make(map[string]*gocb.Bucket, len(confs))
	for i := range confs {
		bucket, err := connectToBucket(cluster, &confs[i])
		if err != nil {
			return nil, err
		}
		buckets[confs[i].Name] = bucket
	}
	return buckets, nil
}

func connectToBucket(cluster *gocb.Cluster, conf *config.BucketConfig) (*gocb.Bucket, error) {
	bucket := cluster.Bucket(conf.Name)
	if err := bucket.WaitUntilReady(_timeout, nil); err != nil {
		return nil, err
	}
	if err := createScopes(cluster, bucket, conf.Scopes); err != nil {
		return nil, err
	}
	if conf.CreatePrimaryIndex {
		if err := createPrimaryIndex(cluster, conf.Name); err != nil {
			return nil, err
		}
	}
	return bucket, nil
}

func createPrimaryIndex(cluster *gocb.Cluster, bucket string) error {
	opts := &gocb.CreatePrimaryQueryIndexOptions{
		IgnoreIfExists: true,
		RetryStrategy:  gocb.NewBestEffortRetryStrategy(nil),
	}

	return cluster.QueryIndexes().CreatePrimaryIndex(bucket, opts)
}

func createScopes(cluster *gocb.Cluster, bucket *gocb.Bucket, scopes []config.ScopeConfig) error {
	existingScopes, err := bucket.Collections().GetAllScopes(nil)
	if err != nil {
		return err
	}

	existingScopesMap := make(map[string]bool, len(existingScopes))
	for _, scope := range existingScopes {
		existingScopesMap[scope.Name] = true
	}

	for _, scopeConf := range scopes {
		if scopeConf.Name != "" && !existingScopesMap[scopeConf.Name] {
			opts := &gocb.CreateScopeOptions{
				RetryStrategy: gocb.NewBestEffortRetryStrategy(nil),
			}

			if err := bucket.Collections().CreateScope(scopeConf.Name, opts); err != nil {
				return err
			}
		}
		if err := createCollections(cluster, bucket, scopeConf.Name, scopeConf.Collections); err != nil {
			return err
		}
	}
	return nil
}

func createCollections(cluster *gocb.Cluster, bucket *gocb.Bucket, scopeName string, collections []config.CollectionConfig) error {
	if scopeName == "" {
		scopeName = bucket.DefaultCollection().ScopeName()
	}

	for _, collection := range collections {
		spec := gocb.CollectionSpec{
			Name:      collection.Name,
			ScopeName: scopeName,
		}

		opts := &gocb.CreateCollectionOptions{
			RetryStrategy: gocb.NewBestEffortRetryStrategy(nil),
		}

		if err := bucket.Collections().CreateCollection(spec, opts); err != nil {
			if !(strings.Contains(err.Error(), "Collection with name") && strings.Contains(err.Error(), "already exists")) {
				return err
			}
		}

		if collection.CreatePrimaryIndex {
			primaryIndexName := fmt.Sprintf("%s`:`%s`.`%s`.`%s", "default", bucket.Name(), scopeName, collection.Name)
			primaryIndexOptions := &gocb.CreatePrimaryQueryIndexOptions{
				IgnoreIfExists: true,
				RetryStrategy:  gocb.NewBestEffortRetryStrategy(nil),
			}

			if err := checkIfKeyspaceIsCreated(cluster, bucket.Name(), scopeName, collection.Name); err != nil {
				return err
			}

			err := cluster.QueryIndexes().CreatePrimaryIndex(primaryIndexName, primaryIndexOptions)
			if err != nil {
				return err
			}
		}

		for _, field := range collection.FieldIndexes {
			err := createCollectionFieldIndex(
				cluster,
				bucket.Name(),
				scopeName,
				collection.Name,
				field,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createCollectionFieldIndex(cluster *gocb.Cluster, bucket, scope, collection, field string) error {
	if err := checkIfKeyspaceIsCreated(cluster, bucket, scope, collection); err != nil {
		return err
	}

	indexName := fmt.Sprintf(
		"%s_%s_%s_%s_FieldIndex",
		bucket,
		scope,
		collection,
		field,
	)

	namespace := fmt.Sprintf(
		"%s`.`%s`.`%s",
		bucket,
		scope,
		collection,
	)

	opts := &gocb.CreateQueryIndexOptions{
		IgnoreIfExists: true,
		RetryStrategy:  gocb.NewBestEffortRetryStrategy(nil),
	}

	return cluster.QueryIndexes().CreateIndex(
		namespace,
		indexName,
		[]string{field},
		opts,
	)
}

func checkIfKeyspaceIsCreated(cluster *gocb.Cluster, bucketName, scopeName, collectionName string) error {
	keyspaceName := fmt.Sprintf("%s:%s.%s.%s", defaultNamespace, bucketName, scopeName, collectionName)
	hasKeyspace, err := hasKeyspace(cluster, keyspaceName)
	if err != nil {
		return err
	}
	if !hasKeyspace {
		return errors.New("cannot create primary index: could not find keyspace")
	}
	return nil
}

func hasKeyspace(cluster *gocb.Cluster, keyspaceName string) (bool, error) {
	type Keyspace struct {
		Path string `json:"path"`
	}
	type KeyspaceResult struct {
		Keyspaces Keyspace `json:"keyspaces"`
	}

	var keyspace KeyspaceResult
	const retryCount = 5
	for i := 0; i < retryCount; i++ {
		result, err := cluster.Query(fmt.Sprintf("SELECT * FROM system:keyspaces where keyspaces.`path`='%s'", keyspaceName), nil)
		if err != nil {
			return false, err
		}

		if result.Next() {
			err = result.Row(&keyspace)
			if err != nil {
				return false, err
			}
		}
		if keyspace.Keyspaces.Path != "" {
			return true, nil
		}

		time.Sleep(1 * time.Second)
	}

	return false, nil
}

func verifyBucketsExists(cluster *gocb.Cluster, confs []config.BucketConfig) (bool, error) {
	buckets, err := cluster.Buckets().GetAllBuckets(nil)
	if err != nil {
		return false, err
	}

	for _, conf := range confs {
		if _, ok := buckets[conf.Name]; !ok {
			return false, nil
		}
	}

	return true, nil
}
