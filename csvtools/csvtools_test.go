package csvtools

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
	"time"
)


func TestLexer_LexColumns(t *testing.T) {
	filename := "testdata/t.csv"
	headers, records := csvData(filename,',') // load test data
	fi, _ := os.Open(filename)
	format := Format{
		sep:    ",",
	}
	l := NewLexer("lexy 1", format, fi)
	for i,v := range l.Headers() {
		if headers[i] != v {
			t.Errorf("headers don't match with file. Got %s, expected %s", v, headers[i])
		}
	}
	xcol := "x"
	ycols := []string{"z", "y"}
	allCols := append(ycols,xcol)
	go l.LexColumns(allCols)
	time.Sleep(time.Nanosecond)
	for line, record := range records {
		for i, c := range l.values {
			num, _ := strconv.ParseFloat(record[i],64)
			select {
			case val := <- c:
				if val.Number() != num {
					t.Errorf("mismatch on line %d. expected %.2f. got %.2f",line, num, val.Number() )
				}
			default:
				continue
			}
		}
	}
	l.Close()
}



// Load all contents into memory separating
// header from contents. header is treated as first
// line read by csv.Read(). Trims leading space.
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
		records = append(records,record)
	}
	return headers, records
}

//func trimSlice(sli []string) []string {
//	for i,s := range sli {
//		sli[i] = strings.TrimSpace(s)
//	}
//	return sli
//}