/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/soypat/decimate/csvtools"
	"github.com/spf13/cobra"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const bufferSize = 3

// flags
var tolerance float64 = 0.1 // default for tests
var xFlag, yFlag, inputSeparator, outputName, outputExtension, floatFormat string
var interp, enforceComma, silent bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "decimate",
	Short: "Reduce # of points of signals and curves.",
	Long: `Decimator is a tool for downsampling
(also known as decimating) numerical data.
Generates decimated files from a token separated
file for use in plotting tools.

Example:

decimate -x time -y "x,y,z"

The code above operates on the time x-column and 3 y-columns
named 'x', 'y' and 'z'.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := checkParameters(args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(args); err != nil {
			fmt.Printf("[ERR] %s", err)
			os.Exit(1)
		}
	},
}

type job struct {
	name         string
	xname, yname string
	tolerance    float64
	stepper
}

func run(args []string) error {
	fi, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("error opening file")
	}
	inputFormat := csvtools.MakeFormat(nil, inputSeparator)
	lex := csvtools.NewLexer("mainLex", inputFormat, fi)

	yColsSplit := splitYcols(yFlag)
	if yColsSplit[0] == "*" && len(yColsSplit) == 1 {
		yColsSplit = []string{}
		for _, h := range lex.Headers() {
			if h != xFlag {
				yColsSplit = append(yColsSplit,h)
			}
		}
	}

	allColumns := append(yColsSplit, xFlag) // x column is last one
	for _, v := range allColumns {
		if !sliceContains(lex.Headers(), v) {
			return fmt.Errorf("%s is not in columns:\n%v", v, lex.Headers())
		}
	}
	valueChans := lex.ValueChannels()
	idxs, err := lex.HeaderIdxs(allColumns)
	if err != nil {
		return fmt.Errorf("Could not find column %s among headers:\n%v", xFlag, lex.Headers())
	}
	var xMux Mux
	xMux.In = valueChans[idxs[len(idxs)-1]] // we will mux x input
	xMux.makeOutputs(len(yColsSplit))
	go lex.LexColumns(allColumns)
	time.Sleep(1)
	go xMux.run()
	var wg sync.WaitGroup
	wg.Add(len(yColsSplit))
	for i, colname := range yColsSplit {
		var a stepper
		if interp {
			a = interpStepper{}
		} else {
			a = inPlaceStepper{}
		}
		var j = job{
			name:      outputName + "-" + colname + "." + outputExtension,
			xname:     xFlag,
			yname:     colname,
			tolerance: tolerance,
			stepper:   a,
		}
		downsample(j, xMux.Out[i], valueChans[idxs[i]], &wg)
	}
	wg.Wait()
	alert("decimate finished")
	return nil
}

func downsample(j job, cx, cy *chan csvtools.Value, wg *sync.WaitGroup) {
	defer wg.Done()
	alert("creating file %s", j.name)
	var xV, yV csvtools.Value
	fo, err := os.Create(j.name)
	if err != nil {
		panic(err)
	}
	defer fo.Close()
	w := csv.NewWriter(fo)
	defer w.Flush()
	if !enforceComma {
		w.Comma = rune(inputSeparator[0])
	}
	err = w.Write([]string{j.xname,j.yname})
	if err != nil {
		panic(err)
	}
	// algorithm asks we print out last value
	jj := 0
	for {
		if len(*cx) == 0 || len(*cy) == 0 {
			time.Sleep(1)
			continue
		}
		jj++
		xV, yV = <-*cx, <-*cy
		if yV.Type() == csvtools.TypeEOF || xV.Type() == csvtools.TypeEOF {
			j.stepper = j.step(math.NaN(),math.NaN()) // false step
			if err := w.Write(j.values(floatFormat)); err != nil {
				panic(err)
			}
			alert("finished file %s", j.name)
			return
		}
		x, y := xV.Number(), yV.Number()
		j.stepper = j.step(x, y)
		if j.ready() {
			if err := w.Write(j.values(floatFormat)); err != nil {
				panic(err)
			}
		}
		w.Flush()
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// actually modifies flag values!
func checkParameters(args []string) error {
	// filename
	if len(args) != 1 {
		return errors.New("requires exactly one argument as input filename")
	}
	// y columns
	ycols := splitYcols(yFlag)
	if len(ycols) < 1 || yFlag == "" {
		return errors.New("found no y-column flag value")
	}
	// Delimiters
	if strings.TrimSuffix(inputSeparator, "s") == "tab" || inputSeparator == "\\t" {
		yFlag = strings.ReplaceAll(yFlag, "\\t", "\t")
		inputSeparator = "\t"
	}
	if len(inputSeparator)> 1 {
		return errors.New("delimiter should be one character. '\\t' and 'tab' work as an option")
	}
	commaReplacement := ";"
	if inputSeparator==commaReplacement {
		commaReplacement = "."
	}
	if enforceComma && inputSeparator != ","{
		yFlag = strings.ReplaceAll(yFlag, ",",commaReplacement)
	}
	// output name forming
	var iname string
	outputName, outputExtension = splitFileExtension(outputName)
	if outputExtension == "" || outputExtension == "<inputExtension>" {
		iname, outputExtension = splitFileExtension(discardPath(args[0]))
		if outputName == "<inputName>-<ycol>" || outputName == "" {
			outputName = iname
		}
	}
	// formatter
	const floatNum = .125
	if num, err := strconv.ParseFloat(fmt.Sprintf(floatFormat, floatNum),64); num != floatNum || err != nil {
		return errors.New("formatting option yielded error. example of usage: \n'%0.2f' for two decimal placed\n'%e' for scientific notation")
	}
	return nil
}
func init() {
 	rootCmd.Flags().StringVarP(&outputName, "output", "o", "", "Output filename. Will have name of y column and input's extension unless otherwise specified.")
	rootCmd.Flags().StringVarP(&floatFormat, "fformat", "f", "%e", "Floating point format")
	rootCmd.Flags().BoolVarP(&enforceComma, "comma", "c", false, "Force output to use comma as delimiter")
	rootCmd.Flags().StringVarP(&inputSeparator, "delimiter", "d", ",", "Delimiter token. Examples: '-d \\t' or '-d=\";\"'")
	rootCmd.Flags().Float64VarP(&tolerance, "tolerance", "t", 0.1, "Downsampling y-value tolerance.")
	rootCmd.Flags().StringVarP(&yFlag, "ycols", "y", "", "Y column names. Pass '*' to process all columns (required). Separate with delimiter token.")
	_ = rootCmd.MarkFlagRequired("ycols")
	rootCmd.Flags().StringVarP(&xFlag, "xcol", "x", "", "X column name (required)")
	_ = rootCmd.MarkFlagRequired("xcol")
	rootCmd.Flags().BoolVarP(&interp, "interp", "i", false, "Use interpolator algorithm. Downsampling is more aggressive at the cost of changing point y values")
	rootCmd.Flags().BoolVarP(&silent, "silent", "s", false, "Silent execution (no printing).")

}

func sliceContains(sli []string, s string) bool {
	for _, v := range sli {
		if v == s {
			return true
		}
	}
	return false
}

func splitYcols(y string) []string {
	s := strings.Split(yFlag, inputSeparator)
	return s
}

func discardPath(fname string) string {
	pathIndex := strings.LastIndex(fname, "/")
	if pathIndex != -1 && pathIndex < len(fname)-1 {
		fname = fname[1+pathIndex:]
	}
	return fname
}

func splitFileExtension(fname string) (string, string) {
	fileTypeIndex := strings.LastIndex(fname, ".")
	if fileTypeIndex == -1 {
		return fname, ""
	}
	return fname[:fileTypeIndex], fname[fileTypeIndex+1:]
}

func alert(format string, args ...interface{}) {
	if !silent {
		msg := fmt.Sprintf(format, args...)
		if args == nil {
			msg = fmt.Sprintf(format)
		}
		fmt.Print("[INFO] ",msg,"\n")
	}
}