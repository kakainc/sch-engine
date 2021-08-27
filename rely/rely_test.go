package rely

import (
	"fmt"
	"sort"
	"testing"
)


var complete chan int = make(chan int)

func loop() {
	fmt.Println("done1")
	complete <- 0 // 执行完毕
}

func nTest(nChan chan int) {
	for i := 1; i < 5; i++ {
		nChan <- i
	}
	close(nChan)
}

func mTest(mChan chan int) {
	for i := 11; i < 15; i++ {
		mChan <- i
	}
	close(mChan)
}

func Test_chan(t *testing.T) {
	nChan := make(chan int, 1000)   // 建立有缓冲区的chan
	mChan := make(chan int, 1000)
	var result []int

	go nTest(nChan)
	go mTest(mChan)

	//go loop()
	//<- complete   // 否则main在此阻塞住阻塞住main线,直至另一个goroutine跑完
	for {
		i, ok := <-mChan
		if !ok {
			break
		}
		result = append(result, i)
	}

	for {
		i, ok := <-nChan
		if !ok {
			break
		}
		result = append(result, i)
	}


	for _,i := range result {
		println(i)
	}
	println("done2")
}

//--------------------------------------------------

func Test_getAsc(t *testing.T) {
	s1 := []rune("都")
	fmt.Println(s1[0])
}

//--------------------------------------------------
func Test_div(t *testing.T) {
	fmt.Println(1.2/0.5)
}

//--------------------------------------------------
func max(vals...float64) float64 {
	var max float64
	for _, val := range vals {
		if val > max {
			max = val
		}
	}
	return max
}


func min(vals...float64) float64 {
	min := 100.0
	for _, val := range vals {
		if val <= min {
			min = val
		}
	}
	return min
}


func Test_mainx(t *testing.T) {
	val := []float64{15,77,57,238,54,108, 0.4, 0.4, 1.3, 0.002}
	fmt.Println(max(val...)) // "238"
	fmt.Println(min(val...)) // "15"

	sort.Slice(val, func(i, j int) bool {
		return val[i] > val[j]
	})
	fmt.Println(val) // "15"
}
