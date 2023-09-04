package ignite

import (
	"context"
	"encoding/base64"
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
	PayloadType string `json:"payload_type"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Data        string `json:"data"`
	Result      string `json:"result"`
}

func Event() {
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
			log.Fatal("AccountRegistry", err)
		}
	}
	address, _ := account.Address("cosmos")
	log.Println("address: ", address)
	res, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody([]byte(fmt.Sprintf(`{"denom":"stake","address":"%s"}`, address))).
		Post(confer.Cfg.FaucetAddress)
	fmt.Println("Faucet:", res.String(), err)
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
					if eventPB.DeviceName == confer.Cfg.AccountName {
						var payload Payload
						// Payload需要base64解密
						payloadBase64, err := base64.StdEncoding.DecodeString(eventPB.Payload)
						if err != nil {
							log.Println("base64.StdEncoding.DecodeString(eventPB.Payload) err:", err)
							break
						}
						log.Println("accept payload:", string(payloadBase64))
						err = json.Unmarshal(payloadBase64, &payload)
						if err != nil {
							log.Println("json.Unmarshal(payloadBase64, &payload):", err)
							break
						}
						switch strings.ToLower(payload.PayloadType) {
						case "publish":
							res, err := publish(payload.Name, payload.Path, payload.Data)
							if err != nil {
								log.Println("publish(payload.Name, payload.Path, payload.Data):", err)
								break
							}
							log.Println("publish:==", payload.Name, payload.Path, res, err)
							//payload.Result = base64.StdEncoding.EncodeToString([]byte(res))
							payload.Result = res
							payloadBytes, _ := json.Marshal(payload)
							response, err := client.BroadcastTx(context.Background(), account, &ttypes.MsgUpdateEventPb{
								Creator:    address,
								Index:      eventPB.Index,
								DeviceName: eventPB.DeviceName,
								Payload:    base64.StdEncoding.EncodeToString(payloadBytes),
							})
							fmt.Println("BroadcastTx publish:==", response.Data, err)
						case "subscribe":
							res, err := subscribe(payload.Name, payload.Path)
							if err != nil {
								log.Println("subscribe(payload.Name, payload.Path):", err)
								break
							}
							log.Println("subscribe:==", payload.Name, payload.Path, res, err)
							payload.Result = res
							payloadBytes, _ := json.Marshal(payload)
							eventpb := &ttypes.MsgUpdateEventPb{
								Creator:    address,
								Index:      eventPB.Index,
								DeviceName: eventPB.DeviceName,
								Payload:    base64.StdEncoding.EncodeToString(payloadBytes),
							}
							response, err := client.BroadcastTx(context.Background(), account, eventpb)
							if err != nil {
								fmt.Println("client.BroadcastTx err ", err)
							}
							fmt.Println("BroadcastTx subscribe:==", response.Data, err)
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
