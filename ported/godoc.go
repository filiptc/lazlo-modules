// Ported from fabioxgn/go-bot
package ported

import (
	"fmt"
	lazlo "github.com/djosephsen/lazlo/lib"
	"github.com/fabioxgn/go-bot/web"
	"net/url"
)

const (
	godocSiteURL    = "http://godoc.org"
	noPackagesFound = "No packages found."
)

type godocResults struct {
	Results []struct {
		Path     string `json:"path"`
		Synopsis string `json:"synopsis"`
	} `json:"results"`
}

var godocSearchURL = "http://api.godoc.org/search"

var GoDoc = &lazlo.Module{
	Name:  `GoDoc`,
	Usage: `godoc <package name>: Searchs godoc.org and displays the first 10 results.`,
	Run:   goDocRun,
}

func goDocRun(b *lazlo.Broker) {
	cb := b.MessageCallback(`(?i)^godoc ([^\s]+)$`, false)
	for {
		pm := <-cb.Chan
		data := &godocResults{}

		url, _ := url.Parse(godocSearchURL)
		q := url.Query()
		q.Set("q", pm.Match[1])
		url.RawQuery = q.Encode()
		err := web.GetJSON(url.String(), data)
		if err != nil {
			pm.Event.Respond(fmt.Sprintf("Error: ", err.Error()))
		}

		if len(data.Results) == 0 {
			pm.Event.Respond(noPackagesFound)
		}
		a := []lazlo.Attachment{
			lazlo.Attachment{
				Color:      "#ff0000",
				Title:      fmt.Sprintf(`Results of "%v"`, pm.Match[1]),
				MarkdownIn: []string{"fields"},
				Fields:     []lazlo.AttachmentField{},
			},
		}

		for i := 0; i < len(data.Results); i++ {
			if i == 10 {
				break
			}
			if data.Results[i].Synopsis == "" {
				data.Results[i].Synopsis = "_No description_"
			}
			a[0].Fields = append(a[0].Fields, lazlo.AttachmentField{
				Title: fmt.Sprintf("%s", data.Results[i].Path),
				Value: fmt.Sprintf("%s\n<%s/%s|full GoDoc>", data.Results[i].Synopsis, godocSiteURL, data.Results[i].Path),
			})
		}

		pm.Event.RespondAttachments(a)
	}
}
