package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/fatih/color"
)

type handler struct {
	listDirHandler       http.Handler
	persisRequestHandler http.Handler
}

func (h *handler) stats(req *http.Request) {
	clientIP := req.RemoteAddr
	if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
		clientIP = clientIP[:colon]
	}

	r := color.New(color.FgHiRed, color.Bold)
	w := color.New(color.FgHiWhite)

	timeFormatted := time.Now().Format("02/Jan/2006 03:04:05")
	requestLine := fmt.Sprintf("%s %s %s", req.Method, req.RequestURI, req.Proto)
	w.Printf("%s - - [%s] %s \"", clientIP, timeFormatted, req.UserAgent())
	r.Print(requestLine)
	w.Println("\"")
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		h.persisRequestHandler.ServeHTTP(rw, r)
	default:
		h.listDirHandler.ServeHTTP(rw, r)
	}
	h.stats(r)
}

func newHandler(listDirHandler http.Handler, persisRequestHandler http.Handler) http.Handler {
	return &handler{
		listDirHandler:       listDirHandler,
		persisRequestHandler: persisRequestHandler,
	}
}

type persistJsonRequestHandler struct {
}

func (h *persistJsonRequestHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Invalid http-method, need 'POST' or 'PUT'"))
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Invalid request content type, need 'application/json'"))
		return
	}

	filePath := r.URL.Path[1:]
	dirPath := path.Dir(filePath)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(fmt.Sprintf("Failed to create folder: %s", dirPath)))
			return
		}
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(fmt.Sprintf("Failed to open file: %s", filePath)))
		return
	}

	defer f.Close()

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Failed to read request body"))
		return
	}

	var i interface{}
	err = json.Unmarshal(bs, &i)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Request is invalid, json required"))
		return
	}
	bs, _ = json.Marshal(i)
	bs = append(bs, byte('\n'))
	if _, err = f.Write(bs); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Failed to write file"))
		return
	}

	rw.Write([]byte(fmt.Sprintf("Write to file %s success", filePath)))
	return
}

func newPersistJsonRequestHandler() http.Handler {
	return &persistJsonRequestHandler{}
}
