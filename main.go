package main

import (
	"encoding/xml"
	"sync"

	"fmt"

	"github.com/piotrjura/darwin/config"
)

func main() {
	conf := config.ReadConfig()

	var wg sync.WaitGroup
	wg.Add(2)

	rC := make(chan []byte)
	tC := make(chan []byte)

	refFile, timetableFile := getReferenceFilenames(conf.Ftp)

	go downloadXML(refFile, &wg, rC, conf.Ftp)
	go downloadXML(timetableFile, &wg, tC, conf.Ftp)

	t := <-tC
	r := <-rC
	var timetable Timetable
	err := xml.Unmarshal(t, &timetable)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(r))

	wg.Wait()
}
