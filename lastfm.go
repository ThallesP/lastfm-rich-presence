package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const BaseEndpoint = "http://ws.audioscrobbler.com/2.0/"

type RecentTrack struct {
	Artist struct {
		Mbid string `json:"mbid"`
		Text string `json:"#text"`
	} `json:"artist"`
	Streamable string `json:"streamable"`
	Image      []struct {
		Size string `json:"size"`
		Text string `json:"#text"`
	} `json:"image"`
	Mbid  string `json:"mbid"`
	Album struct {
		Mbid string `json:"mbid"`
		Text string `json:"#text"`
	} `json:"album"`
	Name string `json:"name"`
	Attr struct {
		Nowplaying string `json:"nowplaying"`
	} `json:"@attr,omitempty"`
	URL  string `json:"url"`
	Date struct {
		Uts  string `json:"uts"`
		Text string `json:"#text"`
	} `json:"date,omitempty"`
}
type Track struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Duration   string `json:"duration"`
	Streamable struct {
		Text      string `json:"#text"`
		Fulltrack string `json:"fulltrack"`
	} `json:"streamable"`
	Listeners string `json:"listeners"`
	Playcount string `json:"playcount"`
	Artist    struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"artist"`
	Album struct {
		Artist string `json:"artist"`
		Title  string `json:"title"`
		URL    string `json:"url"`
		Image  []struct {
			Text string `json:"#text"`
			Size string `json:"size"`
		} `json:"image"`
	} `json:"album"`
	Toptags struct {
		Tag []interface{} `json:"tag"`
	} `json:"toptags"`
}

type NowListeningTrackResponse struct {
	Recenttracks struct {
		Track []RecentTrack `json:"track"`
	} `json:"recenttracks"`
}

type LastFM struct {
	token string
}

func NewLastFM(token string) *LastFM {
	return &LastFM{
		token: token,
	}
}

func (l LastFM) getHTTPClient() *http.Client {
	c := http.Client{
		Timeout: 10 * time.Second,
	}

	return &c
}

func (l LastFM) setDefaultParams(method string, r *http.Request) url.Values {
	params := r.URL.Query()
	params.Add("api_key", l.token)
	params.Add("format", "json")
	params.Add("method", method)

	return params
}

func (l LastFM) GetUserListeningNow(username string) (*RecentTrack, error) {
	c := l.getHTTPClient()

	req, err := http.NewRequest("GET", BaseEndpoint, nil)

	params := l.setDefaultParams("user.getrecenttracks", req)
	params.Add("user", username)

	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = params.Encode()

	resp, err := c.Do(req)

	if err != nil {
		return nil, err
	}

	var nowResp *NowListeningTrackResponse

	bodyB, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyB, &nowResp)

	if err != nil {
		return nil, err
	}

	for _, track := range nowResp.Recenttracks.Track {
		if track.Attr.Nowplaying == "true" {
			return &track, nil
		}
	}

	return nil, errors.New("user is not listening to any track")
}
