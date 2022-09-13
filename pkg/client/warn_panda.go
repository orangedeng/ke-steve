package client

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type APIWarnings struct {
	response http.ResponseWriter
}

func (am APIWarnings) HandleWarningHeader(code int, agent string, message string) {
	logrus.Infof("====HandleWarningHeader===== %s", message)
	am.response.Header().Add("X-API-Warnings", fmt.Sprintf("%d - %s", code, message))
}
