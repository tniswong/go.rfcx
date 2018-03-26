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

func (p *parser) unscan() {
	p.buffer.unscanned = true
}

func (p *parser) scan() (token, string, error) {

	if p.buffer.unscanned {

		p.buffer.unscanned = false
		return p.buffer.token, p.buffer.literal, nil

	}

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

	var (
		keyToken    token
		key         string
		value       string
		valueRead   bool
		hasStar     bool
		quoteOpened bool
		quoteClosed bool
	)

	token, literal, err := p.scanIgnoreWhitespace()

	if err != nil {
		return INVALID, "", "", false, err
	}

	switch token {
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
		keyToken = token
		key = literal
	default:
		return INVALID, "", "", false, ErrInvalidLink
	}

starAndEquals:
	for {

		token, _, err := p.scanIgnoreWhitespace()

		if err != nil {
			return INVALID, "", "", false, err
		}

		switch token {
		case STAR: // optional
			hasStar = true
		case EQ:
			break starAndEquals
		default:
			return INVALID, "", "", false, ErrInvalidLink
		}

	}

valueLoop:
	for {

		token, literal, err := p.scanIgnoreWhitespace()
		if err != nil && value != "" {
			return INVALID, "", "", false, err
		}

		switch token {
		case QUOTE: // optional

			if quoteOpened {
				quoteClosed = true
				valueRead = true // empty counts as read
			}

			quoteOpened = true

		case WORD:

			if !valueRead {
				valueRead = true
				value = literal
				break valueLoop
			} else {
				return INVALID, "", "", false, ErrMissingSemicolon
			}

		default:

			if quoteOpened && !quoteClosed {
				return INVALID, "", "", false, ErrMissingClosingQuote
			} else if !valueRead {
				return INVALID, "", "", false, ErrMissingAttrValue
			} else if err == io.EOF {
				break valueLoop
			} else {
				return INVALID, "", "", false, ErrInvalidLink
			}

		}

	}

	if quoteOpened && !quoteClosed {

		if token, _, err := p.scanIgnoreWhitespace(); err != nil {
			return INVALID, "", "", false, err
		} else if token != QUOTE {
			return INVALID, "", "", false, ErrMissingClosingQuote
		}

		quoteClosed = true

	}

	if token, _, err := p.scanIgnoreWhitespace(); token != SEMICOLON && err != io.EOF {
		return INVALID, "", "", false, ErrMissingSemicolon
	} else if err != nil && err != io.EOF {
		return INVALID, "", "", false, err
	}

	return keyToken, key, value, hasStar, nil

}
