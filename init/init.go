package init

import (
	"cetm_booking/common"
	"cetm_booking/x/config"
	"cetm_booking/x/db/mongodb"
	"cetm_booking/x/fcm"
	"cetm_booking/x/mlog"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path"
	"strconv"
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
	mongodb.MongoConfig.MaxPool = context.IntDefault("mongo.maxPool", 0)
	mongodb.MongoConfig.Path, _ = context.String("mongo.path")
	mongodb.MongoConfig.DBname, _ = context.String("mongo.database")
	mongodb.MongoConfig.DBUser, _ = context.String("mongo.db_user")
	mongodb.MongoConfig.DBPass, _ = context.String("mongo.db_pass")
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
	kmStr, _ := context.String("server.search_km")
	timeSet := context.IntDefault("server.time_set_cache", 0)
	mlog.SolutionDir, _ = context.String("server.folder_log")
	km, _ := strconv.ParseFloat(kmStr, 64)

	startDaySendStr, _ := context.String("server.send_notify_start_day")
	startDaySend, _ := strconv.ParseFloat(startDaySendStr, 64)
	bfHourSendStr, _ := context.String("server.send_notify_bf_hour")
	bfHourSend, _ := strconv.ParseFloat(bfHourSendStr, 64)

	startNearStr, _ := context.String("server.start_near")
	startNear, _ := strconv.ParseFloat(startNearStr, 64)

	endNearStr, _ := context.String("server.end_near")
	endNear, _ := strconv.ParseFloat(endNearStr, 64)

	startOutStr, _ := context.String("server.start_out")
	startOut, _ := strconv.ParseFloat(startOutStr, 64)

	endOutStr, _ := context.String("server.end_out")
	endOut, _ := strconv.ParseFloat(endOutStr, 64)

	cyclePushDay := context.IntDefault("server.cycle_push_day", 0)
	cyclePushTicket := context.IntDefault("server.cycle_push_ticket", 0)

	scanNearStr, _ := context.String("server.scan_near")
	scanNear, _ := strconv.ParseFloat(scanNearStr, 64)

	tkDay := context.IntDefault("server.user_ticket_day", 0)
	cycleDayMissed, _ := context.String("server.cycle_day_missed")

	common.ConfigSystemBooking = common.ConfigSystem{
		LinkCetm:           linkCetm,
		LinkSearchMap:      linkSearchMap,
		PortBooking:        port,
		KmSearch:           km,
		TimeSetCache:       timeSet,
		SendNotifyBfHour:   bfHourSend,
		SendNotifyStartDay: startDaySend,
		StartNear:          startNear,
		EndNear:            endNear,
		StartOut:           startOut,
		EndOut:             endOut,
		ScanNear:           scanNear,
		CyclePushDay:       cyclePushDay,
		CyclePushTicket:    cyclePushTicket,
		CycleDayMissed:     cycleDayMissed,
		UserTicketDay:      tkDay,
	}
}
