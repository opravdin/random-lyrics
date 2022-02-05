package mixer

import (
	"errors"
	"log"
	"math/rand"
	"sync"
)

type Processor struct {
	RecipeBook   RecipeBook
	accessTagged sync.Mutex
	TaggedLines  map[*Tag][]*string
}

func New() Processor {
	rb := InitBook()
	tagged := make(map[*Tag][]*string)
	proc := Processor{
		rb, sync.Mutex{}, tagged,
	}

	return proc
}

func (p *Processor) initTaggedLines() {
	p.accessTagged.Lock()
	for _, tag := range p.RecipeBook.Tags {
		p.TaggedLines[tag] = make([]*string, 0, 100)
	}
	p.accessTagged.Unlock()
}

func (p Processor) ValidateTags() error {
	for tname, tag := range p.RecipeBook.Tags {
		if len(p.TaggedLines[tag]) == 0 {
			return errors.New("No lines for tag " + tname)
		}
	}
	return nil
}

func (p *Processor) ProvideLines(lines []string) {
	for _, line := range lines {
		if line == "" {
			continue
		}
		for _, tag := range p.RecipeBook.Tags {
			if !tag.Re.MatchString(line) {
				continue
			}
			p.accessTagged.Lock()
			//log.Printf("'%s' is good for %s", line, tag.Re.String())
			t := line
			p.TaggedLines[tag] = append(p.TaggedLines[tag], &t)
			p.accessTagged.Unlock()
		}
	}
}

func (p Processor) provideRndByTag(tag *Tag) *string {
	optCount := len(p.TaggedLines[tag])
	getId := rand.Intn(optCount)
	log.Printf("searching line by %s, have %d options, get by key %d", tag.Re.String(), optCount, getId)
	return p.TaggedLines[tag][getId]
}

func (p *Processor) Recipes() []string {
	return p.RecipeBook.RecipeNames
}

// Mix создает случайный куплет
func (p *Processor) Mix() (string, error) {
	recipeName := p.RecipeBook.RecipeNames[rand.Intn(len(p.RecipeBook.RecipeNames))]
	return p.Make(recipeName)
}

func (p *Processor) Make(recipeName string) (string, error) {
	recipe, exists := p.RecipeBook.Recipes[recipeName]
	if !exists {
		return "", errors.New("recipe not found, check /recipes")
	}

	result := ""
	for _, step := range recipe.Steps {
		if rand.Float32() > step.Probability {
			continue
		}
		// TODO: использовать полную логику проверки строк
		result += *p.provideRndByTag(step.Item.MustHave[0]) + "\n"
	}
	return result, nil
}
