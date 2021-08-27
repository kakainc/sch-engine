package main

import (
	"Sch_engine/init_sch"
	"Sch_engine/process"
	"Sch_engine/query_seg"
	"Sch_engine/ranking"
	"Sch_engine/recall"
	"Sch_engine/rely"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

func searchNewHandler() http.Handler {
	rd, mg := init_sch.ServiceDbInit()
	y := init_sch.SegInit()

	cacheContent:=map[int64]string{}
	cacheComment:=map[int64]string{}
	var cLock *sync.Mutex
	var mLock *sync.Mutex

	cacheNInv := map[string][]rely.PidScore{}
	cacheMInv := map[string][]rely.PidScore{}
	cacheTInv := map[string][]rely.PidScore{}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		init_sch.Info.Printf("Remote User(Ip: %s) have connect!", strings.Split(r.RemoteAddr, ":")[0])
		querryMap, _ := ioutil.ReadAll(r.Body)
		init_sch.Info.Println("get query.....")
		var m map[string]interface{}
		_ = json.Unmarshal(querryMap, &m)
		query := m["query"].(string)
		topN := m["num"].(float64)
		recallN := 2*topN

		startTime := time.Now().UnixNano()
		iterms := query_seg.QuerrySeg(y, query)
		endTime1 := time.Now().UnixNano()
		pids := recall.DoDbaRecall(iterms, rd, int(recallN), cacheNInv, cacheMInv, cacheTInv)
		endTime2 := time.Now().UnixNano()
		rPids := ranking.DoRanking(query, pids, rd, mg)
		endTime3 := time.Now().UnixNano()

		//finalPids := process.DoDbProcess(rPids, rd, int(topN))   // 仅返回pid列表
		finalPids := process.DoDbaProcess(rPids, rd, mg, int(topN), cacheContent, cacheComment, cLock, mLock)

		var ret int
		var msg string

		defer func() {
			if err := recover(); err != nil {
				init_sch.Error.Println("func http_fn get error!!!!")
				ret = -1
				msg = err.(string)
				finalPids = "May be meet some accident !"
			} else {
				ret = 1
				msg = ""
			}

			retData := map[string]interface{}{
				"ret": ret,
				"msg": msg,
				"data": finalPids,
			}
			data, err := json.Marshal(retData)
			_, err = w.Write(data)
			if err != nil {
				init_sch.Error.Println("http write json err!!!!!")
			}
		}()

		endTime4 := time.Now().UnixNano()

		milliseconds1:= float64((endTime1 - startTime) / 1e6)
		milliseconds2:= float64((endTime2 - endTime1) / 1e6)
		milliseconds3:= float64((endTime3 - endTime2) / 1e6)
		milliseconds4:= float64((endTime4 - endTime3) / 1e6)
		init_sch.Info.Printf("query_seg cost %0.2f ms\n", milliseconds1)
		init_sch.Info.Printf("recall cost %0.2f ms\n", milliseconds2)
		init_sch.Info.Printf("ranking cost %0.2f ms\n", milliseconds3)
		init_sch.Info.Printf("process cost %0.2f ms\n", milliseconds4)
	})
}


func main() {
	mux := http.NewServeMux()
	th := searchNewHandler()
	mux.Handle("/search", th)
	init_sch.Info.Println("Listening...")

	http.ListenAndServe(":3111", mux)
}

