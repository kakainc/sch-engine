package ranking

import (
	"Sch_engine/init_sch"
	"Sch_engine/rely"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
)

func HotRank(rd *redis.Client, pids []rely.PidScore, hotChan chan map[string]float64) {
	scores := map[string]float64{}
	init_sch.Info.Println("go hotrank......")
	minscore := 100.0
	var maxscore float64
	for _, pid := range pids {
		word := fmt.Sprintf("new_review")
		hotscore := init_sch.ZGet(rd, word, strconv.FormatInt(pid.Pid, 10)) * -1.0
		if hotscore > maxscore {
			maxscore = hotscore
		}
		if hotscore < minscore {
			minscore = hotscore
		}
		if hotscore == 0 {

		}
		scores[strconv.FormatInt(pid.Pid, 10)] = hotscore
	}
	init_sch.Info.Printf("hot score min: %0.2f max: %0.2f", minscore, maxscore)
	for pid, score := range scores {
		scores[pid] = rely.Norm(score, maxscore, minscore)
		init_sch.Info.Printf("%s hotrank score: %f", pid, scores[pid])
	}
	hotChan  <- scores
	close(hotChan)
	//wg.Done()
}
