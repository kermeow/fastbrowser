package gh

import (
	"bufio"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func readSongInfo(f *os.File) map[string]string {
	values := make(map[string]string)
	scanner := bufio.NewScanner(f)
	inSection := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !inSection {
			// .chart has a stupid 3 byte header despite being a text file
			if strings.HasSuffix(strings.ToLower(line), "[song]") {
				inSection = true
			}
			continue
		}
		if line == "}" || (strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]")) {
			break
		}

		tokens := strings.SplitN(line, "=", 2)
		if len(tokens) != 2 {
			continue
		}

		key := strings.TrimSpace(tokens[0])
		val := strings.TrimSpace(tokens[1])
		values[key] = val
	}
	return values
}

func getOrDefault[K, V comparable](m map[K]V, k K, d V) V {
	if val, ok := m[k]; ok {
		return val
	}
	return d
}

type Chart struct {
	// This is only metadata that we're using
	Name       string
	Artist     string
	Album      string
	Genre      string
	Year       string
	Charter    string
	SongLength uint

	path      string
	chartFile string
	midFile   string
	iniFile   string
}

func (c *Chart) ChartFile() string {
	if c.midFile != "" {
		return filepath.Join(c.path, c.midFile)
	}
	return filepath.Join(c.path, c.chartFile)
}

func ReadChart(dirPath string) (*Chart, error) {
	chart := &Chart{
		Name:       "Unknown",
		Artist:     "Unknown",
		Album:      "Unknown",
		Genre:      "Unknown",
		Year:       "Unknown",
		Charter:    "Unknown",
		SongLength: 0,

		path:      dirPath,
		chartFile: "",
		midFile:   "",
		iniFile:   "",
	}

	files, _ := os.ReadDir(dirPath)
	for _, d := range files {
		if d.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(d.Name()))
		if ext == ".chart" {
			chart.chartFile = d.Name()
		}
		if ext == ".mid" {
			chart.midFile = d.Name()
		}
		if strings.ToLower(d.Name()) == "song.ini" {
			chart.iniFile = d.Name()
		}
	}

	if chart.chartFile != "" {
		file, err := os.Open(path.Join(chart.path, chart.chartFile))
		if os.IsNotExist(err) {
			return chart, err
		}

		values := readSongInfo(file)
		chart.Name = getOrDefault(values, "name", "Unknown")
		chart.Artist = getOrDefault(values, "artist", "Unknown")
		chart.Album = getOrDefault(values, "album", "Unknown")
		chart.Genre = getOrDefault(values, "genre", "Unknown")
		chart.Year = getOrDefault(values, "year", "Unknown")
		chart.Charter = getOrDefault(values, "charter", "Unknown")
	}

	if chart.iniFile != "" {
		file, err := os.Open(path.Join(chart.path, chart.iniFile))
		if os.IsNotExist(err) {
			return chart, err
		}

		values := readSongInfo(file)
		chart.Name = getOrDefault(values, "name", chart.Name)
		chart.Artist = getOrDefault(values, "artist", chart.Artist)
		chart.Album = getOrDefault(values, "album", chart.Album)
		chart.Genre = getOrDefault(values, "genre", chart.Genre)
		chart.Year = getOrDefault(values, "year", chart.Year)
		chart.Charter = getOrDefault(values, "charter", chart.Charter)

		songLenStr := getOrDefault(values, "song_length", "")
		if songLenStr != "" {
			if i, err := strconv.ParseUint(songLenStr, 10, 0); err == nil {
				chart.SongLength = uint(i)
			}
		}
	}

	return chart, nil
}
