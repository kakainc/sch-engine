package recall

import (
	"Sch_engine/init_sch"
	"Sch_engine/rely"
	"github.com/go-redis/redis"
	"math"
	"sort"
)

func DoDbaRecall(querryitem []string, rd *redis.Client, recallN int, cacheNInv map[string][]rely.PidScore,
	cacheMInv map[string][]rely.PidScore, cacheTInv map[string][]rely.PidScore ) []rely.PidScore {

	contentChan := make(chan map[rely.PidScore]float64, 1)
	commentChan := make(chan map[rely.PidScore]float64, 1)
	tnameChan := make(chan map[rely.PidScore]float64, 1)
	var resultRecall []rely.PidScore
	go doMatchDbaRecall(querryitem, rd, contentChan, cacheNInv, int64(int(math.Ceil(float64(3*recallN/4)))))
	go doCommentDbaRecall(querryitem, rd, commentChan, cacheMInv, int64(int(math.Ceil(float64(1*recallN/4)))))
	go doTnameDbaRecall(querryitem, rd, tnameChan, cacheTInv, int64(int(math.Ceil(float64(1*recallN/5)))))

	pidRecall := map[rely.PidScore]float64{}

	for {
		i, ok := <-contentChan
		if !ok {
			break
		}
		for k, v := range i {
			if _, ok := pidRecall[k]; ok {
				pidRecall[k] = pidRecall[k] + v * rely.CONTENTWEIGHT
			} else {
				pidRecall[k] = v * rely.CONTENTWEIGHT
			}
		}
	}

	for {
		j, ok := <-commentChan
		if !ok {
			break
		}
		for k, v := range j {
			if _, ok := pidRecall[k]; ok {
				pidRecall[k] = pidRecall[k] + v * rely.COMMENTWEIGHT
			} else {
				pidRecall[k] = v * rely.COMMENTWEIGHT
			}
		}
	}

	for {
		l, ok := <-tnameChan
		if !ok {
			break
		}
		for k, v := range l {
			if _, ok := pidRecall[k]; ok {
				pidRecall[k] = pidRecall[k] + v * rely.TNAMEWEIGHT
			} else {
				pidRecall[k] = v * rely.TNAMEWEIGHT
			}
		}
	}

	minscore := 100.0
	var maxscore float64

	for pidscore, cnt := range pidRecall {
		pidscore.Score = pidscore.Score * cnt * rely.REPEATWEIGHT
		if pidscore.Score > maxscore {
			maxscore = pidscore.Score
		}
		if pidscore.Score < minscore {
			minscore = pidscore.Score
		}
		resultRecall = append(resultRecall, pidscore)
	}

	init_sch.Info.Printf("recall score min: %0.2f max: %0.2f", minscore, maxscore)

	for i, pidsocre := range resultRecall {
		resultRecall[i].Score = rely.Norm(pidsocre.Score, maxscore, minscore)
	}

	sort.Slice(resultRecall, func(i, j int) bool {
		return resultRecall[i].Score > resultRecall[j].Score
	})

	if len(resultRecall) > recallN{
		resultRecall = resultRecall[:recallN]
	}

	init_sch.Info.Printf("recall get %d posts!!", len(resultRecall))
	return resultRecall
}
