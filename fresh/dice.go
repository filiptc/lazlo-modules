package fresh

import (
	"fmt"

	"strings"

	lazlo "github.com/djosephsen/lazlo/lib"
	"github.com/klaidliadon/dice"
)

var Roll = &lazlo.Module{
	Name:  `Roll`,
	Usage: `listen for 'roll nDm', respond with dice roll result`,
	Run: func(b *lazlo.Broker) {
		cb := b.MessageCallback(fmt.Sprintf("roll (%s)+", dice.RollFormat), false)
		for {
			pm := <-cb.Chan
			s := strings.Replace(pm.Match[0], "roll ", "", 1)
			if s == "" {
				return
			}
			p := dice.NewPouch(s)
			p.Roll()
			pm.Event.Respond(fmt.Sprintf("@%s `%s` roll result:\n```%s```", b.SlackMeta.GetUserName(pm.Event.User), s, p.String()))
		}
	},
}
