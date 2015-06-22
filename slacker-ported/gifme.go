package ported

import (
	"encoding/json"
	"fmt"
	lazlo "github.com/djosephsen/lazlo/lib"
	"net/http"
	"net/url"
)

type gifyout struct {
	Meta interface{}
	Data struct {
		Tags                            []string
		Caption                         string
		Username                        string
		Image_width                     string
		Image_frames                    string
		Image_mp4_url                   string
		Image_url                       string
		Image_original_url              string
		Url                             string
		Id                              string
		Type                            string
		Image_height                    string
		Fixed_height_downsampled_url    string
		Fixed_height_downsampled_width  string
		Fixed_height_downsampled_height string
		Fixed_width_downsampled_url     string
		Fixed_width_downsampled_width   string
		Fixed_width_downsampled_height  string
		Rating                          string
	}
}

var Gifme = &lazlo.Module{
	Name:  `Gifme`,
	Usage: `"%BOTNAME% gif me freddie mercury": prints the usage information of every registered plugin`,
	Run:   gifmeRun,
}

func gifmeRun(b *lazlo.Broker) {
	cb := b.MessageCallback(`(?i)gif me (.*)`, true)
	for {
		pm := <-cb.Chan

		search := pm.Match[1]
		q := url.QueryEscape(search)
		myurl := fmt.Sprintf("http://api.giphy.com/v1/gifs/random?rating=pg-13&api_key=dc6zaTOxFJmzC&tag=%s", q)
		g := new(gifyout)
		resp, _ := http.Get(myurl)
		dec := json.NewDecoder(resp.Body)
		dec.Decode(g)
		pm.Event.Respond(g.Data.Image_url)
	}
}
