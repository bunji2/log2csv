package main

import (
	"bufio"
	"encoding/csv"
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
	if len(os.Args) < 2 {
		err = fmt.Errorf(
			"[USAGE] %s config.json logfile.txt > out.csv",
			os.Args[0])
		return
	}

	//al := CommonLog()
	al, e := LoadAL(os.Args[1])
	if e != nil {
		err = e
		return
	}

	w := csv.NewWriter(os.Stdout)

	e = w.Write(al.ItemNames)
	if e != nil {
		err = e
		return
	}

	f, e := os.Open(os.Args[2])
	if e != nil {
		err = e
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		//fmt.Println(line)
		items, e := al.Parse(line)
		if e != nil {
			err = e
			return
		}

		e = w.Write(items)
		if e != nil {
			err = e
			return
		}
	}

	w.Flush()

	e = sc.Err()
	if e != nil {
		err = e
		return
	}

	e = w.Error()
	if e != nil {
		err = e
		//return
	}

	return
}
