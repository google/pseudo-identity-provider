// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package idp

import (
	"bytes"
	"crypto/rand"
	sessionmgmt "customidp/session"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"sync"
	"time"

	"google.golang.org/appengine/v2"
)

// Entry defines a log entry.
type Entry struct {
	Message  string `json:"message"`
	Severity string `json:"severity,omitempty"`
	Trace    string `json:"logging.googleapis.com/trace,omitempty"`

	// Logs Explorer allows filtering and display of this as `jsonPayload.component`.
	Component string `json:"component,omitempty"`
}

// RequestEntry log an incoming Request and its response.
type RequestEntry struct {
	time   time.Time
	input  *sessionmgmt.RequestInput
	action string
	resp   *ResponseEntry
}

// ResponseEntry stores the response data.
type ResponseEntry struct {
	headers http.Header
	body    string
}

// Global mutex protected storage for the Request Log.
var requestLog map[string]RequestEntry
var logMutex sync.Mutex

// init sets up the request log.
func init() {
	requestLog = make(map[string]RequestEntry)
}

// String renders an entry structure to the JSON format expected by Cloud Logging.
func (e Entry) String() string {
	if e.Severity == "" {
		e.Severity = "INFO"
	}
	out, err := json.Marshal(e)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}
	return string(out)
}

// logHandler writes out the request logs. This is protected by authorization since
// the request logs can have codes and tokens.
func logHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(w, r) {
		return
	}
	writeRequestLog(w)
}

// Cloud logging for a Notice.
func logNotice(msg string, r *http.Request) {
	logMsg(msg, r, "NOTICE")
}

// Cloud logging for an Error.
func logError(msg string, r *http.Request) {
	logMsg(msg, r, "Error")
}

// logMsg logs in format used by GCP Logs.
func logMsg(msg string, r *http.Request, sev string) {
	var trace string
	traceHeader := r.Header.Get("X-Cloud-Trace-Context")
	traceParts := strings.Split(traceHeader, "/")
	if appengine.IsAppEngine() && len(traceParts) > 0 && len(traceParts[0]) > 0 {
		trace = fmt.Sprintf("projects/%s/traces/%s", appengine.AppID(appengine.NewContext(r)), traceParts[0])
	}

	log.Println(Entry{
		Severity:  sev,
		Message:   msg,
		Component: "request",
		Trace:     trace,
	})
}

// addRequestLogEntry adds a new request's info to the RequestLog
func addRequestLogEntry(input *sessionmgmt.RequestInput, action string) {
	traceHeader := input.Headers.Get("X-Cloud-Trace-Context")
	time := time.Now()
	logMutex.Lock()
	requestLog[traceHeader] = RequestEntry{time: time, input: input, action: action}
	logMutex.Unlock()
}

// Generate a random base64 url encoded string
func generateBase64ID(size int) (string, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	encoded := base64.URLEncoding.EncodeToString(b)
	return encoded, nil
}

// respLogHandler is a wrapper around another handler that allows us
// to log the Response to the RequestLog.
func respLogHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceHeader := r.Header.Get("X-Cloud-Trace-Context")
		if traceHeader == "" {
			id, err := generateBase64ID(32)
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			}
			r.Header.Add("X-Cloud-Trace-Context", id)
			traceHeader = id
		}

		rec := httptest.NewRecorder()
		fn(rec, r)

		logMutex.Lock()
		req, ok := requestLog[traceHeader]
		if !ok {
			requestLog[traceHeader] = RequestEntry{resp: getResponseEntry(rec.Result())}
		} else {
			req.resp = getResponseEntry(rec.Result())
			requestLog[traceHeader] = req
		}
		logMutex.Unlock()

		for k, v := range rec.Result().Header {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		rec.Body.WriteTo(w)
	}
}

