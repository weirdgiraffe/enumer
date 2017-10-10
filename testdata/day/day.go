// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Simple test: enumeration of type int starting at 0.

package main

import "fmt"

type Day string

const (
	Monday    Day = "Mon"
	Tuesday   Day = "Tue"
	Wednesday Day = "Wen"
	Thursday  Day = "Thu"
	Friday    Day = "Fri"
	Saturday  Day = "Sat"
	Sunday    Day = "Sun"
)

func main() {
	ck(Monday, "Mon")
	ck(Tuesday, "Tue")
	ck(Wednesday, "Wed")
	ck(Thursday, "Thu")
	ck(Friday, "Fri")
	ck(Saturday, "Sat")
	ck(Sunday, "Sun")
}

func ck(day Day, str string) {
	if fmt.Sprint(day) != str {
		panic("day.go: " + str)
	}
}
