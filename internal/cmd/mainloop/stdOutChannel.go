package mainloop

import (
	"fmt"
	"io"
	"sync"
)

// Create a channel that outputs to stdout
func stdOutChannel(wg *sync.WaitGroup, out io.Writer) chan string {
	wg.Add(1)
	ch := make(chan string)
	go func(ch chan string) {
		for chunk := range ch {
			fmt.Fprint(out, chunk)
		}
		wg.Done()
	}(ch)

	return ch
}
