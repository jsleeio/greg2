package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

// exit codes intentionally compatible with GNU grep
const (
	exitStatusMatched    = 0
	exitStatusNotMatched = 1
	exitStatusError      = 2
)

// main
func main() {
	config := configureFromFlags()
	matched := false
	haderr := false
	for _, filename := range flag.Args() {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			haderr = true
			continue
		}
		matches := searchFile(config, file, filename)
		if matches > 0 {
			matched = true
		}
		file.Close()
	}
	if flag.NArg() == 0 {
		matched = searchFile(config, os.Stdin, "stdin") > 0
	}
	if haderr {
		os.Exit(exitStatusError)
	}
	if matched {
		os.Exit(exitStatusMatched)
	}
	os.Exit(exitStatusNotMatched)
}

// Config encapsulates configuration from commandline options and possibly other sources
type Config struct {
	// config derived directly from CLI flags
	Before, After *int
	Match         *string
	IgnoreCase    *bool
	Quiet         *bool
	Filenames     *string

	// config derived after CLI flag parsing is complete
	RE             *regexp.Regexp
	PrintFilenames bool
}

// configureFromFlags derives a full configuration from CLI flags
func configureFromFlags() *Config {
	config := &Config{
		Quiet:      flag.Bool("quiet", false, "do not output any matches"),
		IgnoreCase: flag.Bool("ignorecase", false, "perform case-insensitive matching"),
		Filenames:  flag.String("filenames", "auto", "show filenames in output (valid options: no, auto, yes)"),
		Before:     flag.Int("before", 0, "lines of preceding context to print for each match"),
		After:      flag.Int("after", 0, "lines of following context to print for each match"),
		Match:      flag.String("match", "", "RE2 regular expression to match against the input files"),
	}
	flag.Parse()
	if *config.Before < 0 || *config.After < 0 {
		fmt.Fprintln(os.Stderr, "FATAL: before and after values must not be negative")
		os.Exit(1)
	}
	switch *config.Filenames {
	case "auto":
		// like regular grep: only show filenames when more than one file
		config.PrintFilenames = false
		if flag.NArg() > 1 {
			config.PrintFilenames = true
		}
	case "no":
		config.PrintFilenames = false
	case "yes":
		config.PrintFilenames = true
	default:
		fmt.Fprintln(os.Stderr, "FATAL: invalid value for filenames option")
		os.Exit(1)
	}
	if *config.IgnoreCase {
		config.RE = regexp.MustCompile("(?i)" + *config.Match)
	} else {
		config.RE = regexp.MustCompile(*config.Match)
	}
	return config
}

// printLine prints a single line, with or without filename as required
func printLine(config *Config, line string, filename string) {
	if !*config.Quiet {
		if config.PrintFilenames {
			fmt.Print(filename + ":")
		}
		fmt.Println(line)
	}
}

// searchFile performs a search of a single file and outputs to stdout
func searchFile(config *Config, reader io.Reader, filename string) int {
	scanner := bufio.NewScanner(reader)
	var afterWindow int
	tq := NewTextQueue()
	matches := 0
	inAfterWindow := false
	for scanner.Scan() {
		text := scanner.Text()
		// 1. every line moves through the window
		tq.AddFront(text)
		if config.RE.MatchString(text) {
			// 2. we matched, immediately print the backlog lines, including
			//    the current line. And a separator, if this isn't the first
			//    match
			if matches > 0 && (*config.Before > 0 || *config.After > 0) {
				fmt.Println("----------")
			}
			matches++
			for _, line := range tq.StringSlice() {
				printLine(config, line, filename)
			}
			// 3. ... and purge the queue
			tq.Purge()
			afterWindow = *config.After
			inAfterWindow = true
			continue
		}
		if inAfterWindow && afterWindow > 0 {
			// 4. we didn't match, but we did match recently and the user
			//    would like some trailing context
			printLine(config, text, filename)
			afterWindow--
		}
		if inAfterWindow && afterWindow == 0 {
			// 5. we've finished printing the after-window, if any. Purge
			//    the queue to avoid duplicate-printing anything
			tq.Purge()
			inAfterWindow = false
		}
		// 6. if the window is at maximum width and we haven't matched yet,
		//    discard oldest if necessary
		if qlen := tq.Len(); qlen > *config.Before {
			tq.RemoveBack()
		}
	}
	return matches
}
