package main

import (
	//"fmt"
	//"strings"
	"testing"
	//"bytes"
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

func makeSlice0(n int) []int {
	s := make([]int, 0)
	for i:=0; i < n; i++ {
		s = append(s, i)
	}
	return s
}

func makeSliceWithCap(n int) []int {
	s := make([]int, 0, n)
	for i:=0; i < n; i++ {
		s = append(s, i)
	}
	return s
}


const DefaltCap = 100

func BenchmarkMakeSlice0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		makeSlice0(DefaltCap)
	}
}

func BenchmarkMakeSliceWithCap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		makeSliceWithCap(DefaltCap)
	}
}

func main(){

}