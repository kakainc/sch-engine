package ranking

import (
	"Sch_engine/init_sch"
	"Sch_engine/rely"
	"github.com/go-nlp/bm25"
	"github.com/go-nlp/tfidf"
	"go.mongodb.org/mongo-driver/mongo"
	"sort"
	"strconv"
	"strings"
)

type doc []int

func (d doc) IDs() []int { return []int(d) }

func makeCorpus(a []string) (map[string]int, []string) {
	retVal := make(map[string]int)
	invRetVal := make([]string, 0)
	var id int
	for _, s := range a {
		for _, f := range strings.Fields(s) {
			f = strings.ToLower(f)
			if _, ok := retVal[f]; !ok {
				retVal[f] = id
				invRetVal = append(invRetVal, f)
				id++
			}
		}
	}
	return retVal, invRetVal
}

func makeDocuments(a []string, c map[string]int) []tfidf.Document {
	retVal := make([]tfidf.Document, 0, len(a))
	for _, s := range a {
		var ts []int
		for _, f := range strings.Fields(s) {
			f = strings.ToLower(f)
			id := c[f]
			ts = append(ts, id)
		}
		retVal = append(retVal, doc(ts))
	}
	return retVal
}

func GetBm25 (querry string, pids []rely.PidScore, mg *mongo.Collection, bmChan chan map[string]float64) {
	var contents []string
	var pidsNums []string

	for _,pid := range pids {
		content := init_sch.FindOne(mg, strconv.FormatInt(pid.Pid, 10)).Content
		pidsNums = append(pidsNums, strconv.FormatInt(pid.Pid, 10))
		contents = append(contents, content)
	}

	init_sch.Info.Println("go bm25......")

	corpus, _ := makeCorpus(contents)
	docs := makeDocuments(contents, corpus)
	tf := tfidf.New()
	for _, doc := range docs {
		tf.Add(doc)
	}
	tf.CalculateIDF()
	isha := doc{corpus[querry]}
	ishaScores := bm25.BM25(tf, isha, docs, 1.5, 0.75)

	sort.Sort(sort.Reverse(ishaScores))
	scores := map[string]float64{}
	minscore := 100.0
	var maxscore float64

	for _, d := range ishaScores {
		scores[pidsNums[d.ID]] = d.Score
		//init_sch.Info.Printf("\tID   : %d\n\tScore: %1.3f\n\tDoc  : %q\n", d.ID, d.Score, contents[d.ID])
		if d.Score > maxscore {
			maxscore = d.Score
		}
		if d.Score < minscore {
			minscore = d.Score
		}
	}

	for pid, score := range scores {
		scores[pid] = rely.Norm(score, maxscore, minscore)
		init_sch.Info.Printf("%s bm25rank score: %f", pid, scores[pid])
	}

	bmChan  <- scores
	close(bmChan)
	//wg.Done()
}

