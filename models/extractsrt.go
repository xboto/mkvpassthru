package models

import (
	"os"
	"os/exec"
	"strconv"

	"github.com/gookit/color"
)

/*
ExtractSrt Extract the subtitles from the MKV
*/
func ExtractSrt(Filename string, Subtitles []Track) {
	for _, Sub := range Subtitles {
		args := []string{"tracks", Filename, strconv.Itoa(Sub.ID) + ":" + Filename + strconv.Itoa(Sub.ID) + ".srt"}
		_, err := exec.Command("mkvextract", args...).Output()
		if err != nil {
			color.Red.Println("Error: Can't extract subtitles")
			os.Exit(1)
		}
	}
}
