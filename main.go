package main

import (
	"flag"
	"fmt"
	"git.lenzplace.org/lenzj/pgot/lib"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var Version string

func main() {
	// Get application name as executed from command prompt
	appName := filepath.Base(os.Args[0])

	// Set up formatting for error messages
	log.SetFlags(0)
	log.SetPrefix(appName + ": ")

	// Parse command line
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s [OPTION]... [FILE]...\n"+
				"Process a got (golang template) file and send transformed text to output.\n"+
				"Options:\n", appName)
		flag.PrintDefaults()
	}
	var oflag = flag.String("o", "-", "output `file` path")
	var iflag = flag.String("i", "", "colon separated list of `paths` to search with pgotInclude")
	var dflag = flag.String("d", "", "string of json frontmatter to include")
	var vflag = flag.Bool("v", false, "display "+appName+" version")

	flag.Parse()

	// Display application version if requested
	if *vflag {
		fmt.Println(appName + " " + Version)
		os.Exit(0)
	}

	// Prepare input and output streams
	var (
		fd     string
		input  io.Reader
		output io.Writer
		err    error
	)

	switch flag.NArg() {
	case 0:
		input = os.Stdin
		if fd, err = os.Getwd(); err != nil {
			log.Fatalln(err)
		}
	default:
		// The last argument on the command line is the file to
		// process and send to output.
		fp := flag.Arg(flag.NArg() - 1)
		if input, err = os.Open(fp); err != nil {
			log.Fatalln(err)
		}
		if filepath.IsAbs(fp) {
			fd = filepath.Dir(fp)
		} else {
			pwd, err := os.Getwd()
			if err != nil {
				log.Fatalln(err)
			}
			fd = filepath.Dir(filepath.Join(pwd, fp))
		}
	}

	switch *oflag {
	case "-":
		output = os.Stdout
	default:
		if output, err = os.Create(*oflag); err != nil {
			log.Fatalln(err)
		}
	}

	// Create pgot Parser
	gInclude := strings.Split(*iflag, ":")
	c, err := pgot.NewParser(input, fd, gInclude)
	if err != nil {
		log.Fatalln(err)
	}

	// Include files specified in frontmatter (if any)
	if err := c.ProcessFMInclude(); err != nil {
		log.Fatalln(err)
	}

	// Include files specified on command line
	// Loop through arg0 to argN-2 using c.IncludeFile and filename as the
	// namespace
	for n := 0; n < flag.NArg()-1; n++ {
		base := filepath.Base(flag.Arg(n))
		dir := filepath.Dir(flag.Arg(n))
		name := strings.TrimSuffix(base, filepath.Ext(base))
		if err := c.IncludeFile(flag.Arg(n), dir, name); err != nil {
			log.Fatalln(err)
		}
	}

	// Process any frontmatter specified on command line
	if err := c.ParseFMString(*dflag, fd, "."); err != nil {
		log.Fatalln(err)
	}

	// Execute template
	if err := c.Execute(output); err != nil {
		log.Fatalln(err)
	}
}
