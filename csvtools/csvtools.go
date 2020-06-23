package csvtools

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

const bufferSize = 1

type Value struct {
	str string
	num float64
	typ ValueType
	col int
}

func NewValue(str string, col int) Value {
	return Value{
		str: str,
		num: 0,
		typ: TypeString,
		col: col,
	}
}

type ValueType int

const (
	TypeNil    ValueType = 0
	TypeNumber           = iota
	TypeString
	TypeEOF
	TypeError
)

type lexer struct {
	name   string // for error formatting
	reader *csv.Reader
	file   *os.File
	values []*chan Value // one for each column
	sent   int
	Format
}

type Format struct {
	colstr []string
	typ    []ValueType
	sep    string
}

func NewLexer(name string, format Format, file *os.File) *lexer {
	//r := csv.NewReader(file)
	//r.Comma = rune(format.sep[0])
	//r.TrimLeadingSpace = true
	l := lexer{
		file: file,
		name: name,
		//reader: r,
		Format: format,
		values: make([]*chan Value, bufferSize),
	}
	l.Init()
	return &l
}

func (l *lexer) LexColumns(colstr []string) {
	if hasRepeatedStrings(colstr) {
		l.errorfall("found repeated headers! Cannot continue")
		return
	}
	cols, err := l.HeaderIdxs(colstr)
	if err != nil {
		l.errorfall("%s", err)
		return
	}
	rdr := newCsvReader(l)
	_,_ = rdr.Read() // read header
	numberColumns := l.getNumericalColumns()
	var lexedNumberColumns, lexedStringColumns []int
	for _, c := range numberColumns {
		if containsInt(cols, c) {
			lexedNumberColumns = append(lexedNumberColumns, c)
		}
	}
	stringColumns := l.getStringColumns()
	for _, c := range stringColumns {
		if containsInt(cols, c) {
			lexedStringColumns = append(lexedStringColumns, c)
		}
	}
	for {
		record, errCsv := rdr.Read()
		if errCsv != nil {
			if errCsv.Error() == "EOF" {
				l.sent += len(cols)
				l.sendOnAll(Value{typ: TypeEOF})
			} else {
				l.sent += len(cols)
				l.sendOnAll(Value{typ: TypeError})
			}
			break
		}
		for _, c := range lexedNumberColumns {
			num, err := strconv.ParseFloat(record[c], 64)
			if err != nil {
				l.errorf(c, "unexpected non numerical entry: %s", err)
				return // TODO maybe change this to continue for robustness
			}
			l.sent++
			*l.values[c] <- Value{
				//str: record[c],
				num: num,
				typ: TypeNumber,
				col: c,
			}
		}
		for _, c := range lexedStringColumns {
			l.sent++
			*l.values[c] <- Value{
				str: record[c],
				typ: TypeString,
				col: c,
			}
		}
	}
	l.close()
}

func MakeFormat(colstrings []string, sep string) Format {
	var f Format
	f.colstr = colstrings
	f.sep = sep
	return f
}

func (l *lexer) ValueChannels() []*chan Value {
	return l.values
}

func (l *lexer) close() {
	for _, v := range l.values {
		close(*v)
	}
}

func (l *lexer) errorf(ch int, format string, args ...interface{}) {
	l.sent++
	*l.values[ch] <- Value{
		str: fmt.Sprintf("%s"+format, l.name, args),
		num: 0,
		typ: TypeError,
		col: -1,
	}
	fmt.Printf(format+"\n", args)
}

func (v *Value) Type() ValueType {
	return v.typ
}

func (v *Value) String() string {
	return v.str
}

func (v *Value) Number() float64 {
	return v.num
}
func (v *Value) Column() int {
	return v.col
}

func containsInt(sli []int, i int) bool {
	for _, v := range sli {
		if i == v {
			return true
		}
	}
	return false
}

func containsString(sli []string, i string) bool {
	for _, v := range sli {
		if i == v {
			return true
		}
	}
	return false
}

func hasRepeatedStrings(sli []string) bool {
	for i := 0; i < len(sli)-2; i++ {
		if containsString(sli[i+1:], sli[i]) {
			return true
		}
	}
	return false
}

func (l *lexer) Init() {
	colstr, reader := l.readHeader()
	l.colstr = colstr
	record, err := reader.Read() // reads headers
	l.values = make([]*chan Value, len(l.colstr))
	if err != nil {
		panic("Error initializing lexer") //l.errorf(,"error reading numerical records")
	}
	for i, v := range record {
		if isNumerical(v) {
			l.typ = append(l.typ, TypeNumber)
		} else {
			l.typ = append(l.typ, TypeString)
		}
		c := make(chan Value, bufferSize)
		l.values[i] = &c
	}
}

func (l *lexer) Headers() []string {
	return l.colstr
}

func (l *lexer) readHeader() ([]string, *csv.Reader) {
	rdr := newCsvReader(l) // csv.NewReader(l.file)
	record, err := rdr.Read()
	if err != nil {
		panic("I just wanted to read the header")
	}
	return record, rdr
}

func (l *lexer) HeaderIdxs(colstr []string) (idx []int, err error) {
	for _, v := range colstr {
		for j, w := range l.colstr {
			if v == w {
				idx = append(idx, j)
				break
			} else if j == len(l.colstr)-1 {
				return nil, fmt.Errorf("Could not find '%s' among headers", v)
			}
		}
	}
	return idx, nil
}

func (l *lexer) lexHeader() *csv.Reader {
	headerReader := *l.reader
	record, err := headerReader.Read()
	if err != nil {
		l.errorfall("error ReadHeader: %s", err)
	}
	for i, v := range record {
		if isNumerical(v) {
			l.errorf(i, "expected non-numerical header for column %d. Got %s", i, v)
		}
		*l.values[i] <- Value{
			str: v,
			col: i,
		}
	}
	return &headerReader
}

func (l *lexer) getNumericalColumns() (I []int) {
	for i, v := range l.typ {
		if v == TypeNumber {
			I = append(I, i)
		}
	}
	return I
}

func (l *lexer) getStringColumns() (I []int) {
	for i, v := range l.typ {
		if v == TypeString {
			I = append(I, i)
		}
	}
	return I
}

func isNumerical(s string) bool {
	_, err := strconv.ParseFloat(s, 32)
	return err == nil
}

func newCsvReader(l *lexer) *csv.Reader {
	fi := l.file
	_,_ = fi.Seek(0,io.SeekStart)
	rdr := csv.NewReader(l.file)
	rdr.TrimLeadingSpace = true
	rdr.Comma = rune(l.Format.sep[0])
	return rdr
}
func (l *lexer) sendOnAll(val Value) {
	for _, c := range l.values {
		*c <- val
	}
}

func (l *lexer) errorfall(f string, args ...interface{}) {
	for i, _ := range l.values {
		l.errorf(i, f, args)
	}
}
