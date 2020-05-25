package main

import (
	"fmt"
	"xboto/mkvpassthru/models"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/gookit/color"
	"github.com/theckman/yacspin"
)

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
	nargs := len(args)
	if nargs == 0 {
		color.Gray.Println(appName + " v0.8b")
		fmt.Println("Usage: " + appName + " name1.mkv name2.mkv ...")
		if runtime.GOOS != "windows" {
			fmt.Println("       " + appName + " *.mkv")
		}
		os.Exit(1)
	}

	models.ReadConfig()

	okChar := "✔"
	koChar := "☓"
	spinnN := 7
	if runtime.GOOS == "windows" {
		okChar = "√"
		koChar = "x"
		spinnN = 9
	}

	for k, input := range args {
		cNumber := strconv.Itoa(k + 1)
		if len(cNumber) < 2 && nargs > 9 {
			cNumber = " " + cNumber
		}
		suffix := " "
		if nargs > 1 {
			suffix = " [" + cNumber + "/" + strconv.Itoa(nargs) + "] "
		}
		cfg := yacspin.Config{
			Frequency:       100 * time.Millisecond,
			CharSet:         yacspin.CharSets[spinnN],
			Suffix:          suffix + filepath.Base(input),
			SuffixAutoColon: true,
			Message:         " ",
			StopCharacter:   okChar,
			StopColors:      []string{"fgGreen"},
		}
		Spinner, err := yacspin.New(cfg)

		Videos, Audios, Subtitles, err := models.GetFileInfo(input)
		Spinner.Start()
		if err == nil {
			Spinner.Message("Extracting subtitles ")
			models.ExtractSrt(input, Subtitles)
			Spinner.Message("Cleaning subtitles ")
			models.CleanSubtitles(input, Subtitles)
			Spinner.Message("Muxing tracks ")
			models.Muxing(input, Videos, Audios, Subtitles)

		} else {
			red := color.FgRed.Render
			Spinner.StopCharacter(red(koChar))
			Spinner.StopMessage(red("File type is not supported"))
		}
		Spinner.Stop()
	}

}
