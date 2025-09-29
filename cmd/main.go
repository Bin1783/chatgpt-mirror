package main

import (
	// "bytes"
	"context"
	// "encoding/json"
	"errors"
	"flag"
	"fmt"

	// "io/ioutil"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gpt_mirror/pkg/db"
	"gpt_mirror/pkg/proxy"
	"gpt_mirror/pkg/redis"
	"gpt_mirror/pkg/util"
	"gpt_mirror/routers"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)


func main() {
	InitConfig()
	util.Initialize()
	db.InitMysql()
    _ = db.AutoMigrate()
	redis.InitRedis()
	gin.SetMode(viper.GetString("server.mode"))
	r := gin.Default()
	routers.RegisterRoutes(r)
	go proxy.StartWorkers(5)
	server := &http.Server{
		Addr:           ":" + viper.GetString("server.port"),
		Handler:        r,
		ReadTimeout:    time.Duration(viper.GetInt("server.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(viper.GetInt("server.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			

			os.Exit(1)
		}
	}()
	gracefulExitServer(server)
}

func gracefulExitServer(server *http.Server) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	 <-ch
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		fmt.Println("err", err)
	}
}

func InitConfig() {
	env := flag.String("env", "dev", "Set the application environment (e.g., dev, prod, test)")
	flag.Parse()
	configFileName := "conf_" + *env
	viper.SetConfigName(configFileName)
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./conf")
	viper.AddConfigPath("../conf")
	err := viper.ReadInConfig()
	if err != nil {
		

		panic(err)
	}
}
