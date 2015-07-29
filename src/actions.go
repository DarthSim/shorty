package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func checkServerError(rw http.ResponseWriter, req *http.Request) {
	if err := recover(); err != nil {
		if err == sql.ErrNoRows {
			serverResponse(rw, req, "url not found", 404)
		} else {
			serverError(rw, req, err, 500)
		}
	}
}

func getCode(req *http.Request) string {
	return mux.Vars(req)["code"]
}

func createUrlHandler(rw http.ResponseWriter, req *http.Request) {
	defer checkServerError(rw, req)

	checkErr(req.ParseForm())

	code, err := createUrl(req.Form.Get("url"))
	checkErr(err)

	shortUrl := fmt.Sprintf("http://shorty.com/%s", code)
	serverResponse(rw, req, shortUrl, 200)
}

func redirectHandler(rw http.ResponseWriter, req *http.Request) {
	defer checkServerError(rw, req)

	url, err := getUrl(getCode(req))
	checkErr(err)

	checkErr(hitRedirect(getCode(req)))

	http.Redirect(rw, req, url, 301)
}

func expandHandler(rw http.ResponseWriter, req *http.Request) {
	defer checkServerError(rw, req)

	url, err := getUrl(getCode(req))
	checkErr(err)

	serverResponse(rw, req, url, 200)
}

func statisticsHandler(rw http.ResponseWriter, req *http.Request) {
	defer checkServerError(rw, req)

	count, err := getOpenCount(getCode(req))
	checkErr(err)

	serverResponse(rw, req, fmt.Sprintf("%d", count), 200)
}
