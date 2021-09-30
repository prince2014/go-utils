package concurrency

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// Map type that can be safely shared between
// goroutines that require read/write access to a map
type ConcurrentSafeMap struct {
	sync.RWMutex
	items map[string]interface{}
}

type ConcurrentSafeList struct {
	sync.RWMutex
	items []interface{}
}

func ErrOrDone() error {
	var N = 5000
	quit := make(chan bool)
	errc := make(chan error)
	done := make(chan error)
	for i := 0; i < N; i++ {
		go func(i int) {
			// dummy fetch
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			err := error(nil)
			if i%(rand.Intn(100)+1) == 10 {
				err = fmt.Errorf("goroutine %d's error returned", i)
			}
			ch := done // we'll send to done if nil error and to errc otherwise
			if err != nil {
				ch = errc
			}
			select {
			case ch <- err:
				return
			case <-quit:
				return
			}
		}(i)
	}
	count := 0
	for {
		select {
		case err := <-errc:
			close(quit)
			return err
		case <-done:
			count++
			if count == N {
				return nil // got all N signals, so there was no error
			}
		}
	}
}

func CatchErr(ctx context.Context) error {
	errs, _ := errgroup.WithContext(ctx)

	// run all the http requests in parallel
	for i := 0; i < 4; i++ {
		errs.Go(func() error {
			// pretend this does an http request and returns an error
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			return fmt.Errorf("error in go routine, %v", i)
		})
	}

	// Wait for completion and return the first error (if any)
	return errs.Wait()
}

var countlen = 70
var res []int
var res2 []int

// The result is not correct due to conccurent write
func WaitGroup() {

	var wg sync.WaitGroup

	for i := 0; i < (countlen); i++ {

		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			defer wg.Done()
			res = append(res, i)
		}(&wg, i)
	}

	for i := 0; i < (countlen); i++ {

		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			// dateStr := now.AddDate(0, 0, i).Format("20060103")
			// res = append(res, dateStr)
			res2 = append(res2, i)
		}(i)
	}

	wg.Wait()

	fmt.Println(fmt.Sprintf("res lenght should be %v , but %v", countlen, len(res)))
	fmt.Println(fmt.Sprintf("res2 lenght should be %v , but %v", countlen, len(res2)))

}

type httpPkg struct{}

func (httpPkg) Get(url string) {}

var http httpPkg

// No concurrent write
func WaitGroup2() {
	var res []string

	var wg sync.WaitGroup
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.somename.com/",
	}
	for _, url := range urls {
		// Increment the WaitGroup counter.
		wg.Add(1)
		// Launch a goroutine to fetch the URL.
		go func(url string) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			// Fetch the URL.
			http.Get(url)
			res = append(res, url)
		}(url)
	}
	// Wait for all HTTP fetches to complete.
	wg.Wait()
	fmt.Println("Successfully fetched all URLs.")
	fmt.Printf("res length: %v \n", len(res))
}
