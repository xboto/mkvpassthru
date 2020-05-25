# MKVPassthru
MKVPassthru is a simple CLI application to remove/add ads in video files,
within the subtitles included and in the title of the video.

## Requisitos
This application uses the mkvmerge and mkvextract commands that are part of the package [MKVToolNix](https://mkvtoolnix.download)

## Usage
```bash
mkvpassthru name1.mkv name2.mkv ...
```
o (not available in windows)
```bash
mkvpassthru *.mkv
```
## Configuration file
Configuration file path: `~/mkvpassthru.yml` on Unix-based systems or
`%userprofile%\mkvpassthru.yml` in Windows

## Configuration file example `mkvpassthru.yml`
```yaml
---
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
 your_site_name.net
 ```

## How to compile from source
- Install [Go](https://golang.org/doc/install)
- Requirements
```bash
go get github.com/gookit/color
go get github.com/theckman/yacspin
go get github.com/xboto/mkvpassthru
```
- Go to the source folder and compile
```bash
cd ~/github.com/xboto/mkvpassthru
go build .
```
## License
[MIT](https://choosealicense.com/licenses/mit/)