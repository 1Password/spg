all: build

build:
	go build ./cmd/opgen

generate: agilewords.go agilesyllables.go Makefile

agilewords.go: testdata/agwordlist.txt goify_words.awk
	awk -f goify_words.awk < $< > $@
	gofmt -w $@

agilesyllables.go: testdata/agsyllables.txt goify_syllables.awk
	awk -f goify_syllables.awk < $< > $@
	gofmt -w $@
