package confer

import (
	"github.com/caarlos0/env/v9"
)

var Cfg = Config{}

type Config struct {
	AccountName   string `env:"ACCOUNT_NAME,required" envDefault:"thermometer"`
	NodeAddress   string `env:"NODE_ADDRESS,required" envDefault:"https://cloudx3-iot-rpc.gw103.oneitfarm.com:443"`
	FaucetAddress string `env:"FAUCET_ADDRESS,required" envDefault:"https://cloudx3-iot-faucet.gw103.oneitfarm.com/credit"`
}

func init() {
	if err := env.Parse(&Cfg); err != nil {
		panic(err)
	}
}
