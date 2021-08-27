package ranking

import (
	"Sch_engine/init_sch"
	"Sch_engine/rely"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
)

func CtrRank(rd *redis.Client, pids []rely.PidScore, ctrChan chan map[string]float64) {
	scores := map[string]float64{}
	init_sch.Info.Println("go ctrrank......")
	minscore := 100.0
	var maxscore float64
	for _, pid := range pids {
		word := fmt.Sprintf("new_ctr")
		ctrscore := init_sch.HGet(rd, word, strconv.FormatInt(pid.Pid, 10))
		if ctrscore > maxscore {
			maxscore = ctrscore
		}
		if ctrscore < minscore {
			minscore = ctrscore
		}
		scores[strconv.FormatInt(pid.Pid, 10)] = ctrscore
	}
	init_sch.Info.Printf("ctr score min: %0.2f max: %0.2f", minscore, maxscore)
	for pid, score := range scores {
		scores[pid] = rely.Norm(score, maxscore, minscore)
		init_sch.Info.Printf("%s ctrrank score: %f", pid, scores[pid])
	}
	ctrChan  <- scores
	close(ctrChan)
	//wg.Done()
}
