package goroutin

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type threadNo int
type lastProcessedID int
type threadProcessed struct {
	ThreadNo        threadNo
	LastProcessedID lastProcessedID
}

func main() {

	const threadCount = 50

	updator := threadUpdator(threadCount)
	defer close(updator)

	ids := getIDs() // 中断してたときの途中から開始する処理は省略してます。

	var wg sync.WaitGroup
	for c := 0; c < threadCount; c++ {
		wg.Add(1)
		go func(thread int) {
			defer wg.Done()
			for id := range ids {
				// 更新処理だよー

				rand.Seed(time.Now().UnixNano())
				r := time.Duration(rand.Intn(100) + 1)
				time.Sleep(r * time.Microsecond)

				updator <- threadProcessed{ThreadNo: threadNo(thread), LastProcessedID: lastProcessedID(id)}
			}
		}(c)
	}
	wg.Wait()
	fmt.Println("正常終了")
	time.Sleep(1 * time.Second)
}

func getIDs() <-chan int { // 直にチャネルにしてるけど、元はsliceをチャネルに入れ替えてます。
	const max = 1000
	ids := make(chan int, max)
	for i := 0; i < max; i++ {
		ids <- i + 1
	}
	close(ids)
	return ids
}

func threadUpdator(threadCount threadNo) chan<- threadProcessed { // 名前は適当
	out := make(chan threadProcessed)
	go func() {
		//スレッドごとの処理が完了したIDを保持する。初期化
		processedIDs := map[threadNo]lastProcessedID{}
		for c := threadNo(0); c < threadCount; c++ {
			processedIDs[c] = 0
		}

		var savedID lastProcessedID //DBを更新した値

		//var times int  // 更新した回数
		c := 0
		for tp := range out { // chanに値が入ったら実行される
			processedIDs[tp.ThreadNo] = tp.LastProcessedID

			//処理済みIDの中で最小のもの
			min := processedIDs[0]
			for _, processedID := range processedIDs {
				if min > processedID {
					min = processedID
				}
			}

			//更新処理
			//timeをみて100回ごととかでいいかも
			if savedID != min {
				//更新処理
				savedID = min
				c++
				fmt.Printf("c	%d	￿\n", c)
			}

		}
	}()
	return out
}

// テスト
func idUpdator() chan<- int { // 名前は適当
	out := make(chan int)
	go func() {
		var lastID int // 最後に更新されたid
		var times int  // 更新した回数

		for id := range out {
			if lastID < id { // さっきのより大きい
				lastID = id
				times++
				if times%10 == 0 { // 100回更新したから最後を記録する
					fmt.Println("update last id", times, lastID)
				}
			}
		}
	}()
	return out
}
