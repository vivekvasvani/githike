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

	router.POST("/deletemember", func(ctx *fasthttp.RequestCtx) {
		DeleteMember(ctx, db)
	})

	router.POST("/createrepo", func(ctx *fasthttp.RequestCtx) {
		CreateRepository(ctx, db)
	})

	router.POST("/userdetails", func(ctx *fasthttp.RequestCtx) {
		UserDetails(ctx, db)
	})

	router.POST("/repodetails", func(ctx *fasthttp.RequestCtx) {
		//RepoDetails(ctx, db)
	})

	router.POST("/teamdetails", func(ctx *fasthttp.RequestCtx) {
		//TeamDetails(ctx, db)
	})

	router.POST("/inviteusertohike", func(ctx *fasthttp.RequestCtx) {
		//InviteUserToHike(ctx, db)
	})

	router.PanicHandler = func(ctx *fasthttp.RequestCtx, p interface{}) {
		glog.V(0).Infof("Panic occurred %s", p, ctx.Request.URI().String())
	}
	log.Println("Service Started on port " + "6001")
	glog.Fatal(fasthttp.ListenAndServe(":"+"6001", fasthttp.CompressHandler(router.Handler)))

}
