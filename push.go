package main

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	"github.com/go-stomp/stomp"
)

type PushPort struct {
	TrainSchedule string       `xml:"ts,attr"`
	Version       string       `xml:"version,attr"`
	UpdateOrigin  UpdateOrigin `xml:"uR"`
}

type UpdateOrigin struct {
	UpdateOrigin string        `xml:"updateOrigin,attr"`
	Ts           TrainSchedule `xml:"TS"`
}

type TrainSchedule struct {
	RID string `xml:"rid,attr"`
	SSD string `xml:"ssd,attr"`
	UID string `xml:"uid,attr"`
}

const queryTimetable string = `
	<?xml version="1.0"?>
	<Pport ts="%s" version="11" xmlns="http://www.thalesgroup.com/rtti/PushPort/v11">
		<QueryTimetable />
	</Pport>
`

func listen() {
	conn, err := stomp.Dial("tcp", "datafeeds.nationalrail.co.uk:61613",
		stomp.ConnOpt.Login("d3user", "d3password"))
	if err != nil {
		panic(err)
	}

	sub, err := conn.Subscribe("D399512165-45ca-4082-b02d-3a7796936d3b", stomp.AckClient)
	if err != nil {
		panic(err)
	}

	// q := []byte(fmt.Sprintf(queryTimetable, time.Now().Format(time.RFC3339)))
	// err = conn.Send("", "application/xml", q)
	// if err != nil {
	// 	panic(err)
	// }

	var ok = true
	var msg *stomp.Message

	for ok {
		msg, ok = <-sub.C
		if msg.Err != nil {
			panic(msg.Err)
		}

		err = conn.Ack(msg)
		if err != nil {
			panic(err)
		}

		b := bytes.NewReader(msg.Body)
		r, err := gzip.NewReader(b)
		if err != nil {
			panic(err)
		}

		defer r.Close()
		var xmlBytes []byte
		xmlBytes, err = ioutil.ReadAll(r)
		if err != nil {
			panic(err)
		}

		var pport PushPort
		err = xml.Unmarshal(xmlBytes, &pport)
		if err != nil {
			panic(err)
		}

		fmt.Println(pport)
	}

	err = sub.Unsubscribe()
	if err != nil {
		panic(err)
	}

	err = conn.Disconnect()
	if err != nil {
		panic(err)
	}

	return
}
