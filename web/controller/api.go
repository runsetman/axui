package controller

import (
	"github.com/gin-gonic/gin"
)

type ApiController struct {
	BaseController
	inboundController *InboundController
}

func NewApiController(g *gin.RouterGroup) *ApiController {
	a := &ApiController{}
	a.inboundController = NewInboundController(g)
	a.initRouter(g)
	return a
}

func (a *ApiController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/api")
	g.Use(a.validate)
	g.POST("/inbounds", a.getInbounds)
	g.POST("/inbounds/:name", a.getInboundByName)
}

func (a *ApiController) getInbounds(c *gin.Context) {
	a.inboundController.getAllInboundsStats(c)
}

func (a *ApiController) getInboundByName(c *gin.Context) {
	a.inboundController.getInboundByName(c)
}
