package migrator

import (
	"context"
	"crypto/sha256"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ShoshinNikita/go-disk-buffer"
	"github.com/ShoshinNikita/log/v2"
	"github.com/jessevdk/go-flags"
	"github.com/minio/minio-go"
	"github.com/minio/sio"
	"github.com/pkg/errors"

	"github.com/tags-drive/core/cmd/common"
	"github.com/tags-drive/core/internal/utils"
)

type file struct {
	name    string
	size    int64
	resized bool
	r       io.ReadCloser
}

type app struct {
	config config

	s3 *minio.Client

	logger *clog.Logger
}

type config struct {
	From string `long:"from" required:"true" choice:"disk" choice:"s3"`
	To   string `long:"to" required:"true" choice:"disk" choice:"s3"`

	Disk struct {
		Encrypted        bool   `long:"encrypted"`
		PassPhraseString string `long:"pass-phrase"`
		PassPhrase       [32]byte
	} `group:"disk" namespace:"disk"`

	S3 struct {
		Endpoint        string `long:"endpoint" required:"true"`
		AccessKeyID     string `long:"access-key" required:"true"`
		SecretAccessKey string `long:"secret-key" required:"true"`
		Secure          bool   `long:"secure"`
		BucketLocation  string `long:"bucket-location"`
	} `group:"s3" namespace:"s3"`
}

func newApp(args []string) (*app, error) {
	app := &app{
		// ProdLogger is the default value
		logger: clog.NewProdLogger(),
	}

	// Parse config
	parser := flags.NewParser(&app.config, flags.HelpFlag|flags.PassDoubleDash|flags.IgnoreUnknown)

	_, err := parser.ParseArgs(args)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse flags")
	}

	if app.config.Disk.Encrypted {
		app.config.Disk.PassPhrase = sha256.Sum256([]byte(app.config.Disk.PassPhraseString))
		app.config.Disk.PassPhraseString = ""
	}

	// check --to and --from
	if app.config.From == app.config.To {
		return nil, errors.New("--from and --to can't be equal")
	}

	// Init S3 client
	app.s3, err = minio.New(app.config.S3.Endpoint, app.config.S3.AccessKeyID, app.config.S3.SecretAccessKey, app.config.S3.Secure)
	if err != nil {
		return nil, errors.Wrap(err, "can't init new S3 client")
	}

	// Ping S3 Storage
	defMaxRetry := minio.MaxRetry
	minio.MaxRetry = 1
	_, err = app.s3.ListBuckets()
	minio.MaxRetry = defMaxRetry
	if err != nil {
		return nil, errors.Wrap(err, "can't ping an S3 Storage")
	}

	return app, nil
}

func (app *app) setLogger(l *clog.Logger) {
	app.logger = l
}

func (app *app) prepare() error {
	switch app.config.To {
	case "disk":
		err := os.MkdirAll(common.DataFolder, 0666)
		if err != nil {
			return errors.Wrap(err, "can't create DataFolder")
		}

		err = os.MkdirAll(common.ResizedImagesFolder, 0666)
		if err != nil {
			return errors.Wrap(err, "can't create ResizedImagesFolder")
		}
	case "s3":
		if exist, err := app.s3.BucketExists(common.DataBucket); !exist || err != nil {
			// Have to create a bucket
			err = app.s3.MakeBucket(common.DataBucket, app.config.S3.BucketLocation)
			if err != nil {
				return errors.Wrapf(err, "can't make bucket '%s'", common.DataBucket)
			}
		}

		if exist, err := app.s3.BucketExists(common.ResizedImagesBucket); !exist || err != nil {
			// Have to create a bucket
			err = app.s3.MakeBucket(common.ResizedImagesBucket, app.config.S3.BucketLocation)
			if err != nil {
				return errors.Wrapf(err, "can't make bucket '%s'", common.ResizedImagesBucket)
			}
		}
	default:
		// Can skip
	}

	return nil
}

func (app *app) start() {
	var (
		from From
		to   To
	)

	switch app.config.From {
	case "disk":
		from = app.fromDisk
	case "s3":
		from = app.fromS3
	default:
		// Can skip
	}

	switch app.config.To {
	case "disk":
		to = app.toDisk
	case "s3":
		to = app.toS3
	default:
		// Can skip
	}

	files := from()
	to(files)
}

// From function

type From func() <-chan file

