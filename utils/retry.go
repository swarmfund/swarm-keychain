package utils

func Retry(n int, fn func() error) (err error) {
	for try := 0; try < n; try++ {
		if err = fn(); err == nil {
			return
		}
	}
	return
}
