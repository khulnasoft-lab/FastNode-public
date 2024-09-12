package mockserver

import (
	"net/http/httptest"
	"net/url"

	"github.com/gorilla/mux"
)

// NewFastnodedTestServer returns a new mock instance of fastnoded
func NewFastnodedTestServer() (*FastnodedTestServer, error) {
	router := mux.NewRouter()
	httpd := httptest.NewServer(router)

	url, err := url.Parse(httpd.URL)
	if err != nil {
		return nil, err
	}

	return &FastnodedTestServer{
		server: httpd,
		URL:    url,
		Router: router,
	}, nil
}

// FastnodedTestServer is provides a mocked
type FastnodedTestServer struct {
	server *httptest.Server
	URL    *url.URL
	Router *mux.Router
}

// Close releases the resources used by the fastnoded server
func (t *FastnodedTestServer) Close() {
	t.server.Close()
}
