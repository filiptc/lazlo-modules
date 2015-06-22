package replacements

import (
	lazlo "github.com/djosephsen/lazlo/lib"
	"strings"
)

var Help = &lazlo.Module{
	Name:  `Help`,
	Usage: `"%BOTNAME% help": prints the usage information of every registered plugin`,
	Run:   helpRun,
}

func helpRun(b *lazlo.Broker) {
	cb := b.MessageCallback(`(?i)help`, true)
	for {
		pm := <-cb.Chan
		go getHelp(b, &pm)
	}
}

func getHelp(b *lazlo.Broker, pm *lazlo.PatternMatch) {
	a := []lazlo.Attachment{
		lazlo.Attachment{
			Color:  "#ff0000",
			Title:  "Modules In use",
			Fields: []lazlo.AttachmentField{},
		},
	}
	for _, m := range b.Modules {
		if strings.Contains(m.Usage, `%HIDDEN%`) {
			continue
		}
		usage := strings.Replace(m.Usage, `%BOTNAME%`, b.Config.Name, -1)

		a[0].Fields = append(a[0].Fields, lazlo.AttachmentField{
			Title: m.Name,
			Value: usage,
		})
	}
	pm.Event.RespondAttachments(a)
}
