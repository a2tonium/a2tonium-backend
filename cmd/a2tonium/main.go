package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/a2tonium/a2tonium-backend/internal/app/a2tonium"
	"github.com/a2tonium/a2tonium-backend/internal/app/ipfs"
	jsonGenerator "github.com/a2tonium/a2tonium-backend/internal/app/json_generator"
	"github.com/a2tonium/a2tonium-backend/internal/app/ton"
	"github.com/a2tonium/a2tonium-backend/pkg/config"
	"github.com/a2tonium/a2tonium-backend/pkg/logger"
	"github.com/a2tonium/a2tonium-backend/pkg/ton/crypto"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
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
	var (
		ctx, _            = context.WithCancel(context.Background())
		generatePublicKey = flag.Bool("generatePublicKey", false, "Set this flag to generate public key")
		configFlag        = flag.String("config", "test", "configuration file suffix")
	)
	flag.Parse()
	config.LoadConfig(*configFlag)
	setupLogger()

	if *generatePublicKey {
		keypair, err := crypto.MnemonicToX25519KeyPair(config.GetValue(mnemonicPhrase).String())
		if err != nil {
			fmt.Println("public key generation failed:", err)
			return
		}
		fmt.Println("Your public key:", base64.StdEncoding.EncodeToString(keypair.PublicKey))

		return
	}

	jsonGeneratorService := jsonGenerator.NewJsonGeneratorService()
	logger.Info(ctx, logger.Msg, "jsonGeneratorService created")

	ipfsService, err := ipfs.NewIpfsService(config.GetValue(pinataJwtToken).String())
	if err != nil {
		logger.ErrorKV(ctx, "ipfs.NewIpfsService error", logger.Err, err)
		return
	}
	logger.Info(ctx, logger.Msg, "ipfsService created")

	tonService := ton.NewTonService()
	err = tonService.Init(ctx, config.GetValue(mnemonicPhrase).String())
	if err != nil {
		logger.ErrorKV(ctx, "tonService.Init error", logger.Err, err)
		return
	}
	logger.Info(ctx, logger.Msg, "tonService created and inited")

	a2toniumService := a2tonium.NewA2Tonium(tonService, ipfsService, jsonGeneratorService)
	if err = a2toniumService.Init(ctx); err != nil {
		logger.ErrorKV(ctx, logger.Err, err)
		return
	}
	logger.Info(ctx, logger.Msg, "a2toniumService created and inited")

	logger.Info(ctx, logger.Msg, "Running A2Tonium ...")
	if err = a2toniumService.Run(ctx); err != nil {
		logger.ErrorKV(ctx, logger.Err, err)
		return
	}
}
