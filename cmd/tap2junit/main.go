// Package main represents the program tap2junit, a utility that converts
// a testanything.org's TAP test format into junit format.
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/filmil/tap2junit/pkg/junit"
	"github.com/filmil/tap2junit/pkg/tap"
	"github.com/filmil/tap2junit/pkg/tap/tojunit"
	"github.com/golang/glog"
)

var (
	testName        = flag.String("test_name", "unnamed_test", "Sets the test name to use")
	outputFile      = flag.String("output_file", "", "Sets the output file to use")
	reorderDuration = flag.Bool("reorder_duration", false, "If set, will reorder durations to work around https://github.com/bats-core/bats-core/issues/187")
	reorderAll      = flag.Bool("reorder_all", false, "If set, will reorder all test lines to work around https://github.com/bats-core/bats-core/issues/187")
	singleSuite     = flag.Bool("single_suite", false, "If set, will output only the <testsuite> as top-level tag; not <testsuites>")
)

func run(r io.Reader, w io.Writer, opts tap.ReadOpt, singleSuite bool) error {
	s := ioutil.Discard
	if *outputFile != "" {
		s = w
		nw, err := os.Create(*outputFile)
		if err != nil {
			panic(err)
		}
		w = nw
	}
	t, err := tap.Read(r, s, opts)
	if err != nil {
		return fmt.Errorf("while reading TAP: %v", err)
	}
	j, err := tojunit.FromTAP(t)
	if err != nil {
		return fmt.Errorf("while converting to jUnit: %v", err)
	}
	if err := junit.Write(j, w, singleSuite); err != nil {
		return fmt.Errorf("while writing jUnit: %v", err)
	}
	return nil
}

func init() {
	flag.Parse()
}

func main() {
	opts := tap.ReadOpt{
		Name:            *testName,
		ReorderDuration: *reorderDuration,
		ReorderAll:      *reorderAll,
	}
	if err := run(os.Stdin, os.Stdout, opts, *singleSuite); err != nil {
		glog.Fatalf("unexpected error: %v", err)
	}
}
