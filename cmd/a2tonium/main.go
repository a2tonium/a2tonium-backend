package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	internalTon "github.com/a2tonium/a2tonium-backend/internal/app/ton"
	"github.com/a2tonium/a2tonium-backend/pkg/config"

	//"github.com/a2tonium/a2tonium-backend/pkg/config"
	"github.com/a2tonium/a2tonium-backend/pkg/ton/crypto"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"log"
	"time"
)

func main() {
	var (
		ctx, _   = context.WithCancel(context.Background())
		mnemonic = "artwork vintage physical silk combine faith sketch crisp lion wrestle call credit shell chase donor glare sudden resource edge behave diamond sweet lens fall"
	)
	generatePublicKey := flag.Bool("generatePublicKey", false, "Set this flag to generate public key")
	configFlag := flag.String("config", "test", "configuration file suffix")
	flag.Parse()
	config.LoadConfig(*configFlag)

	if *generatePublicKey {
		keypair, err := crypto.MnemonicToX25519KeyPair(mnemonic)
		if err != nil {
			fmt.Println("public key generation failed:", err)
			return
		}
		fmt.Println("Your public key:", base64.StdEncoding.EncodeToString(keypair.PublicKey))

		return
	}

	client := liteclient.NewConnectionPool()

	configUrl := "https://ton-blockchain.github.io/testnet-global.config.json"
	err := client.AddConnectionsFromConfigUrl(context.Background(), configUrl)
	if err != nil {
		panic(err)
	}
	api := ton.NewAPIClient(client)

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
