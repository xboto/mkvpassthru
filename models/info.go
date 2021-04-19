package models

import (
	"encoding/json"
	"os"
	"os/exec"

	"github.com/gookit/color"
)

/*
Track return structure
*/
type Track struct {
	ID    int
	Lang  string
	Codec string
}

/*
GetFileInfo MKV info JSON structure
*/
func GetFileInfo(Filename string) ([]Track, []Track, []Track, error) {
	type Info struct {
		FileName string `json:"file_name"`
		Tracks   []struct {
			Codec      string `json:"codec"`
			ID         int    `json:"id"`
			Properties struct {
				Language string `json:"language"`
			}
			Type string `json:"type"`
		} `json:"tracks"`
	}
	out, err := exec.Command("mkvmerge", "-F", "json", "-i", Filename).Output()
	if err != nil {
		color.Red.Println("Error: Can't open the file")
		//os.Exit(1)  <-------
	}
	var mkvinfo Info
	err = json.Unmarshal(out, &mkvinfo)
	if err != nil {
		color.Red.Println("Error: Can't extract information from the file")
		os.Exit(1)
	}

	if len(mkvinfo.Tracks) == 0 {
		return nil, nil, nil, exec.ErrNotFound
	}

	if len(mkvinfo.Tracks) == 1 && mkvinfo.Tracks[0].Type == "subtitles" {
		return nil, nil, nil, exec.ErrNotFound
	}

	var Videos []Track
	var Audios []Track
	var Audios2 []Track
	var Subtitles []Track
	var Subtitles2 []Track
	var Tracktmp []Track
	for k := range mkvinfo.Tracks {
		Tracktmp = []Track{{
			ID:    mkvinfo.Tracks[k].ID,
			Lang:  mkvinfo.Tracks[k].Properties.Language,
			Codec: mkvinfo.Tracks[k].Codec,
		}}

		switch mkvinfo.Tracks[k].Type {
		case "video":

			Videos = append(Videos, Tracktmp[0])
		case "audio":
			if mkvinfo.Tracks[k].Properties.Language == CFG.SortLang {
				Audios = append(Audios, Tracktmp[0])
			} else {
				Audios2 = append(Audios2, Tracktmp[0])
			}

		case "subtitles":
			if mkvinfo.Tracks[k].Properties.Language == CFG.SortLang {
				Subtitles = append(Subtitles, Tracktmp[0])
			} else {
				Subtitles2 = append(Subtitles2, Tracktmp[0])
			}
		}

	}
	for k := range Audios2 {
		Audios = append(Audios, Audios2[k])
	}
	for k := range Subtitles2 {
		Subtitles = append(Subtitles, Subtitles2[k])
	}
	if CFG.MaxTracksAudio > 0 && len(Audios) >= CFG.MaxTracksAudio {
		Audios = Audios[:CFG.MaxTracksAudio]
	}
	if CFG.MaxTracksSub > 0 && len(Subtitles) >= CFG.MaxTracksSub {
		Subtitles = Subtitles[:CFG.MaxTracksSub]
	}

	if len(Videos) == 0 {
		return nil, nil, nil, exec.ErrNotFound
	}

	return Videos, Audios, Subtitles, nil

}
