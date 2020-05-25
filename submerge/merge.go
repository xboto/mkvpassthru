/*
Package submerge github.com/gsx95/submerge
*/
package submerge

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/gookit/color"
)

/*
MergeSubs mergue two srt
*/
func MergeSubs(f1 string, f2 string) string {
	file1, err := openFile(f1)
	if err != nil {
		color.Red.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	file2, err := openFile(f2)
	if err != nil {
		color.Red.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	defer closeFile(file1)
	defer closeFile(file2)

	lines, err := parseSubFile(file1)
	if err != nil {
		color.Red.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	lines2, err := parseSubFile(file2)
	if err != nil {
		color.Red.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	lines = append(lines, lines2...)
	sort.Slice(lines, func(i, j int) bool {
		return lines[j].isAfter(lines[i])
	})

	adjustNums(lines)
	return writeLinesToString(lines)
}

func writeLinesToString(lines []*subLine) string {

	w := strings.Builder{}
	for _, line := range lines {
		s := line.toFormat()
		_, err := w.WriteString(s)
		if err != nil {
			panic(err)
		}
	}
	return w.String()
}

func writeLinesToFile(lines []*subLine, outPath string) {
	f, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(f)
	for _, line := range lines {
		s := line.toFormat()
		_, err := w.WriteString(s)
		if err != nil {
			panic(err)
		}
	}
	err = w.Flush()
	if err != nil {
		panic(err)
	}
}

func writeNums(lines []*subLine) {
	for _, line := range lines {
		fmt.Println(line.Num)
	}
}

func writeTimes(lines []*subLine) {
	for _, line := range lines {
		fmt.Println(line.Time)
	}
}

func printMissingNums(lines []*subLine) {
	lastNum := -1
	for _, line := range lines {
		if line == nil {
			continue
		}
		if lastNum == -1 {
			lastNum = line.Num
			continue
		}
		if line.Num != lastNum+1 {
			fmt.Println(lastNum + 1)
		}
		lastNum = line.Num
	}
}
