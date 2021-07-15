package lrc

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"
)

// NOTE: The parser simply ignores all unmatched lines and doen't return any parse errors.

type Duration = time.Duration
type Text = string

// Parse parses a byte slice of LRC lyrics.
func Parse(data []byte) ([]Duration, []Text, error) {
	return NewParser(bytes.NewReader(data)).Parse()
}

// ParseString parses a string of LRC lyrics.
func ParseString(text string) ([]Duration, []Text, error) {
	return NewParser(strings.NewReader(text)).Parse()
}

// Perser is a parser type.
type Parser struct {
	reader io.Reader
}

// NewParser return a new parser from a reader.
func NewParser(r io.Reader) *Parser {
	return &Parser{r}
}

// ll    -> [00:00.00]
// index -> 0123456789

func dtoi(b byte) byte           { return b - '0' }
func ddtoD(b1, b2 byte) Duration { return Duration(10*dtoi(b1) + dtoi(b2)) }
func isdd(b1, b2 byte) bool      { return '0' <= b1 && b1 <= '9' && '0' <= b2 && b2 <= '9' }

// Parse parses the reader according to the LRC format.
// https://en.wikipedia.org/wiki/LRC_(file_format)
func (p *Parser) Parse() ([]Duration, []Text, error) {
	times := make([]Duration, 0)
	lines := make([]Text, 0)
	scnnr := bufio.NewScanner(p.reader)
	//scnnr.Split(bufio.ScanLines)
	for rp := 0; scnnr.Scan(); {
		ll := scnnr.Bytes()
		// [00:00.00][00:00.00]text -> [00:00.00]text -> text
		for len(ll) >= 10 && ll[0] == '[' && isdd(ll[1], ll[2]) && ll[3] == ':' && isdd(ll[4], ll[5]) && ll[6] == '.' && isdd(ll[7], ll[8]) && ll[9] == ']' {
			times = append(times, (ddtoD(ll[1], ll[2])*time.Minute + ddtoD(ll[4], ll[5])*time.Second + ddtoD(ll[7], ll[8])*time.Second/100))
			ll = ll[10:]
			rp++
		}
		for ll := string(ll); rp > 0; rp-- {
			lines = append(lines, ll)
		}
	}
	if err := scnnr.Err(); err != nil {
		return nil, nil, fmt.Errorf("LRC parse: %w", err)
	}
	return times, lines, nil
}
