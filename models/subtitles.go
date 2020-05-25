package models

import (
	"io/ioutil"
	"xboto/mkvpassthru/submerge"
	"os"
	"regexp"
	"strconv"

	"github.com/gookit/color"
)

/*
FindReplaceString remplace string inside str
*/
func FindReplaceString(Filename string, isSRT bool) {
	input, err := ioutil.ReadFile(Filename)
	if err != nil {
		color.Red.Println("Error: Can't read " + Filename)
		os.Exit(1)
	}

	for _, replace := range CFG.ListADStoRemove {
		re := regexp.MustCompile(replace)
		input = re.ReplaceAll([]byte(input), []byte(""))
	}

	if err = ioutil.WriteFile(Filename, input, 0666); err != nil {
		color.Red.Println("Error: Can't write " + Filename)
		os.Exit(1)
	}

	bstr := []byte(CFG.StrAddToStrs)
	if err = ioutil.WriteFile(Filename+"_tmp.srt", bstr, 0666); err != nil {
		color.Red.Println("Error: Can't write temporal .srt")
		os.Exit(1)
	}

	if isSRT {
		input = []byte(submerge.MergeSubs(Filename+"_tmp.srt", Filename))
	}

	err = os.Remove(Filename + "_tmp.srt")
	if err != nil {
		color.Red.Println("Error: Can't remove .srt")
		os.Exit(1)
	}

	if err = ioutil.WriteFile(Filename, input, 0666); err != nil {
		color.Red.Println("Error: Can't write " + Filename)
		os.Exit(1)
	}

}

/*
CleanSubtitles all str files
*/
func CleanSubtitles(Filename string, Subtitles []Track) {
	for _, Sub := range Subtitles {
		namesrt := Filename + strconv.Itoa(Sub.ID) + ".srt"
		isSRT := false
		if Sub.Codec == "SubRip/SRT" {
			isSRT = true
		}
		FindReplaceString(namesrt, isSRT)
	}
}
