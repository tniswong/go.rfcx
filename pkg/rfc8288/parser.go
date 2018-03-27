package rfc8288

import (
	"errors"
	"io"
	"net/url"
)

// Parse Error Types
var (
	ErrInvalidLink         = errors.New("[ERR] invalid link")
	ErrMissingSemicolon    = errors.New("[ERR] invalid link: missing semicolon")
	ErrMissingClosingQuote = errors.New("[ERR] invalid link: missing closing quote")
	ErrMissingAttrValue    = errors.New("[ERR] invalid link: missing attribute value")
)

type parser struct {
	scanner scanner
	buffer  struct {
		token     token
		literal   string
		unscanned bool
	}
}

func newParser(reader io.Reader) parser {
	return parser{scanner: newScanner(reader)}
}

func (p *parser) scan() (token, string, error) {

	token, literal, err := p.scanner.Scan()

	p.buffer.token = token
	p.buffer.literal = literal

	return token, literal, err
}

func (p *parser) scanIgnoreWhitespace() (token, string, error) {

	token, literal, err := p.scan()

	if token == WS {
		token, literal, err = p.scan()
	}

	return token, literal, err

}

func (p parser) Parse() (Link, error) {

	var result = Link{}

	href, err := p.parseHREF()

	if err != nil {
		return Link{}, err
	}

	result.HREF = href

	for {

		token, key, value, hasStar, err := p.parseAttribute()

		if err == io.EOF {
			break
		} else if err != nil {
			return Link{}, err
		}

		switch token {
		case REL:
			result.Rel = value
		case HREFLANG:
			result.HREFLang = value
		case MEDIA:
			result.Media = value
		case TITLE:

			if hasStar {
				result.TitleStar = value
			} else {
				result.Title = value
			}

		case TYPE:
			result.Type = value
		case WORD:
			result.Extend(key, value)
		case EOF:
			return result, nil
		default:
			return Link{}, ErrInvalidLink
		}

	}

	return result, nil

}

func (p *parser) parseHREF() (url.URL, error) {

	var uri *url.URL

	if token, _, err := p.scanIgnoreWhitespace(); err != nil {
		return url.URL{}, err
	} else if token != LT {
		return url.URL{}, ErrInvalidLink
	}

	if token, literal, err := p.scanIgnoreWhitespace(); err != nil {
		return url.URL{}, err
	} else if token != WORD {
		return url.URL{}, ErrInvalidLink
	} else {

		uri, err = url.Parse(literal)

		if err != nil {
			return url.URL{}, err
		}

	}

	if token, _, err := p.scanIgnoreWhitespace(); err != nil {
		return url.URL{}, err
	} else if token != GT {
		return url.URL{}, ErrInvalidLink
	}

	if token, _, err := p.scanIgnoreWhitespace(); token != SEMICOLON && err != io.EOF {
		return url.URL{}, ErrMissingSemicolon
	} else if err != nil && err != io.EOF {
		return url.URL{}, err
	}

	return *uri, nil

}

func (p *parser) parseAttribute() (token, string, string, bool, error) {

	keyToken, key, hasStar, err := p.scanAttributeKey()

	if err != nil {
		return INVALID, "", "", false, err
	}

	value, err := p.scanAttributeValue()

	if err != nil {
		return INVALID, "", "", false, err
	}

	return keyToken, key, value, hasStar, nil

}

func (p parser) isValidAttributeKey(t token) bool {

	switch t {
	case REL:
		fallthrough
	case HREFLANG:
		fallthrough
	case MEDIA:
		fallthrough
	case TITLE:
		fallthrough
	case TYPE:
		fallthrough
	case WORD:
		return true
	default:
		return false
	}

}

func (p *parser) scanAttributeKey() (token, string, bool, error) {

	token, literal, err := p.scanIgnoreWhitespace()

	if err != nil {
		return INVALID, "", false, err
	}

	if !p.isValidAttributeKey(token) {
		return INVALID, "", false, ErrInvalidLink
	}

	var (
		hasStar  bool
		keyToken = token
		key      = literal
	)

	for {

		token, _, err := p.scanIgnoreWhitespace()

		if err != nil {
			return INVALID, "", false, err
		}

		if token == STAR {
			hasStar = true
			continue
		}

		if token != EQ {
			return INVALID, "", false, ErrInvalidLink
		}

		break

	}

	return keyToken, key, hasStar, nil

}

func (p *parser) scanAttributeValue() (string, error) {

	var (
		valueRead   bool
		quoteOpened bool
		quoteClosed bool
		value       string
	)

	for {

		token, literal, err := p.scanIgnoreWhitespace()

		if err != nil && value != "" {
			return "", err
		}

		// optional quote
		if token == QUOTE {

			// if we've already scanned a quote
			if quoteOpened {

				// then this closes the quote
				quoteClosed = true

				// and we've read the value
				valueRead = true

			}

			quoteOpened = true

		} else if token == WORD { // value word

			// if we've not already read the value
			if !valueRead {

				// then we have now
				valueRead = true

				// and this is it
				value = literal

				break

			} else { // otherwise, we're missing a semicolon
				return "", ErrMissingSemicolon
			}

		} else { // anything else

			// if the value hasn't been read
			if !valueRead {

				// we're missing the attribute value
				return "", ErrMissingAttrValue

			} else if err == io.EOF { // if we're at EOF

				// we're done
				break

			} else { // otherwise, the link is invalid
				return "", ErrInvalidLink
			}

		}

	}

	if err := p.verifyQuoteTerminated(quoteOpened, quoteClosed); err != nil {
		return "", err
	}

	if err := p.verifyAttributeTerminatedOrEOF(); err != nil {
		return "", err
	}

	return value, nil

}

func (p *parser) verifyQuoteTerminated(quoteOpened bool, quoteClosed bool) error {

	// if the quote is unterminated
	if quoteOpened && !quoteClosed {

		// scan for the next non-whitespace token
		token, _, err := p.scanIgnoreWhitespace()

		// if there's an error
		if err != nil {
			return err
		} else if token != QUOTE { // otherwise if the token is not a QUOTE, we're missing a closing quote
			return ErrMissingClosingQuote
		}

	}

	return nil

}

func (p *parser) verifyAttributeTerminatedOrEOF() error {

	// scan for the next non-whitespace token
	token, _, err := p.scanIgnoreWhitespace()

	if token != SEMICOLON && err != io.EOF {
		return ErrMissingSemicolon
	} else if err != nil && err != io.EOF {
		return err
	}

	return nil

}
