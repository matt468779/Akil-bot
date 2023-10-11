package main

import (
	"akil_telegram_bot/bootstrap"
	"akil_telegram_bot/gpt"
	"time"

	route "akil_telegram_bot/route"

	"github.com/gin-gonic/gin"
)

func main() {

	app := bootstrap.App()

	env := app.Env

	db := app.Mongo.Database(env.DBName)

	gpt.SetEnv(env)

	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	gin := gin.Default()

	route.Setup(env, timeout, db, gin, app.Bot)

	gin.Run(env.ServerAddress)
}
