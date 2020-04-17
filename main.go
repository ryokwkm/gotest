package main

import (
	"fmt"
	"sync"
)

func main() {
	updator := idUpdator()
	defer close(updator)

	ids := getIDs() // 中断してたときの途中から開始する処理は省略してます。

	var wg sync.WaitGroup
	for c := 0; c < 6; c++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range ids {
				// 更新処理だよー
				updator <- id
			}
		}()
	}
	wg.Wait()
	fmt.Println("正常終了")
}

func idUpdator() chan<- int { // 名前は適当
	out := make(chan int)
	go func() {
		var lastID int // 最後に更新されたid
		var times int  // 更新した回数
		for id := range out {
			if lastID < id { // さっきのより大きい
				lastID = id
				times++
				if times%100 == 0 { // 100回更新したから最後を記録する
					fmt.Println("update last id", times, lastID)
				}
			}
		}
	}()
	return out
}

func getIDs() <-chan int { // 直にチャネルにしてるけど、元はsliceをチャネルに入れ替えてます。
	const max = 100000
	ids := make(chan int, max)
	for i := 0; i < max; i++ {
		ids <- i + 1
	}
	close(ids)
	return ids
}
