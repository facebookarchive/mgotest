package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/facebookgo/mgotest"
)

func main() {
	l := log.New(os.Stdout, "", log.LstdFlags)
	n := flag.Uint("n", 3, "num nodes")
	flag.Parse()
	rs := mgotest.NewReplicaSet(*n, l)
	fmt.Println(rs.Addrs())
	select {}
}
