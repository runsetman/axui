package controller

import (
	"fmt"
	"strconv"
	"x-ui/database/model"
	"x-ui/logger"
	"x-ui/web/global"
	"x-ui/web/service"
	"x-ui/web/session"

	"github.com/gin-gonic/gin"
)

type InboundController struct {
	inboundService service.InboundService
	xrayService    service.XrayService
}

type InboundStats struct {
	Remark   string `json:"string"`
	Up       int64  `json:"up"`
	Down     int64  `json:"down"`
	Expiry   int64  `json:"expiry"`
	Traffic  int64  `json:"traffic"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

func NewInboundController(g *gin.RouterGroup) *InboundController {
	a := &InboundController{}
	a.initRouter(g)
	a.startTask()
	return a
}

func (a *InboundController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/inbound")

	g.POST("/list", a.getInbounds)
	g.POST("/add", a.addInbound)
	g.POST("/del/:id", a.delInbound)
	g.POST("/update/:id", a.updateInbound)
}

func (a *InboundController) startTask() {
	webServer := global.GetWebServer()
	c := webServer.GetCron()
	c.AddFunc("@every 10s", func() {
		if a.xrayService.IsNeedRestartAndSetFalse() {
			err := a.xrayService.RestartXray(false)
			if err != nil {
				logger.Error("restart xray failed:", err)
			}
		}
	})
}

func (a *InboundController) getInbounds(c *gin.Context) {
	user := session.GetLoginUser(c)
	inbounds, err := a.inboundService.GetInbounds(user.Id)
	if err != nil {
		jsonMsg(c, "获取", err)
		return
	}
	jsonObj(c, inbounds, nil)
}

func (a *InboundController) getAllInbounds(c *gin.Context) {
	inbounds, err := a.inboundService.GetAllInbounds()
	if err != nil {
		fmt.Println(err.Error())
		jsonMsg(c, "获取", err)
		return
	}
	jsonObj(c, inbounds, nil)
}

func (a *InboundController) getAllInboundsStats(c *gin.Context) {
	inbounds, err := a.inboundService.GetAllInbounds()
	if err != nil {
		fmt.Println(err.Error())
		jsonMsg(c, "获取", err)
		return
	}

	var inboundsStats []*InboundStats
	for _, inbound := range inbounds {
		stats := &InboundStats{
			Remark:   inbound.Remark,
			Up:       inbound.Up,
			Down:     inbound.Down,
			Traffic:  inbound.Total,
			Expiry:   inbound.ExpiryTime,
			Port:     inbound.Port,
			Protocol: string(inbound.Protocol),
		}
		inboundsStats = append(inboundsStats, stats)
	}

	jsonObj(c, inboundsStats, nil)

}

func (a *InboundController) getInboundByName(c *gin.Context) {
	name := c.Param("name")
	inbound, err := a.inboundService.GetInboundByName(name)
	if err != nil {
		fmt.Println(err.Error())
		jsonMsg(c, "获取", err)
		return
	}

	stats := &InboundStats{
		Remark:   inbound.Remark,
		Up:       inbound.Up,
		Down:     inbound.Down,
		Traffic:  inbound.Total,
		Expiry:   inbound.ExpiryTime,
		Port:     inbound.Port,
		Protocol: string(inbound.Protocol),
	}

	jsonObj(c, stats, nil)
}

func (a *InboundController) addInbound(c *gin.Context) {
	inbound := &model.Inbound{}
	err := c.ShouldBind(inbound)
	if err != nil {
		jsonMsg(c, "添加", err)
		return
	}
	user := session.GetLoginUser(c)
	inbound.UserId = user.Id
	inbound.Enable = true
	inbound.Tag = fmt.Sprintf("inbound-%v", inbound.Port)
	err = a.inboundService.AddInbound(inbound)
	jsonMsg(c, "添加", err)
	if err == nil {
		a.xrayService.SetToNeedRestart()
	}
}

func (a *InboundController) delInbound(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, "删除", err)
		return
	}
	err = a.inboundService.DelInbound(id)
	jsonMsg(c, "删除", err)
	if err == nil {
		a.xrayService.SetToNeedRestart()
	}
}

func (a *InboundController) updateInbound(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, "修改", err)
		return
	}
	inbound := &model.Inbound{
		Id: id,
	}
	err = c.ShouldBind(inbound)
	if err != nil {
		jsonMsg(c, "修改", err)
		return
	}
	err = a.inboundService.UpdateInbound(inbound)
	jsonMsg(c, "修改", err)
	if err == nil {
		a.xrayService.SetToNeedRestart()
	}
}
