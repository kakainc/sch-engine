package init_sch

import (
	"github.com/huichen/sego"
	"os"
)

func SegInit () sego.Segmenter {
	var segmenter sego.Segmenter
	wd, _ := os.Getwd()
	Info.Printf("LoadDictionary from %s/phase_go.txt", wd)
	segmenter.LoadDictionary("./rely/phase_go.txt")

	//Info.Println("Segment init_sch already!")
	return segmenter
}
