package tracks

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// Track holds information about a track
type Track struct {
	// Path points to the actual file
	Path string

	// If the track has been played, the time of its last play is stored here
	LastPlayedAt *time.Time
}

// LoadTracks loads all files within a given folder
// It does not check for file contents so make sure to place
// only wave files there
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
