package confer

import (
	"github.com/caarlos0/env/v9"
)

var Cfg = Config{}

type Config struct {
	AccountName   string `env:"ACCOUNT_NAME,required"`
	NodeAddress   string `env:"NODE_ADDRESS,required"`
	FaucetAddress string `env:"FAUCET_ADDRESS,required"`
}

func init() {
	if err := env.Parse(&Cfg); err != nil {
		panic(err)
	}
}
