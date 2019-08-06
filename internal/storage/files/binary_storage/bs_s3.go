package bs

import (
	"io"
	"os"
	"strconv"
	"time"

	"github.com/minio/minio-go"
	"github.com/pkg/errors"
)

type S3Storage struct {
	client *minio.Client

	config S3StorageConfig
}

type S3StorageConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	Secure          bool

	BucketLocation      string
	DataBucket          string
	ResizedImagesBucket string
}

func NewS3Storage(cnf S3StorageConfig) (*S3Storage, error) {
	var err error

	storage := new(S3Storage)
	storage.config = cnf

	// Init a connection
	storage.client, err = minio.New(cnf.Endpoint, cnf.AccessKeyID, cnf.SecretAccessKey, cnf.Secure)
	if err != nil {
		return nil, errors.Wrap(err, "can't init connection with an S3 storage")
	}

	if err := ping(storage.client); err != nil {
		return nil, errors.Wrap(err, "S3 Storage is unavailable")
	}

	// Create buckets if needed

	if exist, err := storage.client.BucketExists(cnf.DataBucket); !exist || err != nil {
		// Have to create a bucket
		err = storage.client.MakeBucket(cnf.DataBucket, cnf.BucketLocation)
		if err != nil {
			return nil, errors.Wrapf(err, "can't make bucket '%s'", cnf.DataBucket)
		}
	}

	if exist, err := storage.client.BucketExists(cnf.ResizedImagesBucket); !exist || err != nil {
		// Have to create a bucket
		err = storage.client.MakeBucket(cnf.ResizedImagesBucket, cnf.BucketLocation)
		if err != nil {
			return nil, errors.Wrapf(err, "can't make bucket '%s'", cnf.ResizedImagesBucket)
		}
	}

	return storage, nil
}

func ping(client *minio.Client) error {
	// There are no other ways to ping any S3 Storage
	defMaxRetry := minio.MaxRetry
	minio.MaxRetry = 1
	_, err := client.ListBuckets()
	minio.MaxRetry = defMaxRetry
	return err
}

func (s3 S3Storage) GetFile(w io.Writer, fileID int, resized bool) error {
	objectName := strconv.Itoa(fileID)
	bucket := s3.config.DataBucket
	if resized {
		bucket = s3.config.ResizedImagesBucket
	}

	obj, err := s3.client.GetObject(bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return errors.Wrapf(err, "can't get an object '%s/%s'", bucket, objectName)
	}
	defer obj.Close()

	_, err = io.Copy(w, obj)
	if err != nil {
		return errors.Wrapf(err, "can't copy an object '%s/%s'", bucket, objectName)
	}

	return nil
}

func (s3 S3Storage) GetFileStats(fileID int) (os.FileInfo, error) {
	objectName := strconv.Itoa(fileID)
	bucket := s3.config.DataBucket

	stats, err := s3.client.StatObject(bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "can't get stats for an object '%s/%s'", bucket, objectName)
	}

	return newFileInfo(stats), nil
}

func (s3 S3Storage) SaveFile(r io.Reader, fileID int, fileSize int64, resized bool) error {
	objectName := strconv.Itoa(fileID)
	bucket := s3.config.DataBucket
	if resized {
		bucket = s3.config.ResizedImagesBucket
	}

	_, err := s3.client.PutObject(bucket, objectName, r, fileSize, minio.PutObjectOptions{})
	return errors.Wrap(err, "can't put an object")
}

func (s3 S3Storage) DeleteFile(fileID int, resized bool) error {
	objectName := strconv.Itoa(fileID)
	bucket := s3.config.DataBucket
	if resized {
		bucket = s3.config.ResizedImagesBucket
	}

	err := s3.client.RemoveObject(bucket, objectName)
	return errors.Wrapf(err, "can't remove an object '%s/%s'", bucket, objectName)
}

func (s3 *S3Storage) Shutdown() error {
	// We hadn't to shutdown minio.Client{}

	return nil
}

// fileInfo suffice os.FileInfo interface.
type fileInfo struct {
	stats minio.ObjectInfo
}

func newFileInfo(objInfo minio.ObjectInfo) fileInfo {
	return fileInfo{
		stats: objInfo,
	}
}

func (info fileInfo) Name() string {
	return info.stats.Key
}

func (info fileInfo) Size() int64 {
	return info.stats.Size
}

func (info fileInfo) Mode() os.FileMode {
	// TODO: ?
	return 0666
}

func (info fileInfo) ModTime() time.Time {
	return info.stats.LastModified
}

func (info fileInfo) IsDir() bool {
	return false
}

func (info fileInfo) Sys() interface{} {
	return nil
}
