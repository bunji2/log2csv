package main

import (
	"fmt"
	"os"
)

func main() {
	os.Exit(run())
}

func run() int {
	err := process()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

func process() (err error) {
	if len(os.Args) < 3 {
		err = fmt.Errorf(
			"[USAGE] %s config.json logfile1.txt ... logfileN.txt > out.csv",
			os.Args[0])
		return
	}

	conv, e := LoadConf(os.Args[1])
	if e != nil {
		err = e
		return
	}

	// 引数で指定された複数のログファイルをCSVファイルに変換する
	for _, logFile := range os.Args[2:] {
		err = conv.Process(logFile)
		if err != nil {
			break
		}
	}

	return
}
