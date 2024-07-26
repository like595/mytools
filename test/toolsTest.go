package main

import (
	"fmt"
	"github.com/like595/mytools/vtools"
)

func main() {
	data := 65526
	fmt.Println(vtools.BytesToInt16(byte(data/0x100), byte(data)))

	var data1 int
	data1 = -10
	
	fmt.Println(vtools.IntToBytes(data1))
}
