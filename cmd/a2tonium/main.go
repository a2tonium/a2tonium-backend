package main

import (
	"context"
	"fmt"
	internalTon "github.com/a2tonium/a2tonium-backend/internal/app/ton"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"log"
	"time"
)

func main() {
	var (
		ctx, _ = context.WithCancel(context.Background())
	)
	client := liteclient.NewConnectionPool()

	configUrl := "https://ton-blockchain.github.io/testnet-global.config.json"
	err := client.AddConnectionsFromConfigUrl(context.Background(), configUrl)
	if err != nil {
		panic(err)
	}
	api := ton.NewAPIClient(client)

	mnemonic := "artwork vintage physical silk combine faith sketch crisp lion wrestle call credit shell chase donor glare sudden resource edge behave diamond sweet lens fall"
	log.Println("Creating TonService ...")
	tonService := internalTon.NewTonService(api, mnemonic)
	for {
		err = tonService.Init(ctx)
		if err != nil {
			panic(err)
		}
		//tonService.Show()
		fmt.Println("Process Poshel :) ...")

		for i := 0; i < 60; i++ {
			fmt.Println(i)
			time.Sleep(5 * time.Second)
			err = tonService.Run(ctx)
			if err != nil {
				panic(err)
			}
		}
		fmt.Println("Reloading...")
	}
}
