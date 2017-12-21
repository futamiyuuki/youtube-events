package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: ./eventsGen <num_events> <events_output_file> <channels_output_file>")
		os.Exit(2)
	}
	ecnt, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	ef, err := os.OpenFile(os.Args[2], os.O_APPEND|os.O_WRONLY, 0622)
	if err != nil {
		log.Fatal(err)
	}
	defer ef.Close()
	chf, err := os.OpenFile(os.Args[3], os.O_APPEND|os.O_WRONLY, 0622)
	if err != nil {
		log.Fatal(err)
	}
	defer chf.Close()
	fmt.Println("Start generating data for events service...")
	fmt.Printf("Events Output file: %s\n", os.Args[2])
	fmt.Printf("Channels Output file: %s\n\n", os.Args[3])
	st := time.Now()

	genEvents(ecnt, ef, chf)

	fmt.Printf("\nFinished generating events data in %s\n", time.Now().Sub(st).String())
}
