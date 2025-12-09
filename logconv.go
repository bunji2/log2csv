package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"
)

type LogConv struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	LogFmt      string            `json:"log_fmt"`
	ItemNames   []string          `json:"item_names"`
	DateFmts    map[string]string `json:"date_fmts"`
	OutDateFmt  string            `json:"out_date_fmt"`
	re          *regexp.Regexp
}

func LoadConf(filePath string) (r *LogConv, err error) {
	bytes, e := os.ReadFile(filePath)
	if e != nil {
		err = e
		return
	}

	var lc LogConv
	e = json.Unmarshal(bytes, &lc)
	if e != nil {
		err = e
		return
	}

	e = lc.verify()
	if e != nil {
		err = e
		return
	}

	re, e := regexp.Compile(lc.LogFmt)
	if e != nil {
		err = e
		return
	}
	lc.re = re

	r = &lc
	return
}

func (lc *LogConv) verify() (err error) {
	if lc.LogFmt == "" {
		err = errors.New("LogFmt is empty")
		return
	}

	if len(lc.ItemNames) == 0 {
		err = errors.New("ItemNames is empty")
		return
	}

	if len(lc.DateFmts) > 0 && lc.OutDateFmt == "" {
		err = errors.New("OutDateFmt is empty unless DateFmts is not empty")
		//return
	}

	return
}

// 一行分のログから項目をパースする
func (lc *LogConv) Parse(line string) (items []string, err error) {
	if !lc.re.MatchString(line) {
		err = fmt.Errorf("not match")
		return
	}

	match := lc.re.FindStringSubmatch(line)
	if len(match) == len(lc.ItemNames) {
		err = fmt.Errorf("matched number is abnormal")
		return
	}

	for _, itemName := range lc.ItemNames {
		item := match[lc.re.SubexpIndex(itemName)]
		dateFmt, found := lc.DateFmts[itemName]
		if found {
			//if itemName == lc.DateItemName {
			t, e := time.Parse(dateFmt, item)
			if e != nil {
				err = e
				return
			}
			item = t.Format(lc.OutDateFmt)
		}
		items = append(items, item)
	}

	return
}

func (lc *LogConv) Process(filePath string) (err error) {

	w := csv.NewWriter(os.Stdout)

	// ヘッダの出力
	e := w.Write(lc.ItemNames)
	if e != nil {
		err = e
		return
	}

	e = lc.processLogFile(filePath, func(items []string) error {
		return w.Write(items)
	})
	if e != nil {
		err = e
		return
	}

	w.Flush()

	e = w.Error()
	if e != nil {
		err = e
		//return
	}
	return
}

func (lc *LogConv) processLogFile(filePath string, write func([]string) error) (err error) {

	f, e := os.Open(filePath)
	if e != nil {
		err = e
		return
	}
	defer f.Close()

	// ログファイルを1行ずつ読み出して項目をパースする
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		items, e := lc.Parse(line)
		if e != nil {
			err = e
			return
		}

		// パースした項目を出力
		e = write(items)
		if e != nil {
			err = e
			return
		}
	}

	e = sc.Err()
	if e != nil {
		err = e
		// return
	}

	return
}
