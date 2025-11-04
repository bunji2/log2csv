package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"
)

type AL struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	LogFmt      string            `json:"log_fmt"`
	ItemNames   []string          `json:"item_names"`
	DateFmts    map[string]string `json:"date_fmts"`
	OutDateFmt  string            `json:"out_date_fmt"`
	re          *regexp.Regexp
}

func LoadAL(filePath string) (r *AL, err error) {
	bytes, e := os.ReadFile(filePath)
	if e != nil {
		err = e
		return
	}

	var al AL
	e = json.Unmarshal(bytes, &al)
	if e != nil {
		err = e
		return
	}

	e = verifyAL(al)
	if e != nil {
		err = e
		return
	}

	re, e := regexp.Compile(al.LogFmt)
	if e != nil {
		err = e
		return
	}
	al.re = re

	r = &al
	return
}

func verifyAL(al AL) (err error) {
	if al.LogFmt == "" {
		err = errors.New("LogFmt is empty")
		return
	}

	if len(al.ItemNames) == 0 {
		err = errors.New("ItemNames is empty")
		return
	}

	if len(al.DateFmts) > 0 && al.OutDateFmt == "" {
		err = errors.New("OutDateFmt is empty unless DateFmts is not empty")
		//return
	}

	return
}

func (a *AL) Parse(line string) (items []string, err error) {
	if !a.re.MatchString(line) {
		err = fmt.Errorf("not match")
		return
	}

	match := a.re.FindStringSubmatch(line)
	if len(match) == len(a.ItemNames) {
		err = fmt.Errorf("matched number is abnormal")
		return
	}

	for _, itemName := range a.ItemNames {
		item := match[a.re.SubexpIndex(itemName)]
		dateFmt, found := a.DateFmts[itemName]
		if found {
			//if itemName == a.DateItemName {
			t, e := time.Parse(dateFmt, item)
			if e != nil {
				err = e
				return
			}
			item = t.Format(a.OutDateFmt)
		}
		items = append(items, item)
	}

	return
}
