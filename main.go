package main

import (
	"fmt"
	"github.com/fungustt/generator/config"
	"github.com/fungustt/generator/generator"
	o "github.com/fungustt/generator/os"
	"log"
	"time"
)

func main() {
	start := time.Now()
	c, err := config.Get()
	if err != nil {
		log.Fatalf("Error on config init, %s", err)
	}

	os := o.NewWrapper()
	g, err := generator.NewGenerator(c, os)
	if err != nil {
		log.Fatalf("Error on generator init, %s", err)
		return
	}

	if err := g.Generate(); err != nil {
		log.Fatalf("Error on generation, %s", err)
		return
	}

	stop := time.Now()
	fmt.Printf("Successfully generated at \"%s\" in %d seconds \n", c.CsvDir, stop.Sub(start)/time.Second)
}
