package tracks

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// Track xxx
type Track struct {
	Path         string
	LastPlayedAt *time.Time
}

func LoadTracks(path string) ([]*Track, error) {
	tracks := make([]*Track, 0)
	entries, err := os.ReadDir(path)
	if err != nil {
		return tracks, err
	}
	for _, v := range entries {
		if !v.Type().IsRegular() {
			continue
		}
		fullpath := filepath.Join(path, v.Name())
		tracks = append(tracks, &Track{
			Path:         fullpath,
			LastPlayedAt: nil,
		})
	}

	log.Printf("Loaded the following sfx files: %+v", tracks)
	return tracks, nil
}

/*
func LoadTracks(path string) ([]Tracks, error) {
	var tracks Tracks
	if tracks == nil {
		log.Println("making asset map")
		tracks = make(Tracks)
	}
	newUnfilteredEntries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	newEntriesLookup := make(map[string]int, 0)
	for _, v := range newUnfilteredEntries {
		if !v.Type().IsRegular() {
			continue
		}
		newEntriesLookup[filepath.Join(path, v.Name())] = 1
	}

	// delete items
	for k := range tracks {
		if _, ok := newEntriesLookup[k]; !ok {
			delete(tracks, k)
		}
	}

	// do not re-add existing items
	for k := range tracks {
		if _, ok := newEntriesLookup[k]; ok {
			delete(newEntriesLookup, k)
		}
	}

	for k := range newEntriesLookup {
		tracks[k] = &Track{
			Path:         k,
			LastPlayedAt: nil,
		}
	}

	log.Printf("Loaded the following sfx files: %+v", tracks)
	return nil
}
*/
