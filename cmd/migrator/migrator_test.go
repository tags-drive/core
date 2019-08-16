package migrator

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"io/ioutil"
	"log"
	mathRand "math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/minio/minio-go"
	"github.com/minio/sio"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/tags-drive/core/cmd/common"
)

func TestMain(m *testing.M) {
	// Insert values if needed
	// os.Setenv("TEST_STORAGE_S3_ENDPOINT", "")
	// os.Setenv("TEST_STORAGE_S3_ACCESS_KEY_ID", "")
	// os.Setenv("TEST_STORAGE_S3_SECRET_ACCESS_KEY", "")
	// os.Setenv("TEST_STORAGE_S3_SECURE", "")

	// Set Seed for generateTestFiles() function
	mathRand.Seed(time.Now().UnixNano())

	const testFolder = "testdata"

	// Create the test folder
	err := os.Mkdir(testFolder, 0666)
	if err != nil {
		log.Fatalf("can't create the test folder: %s\n", err)
	}

	os.Chdir(testFolder)
	code := m.Run()
	os.Chdir("..")

	// Remove the test folder
	err = os.RemoveAll(testFolder)
	if err != nil {
		log.Fatalf("can't remove the test folder: %s\n", err)
	}

	os.Exit(code)
}

type s3TestConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	Secure          bool
}

func getS3TestConfig() (cnf s3TestConfig, ok bool) {
	cnf.Endpoint = os.Getenv("TEST_STORAGE_S3_ENDPOINT")
	cnf.AccessKeyID = os.Getenv("TEST_STORAGE_S3_ACCESS_KEY_ID")
	cnf.SecretAccessKey = os.Getenv("TEST_STORAGE_S3_SECRET_ACCESS_KEY")
	cnf.Secure = os.Getenv("TEST_STORAGE_S3_SECURE") != "false" // true by default

	if cnf.Endpoint == "" {
		return s3TestConfig{}, false
	}

	return cnf, true
}

