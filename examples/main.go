package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/budougumi0617/msstore-go"
)

const (
	// azureEndpoint is end ponint of Azure Active Directory.
	azureEndpoint = "https://login.microsoftonline.com"
	// storeEndpoint is end point of Microsoft Store API.
	storeEndpoint = "https://manage.devcenter.microsoft.com"
)

func main() {
	var (
		tenantIDFlag  string
		clientIDFlag  string
		secretKeyFlag string
	)
	/* register flag name and shorthand name */
	flag.StringVar(&tenantIDFlag, "tenant", "blank", "Tenant ID")
	flag.StringVar(&tenantIDFlag, "t", "blank", "shorthand of tenantid")
	flag.StringVar(&clientIDFlag, "clientid", "blank", "Client ID")
	flag.StringVar(&clientIDFlag, "c", "blank", "shorthand of clientid")
	flag.StringVar(&secretKeyFlag, "key", "blank", "Secret Key")
	flag.StringVar(&secretKeyFlag, "k", "blank", "shorthand of key")
	flag.Parse()
	tenantID := tenantIDFlag

	u, err := url.Parse(azureEndpoint)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	u.Path = path.Join(u.Path, tenantID, "oauth2/token")
	data := url.Values{
		"grant_type": {"client_credentials"},
		"resource":   {storeEndpoint},
	}
	data.Add("client_id", clientIDFlag)
	data.Add("client_secret", secretKeyFlag)
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	defer resp.Body.Close()

	var r io.Reader = resp.Body
	r = io.TeeReader(r, os.Stderr)

	var result msstore.AADTokenResponese
	err = json.NewDecoder(r).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result.AccessToken)

	// GET Srore information
	storeBase, err := url.Parse(storeEndpoint)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	storeBase.Path = path.Join(storeBase.Path, "v1.0/my/applications")
	req2, err := http.NewRequest("GET", storeBase.String(), nil)
	req2.Header.Set("Authorization", "Bearer "+result.AccessToken)
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	defer resp2.Body.Close()

	var r2 io.Reader = resp2.Body
	r2 = io.TeeReader(r2, os.Stderr)

	var result2 msstore.GetAppDataResponse
	err = json.NewDecoder(r2).Decode(&result2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result2.Value)

}
