package rfc8288

import (
	"errors"
	"net/url"
)

// parse Error Types
var (
	ErrInvalidLink         = errors.New("rfc8288: invalid link")
	ErrMissingSemicolon    = errors.New("rfc8288: invalid link, missing semicolon")
	ErrMissingClosingQuote = errors.New("rfc8288: invalid link, missing closing quote")
	ErrMissingAttrValue    = errors.New("rfc8288: invalid link, missing attribute value")
)

type parser struct {
	scanner scanner
}

func (p *parser) scan() (token, string, error) {
	return p.scanner.Scan()
}

func (p *parser) scanIgnoreWhitespace() (token, string, error) {

	token, literal, err := p.scan()

	if token == WS {
		return p.scan()
	}

	return token, literal, err

}

func (p parser) parse() (Link, error) {

	var result = Link{}
	href, err := p.href()

	if err != nil {
		return Link{}, err
	}

	result.HREF = href

loop:
	for {

		token, key, value, hasStar, err := p.attribute()

		if err != nil {
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
			break loop
		default:
			return Link{}, ErrInvalidLink
		}

	}

	return result, nil

}

func (p *parser) href() (url.URL, error) {

	var uri *url.URL

	token, _, err := p.scanIgnoreWhitespace()

	if err != nil {
		return url.URL{}, err
	}

	if token != LT {
		return url.URL{}, ErrInvalidLink
	}

	token, literal, err := p.scanIgnoreWhitespace()

	if err != nil {
		return url.URL{}, err
	}

	if token != WORD {
		return url.URL{}, ErrInvalidLink
	}

	uri, err = url.Parse(literal)

	if err != nil {
		return url.URL{}, err
	}

	token, _, err = p.scanIgnoreWhitespace()

	if err != nil {
		return url.URL{}, err
	}

	if token != GT {
		return url.URL{}, ErrInvalidLink
	}

	token, _, err = p.scanIgnoreWhitespace()

	if err != nil {
		return url.URL{}, err
	}

	if token != SEMICOLON && token != EOF {
		return url.URL{}, ErrMissingSemicolon
	}

	return *uri, nil

}

func (p *parser) attribute() (token, string, string, bool, error) {

	token, key, hasStar, err := p.attributeKey()

	if err != nil {
		return INVALID, "", "", false, err
	}

	if token == EOF {
		return EOF, "", "", false, nil
	}

	value, err := p.attributeValue()

	if err != nil {
		return INVALID, "", "", false, err
	}

	return token, key, value, hasStar, nil

}

func (p parser) isValidAttributeKey(t token) bool {

	switch t {
	case REL, HREFLANG, MEDIA, TITLE, TYPE, WORD:
		return true
	default:
		return false
	}

}

func (p *parser) attributeKey() (token, string, bool, error) {

	token, literal, err := p.scanIgnoreWhitespace()

	if err != nil {
		return INVALID, "", false, err
	}

	if token == EOF {
		return EOF, literal, false, nil
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

func (p *parser) attributeValue() (string, error) {

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

			}

			// otherwise, we're missing a semicolon
			return "", ErrMissingSemicolon

		} else { // anything else

			// if the value hasn't been read
			if !valueRead {

				// we're missing the attribute value
				return "", ErrMissingAttrValue

			}

			// if we're at EOF
			if token == EOF {
				// we're done
				break
			}

			// otherwise, the link is invalid
			return "", ErrInvalidLink

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
		}

		if token != QUOTE { // otherwise if the token is not a QUOTE, we're missing a closing quote
			return ErrMissingClosingQuote
		}

	}

	return nil

}

func (p *parser) verifyAttributeTerminatedOrEOF() error {

	// scan for the next non-whitespace token
	token, _, err := p.scanIgnoreWhitespace()

	if err != nil {
		return err
	}

	if token != SEMICOLON && token != EOF {
		return ErrMissingSemicolon
	}

	return nil

}
