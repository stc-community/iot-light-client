package script

import (
	"context"
	"encoding/json"
	"log"

	"github.com/urfave/cli"

	ttypes "github.com/stc-community/iot-depin-protocol/x/iotdepinprotocol/types"
	"github.com/stc-community/iot-light-client/pkg/ignite"
	"github.com/stc-community/iot-light-client/pkg/util"
)

type Device struct {
	Name         string `json:"name"`
	Protocol     string `json:"protocol"`
	Address      string `json:"address"`
	DriverSku    string `json:"driver_sku"`
	DriverImage  string `json:"driver_image"`
	Instructions string `json:"instructions"`
	Telemetries  string `json:"telemetries"`
}

func GenerateThermometer(c *cli.Context) {
	device := Device{
		Name:         "thermometer",
		Protocol:     "HTTP",
		Address:      "thermometer.devices:11111",
		DriverSku:    "thermometer",
		DriverImage:  "edgenesis/thermometer:v0.0.1",
		Instructions: "aW5zdHJ1Y3Rpb25TZXR0aW5nczoKICBkZWZhdWx0VGltZW91dFNlY29uZHM6IDgKaW5zdHJ1Y3Rpb25zOgogIGdldF9zdGF0dXM6CiAgcmVhZF92YWx1ZTo=",
		Telemetries:  "dGVsZW1ldHJpZXM6CiAgZGV2aWNlX2hlYWx0aDoKICAgIHByb3BlcnRpZXM6CiAgICAgIGluc3RydWN0aW9uOiBnZXRfc3RhdHVzCiAgICAgIGluaXRpYWxEZWxheU1zOiAxMDAwCiAgICAgIGludGVydmFsTXM6IDEwMDA=",
	}
	mid := util.RandomString(16)
	deviceBytes, _ := json.Marshal(device)
	response, err := ignite.IgniteC.Client.BroadcastTx(context.Background(), ignite.IgniteC.Account, &ttypes.MsgCreateDeviceRegistry{
		Mid:      mid,
		MetaData: string(deviceBytes),
		Creator:  ignite.IgniteC.Address,
	})
	if err != nil {
		log.Println("BroadcastTx err ", err)
		return
	}
	log.Println("register a thermometer success ", response.Code)
	log.Printf("this thermometer mid is : %s", mid)
}

func GenerateAgv(c *cli.Context) {
	device := Device{
		Name:         "agv",
		Protocol:     "HTTP",
		Address:      "agv.devices:11111",
		DriverSku:    "agv",
		DriverImage:  "edgenesis/agv:v0.0.1",
		Instructions: "aW5zdHJ1Y3Rpb25TZXR0aW5nczoKICBkZWZhdWx0VGltZW91dFNlY29uZHM6IDgKaW5zdHJ1Y3Rpb25zOgogIGdldF9wb3NpdGlvbjoKICBnZXRfc3RhdHVzOg==",
		Telemetries:  "dGVsZW1ldHJpZXM6CiAgZGV2aWNlX2hlYWx0aDoKICAgIHByb3BlcnRpZXM6CiAgICAgIGluc3RydWN0aW9uOiBnZXRfc3RhdHVzCiAgICAgIGluaXRpYWxEZWxheU1zOiAxMDAwCiAgICAgIGludGVydmFsTXM6IDEwMDA=",
	}
	mid := util.RandomString(16)
	deviceBytes, _ := json.Marshal(device)
	response, err := ignite.IgniteC.Client.BroadcastTx(context.Background(), ignite.IgniteC.Account, &ttypes.MsgCreateDeviceRegistry{
		Mid:      mid,
		MetaData: string(deviceBytes),
		Creator:  ignite.IgniteC.Address,
	})
	if err != nil {
		log.Println("BroadcastTx err ", err)
		return
	}
	log.Println("register a agv success ", response.Code)
	log.Printf("this agv mid is : %s", mid)
}
