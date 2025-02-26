package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
)

const (
	goroutineNum = 10
	cr           = "\n"
	start        = -1
)

func print(n int) error {
	buf := bytes.Buffer{}
	ic := strconv.Itoa(n)
	if _, err := buf.WriteString(ic); err != nil {
		return err
	}
	buf.WriteString(cr)
	io.WriteString(os.Stdout, buf.String())
	return nil
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(goroutineNum)
	ch := make(chan int, 1)
	ch <- start
	for i := 0; i < 10; i++ {
		go func(i int) {
			for {
				n := <-ch
				if n != i-1 {
					ch <- n
					continue

				}
				if err := print(i); err != nil {
					fmt.Println(err)
					wg.Done()
					return
				}
				ch <- i
				break
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
}