func TestMigrator(t *testing.T) {
	const passPhrase = "test"

	t.Run("from disk to s3", func(t *testing.T) {
		t.Run("without encryption", func(t *testing.T) {
			// ! Warning !
			// We have to use t.Fail() and return to call defer functions!

			assert := assert.New(t)

			// Prepare args

			cnf, ok := getS3TestConfig()
			if !ok {
				t.Skip("Skip test because env vars for connection to an S3 Storage weren't set")
			}

			defer clearDisk()
			defer clearS3(cnf.Endpoint, cnf.AccessKeyID, cnf.SecretAccessKey, cnf.Secure)

			args := []string{
				"--from", "disk",
				"--to", "s3",
				// "--disk.encrypted",
				// "--disk.pass-phrase", "",
				"--s3.endpoint", cnf.Endpoint,
				"--s3.access-key", cnf.AccessKeyID,
				"--s3.secret-key", cnf.SecretAccessKey,
			}
			if cnf.Secure {
				args = append(args, "--s3.secure")
			}

			// Prepare disk

			files := generateTestFiles()
			err := prepareDisk(files, false, [32]byte{})
			if !assert.Nil(err) {
				assert.Fail("can't prepare test files on the disk")
				return
			}

			// Start testing

			app, err := newApp(args)
			if !assert.Nil(err) {
				assert.Fail("can't create a new app")
				return
			}

			err = app.prepare()
			if !assert.Nil(err) {
				assert.Fail("can't prepare the app")
				return
			}

			// Upload files
			app.start()

			// Check files
			err = checkFilesInS3(assert, files, app.s3)
			if !assert.Nil(err) {
				assert.Fail(err.Error())
				return
			}
		})

		t.Run("with encryption", func(t *testing.T) {
			// ! Warning !
			// We have to use t.Fail() and return to call defer function!

			assert := assert.New(t)

			// Prepare args

			cnf, ok := getS3TestConfig()
			if !ok {
				t.Skip("Skip test because env vars for connection to an S3 Storage weren't set")
			}

			defer clearDisk()
			defer clearS3(cnf.Endpoint, cnf.AccessKeyID, cnf.SecretAccessKey, cnf.Secure)

			args := []string{
				"--from", "disk",
				"--to", "s3",
				"--disk.encrypted",
				"--disk.pass-phrase", passPhrase,
				"--s3.endpoint", cnf.Endpoint,
				"--s3.access-key", cnf.AccessKeyID,
				"--s3.secret-key", cnf.SecretAccessKey,
			}
			if cnf.Secure {
				args = append(args, "--s3.secure")
			}

			// Prepare disk

			files := generateTestFiles()
			err := prepareDisk(files, true, sha256.Sum256([]byte(passPhrase)))
			if !assert.Nil(err) {
				assert.Fail("can't prepare test files on the disk")
				return
			}

			// Start testing

			app, err := newApp(args)
			if !assert.Nil(err) {
				assert.Fail("can't create a new app")
				return
			}

			err = app.prepare()
			if !assert.Nil(err) {
				assert.Fail("can't prepare the app")
				return
			}

			// Upload files
			app.start()

			// Check files
			err = checkFilesInS3(assert, files, app.s3)
			if !assert.Nil(err) {
				assert.Fail(err.Error())
				return
			}
		})
	})

	t.Run("from s3 to disk", func(t *testing.T) {
		t.Run("without encryption", func(t *testing.T) {
			// ! Warning !
			// We have to use t.Fail() and return to call defer functions!

			assert := assert.New(t)

			// Prepare args

			cnf, ok := getS3TestConfig()
			if !ok {
				t.Skip("Skip test because env vars for connection to an S3 Storage weren't set")
			}

			defer clearDisk()
			defer clearS3(cnf.Endpoint, cnf.AccessKeyID, cnf.SecretAccessKey, cnf.Secure)

			args := []string{
				"--from", "s3",
				"--to", "disk",
				// "--disk.encrypted",
				// "--disk.pass-phrase", "",
				"--s3.endpoint", cnf.Endpoint,
				"--s3.access-key", cnf.AccessKeyID,
				"--s3.secret-key", cnf.SecretAccessKey,
			}
			if cnf.Secure {
				args = append(args, "--s3.secure")
			}

			// Prepare S3

			files := generateTestFiles()
			err := prepareS3(files, cnf.Endpoint, cnf.AccessKeyID, cnf.SecretAccessKey, cnf.Secure)
			if !assert.Nil(err) {
				assert.Fail("can't prepare test files in an S3 Storage")
				return
			}

			// Start testing

			app, err := newApp(args)
			if !assert.Nil(err) {
				assert.Fail("can't create a new app")
				return
			}

			err = app.prepare()
			if !assert.Nil(err) {
				assert.Fail("can't prepare the app")
				return
			}

			// Upload files
			app.start()

			// Check files
			err = checkFilesOnDisk(assert, files, false, [32]byte{})
			if !assert.Nil(err) {
				assert.Fail(err.Error())
				return
			}
		})

		t.Run("with encryption", func(t *testing.T) {
			// ! Warning !
			// We have to use t.Fail() and return to call defer functions!

			assert := assert.New(t)

			// Prepare args

			cnf, ok := getS3TestConfig()
			if !ok {
				t.Skip("Skip test because env vars for connection to an S3 Storage weren't set")
			}

			defer clearDisk()
			defer clearS3(cnf.Endpoint, cnf.AccessKeyID, cnf.SecretAccessKey, cnf.Secure)

			args := []string{
				"--from", "s3",
				"--to", "disk",
				"--disk.encrypted",
				"--disk.pass-phrase", passPhrase,
				"--s3.endpoint", cnf.Endpoint,
				"--s3.access-key", cnf.AccessKeyID,
				"--s3.secret-key", cnf.SecretAccessKey,
			}
			if cnf.Secure {
				args = append(args, "--s3.secure")
			}

			// Prepare S3

			files := generateTestFiles()
			err := prepareS3(files, cnf.Endpoint, cnf.AccessKeyID, cnf.SecretAccessKey, cnf.Secure)
			if !assert.Nil(err) {
				assert.Fail("can't prepare test files in an S3 Storage")
				return
			}

			// Start testing

			app, err := newApp(args)
			if !assert.Nil(err) {
				assert.Fail("can't create a new app")
				return
			}

			err = app.prepare()
			if !assert.Nil(err) {
				assert.Fail("can't prepare the app")
				return
			}

			// Upload files
			app.start()

			// Check files
			err = checkFilesOnDisk(assert, files, true, sha256.Sum256([]byte(passPhrase)))
			if !assert.Nil(err) {
				assert.Fail(err.Error())
				return
			}
		})
	})
}

type testFile struct {
	name    string
	resized bool

	data []byte
}

func generateTestFiles() []testFile {
	var (
		res               = make([]testFile, 0)
		resizedFileAmount = mathRand.Intn(10) + 10 // from 10 to 20
		usualFileAmount   = mathRand.Intn(10) + 10 // from 10 to 20
	)

	// Usual files
	for i := 0; i < resizedFileAmount; i++ {
		fileSize := mathRand.Intn(1<<20) + 1<<10 // from 1KB to ~1MB
		data := make([]byte, fileSize)
		rand.Read(data)

		file := testFile{
			name:    strconv.Itoa(i),
			resized: false,
			data:    data,
		}
		res = append(res, file)
	}

	// Resized files
	for i := 0; i < usualFileAmount; i++ {
		fileSize := mathRand.Intn(1<<20) + 1<<10 // from 1KB to ~1MB
		data := make([]byte, fileSize)
		rand.Read(data)

		file := testFile{
			name:    strconv.Itoa(i),
			resized: true,
			data:    data,
		}
		res = append(res, file)
	}

	return res
}

// Prepare functions

