package cmd

import (
	"encoding/csv"
	"github.com/soypat/decimate/csvtools"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

const runFilename = "../testdata/ik.tsv"
func TestRun(t *testing.T) {
	yFlag = "*"
	xFlag = "time"
	inputSeparator = "\t"
	go run([]string{runFilename})
	time.Sleep(time.Second*20)
	//if err:= run([]string{runFilename}); err != nil {
	//	t.Errorf("run returned error: %s",err)
	//}
}

const testDataFolder = "../testdata/"
const testDataFilename = "ik.tsv"
const yColumnName = "V(c,n)"
const xColumnName = "time"
//func TestAlgorithm(t *testing.T) {
//	var a inPlaceStepper
//	headers, data := csvData(testDataFolder+testDataFilename, '\t')
//	_, expecteds := csvData(testDataFolder+"ds-"+yColumnName+".tsv", '\t')
//	headerIdx := make([]int, 2)
//	for i, v := range headers {
//		if v == xColumnName {
//			headerIdx[0] = i
//		}
//		if v == yColumnName {
//			headerIdx[1] = i
//		}
//	}
//	// generation of algorithm results
//	var results [][2]float64
//	var steps int
//	for _, w := range data {
//		x, _ := strconv.ParseFloat(w[headerIdx[0]], 64)
//		y, _ := strconv.ParseFloat(w[headerIdx[1]], 64)
//		if stepped := a.step(x, y); stepped {
//			results = append(results, [2]float64{x, y})
//			steps++
//		}
//	}
//	// comparison of results with expected values
//	for i, u := range expecteds {
//		xExpected, _ := strconv.ParseFloat(u[0], 64)
//		yExpected, _ := strconv.ParseFloat(u[1], 64)
//		if xExpected != results[i][0] {
//			t.Errorf("expected x value %f. got %f", xExpected, results[i][0])
//		}
//		if yExpected != results[i][1] {
//			t.Errorf("expected y value %f. got %f", xExpected, results[i][0])
//		}
//	}
//}

const bufSize = 3
const numRecievers = 6
const numTests = 100

// TestMux tests Mux struct by taking a source
// channel (src) and multiplexing it's output
// to multiple receiver (rcvr) channels
func TestMux(t *testing.T) {
	passed := 0
	src := make(chan csvtools.Value, bufSize)
	rcvr := make([]*chan csvtools.Value, numRecievers)
	var wg sync.WaitGroup
	for i, _ := range rcvr {
		c := make(chan csvtools.Value, bufSize)
		rcvr[i] = &c
	}
	m := Mux{
		In:   &src,
		Out:  rcvr,
	}
	data := linearSlice(1, numTests)
	go func(d []int) {
		for _, v := range d {
			src <- csvtools.NewValue(strconv.Itoa(v), 1)
		}
	}(data)
	go m.run() // start muxing
	// reading the mux outputs
	wg.Add(1)
	go func(d []int) {
		for _, v := range d {
			expected := strconv.Itoa(v)
			for ic, c := range m.Out {
				select {
				case val := <-*c:
					if expected != val.String() {
						t.Errorf("expected %s on mux out #%d, got %s", expected, ic, val.String())
					} else {
						passed++
					}
				}
			}
		}
		wg.Done()
	}(data)
	wg.Wait()
	if passed != numTests*numRecievers {
		t.Errorf("Not all subtests passed. %d/%d failed", numTests-passed, numTests)
	}
	return
}

func linearSlice(intStart, intEnd int) []int {
	if r := intEnd - intStart; r < 0 {
		return nil
	}
	I := make([]int, intEnd-intStart+1)
	for i := intStart; i <= intEnd; i++ {
		I[i-intStart] = i
	}
	return I
}

func csvData(filename string, comma rune) ([]string, [][]string) {
	fi, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	reader := csv.NewReader(fi)
	reader.TrimLeadingSpace = true
	reader.Comma = comma
	headers, _ := reader.Read()
	var records [][]string
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		records = append(records, record)
	}
	return headers, records
}
