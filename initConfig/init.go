package initConfig

import (
	"dst-admin-go/api"
	"dst-admin-go/config"
	"dst-admin-go/config/database"
	"dst-admin-go/config/global"
	"dst-admin-go/lobbyServer"
	"dst-admin-go/model"
	"dst-admin-go/schedule"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const logPath = "./dst-admin-go.log"

var f *os.File

func Init() {

	initConfig()
	initLog()
	initDB()
	initCollect()
	initSchedule()
	initLobbyServer()
}

func initDB() {
	db, err := gorm.Open(sqlite.Open(global.Config.Db), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic("failed to connect database")
	}
	database.DB = db
	err = database.DB.AutoMigrate(
		&model.Spawn{},
		&model.PlayerLog{},
		&model.Connect{},
		&model.Proxy{},
		&model.ModInfo{},
		&model.Cluster{},
		&model.JobTask{},
		&api.SystemConfig{},
		&lobbyServer.LobbyHome{},
	)
	if err != nil {
		return
	}
}

func initConfig() {
	yamlFile, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		fmt.Println(err.Error())
	}
	var _config *config.Config
	err = yaml.Unmarshal(yamlFile, &_config)
	if err != nil {
		fmt.Println(err.Error())
	}
	global.Config = _config
}

func initLog() {
	var err error
	f, err = os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}

	// 组合一下即可，os.Stdout代表标准输出流
	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)

	gin.ForceConsoleColor()
	gin.SetMode(gin.DebugMode)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func initCollect() {

	var clusters []model.Cluster
	database.DB.Find(&clusters)
	for _, cluster := range clusters {
		global.CollectMap.AddNewCollect(cluster.ClusterName)
	}

}

func initSchedule() {
	global.Schedule = schedule.NewSchedule()
}

func initLobbyServer() {

	if !global.Config.EnableLobby {
		log.Println("不开启 lobby server")
		return
	}

	// 定义一个定时器，每隔 5 秒执行一次
	lobbyServer2 := lobbyServer.NewLobbyServer2(database.DB)

	global.LobbyServer = lobbyServer2

	ticker := time.NewTicker(60 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				lobbyServer2.SaveLobbyList()
			}
		}
	}()

}