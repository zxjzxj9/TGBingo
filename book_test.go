package main

import (
	"fmt"
	"testing"
)

func TestRequest(t *testing.T) {
	c := make(chan float32)
	wg.Add(1)
	go crawl(2105, c)
	go func(c <- chan float32) {
		for {
			val := <-c
			if int(val*1000) % 50 == 0 {
				fmt.Printf("Progess: %.2f \n", val)
			}
			if val == 1.0 {
				break
			}
		}
	}(c)
	wg.Wait()
}