package main

import (
	"fmt"
)


func main() {
	buzzers := map[int]string{3:"fizz", 5:"buzz"}
	buzz(buzzers)
}

func buzz(b map[int]string) {
	for i:=1;i<=100;i++ {
		var printing bool
		for j,v := range b {
			if i % j == 0 {
				fmt.Printf(v)
				printing = true
			}
		}
		if !printing {
			fmt.Printf("%d", i)
		}
		fmt.Println("")
	}
}
//
//func Downsample(jobname string, l *lexer, tol float64) error {
//
//	if len(f.colstr) != 2 {
//		return fmt.Errorf("Too many/few columns inputted to LinearDownsampler. Need 2.")
//	}
//	xColumn, yColumn, outputSep := f.colstr[0], f.colstr[1], f.sep
//	scanner := bufio.NewScanner(csv.ID)
//	xidx := csv.columns[xColumn]
//	yidx := csv.columns[yColumn]
//
//	var xval, yval, xlast, ylast, Dx, Dy, angle, anglemin, anglemax float64
//	var xstand, ystand float64
//	var rowstr []string
//
//	fo, err := os.Create(outputName)
//	if err != nil {
//		return err
//	}
//	header := format([]string{csv.colstr[xidx],csv.colstr[yidx]}, outputSep)
//	_, err = fo.Write([]byte(header + "\n"))
//	if err != nil {
//		return err
//	}
//	eof := scanner.Scan();eof = scanner.Scan() // skip a line to get to actual values
//	rowstr = strings.Split(scanner.Text(), csv.sep)
//	if !areNumeric(rowstr) {
//		return fmt.Errorf("Non numeric entry in csv during LinearDownsample. Line %d", 2)
//	}
//	_,_ = fo.Write([]byte(format([]string{rowstr[xidx],rowstr[yidx],"\n"},outputSep)))
//	xlast, _ = strconv.ParseFloat(rowstr[xidx], 64)
//	ylast, _ = strconv.ParseFloat(rowstr[yidx], 64)
//	eof = scanner.Scan()
//	rowstr = strings.Split(scanner.Text(), csv.sep)
//	xval, _ = strconv.ParseFloat(rowstr[xidx], 64)
//	yval, _ = strconv.ParseFloat(rowstr[yidx], 64)
//	line := 3
//
//	for eof {
//		xstand, ystand = xlast, ylast
//		Dx, Dy = xval-xstand, yval-ystand
//		permLoAngle, permHiAngle := math.Atan2(Dy-tol, Dx), math.Atan2( Dy+tol, Dx)
//		xlast,ylast = xval,yval
//		for eof  {
//			eof = scanner.Scan()
//			if !eof {
//				break
//			}
//			line++
//			rowstr = strings.Split(scanner.Text(), csv.sep)
//			if !areNumeric(rowstr) {
//				return fmt.Errorf("Non numeric entry in csv during LinearDownsample. Line %d", 2)
//			}
//
//			xval, _ = strconv.ParseFloat(rowstr[xidx], 64)
//			yval, _ = strconv.ParseFloat(rowstr[yidx], 64)
//			Dx, Dy = xval-xstand, yval-ystand
//			angle, anglemin, anglemax = math.Atan2(Dy, Dx), math.Atan2(Dy-tol, Dx), math.Atan2( Dy+tol, Dx)
//			if angle < permLoAngle || angle > permHiAngle {
//				break
//			}
//			if anglemin >= permLoAngle {
//				permLoAngle = anglemin
//			}
//			if anglemax <= permHiAngle {
//				permHiAngle = anglemax
//			}
//			xlast, ylast = xval, yval
//		}
//		_,_ = fo.Write([]byte(fmt.Sprintf("%1.5e%s%1.5e\n", xlast,outputSep, ylast)))
//		//_ = fo.Sync()
//	}
//	return nil
//}
//
//// More aggressive downsampler. Eliminates more points and interpolates on y-axis.
//func (csv *File) InterpolatorDownsample(outputName string, f Format, tol float64) error {
//	_, _ = csv.ID.Seek(0, 0)
//	if len(f.colstr) != 2 {
//		return fmt.Errorf("Too many/few columns inputted to LinearDownsampler. Need 2.")
//	}
//	xColumn, yColumn, outputSep := f.colstr[0], f.colstr[1], f.sep
//	scanner := bufio.NewScanner(csv.ID)
//	xidx := csv.columns[xColumn]
//	yidx := csv.columns[yColumn]
//
//	var xval, yval, xlast, ylast, Dx, Dy, anglemin, anglemax float64
//	var xstand, ystand float64
//	var rowstr []string
//
//	fo, err := os.Create(outputName)
//	if err != nil {
//		return err
//	}
//	header := format([]string{csv.colstr[xidx],csv.colstr[yidx]}, outputSep)
//	_, err = fo.Write([]byte(header + "\n"))
//	if err != nil {
//		return err
//	}
//	eof := scanner.Scan();eof = scanner.Scan() // skip a line to get to actual values
//	rowstr = strings.Split(scanner.Text(), csv.sep)
//	if !areNumeric(rowstr) {
//		return fmt.Errorf("Non numeric entry in csv during LinearDownsample. Line %d", 2)
//	}
//	_,_ = fo.Write([]byte(format([]string{rowstr[xidx],rowstr[yidx],"\n"},outputSep)))
//	xlast, _ = strconv.ParseFloat(rowstr[xidx], 64)
//	ylast, _ = strconv.ParseFloat(rowstr[yidx], 64)
//	eof = scanner.Scan()
//	rowstr = strings.Split(scanner.Text(), csv.sep)
//	xval, _ = strconv.ParseFloat(rowstr[xidx], 64)
//	yval, _ = strconv.ParseFloat(rowstr[yidx], 64)
//	line := 3
//
//	for eof {
//		xstand, ystand = xlast, ylast
//		Dx, Dy = xval-xstand, yval-ystand
//		permLoAngle, permHiAngle := math.Atan2(Dy-tol, Dx), math.Atan2( Dy+tol, Dx)
//		xlast,ylast = xval,yval
//		for eof  {
//			eof = scanner.Scan()
//			if !eof {
//				break
//			}
//			line++
//			rowstr = strings.Split(scanner.Text(), csv.sep)
//			if !areNumeric(rowstr) {
//				return fmt.Errorf("Non numeric entry in csv during LinearDownsample. Line %d", 2)
//			}
//
//			xval, _ = strconv.ParseFloat(rowstr[xidx], 64)
//			yval, _ = strconv.ParseFloat(rowstr[yidx], 64)
//			Dx, Dy = xval-xstand, yval-ystand
//			anglemin, anglemax = math.Atan2(Dy-tol, Dx), math.Atan2( Dy+tol, Dx)
//
//			if permHiAngle < anglemin || permLoAngle > permHiAngle {
//				ylast = ystand + (xlast - xstand) * (math.Tan(permHiAngle) + math.Tan(permLoAngle)) / 2
//				break
//			}
//
//			if anglemin > permLoAngle {
//				permLoAngle = anglemin
//			}
//			if anglemax < permHiAngle {
//				permHiAngle = anglemax
//			}
//			xlast, ylast = xval, yval
//
//		}
//		_,_ = fo.Write([]byte(fmt.Sprintf("%1.5e%s%1.5e\n", xlast,outputSep, ylast)))
//		//_ = fo.Sync()
//	}
//	return nil
//}
//
//