func getResponseEntry(resp *http.Response) *ResponseEntry {
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	// Try to pretty-print if content is JSON.
	jsonBytes := new(bytes.Buffer)
	newStr := ""
	if err := json.Indent(jsonBytes, buf.Bytes(), "", "  "); err != nil {
		newStr = buf.String()
	} else {
		newStr = jsonBytes.String()
	}

	return &ResponseEntry{
		headers: resp.Header,
		body:    newStr,
	}
}

// writeRequestLog outputs the RequestLog as an HTML Table.
func writeRequestLog(w http.ResponseWriter) {
	logMutex.Lock()
	defer logMutex.Unlock()

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// We escape all user provided input, but as a defense-in-depth, only allow style
	// imports from self. Don't allow anything else including inline styles or script.
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'")

	entries := []RequestEntry{}
	for _, v := range requestLog {
		entries = append(entries, v)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].time.Before(entries[j].time)
	})

	fmt.Fprint(w, "<html><head><link rel='stylesheet' href='log.css'></head>")

	// Show log with more recent items first.
	for i := len(entries) - 1; i >= 0; i-- {
		req := entries[i]
		fmt.Fprint(w, "<table>")
		writeRow(w, "Id:", fmt.Sprint(i))
		writeRow(w, "Time:", req.time.Local().Format(time.ANSIC))
		writeRow(w, "Path:", req.input.Path)
		writeRow(w, "Method:", req.input.HTTPMethod)
		writeRow(w, "Proto:", req.input.Proto)
		writeRow(w, "Action Taken:", req.action)
		writeParams(w, req.input)

		if req.input.Session != nil {
			writeRow(w, "Session Code:", req.input.Session.Code)
			writeRow(w, "Session ClientID:", req.input.Session.ClientID)
			writeRow(w, "Session RedirectURI:", req.input.Session.RedirectURI)
			writeRow(w, "Session Nonce:", req.input.Session.Nonce)
			writeRow(w, "Session Challenge:", req.input.Session.CodeChallenge)
			writeRow(w, "Session Challenge Method:", req.input.Session.CodeChallengeMethod)
		}

		if req.resp != nil {
			writeRow(w, "Response:")
			writeRow(w, "Headers:")
			for k, v := range req.resp.headers {
				writeRow(w, append([]string{"", k}, v...)...)
			}

			writeWideRow(w, "Body:", req.resp.body)
		}

		fmt.Fprint(w, "</table>")
	}

	fmt.Fprint(w, "</html>")
}

// writeRow writes one table row.
func writeRow(w http.ResponseWriter, columns ...string) {
	fmt.Fprint(w, "<tr>")
	for _, col := range columns {
		fmt.Fprintf(
			w,
			"<td>%s</td>",
			template.HTMLEscapeString(col))
	}
	fmt.Fprint(w, "</tr>")
}

// writeRow writes one table row.
func writeWideRow(w http.ResponseWriter, columns ...string) {
	fmt.Fprint(w, "<tr>")
	for i, col := range columns {
		if i == 0 {
			fmt.Fprintf(
				w,
				"<td>%s</td>",
				template.HTMLEscapeString(col))
		} else {
			fmt.Fprintf(
				w,
				"<td colspan='2'><div class='widerowcontent'>%s</div></td>",
				template.HTMLEscapeString(col))
		}
	}
	fmt.Fprint(w, "</tr>")
}

// writeParams is a helper to output HTTP parameters.
func writeParams(w http.ResponseWriter, input *sessionmgmt.RequestInput) {
	writeRow(w, "URL Params:")
	for k, v := range input.URLParams {
		writeRow(w, append([]string{"", k}, v...)...)
	}

	if len(input.FormParams) != 0 {
		writeRow(w, "Form Params:")
		for k, v := range input.FormParams {
			writeRow(w, append([]string{"", k}, v...)...)
		}
	}

	writeRow(w, "Headers:")
	for k, v := range input.Headers {
		// Ignore the X-* params.s
		if !strings.HasPrefix(k, "X-") {
			writeRow(w, append([]string{"", k}, v...)...)
		}
	}
}
