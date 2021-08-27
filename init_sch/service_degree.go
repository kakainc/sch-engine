package init_sch

import (
	"Sch_engine/rely"
	"fmt"
	"github.com/go-redis/redis"
)

func DoErrDegree(rd *redis.Client, recallN int64) []rely.PidScore {
	word := fmt.Sprintf("new_review")
	ps := ZRange(rd, word, recallN-1)

	Info.Printf(" Errdegree get %d posts!", len(ps))
	return ps
}


func DoTidDegree(rd *redis.Client, tid string, recallN int64) []rely.PidScore {
	word := fmt.Sprintf("tid_%s", tid)
	ps := ZRange(rd, word, recallN-1)

	Info.Printf("Tiddegree get %d posts!", len(ps))
	return ps
}
