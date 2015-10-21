package fresh

import (
	"bytes"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	lazlo "github.com/djosephsen/lazlo/lib"
)

var Roll = &lazlo.Module{
	Name:  `Roll`,
	Usage: `listen for 'roll nDm', respond with dice roll result`,
	Run:   diceRun,
}

func diceRun(b *lazlo.Broker) {
	cb := b.MessageCallback("roll "+splitter.String(), false)
	for {
		pm := <-cb.Chan
		p := NewPouch(pm.Match[0])
		p.Roll()
		pm.Event.Respond(fmt.Sprintf("Result for `%s`\n```%s```", p.src, p.String()))
	}
}

func NewPouch(s string) *Pouch {
	matches := splitter.FindAllStringSubmatch(strings.Replace(s, " ", "", -1), -1)
	var r = make([]Item, len(matches))
	for i := 0; i < len(matches); i++ {
		m := matches[i]
		if m[3] != "" {
			r[i] = NewBonus(m[0])
		} else {
			r[i] = NewDice(m[0], m[1], m[2])
		}
	}
	return &Pouch{s, r}
}

type Pouch struct {
	src   string
	items []Item
}

func (p *Pouch) Roll() {
	for _, i := range p.items {
		i.Roll()
	}
}

func (p *Pouch) Total() int {
	var t = 0
	for _, i := range p.items {
		t += i.Total()
	}
	return t
}

func (p *Pouch) String() string {
	var b = bytes.NewBuffer(nil)
	for _, i := range p.items {
		fmt.Fprintf(b, "%s\t%v\t%v\n", i, i.Total(), i.Partials())
	}
	fmt.Fprintf(b, "------------------\nTotal\t%d\n", p.Total())
	return b.String()
}

type Item interface {
	Roll()
	Total() int
	Partials() []int
	String() string
}

func NewDice(sign, number, face string) Item {
	f, _ := strconv.Atoi(face)
	q, _ := strconv.Atoi(number)
	if q == 0 {
		q++
	}
	return &Dice{Sign: sign[0] != '-', Face: f, Qty: q}
}

type Dice struct {
	Sign      bool
	Qty, Face int
	Results   []int
}

func (d *Dice) Roll() {
	s := rand.NewSource(time.Now().UnixNano())
	var r = make([]int, d.Qty)
	for i := 0; i < d.Qty; i++ {
		r[i] = 1 + rand.New(s).Intn(d.Face)
	}
	d.Results = r
}

func (d *Dice) Total() int {
	var tot int
	for _, s := range d.Results {
		tot += s
	}
	if !d.Sign {
		return -tot
	}
	return tot
}

func (d *Dice) Partials() []int {
	return d.Results
}

func (d *Dice) String() string {
	var b = bytes.NewBuffer(nil)
	if !d.Sign {
		b.WriteRune('-')
	}
	fmt.Fprintf(b, "%dd%d", d.Qty, d.Face)
	return b.String()
}

func NewBonus(s string) Item {
	v, _ := strconv.Atoi(s)
	return &Bonus{v}
}

type Bonus struct {
	Value int
}

func (b *Bonus) Roll() {}

func (b *Bonus) Total() int { return b.Value }

func (b *Bonus) Partials() []int { return nil }

func (b *Bonus) String() string { return strconv.Itoa(b.Value) }

var splitter = regexp.MustCompile(`\s*[+-]?\s*([0-9]{0,})[dD]([0-9]{1,})|([+-]?[0-9]*)`)
