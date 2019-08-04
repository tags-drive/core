package bs_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/minio/minio-go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	bs "github.com/tags-drive/core/internal/storage/files/binary_storage"
)

const (
	dataBucket          = "test-data"
	resizedImagesBucket = "test-data-resized"
)

func TestMain(m *testing.M) {
	// Insert values if needed
	// os.Setenv("TEST_STORAGE_S3_ENDPOINT", "")
	// os.Setenv("TEST_STORAGE_S3_ACCESS_KEY_ID", "")
	// os.Setenv("TEST_STORAGE_S3_SECRET_ACCESS_KEY", "")
	// os.Setenv("TEST_STORAGE_S3_SECURE", "")

	os.Exit(m.Run())
}

func TestS3Storage_SaveFile(t *testing.T) {
	assert := assert.New(t)

	cnf, ok := getS3Config()
	if !ok {
		t.Skip("Skip test because env vars for connection to an S3 Storage weren't set")
	}

	storage, err := bs.NewS3Storage(cnf)
	if !assert.Nil(err) {
		assert.FailNow("can't init a new S3Storage")
	}

	tests := []struct {
		dataSize int64
		resized  bool
	}{
		{dataSize: 256, resized: false},
		{dataSize: 773, resized: false},
		{dataSize: 1024, resized: false},
		{dataSize: 8192, resized: false},
		//
		{dataSize: 33, resized: true},
		{dataSize: 661, resized: true},
	}

	for id, tt := range tests {
		original := generateRandomData(int(tt.dataSize))
		cp := make([]byte, tt.dataSize)
		copy(cp, original)

		buff := bytes.NewBuffer(cp)

		// Save data
		err := storage.SaveFile(buff, id, tt.dataSize, tt.resized)
		if !assert.Nil(err, "can't save a file with id '%d'", id) {
			continue
		}

		// Read and check data
		buff = &bytes.Buffer{}
		err = storage.GetFile(buff, id, tt.resized)
		if !assert.Nilf(err, "can't get a file with id '%d'", id) {
			continue
		}

		assert.Equal(original, buff.Bytes(), "diffrent data")
	}

	// Clear S3 storage
	err = clearS3(cnf.Endpoint, cnf.AccessKeyID, cnf.SecretAccessKey, cnf.Secure)
	if !assert.Nil(err) {
		assert.FailNow("can't clear S3 storage")
	}
}

func TestS3Storage_GetFile(t *testing.T) {
	assert := assert.New(t)

	cnf, ok := getS3Config()
	if !ok {
		t.Skip("Skip test because env vars for connection to an S3 Storage weren't set")
	}

	storage, err := bs.NewS3Storage(cnf)
	if !assert.Nil(err) {
		assert.FailNow("can't init a new S3Storage")
	}

	tests := []struct {
		id      int
		resized bool
		//
		dataSize     int64
		originalData []byte // should be generated
		//
		exist bool
	}{
		{id: 0, resized: false, dataSize: 128, exist: true},
		{id: 1, resized: true, dataSize: 16380, exist: true},
		//
		{id: 2, resized: false, dataSize: 0, exist: false},
		{id: 3, resized: true, dataSize: 0, exist: false},
	}

	// Save files
	for i, tt := range tests {
		if !tt.exist {
			continue
		}

		tests[i].originalData = generateRandomData(int(tt.dataSize))

		cp := make([]byte, tt.dataSize)
		copy(cp, tests[i].originalData)

		buff := bytes.NewBuffer(cp)
		storage.SaveFile(buff, tt.id, tt.dataSize, tests[i].resized) // ingore errors
	}

	for _, tt := range tests {
		buff := &bytes.Buffer{}
		err := storage.GetFile(buff, tt.id, tt.resized)
		if !tt.exist {
			assert.NotNil(err, "get a non-existing file")
			continue
		}

		if !assert.Nilf(err, "can't get a file with id '%d'", tt.id) {
			continue
		}

		// Check data
		assert.Equal(tt.originalData, buff.Bytes(), "get different data")
	}

	// Clear S3 storage
	err = clearS3(cnf.Endpoint, cnf.AccessKeyID, cnf.SecretAccessKey, cnf.Secure)
	if !assert.Nil(err) {
		assert.FailNow("can't clear S3 storage")
	}
}

