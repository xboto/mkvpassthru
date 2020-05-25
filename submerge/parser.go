package submerge

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func parseSubFile(file *os.File) ([]*subLine, error) {
	var lines []*subLine

	sc := bufio.NewScanner(file)

	nextLine := true
	for nextLine {
		line, notEmpty, err := parseSubLine(sc)
		if err != nil {
			return nil, err
		}
		nextLine = notEmpty
		if line != nil {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

func parseSubLine(sc *bufio.Scanner) (*subLine, bool, error) {
	counter := 0
	var currentLine *subLine
	for sc.Scan() {

		if err := sc.Err(); err != nil {
			return nil, false, err
		}
		line := strings.TrimSpace(sc.Text())
		line = strings.Replace(line, "\ufeff", "", -1)
		if line == "" {
			return currentLine, true, nil
		}

		switch counter {
		case 0:
			num, err := strconv.ParseInt(line, 10, 64)
			if err != nil {
				return nil, false, err
			}
			currentLine = &subLine{Num: int(num)}
		case 1:
			currentLine.Time = line
		case 2:
			currentLine.Text1 = line
		case 3:
			currentLine.Text2 = line
		}
		counter++
	}
	return currentLine, false, nil
}
