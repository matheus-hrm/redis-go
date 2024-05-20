package main

import (
	"errors"
	"strconv"
)

type RespDataType byte

type Command string

type Request struct {
	Command Command
	Args    []string
}

type RespMessage struct {
	Type RespDataType
}

func (datatype RespDataType) IsAggregate() bool {
	return datatype == Array || datatype == BulkString
}

const (
	SimpleString RespDataType = '+'
	Error        RespDataType = '-'
	Integer      RespDataType = ':'
	BulkString   RespDataType = '$'
	Array        RespDataType = '*'
	Status       RespDataType = '='
	Nulls        RespDataType = '_'
	Booleans     RespDataType = '#'
	Doubles      RespDataType = ','
	BulkErrors   RespDataType = '!'
	Maps         RespDataType = '%'
	Sets         RespDataType = '~'
)

const (
	ECHO Command = "ECHO"
	PING Command = "PING"
)

func ParseCommand(b []byte) (Request, error) {
	dataType, n, err := GetCommandType(&b)
	if err != nil {
		return Request{}, err
	}
	switch dataType {
		case SimpleString:
			return Request{
					Command: Command(parseString(&b)),
				}, nil
		case Array:
			cmd := ""
			args := make([]string, 0)
			for i := 0; i < n; i++ {
				dataType, subCount, err := GetCommandType(&b)
				if err != nil {
					return Request{}, err
				}
				if dataType == BulkString {
					val := parseString(&b)
					if len(val) != subCount {
						return Request{}, errors.New("invalid format, incorrect number of arguments")
					}
					if cmd == "" {
						cmd = val
						continue
					}
					args = append(args, val)
				} else {
					return Request{}, errors.New("invalid command")
				}
			}
			return Request{
				Command: Command(cmd),
				Args:    args,
			}, nil
	}	
	return Request{}, errors.New("invalid command")
}

func GetCommandType(b *[]byte) (dataType RespDataType, count int, err error) {
	dataType = RespDataType((*b)[0])
	count = 0
	if dataType.IsAggregate() {
		count, err = strconv.Atoi(string((*b)[1]))
	}
	*b = (*b)[count:]

	return dataType, count, err
}

func parseString(b *[]byte) string {
	read := make([]byte, 0)
	for _, c := range *b {
		if c == '\r' {
			break
		}
		read = append(read, c)
	}
	val := string(read)
	*b = (*b)[len(val)+2:]
	return val
}