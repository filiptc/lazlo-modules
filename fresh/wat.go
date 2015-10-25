package fresh

import (
	"encoding/json"
	"fmt"
	"strings"

	lazlo "github.com/djosephsen/lazlo/lib"
	"github.com/sadbox/mediawiki"
)

var Wat = &lazlo.Module{
	Name:  `Wat`,
	Usage: `Wikipedia queries`,
	Run:   watRun,
}

func watRun(b *lazlo.Broker) {
	cb := b.MessageCallback(`(?i)^wat +is +(.*)$`, false)
	for {
		pm := <-cb.Chan
		search := pm.Match[1]
		client, err := mediawiki.New(`http://en.wikipedia.org/w/api.php`, "lazlo")
		if err != nil {
			pm.Event.Respond(err.Error())
		}
		response, err := readExtract(client, search)
		if err != nil {
			pm.Event.Respond(err.Error())
		}

		response.GenPageList()
		fmt.Printf("%v", response)
		pm.Event.Respond(
			fmt.Sprintf(
				"*%v*\n%v",
				response.Query.PageList[0].Title,
				getFirstParagraph(response.Query.PageList[0].Extract),
			),
		)
	}
}

func readExtract(client *mediawiki.MWApi, search string) (*Response, error) {
	query := map[string]string{
		"action":      "query",
		"prop":        "extracts",
		"titles":      search,
		"redirects": "true",
		"exintro":     "true",
		"explaintext": "true",
		"exsentences": "3",
		"exsectionformat": "plain",
	}
	body, err := client.API(query)
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func getFirstParagraph(str string) string {
	return strings.Split(str, "\n")[0]
}

type Response struct {
	Query struct {
		// The json response for this part of the struct is dumb.
		// It will return something like { '23': { 'pageid': 23 ...
		//
		// As a workaround you can use GenPageList which will create
		// a list of pages from the map.
		Pages    map[string]Page
		PageList []Page
	}
}

// GenPageList generates PageList from Pages to work around the sillyness in
// the mediawiki API.
func (r *Response) GenPageList() {
	r.Query.PageList = []Page{}
	for _, page := range r.Query.Pages {
		r.Query.PageList = append(r.Query.PageList, page)
	}
}

// Page is a mediawiki page and the meatadata about it.
type Page struct {
	Pageid  int
	Ns      int
	Title   string
	Extract string
}
