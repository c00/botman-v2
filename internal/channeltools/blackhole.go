package channeltools

// Create a channel that outputs nothing
func BlackHoleChannel() chan string {
	ch := make(chan string)
	go func(ch chan string) {
		for range ch {
			//do nothing.
		}
	}(ch)

	return ch
}
