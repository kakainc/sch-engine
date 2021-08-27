package query_seg

import (
	"Sch_engine/init_sch"
	"github.com/huichen/sego"
)

func QuerrySeg(segmenter sego.Segmenter, querry string) []string {
	text := []byte(querry)
	segments := segmenter.Segment(text)

	// 普通模式和搜索模式
	items := sego.SegmentsToSlice(segments, true)
	init_sch.Info.Printf("Extract result: %v", items)

	return items
}
