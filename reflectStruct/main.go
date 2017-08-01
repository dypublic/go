// main.go
package main

import (
	"fmt"
	"reflect"
)

func GetTagOfStruct(object interface{}, name string) {
	objectType := reflect.TypeOf(object)
	if objectType.Kind() == reflect.Ptr {
		objectType = objectType.Elem()
	}

	if field, ok := objectType.FieldByName(name); ok {
		tag := field.Tag.Get("json")
		fmt.Println(tag)
	}

}

type Data struct {
	Data string `json:"test_Data"`
}

type test struct {
	test      string   `json:"test_tag"`
	testSlide []string `json:"slide_tag"`
	data      *Data    `json:"pointer"`
}

func main() {
	t := test{
		"123",
		[]string{"234, 234"},
		&Data{"test"},
	}
	GetTagOfStruct(t, "data")
	GetTagOfStruct(t, "testSlide")

}
