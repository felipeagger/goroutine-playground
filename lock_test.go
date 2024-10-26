// go test -bench=. lock_test.go -cpuprofile cpuatm.prof -memprofile mematm.prof -blockprofile blockatm.out -mutexprofile mutexatm.out -trace=traceatm.out -race -benchmem -count=1 -benchtime=10000x -v
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var (
	counterOne int
	mu         sync.Mutex
)

func incrementWithMutex(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		mu.Lock()
		counterOne++
		mu.Unlock()
	}
}

func incrementWithAtomic(wg *sync.WaitGroup, atmCount *atomic.Int32) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		atmCount.Add(1)
	}
}

func incrementWithChannel(ch chan int, w *sync.WaitGroup) {
	defer w.Done()
	for i := 0; i < 100; i++ {
		ch <- 1
	}
}

func BenchmarkLockContentionMutex(b *testing.B) {
	var wg sync.WaitGroup
	now := time.Now()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go incrementWithMutex(&wg)
	}

	wg.Wait()
	dur := time.Now().Sub(now)
	fmt.Printf("\nFinal counterOne value: %v; duration: %d ns\n", counterOne, dur.Nanoseconds())
}

func BenchmarkLockContentionAtomic(b *testing.B) {
	var wg sync.WaitGroup
	var atmCount   atomic.Int32
	now := time.Now()

	atmCount.Store(0)

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go incrementWithAtomic(&wg, &atmCount)
	}

	wg.Wait()
	dur := time.Now().Sub(now)
	fmt.Printf("\nFinal atmCount value: %v; duration: %d ns\n", atmCount.Load(), dur.Nanoseconds())
}

func BenchmarkLockContentionChannels(b *testing.B) {
	var wgTwo sync.WaitGroup
	counterTwo := 0
	ch := make(chan int, 1000000)

	now := time.Now()

	for i := 0; i < b.N; i++ {
		wgTwo.Add(1)
		go incrementWithChannel(ch, &wgTwo)
	}

	go func() {
		wgTwo.Wait()
		close(ch)
	}()

	for val := range ch {
		counterTwo += val
	}

	dur := time.Now().Sub(now)
	fmt.Printf("\nFinal counterTwo value: %v; duration: %d ns\n", counterTwo, dur.Nanoseconds())
}
