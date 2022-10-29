package player

import (
	"log"
	"math/rand"
	"time"

	"github.com/holgerjh/halloween-phone/config"
	"github.com/holgerjh/halloween-phone/tracks"
)

// Loop scheudles tracks for playing until the shutdown channel gets closed
func Loop(micFile string, db *tracks.TrackDB, cfg *config.Config, shutdown <-chan int) {
	if len(*db) == 0 {
		panic("db has no tracks!")
	}
	for {
		log.Printf("Calculating intended playing start time")
		plannedPlayingTime := calculateNextPlayingTime(cfg.MaxWait, cfg.MinWait)
		log.Printf("Want: Start playing in %s", plannedPlayingTime)
		log.Printf("Choosing a track")
		next := getNextTrack(cfg.TrackCooldown, db, plannedPlayingTime)
		if next == nil {
			log.Printf("No track ready in time, choosing from earliest tracks and maximizing delay")
			next = oneOf(db.FilterForEarliestTracks())
			plannedPlayingTime = time.Now().Add(time.Duration(cfg.MaxWait) * time.Second)
		}
		log.Printf("Upcoming track: %v on %s", next, plannedPlayingTime)

		log.Printf("Waiting")
		select {
		case <-time.After(time.Until(plannedPlayingTime)):
			log.Printf("Playing track %s", next.Path)
			playTrack(micFile, next, cfg.SilenceStartOfCall, &cfg.Mic, shutdown)

		case <-shutdown:
			log.Printf("Received shutdown signal, leaving thread")
			return
		}
	}

}

func oneOf(candidates []*tracks.Track) *tracks.Track {
	log.Printf("Choosing from %d sfx items", len(candidates))
	if len(candidates) == 0 {
		return nil
	}
	choiceIdx := rand.Intn(len(candidates))
	log.Printf("-> Picking slot %d", choiceIdx)
	choice := candidates[choiceIdx]
	log.Printf("Chose %+v", choice)
	return choice
}

func calculateNextPlayingTime(maxWait, minWait int) time.Time {
	delay := rand.Intn(maxWait-minWait) + minWait
	return time.Now().Add(time.Duration(delay) * time.Second)
}

func getNextTrack(cooldown int, db *tracks.TrackDB, plannedPlayingTime time.Time) *tracks.Track {
	candidates := make([]*tracks.Track, 0)
	log.Printf("Computing available tracks")
	for _, v := range *db {
		ok := false
		if v.LastPlayedAt == nil {
			log.Printf("-> track %s has not been played, adding", v.Path)
			ok = true
		} else if plannedPlayingTime.Sub(*v.LastPlayedAt) >= (time.Duration(cooldown) * time.Second) {
			log.Printf("-> track %s has been played but became available again", v.Path)
			ok = true
		} else {
			log.Printf("-> track %s is not available", v.Path)
		}
		if ok {
			candidates = append(candidates, v)
		}
	}
	return oneOf(candidates)
}
