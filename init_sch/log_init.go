package init_sch

import (
	"Sch_engine/rely"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	Trace   *log.Logger // 不显示仅记录信息
	Info    *log.Logger // 重要的信息
	Warning *log.Logger // 需要注意的信息
	Error   *log.Logger // 非常严重的问题
)

func init() {

	if _, err := os.Stat(rely.LogFile); os.IsNotExist(err) {
		f, _ := os.Create(rely.LogFile)
		_ = f.Close()
	}

	file, err := os.OpenFile(rely.LogFile,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	Trace = log.New(ioutil.Discard,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(io.MultiWriter(file, os.Stdout),
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(io.MultiWriter(file, os.Stdout),
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(io.MultiWriter(file, os.Stderr),
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

