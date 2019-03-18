package logs_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/junhwong/go-logs"
)

func TestX(t *testing.T) {
	fmt.Println("x")
}

func TestConsole(t *testing.T) {
	logs.Print("HELLO WORLD!")
	logs.Debug("test log2")
	logs.Info("test log3")
	logs.Warn("test log3")
	logs.Prefix("redis.box").Error("test log4")
	logs.Fatal("test log3")
	logs.Std().GracefulStop()
	// time.Sleep(1 * time.Second)
}
func TestFile(t *testing.T) {
	h := logs.FileHandler{
		HandlerBase: logs.HandlerBase{
			FormatType: "text",
			TimeFormat: "%Y-%m-%d %H:%M:%S %z",
		},
		Path:           "/Volumes/sea/docker/fluentd/log",
		PathTimeFormat: "%Y%m%d",
		RotateSize:     256,
	}
	logs.Std().Append(&h)
	logs.Error("test log2")
	logs.Error("test log3")
	logs.Prefix("redis.box").Error("test log4")
	logs.Std().GracefulStop()
}

func TestUDP(t *testing.T) {
	h := logs.NetHandler{
		HandlerBase: logs.HandlerBase{
			FormatType: "json",
			TimeFormat: "timestamp",
		},
		Hostname: "192.168.3.6:9527",
		Protocol: "ipv4.udp",
	}
	logs.Std().Append(&h)
	logs.Error("test log2")
	logs.Error("test log3")
	logs.Prefix("redis.box").Error("test log4")
	time.Sleep(1000 * time.Second)
}
