package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go/extra"
	"github.com/sdjnlh/communal/app"
	"github.com/sdjnlh/communal/log"
	"github.com/sdjnlh/communal/util"
	"github.com/sdjnlh/communal/web"
	"github.com/sdjnlh/op-user/opu"
	userRest "github.com/sdjnlh/op-user/rest"
)

func main() {
	opu.Api.Mount(opu.Service)
	app.RegisterStarter(opu.Api)
	if err := app.Start(); err != nil {
		panic(err)
	}

	flag.Parse()
	log.Slog.Info("starting rest")
	router := gin.Default()
	router.Use(web.CorsHandler(opu.Api))
	//router.Use(session.InitSessionStore(db.Redis))
	userRest.RegisterAPIs(router)
	extra.SetNamingStrategy(util.LowerFirst)

	log.Slog.Fatal(router.Run(":" + opu.Api.Port))
	log.Slog.Info("rest started")
}
