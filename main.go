package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

	fmt.Println(strings.Join(al.ItemNames, ","))
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
		//fmt.Printf("=>[%s]\n", strings.Join(items[1:], "]["))
		fmt.Println(strings.Join(items, ","))
	}

	e = sc.Err()
	if e != nil {
		err = e
	}

	return
}
