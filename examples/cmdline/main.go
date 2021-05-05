package main

import (
	"encoding/base64"
	"fmt"
	"os"

	iid "github.com/theovassiliou/instanceidentification"
)

func main() {

	instanceId := os.Args[1]

	ciid := iid.NewCiid(instanceId)
	fmt.Println(ciid.String())
	if ciid.Miid.Vn == "x" {

		ciid.WithDecoding(func(s string) string {
			b1, _ := base64.StdEncoding.DecodeString(s)
			return string(b1)
		})
		fmt.Println("Decoded: ", ciid.String())
		fmt.Println("As tree: \n", ciid.PrintExtendedCiid())
	}
}
