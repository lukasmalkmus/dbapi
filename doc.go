/*
Package dbapi provides access to the Deutsche Bank API
(https://developer.db.com).
The Deutsche Bank API provides developers with plausible costumer and bank data
to let them build great, highly connected apps.

In order to use the Deutsche Bank API you need to create an account at the
developer portal (https://developer.db.com) and follow the instructions there.
Short version: Create a new application, a new test user and authorize your
application to get the access token. Authentication isn't handled by this
package since there are great packages out there.
If you have a valid access token you can start to use this package.

    // Create a new client.
    AccessToken := "1234567890abcdefghijklmnopqrstuvwxyz"
    api, err := dbapi.New(
        dbapi.SetToken(AccessToken),
    )
    if err != nil {
        log.Fatalln(err)
    }

    // Start to access the Deutsche Bank API (retrieve and print user accounts).
    accounts, response, err := api.Accounts.GetAll()
    if err != nil {
        fmt.Println(response)
        log.Fatalln(err)
    }
    fmt.Printf("%v", accounts)

It is also possible to use a custom http client instead of http.DefaultClient:

    // Create your custom http client.
    client := &http.Client{
        Transport: &http.Transport{
            Proxy: http.ProxyFromEnvironment,
            Dial: (&net.Dialer{
                Timeout:   3 * time.Second,
                KeepAlive: 30 * time.Second,
            }).Dial,
            ExpectContinueTimeout: 1 * time.Second,
            ResponseHeaderTimeout: 3 * time.Second,
            TLSHandshakeTimeout:   3 * time.Second,
        },
    }

    // Use your custom http client.
    api, err := dbapi.New(
        dbapi.SetToken(AccessToken),
        dbapi.SetClient(client),
   )
   // ...
*/
package dbapi
