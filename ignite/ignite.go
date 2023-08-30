package ignite

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-resty/resty/v2"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	ttypes "github.com/stc-community/iot-depin-protocol/x/iotdepinprotocol/types"
	"github.com/stc-community/iot-light-client/pkg/confer"
	"github.com/tendermint/tendermint/types"
)

type Payload struct {
	PayloadType string
	Name        string
	Path        string
	Data        string
}

func Event() {
	client, err := cosmosclient.New(context.Background(),
		cosmosclient.WithNodeAddress(confer.Cfg.NodeAddress),
		cosmosclient.WithUseFaucet(confer.Cfg.FaucetAddress, "", 1000),
	)
	account, _, err := client.AccountRegistry.Create(confer.Cfg.AccountName)
	if err != nil {
		log.Fatal(err)
	}
	_ = client.RPC.Start()
	eventCh, err := client.RPC.Subscribe(context.Background(), "", types.QueryForEvent(types.EventTx).String())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.RPC.UnsubscribeAll(context.Background(), ""); err != nil {
			log.Fatal(err)
		}
		_ = client.RPC.Stop()
	}()
	for {
		event := <-eventCh
		txEvent, ok := event.Data.(types.EventDataTx)
		if ok {
			events := txEvent.Result.GetEvents()
			for _, val := range events {
				if strings.Compare(val.Type, "stccommunity.iotdepinprotocol.iotdepinprotocol.EventPb") == 0 {
					eventPb, _ := sdk.ParseTypedEvent(val)
					eventPB := eventPb.(*ttypes.EventPb)
					if eventPB.Topic == confer.Cfg.AccountName && eventPB.PubType == "request" {
						var payload Payload
						json.Unmarshal([]byte(eventPB.Payload), &payload)
						switch strings.ToLower(payload.PayloadType) {
						case "publish":
							res, _ := publish(payload.Name, payload.Path, payload.Data)
							client.BroadcastTx(context.Background(), account, &ttypes.MsgCreateKv{
								Creator: account.Name,
								Index:   "result",
								Value:   res,
							})
						case "subscribe":
							res, _ := subscribe(payload.Name, payload.Path)
							client.BroadcastTx(context.Background(), account, &ttypes.MsgCreateKv{
								Creator: account.Name,
								Index:   "result",
								Value:   res,
							})
						}
					}
				}
			}
		}
	}
}

func subscribe(name, path string) (result string, err error) {
	res, err := resty.New().R().Get(fmt.Sprintf("http://%s.deviceshifu/%s", name, path))
	return res.String(), err
}

func publish(name, path string, data string) (result string, err error) {
	res, err := resty.New().R().
		//SetHeader("Content-Type", c.Request.Header.Get("Content-Type")).
		SetBody(data).
		Post(fmt.Sprintf("http://%s.deviceshifu/%s", name, path))
	return res.String(), err
}
