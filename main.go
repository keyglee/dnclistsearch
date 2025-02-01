package main

import (
	"dnclist/dnc"
	"dnclist/logger"
	"dnclist/output"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	loadNewCard := flag.String("new", "", "Adds a new card from the text file, requires the file path")
	logLevel := flag.String("log", "error", "Log level, can be error, info, or debug")
	prettyPrint := flag.Bool("pretty", false, "Pretty print the output")
	responseDelimiter := flag.String("delimiter", "none", "Delimiter for the response")
	responseSeparator := flag.String("separator", "newline", "Separator for the response")

	flag.Parse()

	logger := logger.NewLogger(*logLevel)

	dncListFile, err := os.Open("dnc.txt")

	if err != nil {
		logger.Error(err)
		panic(err)
	}

	dncClient, err := dnc.NewClient(dncListFile, logger)

	if err != nil {
		panic(err)
	}

	if *loadNewCard != "" {
		err := dncClient.LoadNewList(*loadNewCard)
		if err != nil {
			logger.Error(err)
			panic(err)
		}
	}

	args := flag.Args()

	var response []string

	for _, arg := range args {

		logger.Debugf("Processing argument: '%s'", arg)
		// Format phone number
		formattedArg, err := dncClient.FormatPhoneNumber(arg)

		logger.Debugf("Formatted argument: '%s'", formattedArg)
		if err != nil {
			logger.Error(err)
			panic(err)
		}
		// Argument needs to be in the form of 'xxx,xxxxxxx'
		found, err := dncClient.Search(formattedArg)

		if err != nil {
			panic(err)
		}

		deliminator := output.GetDelimiter(*responseDelimiter)

		if *prettyPrint == false {
			response = append(response, fmt.Sprintf("%s%t%s", deliminator, found, deliminator))
		} else if found {
			response = append(response, fmt.Sprintf("%s%s Found in the file%s", deliminator, arg, deliminator))
		} else {
			response = append(response, fmt.Sprintf("%s%s Not Found in the file%s", deliminator, arg, deliminator))
		}
	}

	fmt.Println(strings.Join(response, string(output.GetSeparator(*responseSeparator))))
}
