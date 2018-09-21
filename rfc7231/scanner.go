package rfc7231

import (
	"bytes"
	"io"
	"unicode"
)

type token int

// token types
const (
	INVALID token = iota

	// symbols
	SLASH
	SEMICOLON
	EQ
	COMMA

	// multi-character
	WORD
	WS

	// special
	EOF
)

func isWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}

func isSymbol(r rune) bool {
	return r == '/' || r == ';' || r == '=' || r == ','
}

type scanner struct {
	runeScanner io.RuneScanner
	lastRead    token
}

func (s *scanner) scan() (token, string, error) {

	if r, err := s.read(); err != nil { // eof

		return s.scanned(EOF, "", nil)

	} else if isWhitespace(r) { // is whitespace?

		// unread that whitespace, we'll capture it in s.scanWhitespace()
		if err := s.unread(); err != nil {
			return s.scanned(INVALID, string(r), err)
		}

		// scan all the contiguous whitespace
		return s.scanWhitespace()

	} else if isSymbol(r) { // is a Symbol token?

		// which symbol?
		switch r {
		case '/':
			return s.scanned(SLASH, string(r), nil)
		case ';':
			return s.scanned(SEMICOLON, string(r), nil)
		case '=':
			return s.scanned(EQ, string(r), nil)
		case ',':
			return s.scanned(COMMA, string(r), nil)
		}

	}

	// neither whitespace nor comma. unread this rune, we'll capture it in s.scanWord()
	if err := s.unread(); err != nil {
		return s.scanned(INVALID, "", err)
	}

	// keep scanning as a word
	return s.scanWord()

}

func (s *scanner) read() (rune, error) {
	r, _, err := s.runeScanner.ReadRune()
	return r, err
}

func (s *scanner) unread() error {
	return s.runeScanner.UnreadRune()
}

// scanWhitespace scans for contiguous whitespace
func (s *scanner) scanWhitespace() (token, string, error) {

	// buf is a place to store the contiguous whitespace
	var buf bytes.Buffer

	for {

		// read
		if r, err := s.read(); err != nil { // if eof

			// eof
			break

		} else if !isWhitespace(r) { // if not whitespace

			// unread the last rune
			if err := s.unread(); err != nil {
				return INVALID, "", err
			}

			break

		} else { // is whitespace

			// write to buf
			if _, err := buf.WriteRune(r); err != nil {
				return INVALID, "", err
			}

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
		if r, err := s.read(); err != nil { // if eof

			// eof
			break

		} else if isSymbol(r) { // if symbol

			// unread and break
			if err := s.unread(); err != nil {
				return INVALID, "", err
			}

			break

		} else { // otherwise write the rune to buf

			// bytes.Buffer WriteRun always returns nil error.
			// no need to handle err != nil case here
			buf.WriteRune(r)

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
