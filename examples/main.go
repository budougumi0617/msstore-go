package main

import (
	"flag"
	"log"

	"github.com/budougumi0617/msstore-go"
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

	c, err := msstore.NewClient(tenantIDFlag, clientIDFlag, secretKeyFlag)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	c.Init()

	c.GetMyApps()
}
