package ignite

import (
	"context"
	"fmt"
	"log"

	"github.com/ignite/cli/ignite/pkg/cosmosaccount"

	"github.com/go-resty/resty/v2"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/stc-community/iot-light-client/pkg/confer"
)

type Payload struct {
	PayloadType string `json:"payload_type"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Data        string `json:"data"`
	Result      string `json:"result"`
}

type Ignite struct {
	Client  cosmosclient.Client
	Account cosmosaccount.Account
	Address string
}

var IgniteC *Ignite

func init() {
	client, err := cosmosclient.New(context.Background(),
		cosmosclient.WithNodeAddress(confer.Cfg.NodeAddress),
	)
	if err != nil {
		log.Fatal(" cosmosclient.New(context.Background() err", err)
	}
	account, _ := client.Account(confer.Cfg.AccountName)
	if account.Record == nil {
		account, _, err = client.AccountRegistry.Create(confer.Cfg.AccountName)
		if err != nil {
			log.Fatal("AccountRegistry err", err)
		}
	}
	address, _ := account.Address("cosmos")
	log.Println("address: ", address)
	res, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody([]byte(fmt.Sprintf(`{"denom":"stake","address":"%s"}`, address))).
		Post(confer.Cfg.FaucetAddress)
	fmt.Println("Faucet:", res.String(), err)
	IgniteC = &Ignite{
		Client:  client,
		Account: account,
		Address: address,
	}
}
