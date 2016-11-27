// Steve Phillips / elimisteve
// 2016.06.16

package main

import (
	"fmt"
	"log"
	"net/http"
)

const contentTypeJSON = "application/json; charset=utf-8"

func writeErrorStatus(w http.ResponseWriter, errStr string, status int, secretErr error) {
	if Debug {
		log.Printf("Returning HTTP %d w/error: %q;\n  real error: %s\n",
			status, errStr, secretErr)
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"error":%q}`, errStr)
}

func writeError(w http.ResponseWriter, errStr string, secretErr error) {
	writeErrorStatus(w, errStr, http.StatusInternalServerError, secretErr)
}
