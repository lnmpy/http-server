package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
)

type fileHandler struct {
	handler http.Handler
	out     io.Writer
}

func (h *fileHandler) stats(req *http.Request) {
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

func (h *fileHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.handler.ServeHTTP(rw, r)
	h.stats(r)
}

func newFileHandler(handler http.Handler) http.Handler {
	return &fileHandler{handler: handler}
}
