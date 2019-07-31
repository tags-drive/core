package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testFolder = "./test"
	dataFolder = "./var"
	//
	host = "http://localhost"
)

// TestMain setups environment and runs tests
//
// It changes the current directory to 'testFolder'. After finish of all tests it deletes this folder.
func TestMain(m *testing.M) {
	err := os.Mkdir(testFolder, 0666)
	if err != nil {
		log.Fatalf("can't create test folder '%s': %s\n", testFolder, err)
	}

	os.Chdir(testFolder)
	code := m.Run()
	os.Chdir("..")

	err = os.RemoveAll(testFolder)
	if err != nil {
		log.Fatalf("can't delete test folder '%s': %s\n", testFolder, err)
	}

	os.Exit(code)
}

func getCommonEnvVars() map[string]string {
	return map[string]string{
		"DEBUG":               "false",
		"WEB_PORT":            ":80",
		"WEB_TLS":             "false",
		"WEB_LOGIN":           "login",
		"WEB_PASSWORD":        "password",
		"WEB_SKIP_LOGIN":      "false",
		"STORAGE_PASS_PHRASE": "pass_phrase",

		// "WEB_MAX_TOKEN_LIFE": "",
		// "STORAGE_ENCRYPT": "",
		// "STORAGE_TIME_BEFORE_DELETING": "",
	}
}

// prepareAppInstance creates a new App instance and configures it.
func prepareAppInstance(assert *assert.Assertions) *App {
	app, err := PrepareNewApp()
	if !assert.Nil(err) {
		assert.FailNow("can't Prepare a new app")
	}

	err = app.ConfigureServices()
	if !assert.Nil(err) {
		assert.FailNow("can't init services")
	}

	return app
}

func runIntegrationTests(t *testing.T) {
	tests := []struct {
		description string
		test        func(*assert.Assertions)
	}{
		{
			description: "create a new tag",
			test: func(assert *assert.Assertions) {
				const code = http.StatusOK

				body := marshall(
					struct {
						test string
					}{},
				)

				r, err := http.NewRequest("POST", host+"/api/tags", body)
				panicIfNotNil(err)

				resp, err := http.DefaultClient.Do(r)
				panicIfNotNil(err)

				assert.Equal(code, resp.StatusCode)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			tt.test(assert.New(t))
		})
	}

	// TODO: make some requests
	time.Sleep(time.Second * 2)
}

// marshall encode v as json. It panics if an error occurres
func marshall(v interface{}) io.Reader {
	b := &bytes.Buffer{}

	err := json.NewEncoder(b).Encode(v)
	panicIfNotNil(err)

	return b
}

func panicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}

// clear removes the test folder and all subfolders ("./test/var/*").
// Current directory must be 'testFolder'!
func clear() {
	os.RemoveAll(dataFolder)
}

func TestTagsDrive(t *testing.T) {
	t.Run("First launch (no 'var' folder)", func(t *testing.T) {
		defer clear()

		assert := assert.New(t)

		// Setup envs
		envs := getCommonEnvVars()
		envs["STORAGE_ENCRYPT"] = "false"
		for k, v := range envs {
			os.Setenv(k, v)
		}

		app := prepareAppInstance(assert)

		// Start the app and check an error after shutdown
		serverChecked := make(chan struct{})
		go func() {
			if err := app.Start(); err != nil {
				assert.FailNow("app.Start finished with non-nil error:", err)
			}

			close(serverChecked)
		}()

		// Run tests
		runIntegrationTests(t)

		// Shutdown
		app.Shutdown()

		// Wait for server shutdown
		<-serverChecked

	})

	t.Run("Restart ('var' folder exists)", func(t *testing.T) {
		// TODO
	})

	t.Run("Without encryption", func(t *testing.T) {
		// TODO
	})

	t.Run("With encryption", func(t *testing.T) {
		// TODO
	})
}
