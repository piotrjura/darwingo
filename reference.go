package main

import (
	"compress/gzip"
	"io/ioutil"
	"sync"

	"github.com/jlaffaye/ftp"
	"github.com/piotrjura/darwin/config"
)

type CallingPoint struct {
	Tpl string `xml:"tpl,attr"`
	Act string `xml:"act,attr"`
}

type ArrivingPoint struct {
	PlannedArrival string `xml:"pta,attr"`
	ActualArrival  string `xml:"wta,attr"`
}

type DepartingPoint struct {
	PlannedDeparture string `xml:"ptd,attr"`
	ActualDeparture  string `xml:"wtd,attr"`
}

type IntermediatePoint struct {
	*CallingPoint
	*ArrivingPoint
	*DepartingPoint
}

type OriginPoint struct {
	*CallingPoint
	*DepartingPoint
}

type DestinationPoint struct {
	*CallingPoint
	*ArrivingPoint
}

type Journey struct {
	Origin             OriginPoint         `xml:"OR"`
	Destination        DestinationPoint    `xml:"DT"`
	IntermediatePoints []IntermediatePoint `xml:"IP"`
}

type Timetable struct {
	Journeys []Journey `xml:"Journey"`
}

func downloadXML(file string, wg *sync.WaitGroup, c chan []byte, conf config.FtpConfig) {
	defer wg.Done()
	conn := connect(conf)

	ftpFile, err := conn.Retr(file)
	if err != nil {
		panic(err)
	}
	defer ftpFile.Close()

	r, err := gzip.NewReader(ftpFile)
	defer r.Close()

	xmlBytes, _ := ioutil.ReadAll(r)
	c <- xmlBytes
	close(c)
}

func getReferenceFilenames(conf config.FtpConfig) (string, string) {
	conn := connect(conf)
	files, err := conn.List("")
	if err != nil {
		panic(err)
	}

	return files[0].Name, files[1].Name
}

func connect(config config.FtpConfig) *ftp.ServerConn {
	conn, err := ftp.Connect(config.URL)
	if err != nil {
		panic(err)
	}

	err = conn.Login(config.User, config.Password)
	if err != nil {
		panic(err)
	}

	return conn
}
