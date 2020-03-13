package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/urfave/negroni"

	"github.com/gorilla/schema"
	"github.com/leeif/go-web-template/datatype/error"
	"github.com/leeif/go-web-template/datatype/request"
	resp "github.com/leeif/go-web-template/datatype/response"
	"github.com/leeif/go-web-template/utils/view"
)

func getBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

func getBody(r *http.Request, reciver interface{}) *error.Error {

	contentType := r.Header.Get("Content-type")
	if strings.Contains(contentType, "application/json") {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return error.ServerError.Wrapper(errors.New("Read body failed: " + err.Error()))
		}
		err = json.Unmarshal(body, &reciver)
		if err != nil {
			return error.BadRequest
		}
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		err := r.ParseForm()
		if err != nil {
			return error.BadRequest
		}
		decoder := schema.NewDecoder()
		err = decoder.Decode(reciver, r.PostForm)
		if err != nil {
			return error.BadRequest.Wrapper(err)
		}
	}
	pr, ok := reciver.(request.Request)
	// check request body validation
	if ok && !pr.Validation() {
		return error.BadRequest
	}
	return nil
}

func getQuery(r *http.Request, reciver interface{}) *error.Error {
	decoder := schema.NewDecoder()
	if err := decoder.Decode(reciver, r.URL.Query()); err != nil {
		return error.BadRequest.Wrapper(err)
	}

	pr, ok := reciver.(request.Request)
	// check request body validation
	if ok && !pr.Validation() {
		return error.BadRequest
	}
	return nil
}

func responseOK(body interface{}, w http.ResponseWriter) *error.Error {
	response := resp.ReponseOK{}
	response.Status = resp.STATUSOK
	response.Body = body
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(response)
	if err != nil {
		return error.ServerError.Wrapper(err)
	}
	_, err = w.Write(b)
	if err != nil {
		return error.ServerError.Wrapper(err)
	}
	return nil
}

func responseError(e *error.Error, w http.ResponseWriter) *error.Error {
	response := resp.ReponseError{}
	response.Status = resp.STATUSERROR
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(e.HTTPCode)

	m := make(map[string]interface{})
	m["code"] = e.Code
	m["message"] = e.HTTPError.Error()
	response.Error = m
	b, err := json.Marshal(response)
	if err != nil {
		return error.ServerError.Wrapper(err)
	}
	_, err = w.Write(b)
	if err != nil {
		return error.ServerError.Wrapper(err)
	}
	return nil
}

func responseHTMLOK(file string, data interface{}, w http.ResponseWriter) *error.Error {
	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	t, err := view.Parse(file)
	if err != nil {
		return error.ServerError.Wrapper(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		return error.ServerError.Wrapper(err)
	}
	return nil
}

func responseHTMLError(file string, data interface{}, w http.ResponseWriter, status int) *error.Error {
	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(status)
	t, err := view.Parse(file)
	if err != nil {
		return error.ServerError.Wrapper(err)
	}
	if err != nil {
		return error.ServerError.Wrapper(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		return error.ServerError.Wrapper(err)
	}
	return nil
}

func (route *Route) handlerWrapper(handler func(http.ResponseWriter, *http.Request) *error.Error) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if err := handler(w, r); err != nil {
			responseError(err, w)
			route.log(err, r)
			return
		}

		next(w, r)
	}
}

func (route *Route) webHandlerWrapper(handler func(http.ResponseWriter, *http.Request) *error.Error) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if err := handler(w, r); err != nil {
			responseHTMLError("error.html", nil, w, http.StatusInternalServerError)
			route.log(err, r)
			return
		}

		next(w, r)
	}
}

func (route *Route) log(e *error.Error, request *http.Request) {
	url := request.URL.String()
	if e.LogError != nil {
		route.logger.Error(fmt.Sprintf("[%s]:%s", url, e.LogError.Error()))
	}
	if e.HTTPError != nil {
		route.logger.Debug(fmt.Sprintf("[%s]:%s", url, e.HTTPError.Error()))
	}
}

// Redirects to a new path while keeping current request's query string
func redirectWithQueryString(to string, query url.Values, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("%s%s", to, getQueryString(query)), http.StatusFound)
}

// Redirects to a new path with the query string moved to the URL fragment
func redirectWithFragment(to string, query url.Values, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("%s#%s", to, query.Encode()), http.StatusFound)
}

// Returns string encoded query string of the request
func getQueryString(query url.Values) string {
	encoded := query.Encode()
	if len(encoded) > 0 {
		encoded = fmt.Sprintf("?%s", encoded)
	}
	return encoded
}
