package main

import (
	"context"
	"fmt"

	"github.com/prince2014/go-utils/concurrency"
)

func main() {
	// rand.Seed(time.Now().UnixNano())
	fmt.Println(concurrency.ErrOrDone())
	fmt.Println(concurrency.CatchErr(context.TODO()))
	concurrency.WaitGroup()
	concurrency.WaitGroup2()

}
