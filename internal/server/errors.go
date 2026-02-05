package server

import (
	"errors"
	"net/http"
	"time"
)

func httpError(w http.ResponseWriter, err error, code int) {
	errDTO := ErrorDTO{
		Message: err.Error(),
		Time:    time.Now(),
	}

	http.Error(w, errDTO.ToString(), code)
}

func httpErrorIs(target, err error, w http.ResponseWriter) {
	errDTO := ErrorDTO{
		Message: err.Error(),
		Time:    time.Now(),
	}

	if errors.Is(err, target) {
		http.Error(w, errDTO.ToString(), http.StatusNotFound)
	} else {
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
	}
}
