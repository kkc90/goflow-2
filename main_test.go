package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"testing"
)

const (
	connString = "localhost:2055"
)

var recordTest = []NetFlow5Record{
	{3232237668, 3232237569, 0, 0, 0, 1, 66, 1048515346, 1048515346, 11065, 53, 0, 0, 17, 0, 0, 0, 0, 0, 0},
	{3232237668, 2539981957, 0, 0, 0, 1, 40, 1048513581, 1048513581, 51673, 443, 0, 16, 6, 0, 0, 0, 0, 0, 0},
	{2899907950, 3232237668, 0, 0, 0, 1, 64, 1048515121, 1048515121, 443, 51674, 0, 16, 6, 0, 0, 0, 0, 0, 0},
	{3232237569, 3232237668, 0, 0, 0, 1, 159, 1048515384, 1048515384, 53, 11065, 0, 0, 17, 0, 0, 0, 0, 0, 0},
	{3232237668, 2899907950, 0, 0, 0, 1, 40, 1048515092, 1048515092, 51674, 443, 0, 16, 6, 0, 0, 0, 0, 0, 0},
	{3232237668, 2726109825, 0, 0, 0, 3, 1713, 1048513409, 1048513412, 50986, 443, 0, 24, 6, 0, 0, 0, 0, 0, 0},
	{2726109825, 3232237668, 0, 0, 0, 3, 413, 1048513409, 1048513451, 443, 50986, 0, 24, 6, 0, 0, 0, 0, 0, 0},
	{2637660846, 3232237668, 0, 0, 0, 1, 49, 1048511795, 1048511795, 40020, 33683, 0, 0, 17, 0, 0, 0, 0, 0, 0},
	{2637687698, 3232237668, 0, 0, 0, 1, 49, 1048511795, 1048511795, 40014, 33683, 0, 0, 17, 0, 0, 0, 0, 0, 0},
	{3232237668, 2637660830, 0, 0, 0, 1, 177, 1048511514, 1048511514, 33683, 40023, 0, 0, 17, 0, 0, 0, 0, 0, 0},
	{3232237668, 2637660840, 0, 0, 0, 1, 174, 1048511513, 1048511513, 33683, 40004, 0, 0, 17, 0, 0, 0, 0, 0, 0},
	{3232237668, 2637660846, 0, 0, 0, 1, 192, 1048511663, 1048511663, 33683, 40020, 0, 0, 17, 0, 0, 0, 0, 0, 0},
	{3232237668, 2637687698, 0, 0, 0, 1, 176, 1048511726, 1048511726, 33683, 40014, 0, 0, 17, 0, 0, 0, 0, 0, 0},
	{1094180650, 3232237668, 0, 0, 0, 1, 910, 1048511660, 1048511660, 40018, 33683, 0, 0, 17, 0, 0, 0, 0, 0, 0},
	{3232237668, 2637641885, 0, 0, 0, 1, 66, 1048511514, 1048511514, 33683, 40008, 0, 0, 17, 0, 0, 0, 0, 0, 0},
	{2637660840, 3232237668, 0, 0, 0, 1, 104, 1048511843, 1048511843, 40004, 33683, 0, 0, 17, 0, 0, 0, 0, 0, 0},
	{3232237668, 1094180650, 0, 0, 0, 1, 66, 1048511514, 1048511514, 33683, 40018, 0, 0, 17, 0, 0, 0, 0, 0, 0},
	{676829863, 3232237668, 0, 0, 0, 1, 78, 1048511761, 1048511761, 33033, 33683, 0, 0, 17, 0, 0, 0, 0, 0, 0},
	{2637641885, 3232237668, 0, 0, 0, 1, 48, 1048511674, 1048511674, 40008, 33683, 0, 0, 17, 0, 0, 0, 0, 0, 0},
	{3232237668, 676829863, 0, 0, 0, 1, 193, 1048511663, 1048511663, 33683, 33033, 0, 0, 17, 0, 0, 0, 0, 0, 0},
	{2637660830, 3232237668, 0, 0, 0, 1, 49, 1048511660, 1048511660, 40023, 33683, 0, 0, 17, 0, 0, 0, 0, 0, 0}}

func TestServer(t *testing.T) {
	c := make(chan string)
	go setupUDPServer(connString, c)
	conn, err := net.Dial("udp4", connString)
	if err != nil {
		t.Error("Could not connect to the server")
	}
	for _, record := range recordTest {
		h := &NetFlow5Header{Version: uint16(5), Count: uint16(1), SamplingInterval: uint16(3)}
		writer := new(bytes.Buffer)
		err = binary.Write(writer, binary.BigEndian, h)
		if err != nil {
			t.Error("Could not encode the header")
		}
		rec := &record
		err = binary.Write(writer, binary.BigEndian, rec)
		if err != nil {
			t.Error("Could not encode the record")
		}
		_, err = conn.Write(writer.Bytes())
		if err != nil {
			t.Error("Could not send the record")
		}
		sent := strings.Trim(fmt.Sprintf("%s", record.String()), "\n")
		got := strings.Trim(<-c, "\n")
		if sent != got {
			t.Errorf("Expected %s to be equal to %s", sent, got)
		}
	}
}
