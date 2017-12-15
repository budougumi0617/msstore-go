package msstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"
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
	URL                           *url.URL
	authorization                 string
	tenantID, clientID, secretKey string
}

// TokenResponse defines response from Azure Active Directory
type TokenResponse struct {
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
func NewClient(tenantID, clientID, secretKey string) (*Client, error) {
	if len(tenantID) == 0 {
		return nil, errors.New("tenantID is nil")
	}
	if len(clientID) == 0 {
		return nil, errors.New("clientID is nil")
	}
	if len(secretKey) == 0 {
		return nil, errors.New("secretKey is nil")
	}
	u, err := url.Parse(storeEndpoint)
	if err != nil {
		return nil, err
	}
	return &Client{
		URL:       u,
		tenantID:  tenantID,
		clientID:  clientID,
		secretKey: secretKey,
	}, nil
}

// Init gets Token from Azure Active Directory.
func (c *Client) Init() error {
	u, err := url.Parse(azureEndpoint)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, c.tenantID, "oauth2/token")
	data := url.Values{
		"grant_type": {"client_credentials"},
		"resource":   {storeEndpoint},
	}
	data.Add("client_id", c.clientID)
	data.Add("client_secret", c.secretKey)
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Set("User-Agent", userAgent)
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r io.Reader = resp.Body
	r = io.TeeReader(r, os.Stderr)

	var result TokenResponse
	err = json.NewDecoder(r).Decode(&result)
	if err != nil {
		return err
	}

	fmt.Println(result.AccessToken)
	c.authorization = result.AccessToken
	return nil
}

// GetMyApps gets Windows applications infomation from Windows Developer Center.
func (c *Client) GetMyApps() error {
	// GET Srore information
	u := c.URL
	u.Path = path.Join(c.URL.Path, "v1.0/my/applications")
	req, err := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("Authorization", "Bearer "+c.authorization)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	defer resp.Body.Close()

	var r io.Reader = resp.Body
	r = io.TeeReader(r, os.Stderr)

	var result GetAppDataResponse
	err = json.NewDecoder(r).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result.Value)
	return nil
}
