// Ported from fabioxgn/go-bot
package ported

import (
	"fmt"
	lazlo "github.com/djosephsen/lazlo/lib"
	"github.com/fabioxgn/go-bot/web"
)

const (
	pattern   = "(?i)\\b(cat|gato|miau|meow|garfield|lolcat)[s|z]{0,1}\\b"
	msgPrefix = "I love cats! Here's a fact: %s"
)

type facts struct {
	Facts   []string `json:"facts"`
	Success string   `json:"success"`
}

var catFactsURL = "http://catfacts-api.appspot.com/api/facts?number=1"

var CatFacts = &lazlo.Module{
	Name:  `Cat Facts`,
	Usage: `If you talk about cats, %BOTNAME% will give you a fun fact.`,
	Run:   catFactsRun,
}

func catFactsRun(b *lazlo.Broker) {
	cb := b.MessageCallback(`(?i)\b(cat|gato|miau|meow|garfield|lolcat)[s|z]?\b`, false)
	for {
		pm := <-cb.Chan
		data := &facts{}
		err := web.GetJSON(catFactsURL, data)
		if err != nil || len(data.Facts) == 0 {
			return
		}
		pm.Event.Respond(fmt.Sprintf("I love cats! Here's a fact: %v", data.Facts[0]))
	}
}
