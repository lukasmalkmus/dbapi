package main

import (
	"fmt"
	"os"
	"time"

	"github.com/LukasMa/dbapi"
	"github.com/urfave/cli"
)

var (
	// CommitHash is the git commit hast that describes the state of the
	// repository the application was build from.
	CommitHash = "Unknown"
	// CompileTime is the build date/ime of the application.
	CompileTime time.Time
	// Version is the applications version.
	Version = "Unknown"
)

var (
	// AccessToken is the OAuth2 AccessToken.
	AccessToken string
)

var api *dbapi.Client

func main() {
	api, _ := dbapi.New()

	app := cli.NewApp()
	app.Name = "dbcli"
	app.Usage = "A small command line application to test the LukasMa/dbapi client library."
	app.Compiled = CompileTime
	app.Version = Version
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Lukas Malkmus",
			Email: "mail@lukasmalkmus.com",
		},
	}
	app.Copyright = "(c) 2016 Lukas Malkmus"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "token, t",
			Usage:       "access token `TOKEN`",
			EnvVar:      "TOKEN",
			Destination: &AccessToken,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "accounts",
			Usage: "access the /cashAccounts endpoint",
			Action: func(c *cli.Context) error {
				var res *dbapi.Accounts
				var err error

				api.Options(dbapi.SetToken(AccessToken))
				if !api.Authentication.HasAuth() {
					fmt.Println("No AccessToken provided!")
					return nil
				}

				if iban := c.Args().First(); iban != "" {
					res, _, err = api.Accounts.Get(iban)
				} else {
					res, _, err = api.Accounts.GetAll()
				}
				if err != nil {
					return err
				}

				fmt.Printf("\nAccount(s):\n\n")
				for _, acc := range *res {
					fmt.Printf("\t%s\n\t%.2f€\n\t%s\n\n", acc.Iban, acc.Balance, acc.ProductDescription)
				}

				return nil
			},
		},
		{
			Name:  "addresses",
			Usage: "access the /addresses endpoint",
			Action: func(c *cli.Context) error {
				api.Options(dbapi.SetToken(AccessToken))
				if !api.Authentication.HasAuth() {
					fmt.Println("No AccessToken provided!")
					return nil
				}

				res, _, err := api.Addresses.Get()
				if err != nil {
					return err
				}

				fmt.Printf("\nAddress(es):\n\n")
				for _, addr := range *res {
					fmt.Printf("\t%s\n\t%s %d\n\t%d %s\n\t%s\n\n", addr.Type, addr.Street, addr.HouseNumber, addr.ZipCode, addr.City, addr.Country)
				}

				return nil
			},
		},
		{
			Name:  "transactions",
			Usage: "access the /transactions endpoint",
			Action: func(c *cli.Context) error {
				var res *dbapi.Transactions
				var err error

				api.Options(dbapi.SetToken(AccessToken))
				if !api.Authentication.HasAuth() {
					fmt.Println("No AccessToken provided!")
					return nil
				}

				if iban := c.Args().First(); iban != "" {
					res, _, err = api.Transactions.Get(iban)
				} else {
					res, _, err = api.Transactions.GetAll()
				}
				if err != nil {
					return err
				}

				fmt.Printf("\nTransaction(s):\n\n")
				for _, trans := range *res {
					fmt.Printf("\t%.2f€\n\t%s\n\t%s\n\tFrom/To: %s <%s>\n\n", trans.Amount, trans.Date, trans.Usage, trans.CounterPartyName, trans.CounterPartyIBAN)
				}

				return nil
			},
		},
		{
			Name:  "userinfo",
			Usage: "access the /userInfo endpoint",
			Action: func(c *cli.Context) error {
				api.Options(dbapi.SetToken(AccessToken))
				if !api.Authentication.HasAuth() {
					fmt.Println("No AccessToken provided!")
					return nil
				}

				res, _, err := api.UserInfo.Get()
				if err != nil {
					return err
				}

				fmt.Printf("\nUser Info:\n\n")
				fmt.Printf("\t%s\n\t%s %s\n\t%s\n", res.Gender, res.FirstName, res.LastName, res.DateOfBirth)

				return nil
			},
		},
	}

	app.Run(os.Args)

	/*
		    app.Action = func(c *cli.Context) error {
				// Create a new client.
				AccessToken := "1234567890abcdefghijklmnopqrstuvwxyz"
				api, err := dbapi.New(
					dbapi.SetToken(AccessToken),
				)
				if err != nil {
					log.Fatalln(err)
					return err
				}

				// Start to access the Deutsche Bank API (retrieve and print user accounts).
				accounts, response, err := api.Accounts.GetAll()
				if err != nil {
					fmt.Println(response.Response)
					log.Fatalln(err)
					return err
				}
				fmt.Printf("%v", accounts)
				return nil
			}
	*/
}
