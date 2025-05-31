package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	internalTon "github.com/a2tonium/a2tonium-backend/internal/app/ton"
	"github.com/a2tonium/a2tonium-backend/pkg/config"
	"github.com/a2tonium/a2tonium-backend/pkg/logger"
	"github.com/a2tonium/a2tonium-backend/pkg/ton/crypto"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"time"
)

func getLogLevel(lvl string) zapcore.Level {
	res := zapcore.PanicLevel
	switch lvl {
	case "debug":
		res = zapcore.DebugLevel
	case "info":
		res = zapcore.InfoLevel
	case "warn":
		res = zapcore.WarnLevel
	case "error":
		res = zapcore.ErrorLevel
	}
	return res
}

func setupLogger() {
	if config.GetValue(logIntoFile).Boolean() {
		logToFile := &lumberjack.Logger{
			Filename:   config.GetValue(logFilePath).String(),
			MaxAge:     config.GetValue(logFileMaxAge).Int(),
			MaxSize:    config.GetValue(logFileMaxSize).Int(),
			MaxBackups: config.GetValue(logFileBackups).Int(),
		}
		logger.SetLogger(logger.NewWithOutput(getLogLevel(config.GetValue(logLevel).String()), logToFile))
	} else {
		logger.SetLogger(logger.New(getLogLevel(config.GetValue(logLevel).String())))
	}
}

func main() {
	setupLogger()
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
