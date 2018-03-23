package rfc8288

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

// Token is an enum type to represent known lexer tokens for RFC8288
type Token int

const (
	INVALID Token = iota

	// delimiters
	QUOTE
	SEMICOLON
	LT
	GT
	EQ

	// special
	EOF
	STAR

	// multicharacter
	WORD
	WS

	// reserved attribute names
	REL
	HREFLANG
	MEDIA
	TITLE
	TYPE
)

// Scanner is a lexer for rfc8288
type Scanner struct {
	reader   bufio.Reader
	lastRead Token

	quoteOpen   bool
	bracketOpen bool
}

// NewScanner
func NewScanner(reader io.Reader) Scanner {
	return Scanner{reader: *bufio.NewReader(reader)}
}

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

// read the next rune from the buffered reader (io.EOF is returned as error on attempting to read EOF)
func (s *Scanner) read() (rune, error) {
	r, _, err := s.reader.ReadRune()
	return r, err
}

// unread the last rune from the reader
func (s *Scanner) unread() error {
	return s.reader.UnreadRune()
}

// Scan returns the next token and literal, or error
func (s *Scanner) Scan() (token Token, literal string, err error) {

	// read
	if r, err := s.read(); err != nil { // eof

		return EOF, "", io.EOF

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

// scanWhitespace scans for contiguous whitespace
func (s *Scanner) scanWhitespace() (token Token, literal string, err error) {

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
func (s *Scanner) scanWord() (token Token, literal string, err error) {

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

		} else if isStar(r) && s.lastRead != QUOTE && s.lastRead != LT { // if star, and not in quotes or chevrons

			// unread and break
			if err := s.unread(); err != nil {
				return INVALID, "", err
			}

			break

		} else { // otherwise write the rune to buf

			if _, err := buf.WriteRune(r); err != nil {
				return INVALID, "", err
			}

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
func (s *Scanner) scanned(t Token, literal string, err error) (Token, string, error) {
	s.lastRead = t
	return t, literal, err
}
