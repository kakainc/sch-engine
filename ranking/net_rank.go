package ranking

import (
	"Sch_engine/init_sch"
	"Sch_engine/rely"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"sort"
	"strconv"
)

func DoRanking(querry string, pids []rely.PidScore, rd *redis.Client, mg *mongo.Collection) []rely.PidScore {
	hotChan := make(chan map[string]float64, 1)
	ctrChan := make(chan map[string]float64, 1)
	bmChan := make(chan map[string]float64, 1)
	rkChan := make(chan map[string]float64, 1)
	//wg := sync.WaitGroup{}
	//wg.Add(2)

	go HotRank(rd, pids, hotChan)
	go CtrRank(rd, pids, ctrChan)
	go GetBm25(querry, pids, mg, bmChan)
	go post(querry, pids, mg, rkChan)
	//wg.Wait()

	hotscores, _ := <-hotChan
	ctrscores, _ := <-ctrChan
	bmscores, _ := <-bmChan
	rkscores, _ := <-rkChan
	hotpids := rankingResult(pids, hotscores, rely.HOTWEIGHT)
	ctrpids := rankingResult(hotpids, ctrscores, rely.CTRWEIGHT)
	bmpids := rankingResult(ctrpids, bmscores, rely.BMWEIGHT)
	rkpids := rankingResult(bmpids, rkscores, rely.MODELWEIGHT)

	sort.Slice(rkpids, func(i, j int) bool {
		return rkpids[i].Score > rkpids[j].Score
	})

	init_sch.Info.Printf("ranking get %d posts!", len(rkpids))
	return rkpids
}

func rankingResult(pids []rely.PidScore, rkscore map[string]float64, weight float64) []rely.PidScore {
	for _, pid := range pids {
		pid.Score = pid.Score + rkscore[strconv.FormatInt(pid.Pid, 10)] * weight
	}
	return pids
}
