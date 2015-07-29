package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ActionsTestSuite struct {
	suite.Suite
	Response *http.Response
}

func (suite *ActionsTestSuite) SetupSuite() {
	os.Setenv("DB_CONN", "dbname=shorty_test sslmode=disable")
	os.Setenv("HOSTNAME", "the-custom-domain.shorty.com")
	os.Setenv("RESET_DB", "1")
	os.Setenv("ADDRESS", "localhost:8088")

	initDB(false)
	go startServer()
}

func (suite *ActionsTestSuite) TearDownSuite() {
	db.Exec("DROP TABLE urls;")
	closeDB()
}

func (suite *ActionsTestSuite) SetupTest() {
}

func (suite *ActionsTestSuite) TearDownTest() {
	db.Exec("DELETE FROM urls;")
}

func (suite *ActionsTestSuite) SendRequest(method, path string, body ...string) (err error) {
	var req *http.Request

	if method == "POST" {
		reqBody := strings.NewReader(body[0])
		req, err = http.NewRequest(method, "http://localhost:8088"+path, reqBody)
	} else {
		req, err = http.NewRequest(method, "http://localhost:8088"+path, nil)
	}

	if err != nil {
		return
	}

	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	suite.Response, err = http.DefaultTransport.RoundTrip(req)

	return
}

func (suite *ActionsTestSuite) ResponseBody() string {
	body, _ := ioutil.ReadAll(suite.Response.Body)
	return string(body)
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
	suite.Equal(0, int(openCount))

	suite.Equal(200, suite.Response.StatusCode)
	suite.Equal(
		fmt.Sprintf("http://the-custom-domain.shorty.com/%s", code),
		suite.ResponseBody(),
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

	suite.Equal(200, suite.Response.StatusCode)
	suite.Equal(
		"http://google.com/",
		suite.ResponseBody(),
	)
}

func (suite *ActionsTestSuite) TestRedirectToUrl() {
	db.Exec("INSERT INTO urls (url, code) VALUES ('http://google.com/', 'abcd')")

	suite.Nil(suite.SendRequest(
		"GET",
		"/abcd",
	))

	suite.Equal(301, suite.Response.StatusCode)
	suite.Equal(
		"http://google.com/",
		string(suite.Response.Header.Get("Location")),
	)
}

func (suite *ActionsTestSuite) TestStatistics() {
	db.Exec("INSERT INTO urls (url, code, open_count) VALUES ('http://google.com/', 'abcd', 1234)")

	suite.Nil(suite.SendRequest(
		"GET",
		"/statistics/abcd",
	))

	suite.Equal(200, suite.Response.StatusCode)
	suite.Equal(
		"1234",
		suite.ResponseBody(),
	)
}

func TestActions(t *testing.T) {
	suite.Run(t, new(ActionsTestSuite))
}