func (app *app) fromDisk() <-chan file {
	files := make(chan file, 256)

	go func() {
		walkFunction := func(root string, resized bool) filepath.WalkFunc {
			return func(path string, info os.FileInfo, err error) error {
				if info.IsDir() {
					// Skip folders
					if path != root {
						return filepath.SkipDir
					}
					return nil
				}

				f, err := os.Open(path)
				if err != nil {
					app.logger.Errorf("can't open file '%s': %s\n", path, err)
					return nil
				}

				var (
					src      io.ReadCloser = f
					fileSize int64         = info.Size()
				)

				if app.config.Disk.Encrypted {
					// Have to decrypt a file. File must be closed after decryption!

					buff := buffer.NewBuffer(nil)
					_, err := sio.Decrypt(buff, f, sio.Config{Key: app.config.Disk.PassPhrase[:]})
					if err != nil {
						app.logger.Errorf("can't decrypt a file '%s': %s\n", path, err)
						return nil
					}

					r, size := utils.GetReaderSize(buff)

					// Update with decrypted file propertirs
					fileSize = size
					src = readCloserWrapper{r: r}

					// Have to close a file
					f.Close()
				}

				file := file{
					name:    info.Name(),
					size:    fileSize,
					resized: resized,
					r:       src,
				}
				files <- file

				return nil
			}
		}

		// Run every Walk function in a goroutine.

		var wg sync.WaitGroup

		// Upload files from "./var/data"
		wg.Add(1)
		go func() {
			defer wg.Done()

			filepath.Walk(common.DataFolder, walkFunction(common.DataFolder, false))
		}()

		// Upload files from "./var/data/resized"
		wg.Add(1)
		go func() {
			defer wg.Done()

			filepath.Walk(common.ResizedImagesFolder, walkFunction(common.ResizedImagesFolder, true))
		}()

		wg.Wait()

		// Close the channel
		close(files)
	}()

	return files
}

func (app *app) fromS3() <-chan file {
	files := make(chan file, 256)

	go func() {
		getFunction := func(bucket string, resized bool) {
			done := make(chan struct{})

			objects := app.s3.ListObjects(bucket, "", false, done)
			for object := range objects {
				obj, err := app.s3.GetObject(bucket, object.Key, minio.GetObjectOptions{})
				if err != nil {
					app.logger.Errorf("can't get an object '%s/%s': %s\n", bucket, object.Key, err)
					continue
				}

				file := file{
					name:    object.Key,
					size:    object.Size,
					resized: resized,
					r:       obj,
				}
				files <- file
			}

			close(done)
		}

		var wg sync.WaitGroup

		// Get files from DataBucket
		wg.Add(1)
		go func() {
			defer wg.Done()

			getFunction(common.DataBucket, false)
		}()

		// Get files from ResizedImagesFolder
		wg.Add(1)
		go func() {
			defer wg.Done()

			getFunction(common.ResizedImagesBucket, true)
		}()

		wg.Wait()

		// Close the channel
		close(files)
	}()

	return files
}

// To functions

type To func(<-chan file)

func (app *app) toDisk(files <-chan file) {
	const maxWorkers = 3

	var wg sync.WaitGroup

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for file := range files {
				path := common.DataFolder
				if file.resized {
					path = common.ResizedImagesFolder
				}
				path += "/" + file.name

				f, err := os.Create(path)
				if err != nil {
					app.logger.Errorf("can't create a file: %s\n", err)
					continue
				}

				var src io.Reader = file.r
				if app.config.Disk.Encrypted {
					// Have to encrypt a file
					buff := buffer.NewBuffer(nil)
					_, err := sio.Encrypt(buff, file.r, sio.Config{Key: app.config.Disk.PassPhrase[:]})
					if err != nil {
						app.logger.Errorf("can't encrypt file '%s': %s\n", file.name, err)
						continue
					}

					// Update the src
					src = buff
				}
				io.Copy(f, src)

				// Close created file
				f.Close()
				// Close the reader
				file.r.Close()
			}
		}()
	}

	wg.Wait()
}

func (app *app) toS3(files <-chan file) {
	const maxWorkers = 3

	var wg sync.WaitGroup

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for file := range files {
				bucket := common.DataBucket
				if file.resized {
					bucket = common.ResizedImagesBucket
				}
				key := file.name

				_, err := app.s3.PutObject(bucket, key, file.r, file.size, minio.PutObjectOptions{})
				if err != nil {
					app.logger.Errorf("can't put a file: %s\n", err)
					continue
				}

				// Close reader
				file.r.Close()
			}
		}()
	}

	wg.Wait()
}

func StartMigrator(version string) <-chan struct{} {
	logger := clog.NewProdConfig().PrintTime(false).Build()

	logger.Infoln("init Migrator")

	app, err := newApp(os.Args[1:])
	if err != nil {
		logger.Fatalf("can't init a new app: %s\n", err)
	}

	app.setLogger(logger)

	err = app.prepare()
	if err != nil {
		logger.Fatalf("can't prepare app: %s\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		app.start()

		cancel()
	}()
	printWaitingMessage(ctx, logger, "migration")

	logger.Infoln("migration is finished")

	done := make(chan struct{})
	close(done)
	return done
}

func printWaitingMessage(ctx context.Context, l *clog.Logger, msg string) {
	usedSpace := len(msg) + 9 // "[INF] " + "..." = 9

	var clearingString string
	for i := 0; i < usedSpace; i++ {
		clearingString += " "
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	// Print for the first time
	l.Info(msg)

	n := 0
	for {
		select {
		case <-ticker.C:
			if n > 3 {
				n = 0
				// Clear space and update the message
				l.Print("\r", clearingString, "\r")
				l.Info(msg)
			}

			if n != 0 {
				l.Print(".")
			}
			n++
		case <-ctx.Done():
			l.Print("\n")
			return
		}
	}
}
