package main

import (
	"testing"
)

func TestParsePing(t *testing.T) {

	raw := "+PING\r\n"
	request, err :=  ParseCommand([]byte(raw))

	if err != nil {
		t.Errorf("Unexpected error has occured %s", err.Error())
		t.FailNow()
	}
	if request.Command != PING {
		t.Errorf("Expected PING, got %s", request.Command)
		t.FailNow()
	}
}

func TestParseEcho(t *testing.T) {
	
	raw := "*2\r\n$4\r\nECHO\r\n$7\r\ntesting\r\n"
	request, err := ParseCommand([]byte(raw))

	if err != nil {
		t.Errorf("Unexpected error has occured %s", err.Error())
		t.FailNow()
	}
	if request.Command != ECHO {
		t.Errorf("Expected ECHO, got %s", request.Command)
		t.Fail()
	}
	if len(request.Args) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(request.Args))
		t.FailNow()
	}
	if request.Args[0] != "testing" {
		t.Errorf("Expected ECHO argument testing, got %s", request.Args[0])
		t.FailNow()
	}
}