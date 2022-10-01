package main

import (
	"context"
	"fmt"
	"github.com/chinaDL/dictGenerate"
)

func main() {
	i := 1
	dictGenerate.GenerateDo(dictGenerate.HexdigitsLowercase, 32, func(s string, cancelFunc context.CancelFunc) {
		fmt.Println(s)
		i++
		if i > 5 {
			cancelFunc()
		}
	})
}
