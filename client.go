package msstore

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

const (
	// azureEndpoint is end ponint of Azure Active Directory.
	azureEndpoint = "https://login.microsoftonline.com"
	// storeEndpoint is end point of Microsoft Store API.
	storeEndpoint = "https://manage.devcenter.microsoft.com"

	defaultHTTPTimeout = 120 * time.Second
)

var (
	httpClient = &http.Client{Timeout: defaultHTTPTimeout}
	userAgent  = fmt.Sprintf("msstore-go (%s)", runtime.Version())
)

// Client is custom httpclient.
type Client struct {
	URL           *url.URL
	client        *http.Client
	Authorization string
	TenantID      string
}

// AADTokenResponese defines response from "/oauth2/token".
type AADTokenResponese struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
	ExtExpiresIn string `json:"ext_expires_in"`
	ExpiresOn    string `json:"expires_on"`
	NotBefore    string `json:"not_before"`
	Resource     string `json:"resource"`
	AccessToken  string `json:"access_token"`
}

// GetAppDataResponse defines response from "/v1.0/my/applications"
type GetAppDataResponse struct {
	Value []struct {
		ID                           string    `json:"id"`
		PrimaryName                  string    `json:"primaryName"`
		PackageFamilyName            string    `json:"PackageFamilyName"`
		PackageIdentityName          string    `json:"packageIdentityName"`
		PublisherName                string    `json:"publisherName"`
		FirstPublishedDate           time.Time `json:"firstPublishedDate"`
		PendingApplicationSubmission struct {
			ID               string `json:"id"`
			ResourceLocation string `json:"resourceLocation"`
		} `json:"pendingApplicationSubmission,omitempty"`
		HasAdvancedListingPermission       bool `json:"hasAdvancedListingPermission"`
		LastPublishedApplicationSubmission struct {
			ID               string `json:"id"`
			ResourceLocation string `json:"resourceLocation"`
		} `json:"lastPublishedApplicationSubmission,omitempty"`
	} `json:"value"`
	TotalCount int `json:"totalCount"`
}

// NewClient returns Intial Client.
func NewClient(tenantID string) (*Client, error) {
	if len(tenantID) == 0 {
		return nil, errors.New("tenantID is nil")
	}
	return &Client{TenantID: tenantID}, nil
}
