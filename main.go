package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/xboto/mkvpassthru/models"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/gookit/color"
	"github.com/theckman/yacspin"
)

func widthterm() int {
	fd := int(os.Stdout.Fd())
	termWidth, _, _ := terminal.GetSize(fd)
	return termWidth
}

func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 1 {
			//num -= 1
			num--
		}
		bnoden = str[0:num] + "…"
		//bnoden = str[0:num] + "..."
	}
	return bnoden
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if filepath.Base(path)[0:1] != "." {
			if info, err := os.Stat(path); err == nil && !info.IsDir() {
				*files = append(*files, path)
			}
		}
		return nil
	}
}

func main() {

	_, err := exec.LookPath("mkvmerge")
	if err != nil {
		color.Red.Println("Error: \"mkvmerge\" executable file not found in $PATH")
		fmt.Println("Please install the latest package \"MKVtoolnix\" or fix $PATH")
		os.Exit(1)
	}
	_, err = exec.LookPath("mkvextract")
	if err != nil {
		color.Red.Println("Error: \"mkvextract\" executable file not found in $PATH")
		fmt.Println("Please install the latest package \"MKVtoolnix\" or fix $PATH")
		os.Exit(1)
	}

	appName := os.Args[0]
	args := os.Args[1:]
	if len(args) == 0 {
		color.Gray.Println(appName + " v0.9b by xboto")
		fmt.Println("Usage: " + appName + " dir1 dir2 name1.mkv name2.mp4 …")
		if runtime.GOOS != "windows" {
			fmt.Println("       " + appName + " *.mkv")
		}
		os.Exit(1)
	}

	//fmt.Println(args)

	var files []string
	for _, input := range args {
		if info, err := os.Stat(input); err == nil && info.IsDir() {
			err = filepath.Walk(input, visit(&files))
			if err != nil {
				panic(err)
			}
			//input = input + "/*.mkv"
			//fmt.Print("dir:")
			//args[k] = filepath.Join(args[k], "*")
			//args[k] = args[k] + "/*"
		} else {
			files = append(files, input)
		}
		//fmt.Println(input)
		//fmt.Println(k)
		//fmt.Println(args[k])

	}
	//fmt.Println(files)
	nfiles := len(files)
	//os.Exit(1)

	models.ReadConfig()

	okChar := "✔"
	koChar := "☓"
	spinnN := 7
	if runtime.GOOS == "windows" {
		okChar = "√"
		koChar = "x"
		spinnN = 9
	}

	for k, input := range files {
		//fmt.Println(input)
		cNumber := strconv.Itoa(k + 1)
		if len(cNumber) < 2 && nfiles > 9 {
			cNumber = " " + cNumber
		}

		suffix := " "
		if nfiles > 1 {
			suffix = "[" + cNumber + "/" + strconv.Itoa(nfiles) + "] "
		}
		if widthterm() < 24 {
			suffix = ""
		} else {
			suffix = truncateString(suffix+filepath.Base(input), widthterm()-24)
		}

		cfg := yacspin.Config{
			Frequency:       100 * time.Millisecond,
			CharSet:         yacspin.CharSets[spinnN],
			Suffix:          suffix,
			SuffixAutoColon: true,
			Message:         " ",
			StopCharacter:   okChar,
			StopColors:      []string{"fgGreen"},
		}

		Spinner, err := yacspin.New(cfg)

		Videos, Audios, Subtitles, err := models.GetFileInfo(input)
		Spinner.Start()
		if err == nil {
			Spinner.Message(truncateString("Extracting subtitles ", widthterm()-9))
			models.ExtractSrt(input, Subtitles)
			Spinner.Message(truncateString("Cleaning subtitles ", widthterm()-9))
			models.CleanSubtitles(input, Subtitles)
			Spinner.Message(truncateString("Muxing tracks ", widthterm()-6))
			models.Muxing(input, Videos, Audios, Subtitles)

		} else {
			red := color.FgRed.Render
			Spinner.StopCharacter(red(koChar))
			Spinner.StopMessage(red("File type is not supported"))
		}
		Spinner.Stop()
	}

}
