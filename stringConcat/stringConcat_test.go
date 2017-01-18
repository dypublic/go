package main

import (
	"testing"
)
/*
go test -bench=. stringConcat_test.go
*/


func BenchmarkStringAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stringAdd(xml)
	}
}

func BenchmarkStringJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stringJoin(xml)
	}
}
func BenchmarkStringFmt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stringFmt(xml)
	}
}
func BenchmarkStringBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stringBuffer(xml)
	}
}

func BenchmarkStringBufferBuf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stringBufferBuf(xml)
	}
}
