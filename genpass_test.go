package spg

import (
	"testing"
)

// Toy example of how any recipe can satisfy Generator
func TestGeneratorInterface(t *testing.T) {
	var (
		g Generator
		p *Password
	)

	g = createSizeFourGenerator("chars")
	p, err := g.Generate()
	if err != nil {
		t.Errorf("Generation error: %v", err)
	}

	if len(p.String()) != 4 {
		t.Errorf("Could not create a char passsword using the Generator interface: %s", p.String())
	}

	g = createSizeFourGenerator("words")
	p, _ = g.Generate()
	if len(p.Tokens().Atoms()) != 4 {
		t.Errorf("Could not create a word passsword using the Generator interface: %s", p.String())
	}
}

func createSizeFourGenerator(name string) Generator {
	if name == "words" {
		wl, _ := NewWordList([]string{"one", "two", "three"})
		r := NewWLRecipe(4, wl)
		r.SeparatorChar = "-"
		return r
	}

	r := &CharRecipe{Length: 4, Allow: Lowers}
	return r
}
