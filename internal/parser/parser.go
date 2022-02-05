package parser

import (
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Song struct {
	Lines        []string
	Name, Artist string
}

type dataColumns struct {
	artist, line, songName int8
}

func ReadSongs() []Song {
	targetDir := "./assets"
	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		log.Fatal(err)
	}

	songs := make([]Song, 0, 16)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		if ext := filepath.Ext(filename); ext != ".csv" {
			log.Printf("skip file %s cause of extension %s", filename, ext)
			continue
		}
		res, err := readFile(filepath.Join(targetDir, filename))
		if err != nil {
			log.Fatal(err)
			continue
		}
		songs = append(songs, res...)
	}
	log.Printf("done parsing songs, got %d", len(songs))
	return songs
}

func readFile(filePath string) ([]Song, error) {
	res := make([]Song, 0, 16)

	fileReader, err := os.Open(filePath)
	if err != nil {
		return res, err
	}

	reader := csv.NewReader(fileReader)
	cols, err := describeFile(reader)
	if err != nil {
		return res, err
	}

	curSong := Song{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
			continue
		}

		// новая песня
		if record[cols.songName] != curSong.Name {
			if curSong.Name != "" {
				res = append(res, curSong)
			}
			curSong = Song{
				Name:   record[cols.songName],
				Artist: record[cols.artist],
				Lines:  make([]string, 0, 64),
			}
		}
		curSong.Lines = append(curSong.Lines, record[cols.line])
	}

	return res, nil
}

// describeFile читает из файла строку заголовка и определяет колонки для работы
func describeFile(reader *csv.Reader) (dataColumns, error) {
	res := dataColumns{-1, -1, -1}
	line, err := reader.Read()
	if err == io.EOF {
		return res, errors.New("unable to describe file: EOF on first line")
	} else if err != nil {
		return res, err
	}

	for index, val := range line {
		if index > 100 {
			break
		}
		switch val {
		case "artist":
			res.artist = int8(index)
		case "track_title":
			res.songName = int8(index)
		case "lyric":
			res.line = int8(index)
		}
	}

	if res.artist == -1 || res.line == -1 || res.songName == -1 {
		log.Println(res)
		return res, errors.New("unable to describe file: not all cols found")
	}
	return res, nil
}
