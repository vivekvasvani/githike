package server

import (
	"log"

	"database/sql"

	"github.com/buaazp/fasthttprouter"
	"github.com/golang/glog"
	"github.com/valyala/fasthttp"
)

func NewServer(db *sql.DB) {

	router := fasthttprouter.New()

	router.POST("/gitoptions", func(ctx *fasthttp.RequestCtx) {
		SendGitHikeOptions(ctx)
	})

	router.POST("/getrequests", func(ctx *fasthttp.RequestCtx) {
		HandleAppRequests(ctx)
	})

	router.POST("/addhiketeammembership", func(ctx *fasthttp.RequestCtx) {
		AddHikeTeamMembership(ctx, db)
	})

	router.POST("/addorupdateteam", func(ctx *fasthttp.RequestCtx) {
		AddOrUpdateTeam(ctx, db)
	})

	router.PanicHandler = func(ctx *fasthttp.RequestCtx, p interface{}) {
		glog.V(0).Infof("Panic occurred %s", p, ctx.Request.URI().String())
	}
	log.Println("Service Started on port " + "6001")
	glog.Fatal(fasthttp.ListenAndServe(":"+"6001", fasthttp.CompressHandler(router.Handler)))

}
