package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type ActionsTestSuite struct {
	suite.Suite
	Router   *mux.Router
	Response *httptest.ResponseRecorder
}

func (suite *ActionsTestSuite) SetupSuite() {
	config = Config{}

	config.Url.Domain = "shorty.test"

	config.Database.Host = "localhost"
	config.Database.Database = "shorty_test"
	config.Database.User = os.Getenv("DBUSER")
	config.Database.Password = os.Getenv("DBPASS")

	config.Database.MaxOpenConnections = 5
	config.Database.MaxIdleConnections = 5

	config.Database.InitSchema = true

	config.Log.Path = "test.log"

	initLogger()

	initDB()

	suite.Router = setupRouter()
}

func (suite *ActionsTestSuite) TearDownSuite() {
	db.Exec("DROP TABLE urls;")
	closeDB()

	closeLogger()
	os.Remove(absPathToFile(config.Log.Path))
}

func (suite *ActionsTestSuite) SetupTest() {
}

func (suite *ActionsTestSuite) TearDownTest() {
	db.Exec("DELETE FROM urls;")
}

func (suite *ActionsTestSuite) SendRequest(method, path string, body ...string) error {
	suite.Response = httptest.NewRecorder()

	var (
		req *http.Request
		err error
	)

	if method == "POST" {
		reqBody := strings.NewReader(body[0])
		req, err = http.NewRequest(method, "http://shorty.test"+path, reqBody)
	} else {
		req, err = http.NewRequest(method, "http://shorty.test"+path, nil)
	}

	if err != nil {
		return err
	}

	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	suite.Router.ServeHTTP(suite.Response, req)
	return nil
}

func (suite *ActionsTestSuite) TestCreateUrl() {
	suite.Nil(suite.SendRequest(
		"POST",
		"/shorten",
		"url=http%3A%2F%2Fgoogle.com%2F",
	))

	var (
		url, code string
		openCount int64
	)

	suite.Nil(
		db.QueryRow("SELECT url, code, open_count FROM urls").
			Scan(&url, &code, &openCount),
	)

	suite.Equal("http://google.com/", url)
	suite.NotEmpty(code)
	suite.Equal(0, openCount)

	suite.Equal(200, suite.Response.Code)
	suite.Equal(
		fmt.Sprintf("http://%s/%s", config.Url.Domain, code),
		string(suite.Response.Body.Bytes()),
	)
}

func (suite *ActionsTestSuite) TestCreateUrlTwice() {
	suite.Nil(suite.SendRequest(
		"POST",
		"/shorten",
		"url=http%3A%2F%2Fgoogle.com%2F",
	))

	suite.Nil(suite.SendRequest(
		"POST",
		"/shorten",
		"url=http%3A%2F%2Fgoogle.com%2F",
	))

	rows, err := db.Query("SELECT code FROM urls")
	defer rows.Close()

	suite.Nil(err)

	var code1, code2 string

	rows.Next()
	rows.Scan(&code1)

	rows.Next()
	rows.Scan(&code2)

	suite.NotEqual(code1, code2, "Codes should not be equal even for equal urls")
}

func (suite *ActionsTestSuite) TestExpandUrl() {
	db.Exec("INSERT INTO urls (url, code) VALUES ('http://google.com/', 'abcd')")

	suite.Nil(suite.SendRequest(
		"GET",
		"/expand/abcd",
	))

	suite.Equal(200, suite.Response.Code)
	suite.Equal(
		"http://google.com/",
		string(suite.Response.Body.Bytes()),
	)
}

func (suite *ActionsTestSuite) TestRedirectToUrl() {
	db.Exec("INSERT INTO urls (url, code) VALUES ('http://google.com/', 'abcd')")

	suite.Nil(suite.SendRequest(
		"GET",
		"/abcd",
	))

	suite.Equal(301, suite.Response.Code)
	suite.Equal(
		"http://google.com/",
		string(suite.Response.Header().Get("Location")),
	)
}

func (suite *ActionsTestSuite) TestStatistics() {
	db.Exec("INSERT INTO urls (url, code, open_count) VALUES ('http://google.com/', 'abcd', 1234)")

	suite.Nil(suite.SendRequest(
		"GET",
		"/statistics/abcd",
	))

	suite.Equal(200, suite.Response.Code)
	suite.Equal(
		"1234",
		string(suite.Response.Body.Bytes()),
	)
}

func TestActions(t *testing.T) {
	suite.Run(t, new(ActionsTestSuite))
}