package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main()  {

	fmt.Println("simpleAsyncCalls")
	simpleAsyncCalls()

	fmt.Println("\nasyncCallsWithResults")
	asyncCallsWithResults()

}

func simpleAsyncCalls () {
	var wg sync.WaitGroup
	var calls int = 5

	wg.Add(calls)

	for i := 1; i <= calls; i++ {

		go func(wg *sync.WaitGroup, idx int) {
			defer wg.Done()

			//Your process here
			time.Sleep(2 * time.Second)
			fmt.Printf("execution: %v\n", idx)

		}(&wg, i)
	}

	wg.Wait()
}

func asyncCallsWithResults () {
	var wg sync.WaitGroup
	var results []chanResp
	users := []string{"felipeagger", "maria", "joao", "mateus", "marcos"}

	wg.Add(len(users))

	channel := make(chan chanResp, len(users))
	for i, username := range users {

		go processGetUser(&wg, channel, username, i)
	}

	//Wait only if you want to wait for all responses... if you want to stream them, comment the wait()
	//wg.Wait()

	for range users {
		select {
			case resp := <-channel:
				results = append(results, resp)
				fmt.Printf("\nReceived user: %v\n", resp.User.Name)
		}
	}

	close(channel)
	fmt.Println(results)
}

func processGetUser(wg *sync.WaitGroup, channel chan chanResp, username string, idx int) {
	defer wg.Done()

	var err error

	//external request simulation
	time.Sleep(time.Duration(idx) * time.Second)

	user := UserInfo{
		Name: username,
		ID: rand.Intn(10000),
	}

	//real call to github
	//user, err := GetUserInfo(username)

	channel <- chanResp{user, err}
}