package recall

import (
	"Sch_engine/init_sch"
	"Sch_engine/rely"
	"fmt"
	"github.com/go-redis/redis"
)

func doMatchDbaRecall(querryitem []string, rd *redis.Client, contentChan chan map[rely.PidScore]float64,
	cacheNInv map[string][]rely.PidScore, recallN int64) {
	pidMap := map[rely.PidScore]float64{}

	for _, word := range querryitem {
		var ps []rely.PidScore
		word := fmt.Sprintf("内容_"+word)
		if cn, ok := cacheNInv[word]; ok && len(cn) >= int(recallN) {
			ps = cacheNInv[word][:recallN]
		} else {
			ps = init_sch.ZRange(rd, word, recallN)
			cacheNInv[word] = ps
		}

		if len(ps)>0 {
			for _, p := range ps {
				p.Score = p.Score * -1.0
				if _, ok := pidMap[p]; ok {
					pidMap[p]+=1
				} else {
					pidMap[p] = 1
				}
			}
		}
	}

	init_sch.Info.Printf("content get %d posts!!", len(pidMap))
	contentChan  <- pidMap
	close(contentChan)
}

func doCommentDbaRecall(querryitem []string, rd *redis.Client, commentChan chan map[rely.PidScore]float64,
	cacheMInv map[string][]rely.PidScore, recallN int64) {

	pidMap := map[rely.PidScore]float64{}
	for _, word := range querryitem {
		var ps []rely.PidScore
		word := fmt.Sprintf("评论_"+word)
		if cn, ok := cacheMInv[word]; ok && len(cn) >= int(recallN) {
			ps = cacheMInv[word][:recallN]
		} else {
			ps = init_sch.ZRange(rd, word, recallN)
			cacheMInv[word] = ps
		}

		if len(ps)>0 {
			for _, p := range ps {
				p.Score = p.Score * -1.0
				if _, ok := pidMap[p]; ok {
					pidMap[p]+=1
				} else {
					pidMap[p] = 1
				}
			}
		}
	}

	init_sch.Info.Printf("comment get %d posts!!", len(pidMap))
	commentChan <- pidMap
	close(commentChan)
}

func doTnameDbaRecall(querryitem []string, rd *redis.Client, tnameChan chan map[rely.PidScore]float64,
	cacheTInv map[string][]rely.PidScore, recallN int64) {

	pidMap := map[rely.PidScore]float64{}
	for _, word := range querryitem {
		var ps []rely.PidScore
		word := fmt.Sprintf("话题_"+word)
		if cn, ok := cacheTInv[word]; ok && len(cn) >= int(recallN) {
			ps = cacheTInv[word][:recallN]
		} else {
			ps = init_sch.ZRange(rd, word, recallN)
			cacheTInv[word] = ps
		}

		if len(ps)>0 {
			for _, p := range ps {
				p.Score = p.Score * -1.0
				if _, ok := pidMap[p]; ok {
					pidMap[p]+=1
				} else {
					pidMap[p] = 1
				}
			}
		}
	}

	init_sch.Info.Printf("tname get %d posts!!", len(pidMap))
	tnameChan <- pidMap
	close(tnameChan)
}