func prepareDisk(files []testFile, encrypted bool, passPhrase [32]byte) error {
	err := os.MkdirAll(common.VarFolder, 0666)
	if err != nil {
		return errors.Wrap(err, "can't create VarFolder")
	}
	err = os.MkdirAll(common.DataFolder, 0666)
	if err != nil {
		return errors.Wrap(err, "can't create DataFolder")
	}
	err = os.MkdirAll(common.ResizedImagesFolder, 0666)
	if err != nil {
		return errors.Wrap(err, "can't create ResizedImagesFolder")
	}

	for _, f := range files {
		path := common.DataFolder
		if f.resized {
			path = common.ResizedImagesFolder
		}
		path += "/" + f.name

		file, err := os.Create(path)
		if err != nil {
			return errors.Wrap(err, "can't create a new test file")
		}

		// Make the copy of f.data
		cp := make([]byte, len(f.data))
		copy(cp, f.data)
		var src io.Reader = bytes.NewBuffer(cp)
		if encrypted {
			buff := bytes.NewBuffer(nil)
			sio.Encrypt(buff, src, sio.Config{Key: passPhrase[:]})
			src = buff
		}

		_, err = io.Copy(file, src)
		if err != nil {
			return errors.Wrap(err, "can't Copy data into a test file")
		}

		file.Close()
	}

	return nil
}

func prepareS3(files []testFile, endpoint string, accessKeyID string, secretAccessKey string, secure bool) error {
	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, secure)
	if err != nil {
		return errors.Wrap(err, "can't connect to the S3 Storage")
	}

	// Create buckets (buckets must not exist)
	err = client.MakeBucket(common.DataBucket, "")
	if err != nil {
		return errors.Wrap(err, "can't create DataBucket")
	}
	err = client.MakeBucket(common.ResizedImagesBucket, "")
	if err != nil {
		return errors.Wrap(err, "can't create ResizedImagesBucket")
	}

	// Put files
	for _, f := range files {
		bucket := common.DataBucket
		if f.resized {
			bucket = common.ResizedImagesBucket
		}
		key := f.name

		// Make the copy of f.data
		cp := make([]byte, len(f.data))
		copy(cp, f.data)
		buff := bytes.NewBuffer(cp)

		_, err = client.PutObject(bucket, key, buff, int64(buff.Len()), minio.PutObjectOptions{})
		if err != nil {
			return errors.Wrap(err, "can't put an object")
		}
	}

	return nil
}

// Check functions

func checkFilesOnDisk(assert *assert.Assertions, files []testFile, encrypted bool, passPhrase [32]byte) error {
	for _, f := range files {
		path := common.DataFolder
		if f.resized {
			path = common.ResizedImagesFolder
		}
		path += "/" + f.name

		file, err := os.Open(path)
		if !assert.Nil(err) {
			return errors.Wrap(err, "can't open a file")
		}
		defer file.Close()

		var src io.Reader = file
		if encrypted {
			buff := bytes.NewBuffer(nil)
			_, err = sio.Decrypt(buff, src, sio.Config{Key: passPhrase[:]})
			if !assert.Nil(err) {
				return errors.Wrap(err, "can't decrypt a file")
			}
			src = buff
		}

		data, err := ioutil.ReadAll(src)
		if !assert.Nil(err) {
			return errors.Wrap(err, "can't read data from a file")
		}

		if !assert.Equal(f.data, data) {
			return errors.New("content of original file and file on a disk are different")
		}
	}

	return nil
}

func checkFilesInS3(assert *assert.Assertions, files []testFile, client *minio.Client) error {
	for _, f := range files {
		bucket := common.DataBucket
		if f.resized {
			bucket = common.ResizedImagesBucket
		}
		key := f.name

		obj, err := client.GetObject(bucket, key, minio.GetObjectOptions{})
		if !assert.Nil(err) {
			return errors.Wrapf(err, "can't get file '%s/%s'", bucket, key)
		}

		data, err := ioutil.ReadAll(obj)
		if !assert.Nil(err) {
			return errors.Wrapf(err, "can't read data from file '%s/%s'", bucket, key)
		}
		obj.Close()

		if !assert.Equal(f.data, data) {
			return errors.New("content of original file and file in S3 Storage are different")
		}
	}

	return nil
}

// Clear functions

func clearDisk() error {
	// Can remove only VarFolder
	err := os.RemoveAll(common.VarFolder)
	return errors.Wrap(err, "can't remove VarFolder")
}

// clearS3 removes test buckets. If an error occurred, it panics
func clearS3(endpoint string, accessKeyID string, secretAccessKey string, secure bool) error {
	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, secure)
	if err != nil {
		return errors.Wrap(err, "can't connect to the S3 Storage")
	}

	buckets := [...]string{common.DataBucket, common.ResizedImagesBucket}
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
		err := client.RemoveBucket(bucket)
		if err != nil {
			return errors.Wrap(err, "can't delete a bucket")
		}

		close(done)
	}

	return nil
}
