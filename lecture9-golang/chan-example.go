package main

import "fmt"

func main(){
	i := make(chan int)
	go func(){
		i <- 1
	}()
	fmt.Println(<-i)
}
