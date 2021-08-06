package handler

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rancher/apiserver/pkg/types"
	"github.com/rancher/steve/pkg/attributes"
	"github.com/rancher/steve/pkg/schema"
)

func k8sAPI(sf schema.Factory, apiOp *types.APIRequest) {
	vars := mux.Vars(apiOp.Request)
	group := vars["group"]
	if group == "core" {
		group = ""
	}

	apiOp.Name = vars["name"]
	apiOp.Type = vars["type"]

	nOrN := vars["nameorns"]
	if nOrN != "" {
		schema := apiOp.Schemas.LookupSchema(apiOp.Type)
		if attributes.Namespaced(schema) {
			vars["namespace"] = nOrN
		} else {
			vars["name"] = nOrN
		}
	}

	if namespace := vars["namespace"]; namespace != "" {
		apiOp.Namespace = namespace
	}
}

func apiRoot(sf schema.Factory, apiOp *types.APIRequest) {
	apiOp.Type = "apiRoot"
}

func rewriteLocalCluster(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if strings.HasPrefix(req.URL.Path, "/k8s/clusters/local") {
			req.URL.Path = strings.TrimPrefix(req.URL.Path, "/k8s/clusters/local")
			if req.URL.Path == "" {
				req.URL.Path = "/"
			}
		}
		next.ServeHTTP(rw, req)
	})
}
