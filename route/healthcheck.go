package route

import (
	"net/http"

	"github.com/leeif/go-web-template/datatype/error"
)

func (route *Route) healthcheck(w http.ResponseWriter, r *http.Request) *error.Error {
	err := route.manager.Healthcheck()
	if err != nil {

	}

	respBody := make(map[string]interface{})
	respBody["version"] = route.config.Version

	if err := responseOK(respBody, w); err != nil {
		route.logger.Error(err.LogError)
	}
	return nil
}
