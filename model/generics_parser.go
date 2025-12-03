package model

import (
	"fmt"
	"strings"
	"unicode"
)

type genericsParser struct {
	input string
	pos   int
}

func newGenericsParser(input string) *genericsParser {
	return &genericsParser{
		input: input,
		pos:   0,
	}
}

func (p *genericsParser) eof() bool {
	return p.pos >= len(p.input)
}

func (p *genericsParser) peek() rune {
	if p.eof() {
		return 0
	}
	r, _ := runeAt(p.input, p.pos)
	return r
}

func (p *genericsParser) next() rune {
	if p.eof() {
		return 0
	}
	r, size := runeAt(p.input, p.pos)
	p.pos += size
	return r
}

func (p *genericsParser) remaining() string {
	if p.eof() {
		return ""
	}
	return p.input[p.pos:]
}

func (p *genericsParser) skipSpaces() {
	for !p.eof() {
		r, size := runeAt(p.input, p.pos)
		if !unicode.IsSpace(r) {
			return
		}
		p.pos += size
	}
}

// parseType := identifier [ "<" type { "," type } ">" ]
func (p *genericsParser) parseType() (*TypeNode, error) {
	p.skipSpaces()
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	node := &TypeNode{Name: name}

	p.skipSpaces()
	if p.peek() != '<' {
		// no generics
		return node, nil
	}

	// parse generic arguments
	p.next() // consume '<'
	p.skipSpaces()

	for {
		arg, err := p.parseType()
		if err != nil {
			return nil, err
		}
		node.Args = append(node.Args, arg)

		p.skipSpaces()
		ch := p.peek()
		if ch == ',' {
			p.next()
			p.skipSpaces()
			continue
		}
		if ch == '>' {
			p.next() // consume '>'
			break
		}
		if ch == 0 {
			return nil, fmt.Errorf("unexpected end of input while parsing generic args")
		}
		return nil, fmt.Errorf("unexpected character %q at position %d", ch, p.pos)
	}

	return node, nil
}

// identifier is any run of non-space, non-delimiter characters
// delimiters: '<', '>', ',', whitespace
func (p *genericsParser) parseIdentifier() (string, error) {
	p.skipSpaces()
	start := p.pos
	for !p.eof() {
		r, size := runeAt(p.input, p.pos)
		if unicode.IsSpace(r) || strings.ContainsRune("<>,", r) {
			break
		}
		p.pos += size
	}
	if p.pos == start {
		return "", fmt.Errorf("expected identifier at position %d", p.pos)
	}
	return p.input[start:p.pos], nil
}

func runeAt(s string, i int) (rune, int) {
	if i >= len(s) {
		return 0, 0
	}
	r, size := rune(s[i]), 1
	if r >= 0x80 {
		// fall back to proper rune decoding
		r, size = utf8DecodeRuneInString(s[i:])
	}
	return r, size
}

// minimal utf8 decoding (we could also just use utf8.DecodeRuneInString from "unicode/utf8")
func utf8DecodeRuneInString(s string) (rune, int) {
	// In real code, just use utf8.DecodeRuneInString.
	// Here we keep it explicit.
	for i := 1; i <= len(s); i++ {
		r := []rune(s[:i])
		if len(r) == 1 {
			return r[0], i
		}
	}
	return rune(s[0]), 1
}
