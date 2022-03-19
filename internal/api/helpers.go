package api

import (
	"encoding/json"
	"fmt"
	"github.com/LasseJacobs/go-starter-kit/pkg/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

var statusMap = map[int]string{
	100: "Continue",
	101: "Switching Protocols",
	200: "OK",
	201: "Created",
	202: "Accepted",
	203: "Non-Authoritative Information",
	204: "No Content",
	205: "Reset Content",
	206: "Partial Content",
	300: "Multiple Choices",
	301: "Moved Permanently",
	302: "Found",
	303: "See Other",
	304: "Not Modified",
	305: "Use Proxy",
	307: "Temporary Redirect",
	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	409: "Conflict",
	410: "Gone",
	411: "Length Required",
	412: "Precondition Failed",
	413: "Request Entity Too Large",
	414: "Request-URI Too Long",
	415: "Unsupported Media Type",
	416: "Requested Range Not Satisfiable",
	417: "Expectation Failed",
	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
	505: "HTTP Version Not Supported",
}

func readModel(r *http.Request, m interface{}) error {
	//todo: why is this defer here?
	defer io.Copy(ioutil.Discard, r.Body) //nolint:errcheck
	return json.NewDecoder(r.Body).Decode(m)
}

func sendJSON(w http.ResponseWriter, status int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, errors.Wrap(err, fmt.Sprintf("Error encoding json response: %v", obj)).Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	if _, err = w.Write(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func sendError(w http.ResponseWriter, err error) {
	//var errorID = "A011"
	switch e := err.(type) {
	case *model.StoryNotFoundError:
		sendHttpError(w, http.StatusNotFound, e.Error())
	default:
		// hide real error details from response to prevent info leaks
		sendHttpError(w, http.StatusInternalServerError, e.Error())
	}
}

func sendHttpError(w http.ResponseWriter, status int, errorId string) {
	w.WriteHeader(status)
	if _, writeErr := w.Write([]byte(fmt.Sprintf(`{"code":%d,"msg":"%s","error_message":"%s"}`, status, statusMap[status], errorId))); writeErr != nil {
		logrus.WithError(writeErr).Error("Error writing generic error message")
	}
}
