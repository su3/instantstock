package main

import "instantstock/ctrl"

func main() {
	codeList := []string{"600000", "601398", "000972"}
	f := ctrl.Fetch{}
	f.FetchCodes(codeList)
}
