package web

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/jpeg"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
	"github.com/kbinani/screenshot"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
	"hsott.cn/next-monitor/middlewares"
)

type updateData struct {
	logs chan log
	lock sync.Mutex
}

type log struct {
	Time string `json:"time"`
	Msg  string `json:"msg"`
}

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var data updateData
var img_byte []byte

//go:embed templates\index.html
var index []byte

func init() {
	data.logs = make(chan log, 1024)
}

func WebRun(sendMsg func(string) error, log *logrus.Logger) {
	r := gin.Default()
	r.Use(middlewares.Cors())
	r.GET("/", func(c *gin.Context) {
		if _, err := c.Writer.Write(index); err != nil {
			log.Warn(err)
		}
	})
	r.GET("/video_feed", gen)
	r.GET("/status", status)
	r.GET("/info", info)
	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Warn("Fail to upgrade:", err)
			return
		}
		go wsRec(conn, sendMsg)
		go WsWri(conn)
	})

	go get_screen()

	r.Run()
}

func gen(c *gin.Context) {
	c.Header("Content-Type", "multipart/x-mixed-replace; boundary=frame")

	for {
		c.Writer.Write([]byte("--frame\r\nContent-Type: image/jpeg\r\n\r\n"))
		c.Writer.Write(img_byte)
		c.Writer.Write([]byte("\r\n"))
		time.Sleep(500 * time.Millisecond)
	}
}

func status(c *gin.Context) {
	mem, _ := mem.VirtualMemory()
	// cpu_used, _ := cpu.Percent(time.Duration(time.Second), false)
	c.JSON(200, gin.H{"mem_total": mem.Total, "mem_free": mem.Free, "mem_used": mem.UsedPercent})
	// c.JSON(200, gin.H{"mem_total": mem.Total, "mem_free": mem.Free, "mem_used": mem.UsedPercent, "cpu_used": cpu_used})
}

func info(c *gin.Context) {
	info, err := host.Info()
	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, gin.H{"data": info})
}

func get_screen() {
	for {
		img, err := screenshot.CaptureDisplay(0)
		if err != nil {
			fmt.Println(err)
			return
		}

		frame := new(bytes.Buffer)
		err = jpeg.Encode(frame, img, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		img_byte = frame.Bytes()
		time.Sleep(500 * time.Millisecond)
	}
}

func wsRec(c *ws.Conn, sendMsg func(string) error) {
	for {
		_, data, err := c.ReadMessage()
		if err != nil {
			return
		}
		sendMsg(string(data))
	}
}

func WsWri(c *ws.Conn) {
	for l := range data.logs {
		err := c.WriteJSON(l)
		if err != nil {
			return
		}
	}
}

func AddLog(msg string, t time.Time) {
	defer data.lock.Unlock()
	data.lock.Lock()
	data.logs <- log{
		Time: t.Format("15:04:05"),
		Msg:  msg,
	}
}
