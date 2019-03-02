package logs

import (
	"errors"
	"net"
)

type NetHandler struct {
	HandlerBase
	Hostname string // 主机名，包含端口
	Protocol string // 协议，支持类型：ipv4.udp(默认),ipv4.tcp,ipv6.udp,ipv6.udp,http
	conn     net.Conn
}

func (h *NetHandler) Init() (err error) {
	if h.Protocol == "" {
		h.Protocol = "ipv4.udp"
	}
	if h.FormatType == "" {
		h.FormatType = "json"
	}
	h.TimeFormat = replaceTimeFormat(h.TimeFormat)

	switch h.Protocol {
	case "ipv4.udp":
		fallthrough
	case "ipv6.udp":
		h.conn, err = net.Dial("udp", h.Hostname)
	default:
		err = errors.New("logs.NetHandler: not support protocol")
	}

	if err != nil {
		return
	}

	return nil
}

func (h *NetHandler) Handle(v Entry) (err error) {
	switch h.FormatType {
	case "json":
		var data []byte
		data, err = formatJson(v, h.TimeFormat)
		if err != nil {
			return
		}
		_, err = h.conn.Write(data) // TODO: 循环发送完整
	}
	return
}
