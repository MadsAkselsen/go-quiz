package main

import "flag"

func main() {
	csvFilename := flag.String("Csv", "problems.csv", "a csv file in the format of 'question,answer'")
	flag.Parse()
	_ = csvFilename
}