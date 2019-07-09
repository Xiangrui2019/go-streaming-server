package routers

import (
	"go-streaming-server/controllers"
	"go-streaming-server/limiters"
	"go-streaming-server/response"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Router struct {
	context *httprouter.Router
	limiter *limiters.ConnectionLimiter
	video   *controllers.VideoController
}

func NewRouter(ctx *httprouter.Router, maxconnection int) *Router {
	return &Router{
		context: ctx,
		limiter: limiters.NewConnectionLimiter(maxconnection),
		video:   controllers.NewVideoController(),
	}
}

func (router *Router) ConfigureRouter() error {
	router.context.GET("/video/:vid", router.video.Streaming)
	router.context.POST("/video/:vid", router.video.Upload)

	return nil
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !router.limiter.GetConnection() {
		response.SendResponse(w, http.StatusTooManyRequests, &response.ErrorResponse{
			Code:    http.StatusTooManyRequests,
			Message: "Too many requests.",
		})
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

	router.context.ServeHTTP(w, r)

	defer router.limiter.FreeConnection()
}
