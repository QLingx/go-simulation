package main

import (
	"log"
	"strings"
	"bufio"
	"fmt"
	"net"
	"go-simulation/pubsub"
)


func main(){
	ps := pubsub.New()

	l,err := net.Listen("tcp",":5555")

	if err != nil{
		fmt.Println(err)
	}

	fmt.Println("Listening",l.Addr().String())
	fmt.Println("Clients can connect to this server like follow:")
	log.Println(" $ telnet server:5555")

	for{
		c,err := l.Accept()
		if err != nil{
			fmt.Println(err)
		}
		go func(c net.Conn){
			buf := bufio.NewReader(c)
			ps.Sub(func(t string){
				fmt.Println(t)
				c.Write([]byte(t + "\n"))
			})
			fmt.Println("Subscribed",c.RemoteAddr().String())
			for{
				b,_,err := buf.ReadLine()
				if err != nil{
					fmt.Println("Closed",c.RemoteAddr().String())
					break
				}
				ps.Pub(strings.TrimSpace(string(b)))
			}
			ps.Leave(nil)
		}(c)
	}
}