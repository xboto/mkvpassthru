package models

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gookit/color"
)

/*
Muxing all tracks
*/
func Muxing(Filename string, Videos []Track, Audios []Track, Subtitles []Track) {
	now := time.Now().Format("20060102150405")
	args := []string{"-q", "-o", Filename + now}
	argcomsrt := []string{"--no-global-tags", "--no-chapters", "--disable-track-statistics-tags", "--no-attachments"}
	argcommon := append(argcomsrt, Filename)
	argsvideo := []string{
		"--track-name", strconv.Itoa(Videos[0].ID) + ":",
		"--default-track", strconv.Itoa(Videos[0].ID) + ":yes", "--forced-track",
		strconv.Itoa(Videos[0].ID) + ":no", "-d",
		strconv.Itoa(Videos[0].ID), "-A", "-S", "-T",
	}
	argsvideo = append(argsvideo, argcommon...)
	args = append(args, argsvideo...)
	torder := "0:" + strconv.Itoa(Videos[0].ID)
	var argsdaudio = []string{}
	for k := range Audios {
		var defaultstr string
		if k == 0 {
			defaultstr = "yes"
		} else {
			defaultstr = "no"
		}
		argsk := []string{
			"--language", strconv.Itoa(Audios[k].ID) + ":" + Audios[k].Lang, "--track-name", strconv.Itoa(Audios[k].ID) + ":",
			"--default-track", strconv.Itoa(Audios[k].ID) + ":" + defaultstr, "--forced-track", strconv.Itoa(Audios[k].ID) + ":no",
			"-a", strconv.Itoa(Audios[k].ID), "-D", "-S", "-T",
		}
		argsk = append(argsk, argcommon...)
		argsdaudio = append(argsdaudio, argsk...)
		torder += "," + strconv.Itoa(k+1) + ":" + strconv.Itoa(Audios[k].ID)
	}

	tnums := len(Audios) + 1
	var argssubtitles = []string{}
	for k := range Subtitles {
		var defaultstr string
		if k == 0 {
			defaultstr = "yes"
		} else {
			defaultstr = "no"
		}
		argsk := []string{
			"--language", "0:" + Subtitles[k].Lang, "--track-name", "0:",
			"--default-track", "0:" + defaultstr, "--forced-track", "0:no",
			"-s", "0", "-D", "-A", "-T",
		}
		argsk = append(argsk, argcomsrt...)
		argsk = append(argsk, Filename+strconv.Itoa(Subtitles[k].ID)+".srt")
		argssubtitles = append(argssubtitles, argsk...)

		torder += "," + strconv.Itoa(k+tnums) + ":" + strconv.Itoa(0)

	}
	args = append(args, argsdaudio...)
	args = append(args, argssubtitles...)
	title := filepath.Base(Filename)[:len(filepath.Base(Filename))-4]
	customTitle := ""
	if CFG.AddToTitle != "" {
		customTitle = " - " + CFG.AddToTitle
	}
	args = append(args, "--track-order", torder, "--title", title+customTitle)

	_, err := exec.Command("mkvmerge", args...).Output()
	if err != nil {
		color.Red.Println("\nError Muxing: Invalid file type")
		os.Exit(1)
	}

	for _, Sub := range Subtitles {
		err := os.Remove(Filename + strconv.Itoa(Sub.ID) + ".srt")
		if err != nil {
			color.Red.Println("Error: Can't remove temporal .srt")
			os.Exit(1)
		}
	}

	err = os.Rename(Filename, Filename+".org")
	if err != nil {
		color.Red.Println("Error: Can't rename to: " + Filename + ".org")
		os.Exit(1)
	}

	err = os.Rename(Filename+now, Filename)
	if err != nil {
		color.Red.Println("Error: Can't rename to: " + Filename)
		os.Exit(1)
	}

	if !CFG.KeepOrgFile {
		err := os.Remove(Filename + ".org")
		if err != nil {
			color.Red.Println("Error: Can't remove " + Filename + ".org")
			os.Exit(1)
		}

	}

}
