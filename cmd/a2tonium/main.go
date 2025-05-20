package main

import (
	"context"
	"fmt"
	internalTon "github.com/a2tonium/a2tonium-backend/internal/app/ton"
	"github.com/a2tonium/a2tonium-backend/pkg/config"
	"github.com/a2tonium/a2tonium-backend/pkg/logger"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
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
		ctx, _ = context.WithCancel(context.Background())
	)
	client := liteclient.NewConnectionPool()

	configUrl := "https://ton-blockchain.github.io/testnet-global.config.json"
	err := client.AddConnectionsFromConfigUrl(context.Background(), configUrl)
	if err != nil {
		panic(err)
	}
	api := ton.NewAPIClient(client)

	mnemonic := "erupt unfair flee rent inquiry nerve enlist swamp report lucky witness donkey race task evoke shed cave mercy puzzle come slight limb fun prefer"
	log.Println("Creating TonService ...")
	tonService := internalTon.NewTonService(api, mnemonic)
	err = tonService.Init(ctx)
	if err != nil {
		panic(err)
	}
	//tonService.Show()
	fmt.Println("Process Poshel :) ...")
	err = tonService.Run(ctx)
	if err != nil {
		panic(err)
	}
}
