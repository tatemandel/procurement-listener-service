package main

import (
	"flag"
	"log"
	"procurementlistenerservice/inmemory"
	"procurementlistenerservice/server"
)

// Options contains the options for the service.
type Options struct {
	Port         int
	MetadataFile string
}

var options Options

func init() {
	flag.IntVar(&options.Port, "port", 11000, "use '--port' option to specify the port for service to listen on")
	flag.StringVar(&options.MetadataFile, "metadataFile", "metadata.json", "use '--metadataFile'"+
		"option to specify the metadata file that contains service definitions")
	flag.Parse()
}

func main() {

	metadata, err := inmemory.ReadMetadataFile(options.MetadataFile)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Loaded metadata:")
	log.Printf("%+v\n", metadata)

	s, err := server.CreateServer(options.Port, inmemory.CreateService(metadata))
	if err != nil {
		log.Fatalf("Error creating server: '%v'\n", err)
	}
	s.Start()
}
