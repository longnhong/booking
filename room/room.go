package room

import (
	"cetm_booking/x/rest"
	"cetm_booking/x/socket"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type RoomServer struct {
	*gin.RouterGroup
	rest.JsonRender
	box *socket.Box
}

func NewRoomServer(parent *gin.RouterGroup) *RoomServer {
	var s = RoomServer{
		RouterGroup: parent,
		box:         socket.NewBox(),
	}
	s.box.Handle("/chat", s.handleChat)
	s.GET("/join", s.handleJoin)
	return &s
}
func (s *RoomServer) handleJoin(ctx *gin.Context) {
	s.serveWS(ctx.Writer, ctx.Request)
}
func (s *RoomServer) handleChat(r *socket.Request) {
	var msg = struct {
		To   string
		Text string
	}{}
	r.MustDecodeBody(&msg)
	glog.Info(msg)
	s.box.SendToOther(msg.To, []byte(msg.Text))
}
func (s *RoomServer) serveWS(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		glog.Infof("error upgrading %s", err)
		return
	}
	s.box.Accept(id, conn)
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
