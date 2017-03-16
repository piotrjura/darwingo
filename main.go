package main

import (
	"sync"

	"github.com/piotrjura/darwin/config"
)

func main() {
	conf := config.ReadConfig()

	var wg sync.WaitGroup
	wg.Add(4)

	rC := make(chan []byte)
	tC := make(chan []byte)

	refFile, timetableFile := getReferenceFilenames(conf.Ftp)

	go downloadXML(refFile, &wg, rC, conf.Ftp)
	go downloadXML(timetableFile, &wg, tC, conf.Ftp)

	go parseReference(rC, &wg)
	go parseTimetables(tC, &wg)

	wg.Wait()
}