func TestS3Storage_GetFileStats(t *testing.T) {
	assert := assert.New(t)

	cnf, ok := getS3Config()
	if !ok {
		t.Skip("Skip test because env vars for connection to an S3 Storage weren't set")
	}

	storage, err := bs.NewS3Storage(cnf)
	if !assert.Nil(err) {
		assert.FailNow("can't init a new S3Storage")
	}

	tests := []struct {
		originalData []byte // should be generated
		dataSize     int64
	}{
		{dataSize: 1},
		{dataSize: 256},
		{dataSize: 16384},
	}

	// Save files
	for i := range tests {
		tests[i].originalData = generateRandomData(int(tests[i].dataSize))

		cp := make([]byte, tests[i].dataSize)
		copy(cp, tests[i].originalData)

		buff := bytes.NewBuffer(cp)
		storage.SaveFile(buff, i, tests[i].dataSize, false) // ingore errors
	}

	// Get stats
	for i, tt := range tests {
		stats, err := storage.GetFileStats(i)
		if !assert.Nilf(err, "can't get stats for file with id '%d'", i) {
			continue
		}

		// Only size matters
		assert.Equal(tt.dataSize, stats.Size(), "different sized got")
	}

	// Clear S3 storage
	err = clearS3(cnf.Endpoint, cnf.AccessKeyID, cnf.SecretAccessKey, cnf.Secure)
	if !assert.Nil(err) {
		assert.FailNow("can't clear S3 storage")
	}
}

func TestS3Storage_DeleteFile(t *testing.T) {
	assert := assert.New(t)

	cnf, ok := getS3Config()
	if !ok {
		t.Skip("Skip test because env vars for connection to an S3 Storage weren't set")
	}

	storage, err := bs.NewS3Storage(cnf)
	if !assert.Nil(err) {
		assert.FailNow("can't init a new S3Storage")
	}

	tests := []struct {
		dataSize int64
		resized  bool
	}{
		{dataSize: 115, resized: false},
		{dataSize: 16384, resized: false},
		//
		{dataSize: 256, resized: true},
	}

	// Save files
	for id, tt := range tests {
		data := generateRandomData(int(tt.dataSize))
		buff := bytes.NewBuffer(data)
		storage.SaveFile(buff, id, tt.dataSize, tt.resized) // ingore errors
	}

	// Delete files
	for id, tt := range tests {
		err := storage.DeleteFile(id, tt.resized)
		assert.Nilf(err, "can't delete a file with id '%d'", id)
	}

	// Clear S3 storage
	err = clearS3(cnf.Endpoint, cnf.AccessKeyID, cnf.SecretAccessKey, cnf.Secure)
	if !assert.Nil(err) {
		assert.FailNow("can't clear S3 storage")
	}
}

// List of test env variables:
//  - TEST_STORAGE_S3_ENDPOINT
//  - TEST_STORAGE_S3_ACCESS_KEY_ID
//  - TEST_STORAGE_S3_SECRET_ACCESS_KEY
//  - TEST_STORAGE_S3_SECURE
//  - TEST_STORAGE_S3_BUCKET_LOCATION
//
func getS3Config() (cnf bs.S3StorageConfig, ok bool) {
	var (
		endpoint        string = os.Getenv("TEST_STORAGE_S3_ENDPOINT")
		accessKeyID     string = os.Getenv("TEST_STORAGE_S3_ACCESS_KEY_ID")
		secretAccessKey string = os.Getenv("TEST_STORAGE_S3_SECRET_ACCESS_KEY")
		secure          bool   = os.Getenv("TEST_STORAGE_S3_SECURE") != "false" // true by default
		bucketLocation  string = os.Getenv("TEST_STORAGE_S3_BUCKET_LOCATION")   // can be empty
	)

	if isAnyEmpty(endpoint, accessKeyID, secretAccessKey) {
		return cnf, false
	}

	return bs.S3StorageConfig{
		Endpoint:            endpoint,
		AccessKeyID:         accessKeyID,
		SecretAccessKey:     secretAccessKey,
		Secure:              secure,
		BucketLocation:      bucketLocation,
		DataBucket:          dataBucket,
		ResizedImagesBucket: resizedImagesBucket,
	}, true
}

func isAnyEmpty(vars ...string) bool {
	for _, v := range vars {
		if v == "" {
			return true
		}
	}

	return false
}

// clearS3 removes test buckets. If an error occurred, it panics
func clearS3(endpoint string, accessKeyID string, secretAccessKey string, secure bool) error {
	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, secure)
	if err != nil {
		return errors.Wrap(err, "can't init connection")
	}

	buckets := [...]string{dataBucket, resizedImagesBucket}
	for _, bucket := range buckets {
		// Get list of all objects
		done := make(chan struct{})
		objects := client.ListObjects(bucket, "", true, done)

		// Remove all objects
		for obj := range objects {
			err := client.RemoveObject(bucket, obj.Key)
			if err != nil {
				return errors.Wrap(err, "can't delete an object")
			}
		}

		// Remove a bucket
		err = client.RemoveBucket(bucket)
		if err != nil {
			return errors.Wrap(err, "can't delete a bucket")
		}

		close(done)
	}

	return nil
}
