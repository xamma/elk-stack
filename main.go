package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/xamma/elk-stack/config"
)

var sugar *zap.SugaredLogger

func hello(ctx *gin.Context) {
ctx.JSON(http.StatusOK, gin.H{"message": "Hello world"})
sugar.Infow("Hello world endpoint accessed",
		"remote_address", ctx.Request.RemoteAddr,
		"user_agent", ctx.Request.UserAgent(),
	)
}

func secret(ctx *gin.Context) {
	user, _ := ctx.Get(gin.AuthUserKey)
	ctx.JSON(http.StatusOK, gin.H{"message": "Damn this is so secret"})
	sugar.Infow("Secret endpoint accessed",
		"remote_address", ctx.Request.RemoteAddr,
		"user_agent", ctx.Request.UserAgent(),
		"user", user,
	)
}

func setupLogger() (*zap.SugaredLogger, func(), error) {
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleDebugging := zapcore.Lock(os.Stdout)

	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		return nil, nil, err
	}

	file, err := os.Create("logs/app.log")
	if err != nil {
		return nil, nil, err
	}

	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	fileDebugging := zapcore.Lock(file)

	consoleCore := zapcore.NewCore(consoleEncoder, consoleDebugging, zap.DebugLevel)
	fileCore := zapcore.NewCore(fileEncoder, fileDebugging, zap.InfoLevel)

	logger := zap.New(zapcore.NewTee(consoleCore, fileCore), zap.AddCaller())

	sugar = logger.Sugar()
	sugar.Info("Logger initialized")

	// Needed to keep the file open
	closeFile := func() {
		err := file.Close()
		if err != nil {
			sugar.Error("Failed to close the log file:", err)
		}
	}

	return sugar, closeFile, nil
}

func main() {

	sugar, closeFile, err := setupLogger()
	if err != nil {
		panic("Something bad happened to logging")
	}
	defer closeFile()

	cfg, err := config.LoadConfig()
	if err != nil {
		sugar.Errorf("Error loading Config: %v", err)
		panic("Failed to load configuration")
	}

	user := cfg.User
	password := cfg.Password

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	
	v1 := router.Group("/api/v1")
	v1.GET("/hello", hello)
	
	authorized := router.Group("/api/v1", gin.BasicAuth(gin.Accounts{
		user: password,
		// "admin": "0000",
	}))

	authorized.GET("/secret", secret)
	router.Run(":9090")
}
