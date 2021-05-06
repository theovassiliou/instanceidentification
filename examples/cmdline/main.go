package main

import (
	"fmt"
	"os"

	iid "github.com/theovassiliou/instanceidentification"
)

func main() {

	instanceId := os.Args[1]

	ciid := iid.NewStdCiid(instanceId)
	fmt.Println(ciid.String())
	fmt.Println(ciid.TreePrint())
}
