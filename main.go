package main

import (
	"fmt"

	"github.com/robertokbr/iago/pkg"
)

func main() {
	word, timesSaid := pkg.GetMostSaidWord()
	fmt.Printf("the word %s was said %d times", word, timesSaid)
}
