package main

import (
	"fmt"
	"github.com/hugolgst/rich-go/client"
	"os"
	"time"
)

var token = os.Getenv("LASTFM_API_KEY")
var username = "thallesp"
var pastTrack *RecentTrack

func main() {
	ticker := time.NewTicker(5 * time.Second)

	for {
		<-ticker.C
		track := GetCurrentPlayingTrack()

		if track == nil {
			continue
		}

		if track.Name == pastTrack.Name && track.Artist.Text == pastTrack.Artist.Text {
			continue
		}

		SetStatus(track)
	}
}

func GetCurrentPlayingTrack() *RecentTrack {
	lastClient := NewLastFM(token)

	track, err := lastClient.GetUserListeningNow(username)

	if err != nil {
		return nil
	}

	return track
}

func SetStatus(track *RecentTrack) {
	err := client.Login("1077561394809020416")
	defer client.Logout()

	image := track.Image[3].Text

	if image == "" {
		image = "https://pngimg.com/d/vinyl_PNG35.png"
	}

	err = client.SetActivity(client.Activity{
		State:      track.Name,
		Details:    track.Artist.Text,
		LargeImage: image,
		Buttons: []*client.Button{
			{
				Label: "Last.fm Profile",
				Url:   fmt.Sprintf("https://www.last.fm/user/%s", username),
			},
		},
		LargeText: track.Album.Text,
	})

	if err != nil {
		fmt.Errorf("failed to set activity. Error: %s", err.Error())
	}
}
