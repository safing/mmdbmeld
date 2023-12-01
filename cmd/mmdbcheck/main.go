package main

import "log"

func main() {
	if err := getRootCmd().Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
