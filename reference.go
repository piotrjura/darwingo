package main

import (
	"compress/gzip"
	"encoding/xml"
	"io/ioutil"
	"sync"

	"fmt"

	"github.com/fatih/color"
	"github.com/jlaffaye/ftp"
	"github.com/piotrjura/darwingo/config"
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

type Location struct {
	Tpl  string `xml:"tpl,attr"`
	Name string `xml:"locname,attr"`
	CRS  string `xml:"crs,attr"`
	TOC  string `xml:"toc,attr"`
}

type Company struct {
	TOC  string `xml:"toc,attr"`
	Name string `xml:"tocname,attr"`
	URL  string `xml:"url,attr"`
}

type LateReason struct {
	Code int    `xml:"code,attr"`
	Text string `xml:"reasontext,attr"`
}

type TimetableReference struct {
	Locations   []Location   `xml:"LocationRef"`
	Companies   []Company    `xml:"Company"`
	LateReasons []LateReason `xml:"LateRunningReasons>Reason"`
}

func downloadXML(file string, wg *sync.WaitGroup, c chan []byte, conf config.FtpConfig) {
	color.Blue("Downloading XML: %s...\n", file)
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
	color.Green("Downloading %s completed\n", file)
	c <- xmlBytes
	close(c)
}

func getReferenceFilenames(conf config.FtpConfig) (string, string) {
	color.Blue("Fetching timetable & reference file names...")
	conn := connect(conf)
	files, err := conn.List("")
	if err != nil {
		panic(err)
	}

	color.Green("Timetable & reference file names fetched")
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

func parseTimetables(x chan []byte, wg *sync.WaitGroup) Timetable {
	defer wg.Done()
	d := <-x
	color.Blue("Parsing timetable data...")
	var timetable Timetable
	err := xml.Unmarshal(d, &timetable)
	if err != nil {
		panic(err)
	}

	color.Green("Timetable data parsed")
	return timetable
}

func parseReference(x chan []byte, wg *sync.WaitGroup) TimetableReference {
	defer wg.Done()

	d := <-x
	color.Blue("Parsing reference data...")
	var ref TimetableReference
	err := xml.Unmarshal(d, &ref)
	if err != nil {
		panic(err)
	}

	for _, reason := range ref.LateReasons {
		fmt.Println(reason.Text)
	}

	color.Green("Reference data parsed")
	return ref
}
