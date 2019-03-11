package main

import "time"

func main() {
	var timer *time.Timer
	timer = time.NewTimer(time.Second)
	go func() {
		for {
			select {
			case <-timer.C:
				println("Timer expired")
			}
		}

	}()
	println("XXXXX")
	stop := timer.Stop()
	println("Timer cancelled:", stop)

	timer.Reset(time.Second * 10)

	time.Sleep(time.Second * 20)
	println("XXXXX")
}
