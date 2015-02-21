package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

func checkServerError(rw http.ResponseWriter, err error) bool {
	switch {
	case err == sql.ErrNoRows:
		serverResponse(rw, "url not found", 404)
		return true
	case err != nil:
		serverError(rw, err, 500)
		return true
	}

	return false
}

func createUrlHandler(rw http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if checkServerError(rw, err) {
		return
	}

	url := req.Form.Get("url")

	if url == "" {
		serverResponse(rw, "url should be defined", 422)
		return
	}

	code, err := createUrl(url)
	if checkServerError(rw, err) {
		return
	}

	shortUrl := fmt.Sprintf("http://%s/%s", config.Url.Domain, code)

	serverResponse(rw, shortUrl, 200)
}

func redirectHandler(rw http.ResponseWriter, req *http.Request) {
	code := requestVars(req)["code"]

	url, err := getUrl(code)
	if checkServerError(rw, err) {
		return
	}

	err = hitRedirect(code)
	if checkServerError(rw, err) {
		return
	}

	http.Redirect(rw, req, url, 301)
}

func expandHandler(rw http.ResponseWriter, req *http.Request) {
	code := requestVars(req)["code"]

	url, err := getUrl(code)
	if checkServerError(rw, err) {
		return
	}

	serverResponse(rw, url, 200)
}

func statisticsHandler(rw http.ResponseWriter, req *http.Request) {
	code := requestVars(req)["code"]

	count, err := getOpenCount(code)
	if checkServerError(rw, err) {
		return
	}

	serverResponse(rw, fmt.Sprintf("%d", count), 200)
}
