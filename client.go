package msstore

import "net/http"
import "net/url"

const (
	apiEndpoint = "https://login.microsoftonline.com"
)

// Client is custom httpclient.
type Client struct {
	URL           *url.URL
	client        *http.Client
	Authorization string
	TenantID      string
}
