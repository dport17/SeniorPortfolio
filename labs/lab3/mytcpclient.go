package main

import "pkt"
import "net"

func main(){
	conn, err := net.Dial("tcp", "golag.org:80")
	if err != nil{
		//handle error
	}
 
}