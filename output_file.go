package logs

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
)

type FileHandler struct {
	HandlerBase
	Path            string // 日志存放的路径
	PathTimeFormat  string
	RotateSize      int64 // 单个文件的大小
	AppendDelimiter string
	stream          chan Entry
	file            *os.File
}

func (h *FileHandler) Init() error {
	dir, err := testDir(h.Path)
	if err != nil {
		return err
	}
	if h.FormatType == "" {
		h.FormatType = "json"
	}
	if h.AppendDelimiter == "" {
		h.AppendDelimiter = "\n"
	}
	h.TimeFormat = replaceTimeFormat(h.TimeFormat)
	h.Path = dir
	h.stream = make(chan Entry) //, 1000
	h.PathTimeFormat = replaceTimeFormat(h.PathTimeFormat)

	return nil
}

func (h *FileHandler) Handle(v Entry) (err error) {
	file, err := openFile(h, v)
	if err != nil {
		return
	}

	if h.FormatType == "text" {
		msg := formatText(v, h.TimeFormat, h.AppendDelimiter)
		_, err = file.WriteString(msg)
	} else {
		var data []byte
		data, err = formatJson(v, h.TimeFormat)
		if err == nil {
			_, err = file.Write(data)
		}
	}

	if err == nil {
		h.file = file
	}
	return
}

func testDir(d string) (string, error) {
	s, err := os.Stat(d)
	if err != nil {
		return "", err
	}
	if !s.IsDir() {
		return "", errors.New("不是一个有效的目录")
	}
	d, err = filepath.Abs(d)
	if err != nil {
		return "", err
	}
	err = os.Mkdir(d+"/test.tmp000000", os.ModeDir)
	if err != nil {
		return "", err
	}
	err = os.RemoveAll(d + "/test.tmp000000")
	if err != nil {
		return "", err
	}
	// TODO: 测试写入文件
	return d, nil
}

func openFile(ctx *FileHandler, v Entry) (file *os.File, err error) {
	var s os.FileInfo
	if ctx.file != nil {
		file = ctx.file
		s, err = file.Stat()
		if err == nil && ctx.RotateSize > 0 && s.Size() > ctx.RotateSize {
			return
		}
		file.Close()
		ctx.file = nil
		file = nil
	}
	dir := ctx.Path
	cat := ""
	if ctx.PathTimeFormat != "" {
		cat = v.GetTimestamp().Format(ctx.PathTimeFormat)
		dir += "/" + cat

		if s, derr := os.Stat(dir); derr != nil || !s.IsDir() {
			err = os.Mkdir(dir, os.ModePerm)
			if err != nil {
				return
			}
		}
	}
	suffix := ".log"
	prefix := v.GetPrefix()
	num := 0
	for {
		filename := ""
		if prefix != "" {
			filename += prefix + "."
		}
		filename += v.GetTimestamp().Format("20060102")
		if num != 0 {
			filename += "." + strconv.Itoa(num)
		}
		num++
		filename += suffix
		filename = dir + "/" + filename
		s, err = os.Stat(filename)
		if err != nil {
			return os.Create(filename)
		}
		if ctx.RotateSize > 0 && s.Size() > ctx.RotateSize {
			continue
		}
		return os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	}
}
