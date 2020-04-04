package main_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	neturl "net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-transaction-example/internal/pkg/must"
	"go-transaction-example/internal/platform/service"
)

var config service.Config
var wg sync.WaitGroup

func TestMain(m *testing.M) {
	var err error

	config, err = service.ParseConfig()
	must.NotFail(err)

	wg.Add(1)
	go func() {
		wg.Done()
		must.NotFail(service.Run(config))
	}()

	os.Exit(m.Run())
}

func TestPostUsers(t *testing.T) {
	start := time.Now()

	name := "dani"
	age := 25

	url := fmt.Sprintf("http://%s:%d/users?name=%s&age=%d",
		"localhost", config.Server.Port,
		name, age,
	)

	res, err := http.Post(url, "application/x-www-urlencoded", nil)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	url = fmt.Sprintf("http://%s:%d/users/history?from=%s",
		"localhost", config.Server.Port,
		neturl.QueryEscape(start.Format(time.RFC3339)),
	)

	res, err = http.Get(url)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)
	defer res.Body.Close()

	var body []map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&body)
	if !assert.NoError(t, err) {
		return
	}

	t.Log(body)
}
