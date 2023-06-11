package main

import (
	rs "github.com/vivint/infectious"
	"time"
)

func main() {
	fec, _ := rs.NewFEC(8, 20)
	messageSize := 126 * 2 * 2

	good := true
	iterations := 0

	go func() {
		time.Sleep(time.Second * 10)
		good = false
	}()

	for good {
		iterations += 1
		paddedData := make([]byte, messageSize)
		for i := range paddedData {
			paddedData[i] = 2
		}

		shares := make([]rs.Share, 0)

		output := func(s rs.Share) {
			shares = append(shares, s.DeepCopy())
		}

		err := fec.Encode(paddedData, output)
		if err != nil {
			panic(err)
		}

		shares = shares[6:]

		err = fec.Rebuild(shares, func(share rs.Share) {

		})
		if err != nil {
			return
		}
		if err != nil {
			panic(err)
		}
	}

	println(iterations)
}
