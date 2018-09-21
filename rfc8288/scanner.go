package rfc8288

import (
	"bytes"
	"io"
	"unicode"
)

// token is an enum type to represent known lexer tokens for RFC8288
type token int

// token types
const (
	INVALID token = iota

	// delimiters
	QUOTE
	SEMICOLON
	LT
	GT
	EQ

	// special
	EOF
	STAR

	// multi-character
	WORD
	WS

	// reserved attribute names
	REL
	HREFLANG
	MEDIA
	TITLE
	TYPE
)

// isWhitespace returns true if rune a unicode whitespace character
func isWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}

// isSymbol returns true if rune is a link symbol
func isSymbol(r rune) bool {
	return r == '"' || r == ';' || r == '<' || r == '>' || r == '='
}

// isStar returns true if rune is an asterisk
func isStar(r rune) bool {
	return r == '*'
}

// scanner is a lexer for rfc8288
type scanner struct {
	runeScanner io.RuneScanner
	lastRead    token

	quoteOpen   bool
	bracketOpen bool
}

// scan returns the next token and literal, or error
func (s *scanner) Scan() (token, string, error) {

	// read
	if r, err := s.read(); err != nil { // eof

		return s.scanned(EOF, "", nil)

	} else if isWhitespace(r) { // is whitespace?

		// unread that whitespace, we'll capture it in s.scanWhitespace()
		if err := s.unread(); err != nil {
			return INVALID, "", err
		}

		// scan all the contiguous whitespace
		return s.scanWhitespace()

	} else if isStar(r) && s.lastRead != QUOTE && s.lastRead != LT { // if r is '*', but the lastRead tokens aren't '"' or '<'

		// then we scanned a STAR token
		return s.scanned(STAR, string(r), nil)

	} else if isSymbol(r) { // is a Symbol token?

		// which symbol?
		switch r {
		case '"':
			s.quoteOpen = !s.quoteOpen
			return s.scanned(QUOTE, string(r), nil)
		case ';':
			return s.scanned(SEMICOLON, string(r), nil)
		case '<':
			s.bracketOpen = true
			return s.scanned(LT, string(r), nil)
		case '>':
			s.bracketOpen = false
			return s.scanned(GT, string(r), nil)
		case '=':
			return s.scanned(EQ, string(r), nil)
		}

	}

	// neither whitespace, star, nor symbol. unread this rune, we'll capture it in s.scanWord()
	if err := s.unread(); err != nil {
		return INVALID, "", err
	}

	// keep scanning as a word
	return s.scanWord()

}

// read the next rune from the buffered runeScanner
func (s *scanner) read() (rune, error) {
	r, _, err := s.runeScanner.ReadRune()
	return r, err
}

// unread the last rune from the runeScanner
func (s *scanner) unread() error {
	return s.runeScanner.UnreadRune()
}

// scanWhitespace scans for contiguous whitespace
func (s *scanner) scanWhitespace() (token, string, error) {

	// buf is a place to store the contiguous whitespace
	var buf bytes.Buffer

	for {

		// read
		r, err := s.read()

		if err != nil { // if eof

			// eof
			break

		}

		if !isWhitespace(r) { // if not whitespace

			// unread the last rune
			err := s.unread()

			if err != nil {
				return INVALID, "", err
			}

			break

		}

		// its whitespace, write to buf
		_, err = buf.WriteRune(r)

		if err != nil {
			return INVALID, "", err
		}

	}

	// scanned WS
	return s.scanned(WS, buf.String(), nil)

}

// scanWord scans for continuous word runes
func (s *scanner) scanWord() (token, string, error) {

	// buf is a place to store the contiguous word runes
	var buf bytes.Buffer

	for {

		// read
		r, err := s.read()

		if err != nil {

			// eof
			break

		}

		if isSymbol(r) {

			err := s.unread()

			if err != nil {
				return INVALID, "", err
			}

			break

		}

		if isStar(r) && s.lastRead != QUOTE && s.lastRead != LT { // if star, and not in quotes or chevrons

			err := s.unread()

			if err != nil {
				return INVALID, "", err
			}

			break

		}

		// otherwise write the rune to buf
		_, err = buf.WriteRune(r)

		if err != nil {
			return INVALID, "", err
		}

	}

	// as long as we're not in quotes or chevrons
	if s.lastRead != QUOTE && s.lastRead != LT {

		// match for known attribute names then return the corresponding token.
		switch buf.String() {
		case "rel":
			return s.scanned(REL, buf.String(), nil)
		case "hreflang":
			return s.scanned(HREFLANG, buf.String(), nil)
		case "media":
			return s.scanned(MEDIA, buf.String(), nil)
		case "title":
			return s.scanned(TITLE, buf.String(), nil)
		case "type":
			return s.scanned(TYPE, buf.String(), nil)
		}

	}

	// scanned a WORD.
	return s.scanned(WORD, buf.String(), nil)

}

// scanned tells the scanner what we've just scanned. the error parameter is passthrough as a convenience
func (s *scanner) scanned(t token, literal string, err error) (token, string, error) {
	s.lastRead = t
	return t, literal, err
}
