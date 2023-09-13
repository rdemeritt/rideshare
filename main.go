package main

import (
	"flag"
)

func init() {
	initLog()
	parseCommandLineFlags()
}

func parseCommandLineFlags() {
	// Define command line flags
	logLevel := flag.String("log-level", "info", "Logging level (debug, info, warn, error)")
	startCoord := flag.String("start", "37.7749,-122.4194", "Starting coordinates")
	endCoord := flag.String("end", "37.3352,-121.8811", "Ending coordinates")

	// Parse command line flags
	flag.Parse()

	// Set the logging level
	setLogLevel(*logLevel)
	// Set the start and end coordinates
	setStart(*startCoord)
	setEnd(*endCoord)
}

func main() {
    // Get the driving directions
    directions := getDirections(SourceCoordinates, DestinationCoordinates)

	printDirections(directions)
}
