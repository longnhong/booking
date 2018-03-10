package init

import (
	"cetm_booking/common"
	"cetm_booking/x/config"
	"cetm_booking/x/db/mongodb"
	"cetm_booking/x/fcm"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path"
)

func init() {
	loadConfig()
	initLog()
	initDB()
	initFcm()
	initConfigSytem()
}

var context *config.Context

func loadConfig() {
	context, _ = config.LoadContext("app.conf", []string{""})
}

func initLog() {
	//config for gin request log
	{
		f, _ := os.Create(path.Join("log", "gin.log"))
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	}
	//config for app log use glog
	{
		logDir, _ := context.String("glog.log_dir")
		logStd, _ := context.String("glog.alsologtostderr")
		flag.Set("alsologtostderr", logStd)
		flag.Set("log_dir", logDir)
		flag.Parse()
	}
}
func initDB() {
	fmt.Println("init db")
	// Read configuration.
	mongodb.MaxPool = context.IntDefault("mongo.maxPool", 0)
	mongodb.PATH, _ = context.String("mongo.path")
	mongodb.DBNAME, _ = context.String("mongo.database")
	mongodb.CheckAndInitServiceConnection()
}

func initFcm() {
	fcm.FCM_SERVER_KEY_CUSTOMER, _ = context.String("fcm.serverkey.customer")
	fcm.NewFcmApp(fcm.FCM_SERVER_KEY_CUSTOMER)
	fmt.Print("Qua fcm")
}

func initConfigSytem() {
	linkCetm, _ := context.String("server.cetm")
	linkSearchMap, _ := context.String("server.map_search")
	port, _ := context.String("server.port")
	common.ConfigSystemBooking = common.ConfigSystem{
		LinkCetm:      linkCetm,
		LinkSearchMap: linkSearchMap,
		PortBooking:   port,
	}
}
