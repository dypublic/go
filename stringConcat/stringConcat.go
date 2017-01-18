package main

import (
	"strings"
	"fmt"
	"bytes"
	"runtime/pprof"
	"os"
	"log"
)
/*
go test -bench=. stringConcat_test.go
*/
const SOAP_HEAD = "<SOAP-ENV:Envelope xmlns:SOAP-ENV=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:SOAP-ENC=\"http://schemas.xmlsoap.org/soap/encoding/\" xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns=\"http://www.derbysoft.com/doorway\"><SOAP-ENV:Body>"
const SOAP_FOOT = "</SOAP-ENV:Body></SOAP-ENV:Envelope>"

const xml = "<Hello>world</Hello>"

//func stringAdd(xml string) string {
//
//	return SOAP_HEAD + xml + SOAP_FOOT + SOAP_HEAD + xml + SOAP_FOOT + SOAP_HEAD + xml + SOAP_FOOT
//}

func stringAdd(xml string) string {
	build := SOAP_HEAD
	build += xml
	build += SOAP_FOOT
	build += SOAP_HEAD
	build += xml
	build += SOAP_FOOT
	build += SOAP_HEAD
	build += xml
	build += SOAP_FOOT
	return build
}

func stringJoin(xml string) string {

	return strings.Join([]string{SOAP_HEAD, xml, SOAP_FOOT},"")
}

func stringFmt(xml string) string {

	return fmt.Sprintf("%s%s%s",SOAP_HEAD , xml , SOAP_FOOT)
}

func stringBuffer(xml string) string {
	buf := make([]byte, 0, 3000)
	buf = append(buf, SOAP_HEAD...)
	buf = append(buf, xml...)
	buf = append(buf, SOAP_FOOT...)
	buf = append(buf, SOAP_HEAD...)
	buf = append(buf, xml...)
	buf = append(buf, SOAP_FOOT...)
	buf = append(buf, SOAP_HEAD...)
	buf = append(buf, xml...)
	buf = append(buf, SOAP_FOOT...)
	return string(buf)
}

func stringBufferBuf(xml string) string {
	buf:=new(bytes.Buffer)

	buf.WriteString(SOAP_HEAD)
	buf.WriteString(xml)
	buf.WriteString(SOAP_FOOT)
	buf.WriteString(SOAP_HEAD)
	buf.WriteString(xml)
	buf.WriteString(SOAP_FOOT)
	buf.WriteString(SOAP_HEAD)
	buf.WriteString(xml)
	buf.WriteString(SOAP_FOOT)
	return buf.String()
}

func main(){
	f, err := os.Create("cpuprofile")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	var str = ""
	for i:=0; i < 1000000; i++ {

		str = stringAdd(xml)
		//println(str)
		str = stringJoin(xml)
		//println(str)
		str = stringFmt(xml)
		//println(str)
		str = stringBuffer(xml)
		//println(str)
	}
	println(str)
}

