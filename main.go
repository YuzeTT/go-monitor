package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"hsott.cn/next-monitor/hook"
	"hsott.cn/next-monitor/web"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

var log *logrus.Logger

func get_info() *host.InfoStat {
	info, err := host.Info()
	if err != nil {
		fmt.Println(err)
		// return
	}

	return info
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"info": gin.H{
			"hostname": get_info().Hostname,
			"platform": get_info().Platform,
		},
	})
}

func info(c *gin.Context) {
	info, err := host.Info()
	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, gin.H{"data": info})
}

func status(c *gin.Context) {
	mem, _ := mem.VirtualMemory()
	cpu_used, _ := cpu.Percent(time.Duration(time.Second), false)
	c.JSON(200, gin.H{"mem_total": mem.Total, "mem_free": mem.Free, "mem_used": mem.UsedPercent, "cpu_used": cpu_used})
}

func main() {
	log = logrus.New()
	// log.SetOutput(colorable.NewColorableStdout())
	hook.InitHook(log)
	web.WebRun(sendMsg, log)
}

func sendMsg(str string) error {
	// if str == "/throw" {
	// if err := useItem(); err != nil {
	// 	return err
	// }
	// }
	// if err := c.Conn.WritePacket(pk.Marshal(packetid.ChatServerbound, pk.String(str))); err != nil {
	// 	return err
	// }
	// log.Info(str)
	return nil
}
