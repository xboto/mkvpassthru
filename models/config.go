package models

import (
	"io/ioutil"
	"os"

	"github.com/gookit/color"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/mitchellh/go-homedir"
)

/*
Config extract config file
*/
type Config struct {
	ListADStoRemove []string `yaml:"list_ads_to_remove"`
	SortLang        string   `yaml:"sort_lang"`
	MaxTracksAudio  int      `yaml:"max_tracks_audio"`
	MaxTracksSub    int      `yaml:"max_tracks_sub"`
	KeepOrgFile     bool     `yaml:"keep_org_file"`
	AddToTitle      string   `yaml:"add_to_title"`
	StrAddToStrs    string   `yaml:"str_add_to_srts"`
}

/*
CFG global variable from config mkvpassthrucfg.yml
*/
var CFG Config

/*
ReadConfig is read yaml file config
*/
func ReadConfig() {
	fileconfig, err := homedir.Expand("~/mkvpassthru.yml")
	if err != nil {
		color.Red.Println("Error: Invalid path")
		os.Exit(1)
	}

	_, err = os.Stat(fileconfig)
	if err != nil {
		yamlFile := []byte(`---
# List of words/phrases to remove in subtitles
# supports regex, one words/phrases per line
# (?mi) is case insensitive in regex
list_ads_to_remove:
- ADS Web site
- (?mi)jhon ripper
- website.com

# Default language to sort tracks
sort_lang: eng

# Maximum number of audio tracks
# 0 are all tracks
max_tracks_audio: 0

# Maximum number of subtitles tracks
# 0 are all tracks
max_tracks_sub: 0

# Keep the original file, finish add .org
# false deletes the original file
keep_org_file: true 

# Custom video title
add_to_title: Custom Title

# Custom srt to add to each srt
str_add_to_srts: |-
 1
 00:00:05,000 --> 00:00:10,000
 your_site_name.net

 2
 00:20:00,000 --> 00:20:01,500
 your_site_name.net`)
		if err = ioutil.WriteFile(fileconfig, yamlFile, 0666); err != nil {
			color.Red.Println("Error: Can't write " + fileconfig)
			os.Exit(1)
		} else {
			color.Normal.Println("Sample configuration file created, customize content:")
			color.BgCyan.Println(fileconfig)
			os.Exit(0)
		}
	}

	err = cleanenv.ReadConfig(fileconfig, &CFG)
	if err != nil {
		color.Red.Println("Error: Invalid configuration file " + fileconfig)
		color.Normal.Println("Check syntax in YAML file")
		os.Exit(1)
	}

}
