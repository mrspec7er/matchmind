package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/mrspec7er/matchmind/server/internal/handler"
	"github.com/mrspec7er/matchmind/server/internal/service"
)

type Controller struct {
	Service  service.Service
	Response handler.ResponseJSON
	Upgrader websocket.Upgrader
}

func (c *Controller) SendResponse(w http.ResponseWriter, r *http.Request) {
	roomId := chi.URLParam(r, "roomId")

	wsConn := c.Service.WebsocketConnection()
	newConn, err := wsConn.Upgrade(w, r, nil)
	if err != nil {
		c.Response.GeneralErrorHandler(w, 500, err)
		return
	}

	c.Service.ProcessMessage(newConn, roomId)
}

func (c *Controller) CreateRoom(w http.ResponseWriter, r *http.Request) {

	room, err := c.Service.CreateRoom()
	if err != nil {
		c.Response.GeneralErrorHandler(w, 500, err)
		return
	}

	c.Response.QuerySuccessResponse(w, nil, room, nil)
}
