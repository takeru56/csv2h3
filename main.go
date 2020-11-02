package main

import (
	"encoding/csv"
	"log"
	"os"
	"os/exec"
	"strings"
)

// [経度, 緯度]のcsvファイル → [h3]のcsvファイル

// *********************************************
var filePath = "./ovd_result_09_bmw.csv"
var resolution = "8"
var outFilePath = "ovd_hex_09_bmw.csv"

// *********************************************

func main() {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var line []string
	var columns [][]string
	columns = append(columns, []string{"id", "created_at"})

	for {
		line, err = reader.Read()
		if err != nil {
			break
		}

		lat := line[0]
		lng := line[1]
		createdAt := line[2]

		if lat == "" || lng == "" {
			continue
		}

		if lat == "lat" {
			continue
		}

		out, _ := exec.Command("geoToH3", "--resolution", resolution, "--latitude", lat, "--longitude", lng).Output()
		columns = append(columns, []string{strings.TrimRight(string(out), "\n"), createdAt})
	}

	f, err := os.Create(outFilePath)
	w := csv.NewWriter(f)

	for _, column := range columns {
		if err := w.Write(column); err != nil {
			log.Fatal(err)
		}
	}

	w.Flush()
}
