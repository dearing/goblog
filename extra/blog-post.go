package main

import (
	"flag"
	"fmt"
	"log"

	store "github.com/dearing/blog/storage/redis"
	
)

var host = flag.String("host", "192.168.1.150:6379", "Redis host")
var pass = flag.String("pass", "", "Redis pass")
var db = flag.Int64("db", -1, "Redis DB index")

func main() {

	flag.Parse()

	store.Connect(*host, *pass, *db)
	defer store.Close()
	log.Println("conntected to redis", *host, *db)

	for i, k := range store.GetPosts().Val() {
		fmt.Printf("[%v]:%s\n", i, k)
	}

	input := ""
	for {
		input, _ = readLine()
		log.Println(input)
	}
}

func readLine() (l string, e error) {
	_, e = fmt.Scanln(&l)
	if e != nil {
		log.Println(e)
	}
	return l, e
}
