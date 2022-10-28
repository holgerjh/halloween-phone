package tracks

import "time"

type TrackDB map[string]*Track

func NewDB() TrackDB {
	return make(TrackDB)
}

func (db TrackDB) HasTrack(track string) bool {
	_, ok := db[track]
	return ok
}

// TrackDB xxx

func (db TrackDB) InsertFiltered(tracks []*Track) {
	lookup := TracksToLookupTable(tracks)

	for k := range db {
		if _, ok := lookup[db[k].Path]; !ok {
			delete(db, k)
		}
	}
	for _, v := range tracks {
		if _, ok := db[v.Path]; !ok {
			db[v.Path] = v
		}
	}
}

func (db TrackDB) FilterForEarliestTracks() []*Track {
	var min *time.Time
	var best []*Track
	for _, v := range db {
		if min == nil || min.Before(*v.LastPlayedAt) {
			min = v.LastPlayedAt
			best = []*Track{v}
		} else if min != nil && min.Equal(*v.LastPlayedAt) {
			best = append(best, v)
		}
	}
	return best
}

func TracksToLookupTable(tracks []*Track) map[string]*Track {
	lookup := make(map[string]*Track)
	for _, v := range tracks {
		lookup[v.Path] = v
	}
	return lookup
}

func LoadTracksIntoNewDB(folder string) (*TrackDB, error) {
	db := NewDB()
	t, err := LoadTracks(folder)
	if err != nil {
		return nil, err
	}
	db.InsertFiltered(t)
	return &db, nil
}
