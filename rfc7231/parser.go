package rfc7231

import (
	"errors"
	"io"
	"strconv"
)

// Parsing Errors for Accept
var (
	ErrInvalidMediaRange         = errors.New("rfc7231: invalid media range")
	ErrQMustBeNumberBetween0And1 = errors.New("rfc7231: invalid media range: Q must be a number between 0 and 1")
)

type parser struct {
	scanner scanner
	buffer  struct {
		token     token
		literal   string
		unscanned bool
	}
}

func (p *parser) scan() (token, string, error) {

	if p.buffer.unscanned {
		p.buffer.unscanned = false
		return p.buffer.token, p.buffer.literal, nil
	}

	token, literal, err := p.scanner.scan()

	p.buffer.token = token
	p.buffer.literal = literal

	return token, literal, err

}

func (p *parser) scanIgnoreWhitespace() (token, string, error) {

	token, literal, err := p.scan()

	if token == WS {
		return p.scan()
	}

	return token, literal, err

}

func (p *parser) unscan() {
	p.buffer.unscanned = true
}

func (p parser) parse() (Accept, error) {

	mediaRanges, err := p.parseMediaRanges()

	if err != nil {
		return Accept{}, err
	}

	var result = Accept{
		MediaRanges: mediaRanges,
	}

	return result, nil

}

func (p *parser) parseMediaRanges() ([]MediaRange, error) {

	var result []MediaRange

	for {

		typeName, subtypeName, err := p.mediaRange()

		if err == io.EOF {
			break
		} else if err != nil {
			return []MediaRange{}, err
		}

		mr := MediaRange{
			TypeName:    typeName,
			SubtypeName: subtypeName,
		}

		params, err := p.params()

		if err != nil {
			return []MediaRange{}, err
		}

		if q, ok := params["q"]; ok {

			qf, err := strconv.ParseFloat(q, 64)

			if err != nil {
				return []MediaRange{}, ErrQMustBeNumberBetween0And1
			}

			mr.Q = qf
			delete(params, "q")

		}

		mr.Params = params
		result = append(result, mr)

	}

	return result, nil

}

func (p *parser) mediaRange() (string, string, error) {

	token, typeName, err := p.scanIgnoreWhitespace()

	if err != nil {
		return "", "", err
	}

	if token == EOF {
		return "", "", io.EOF
	}

	if token == COMMA {

		token, typeName, err = p.scanIgnoreWhitespace()

		if err != nil {
			return "", "", err
		}

	}

	if token != WORD {
		return "", "", ErrInvalidMediaRange
	}

	token, _, err = p.scanIgnoreWhitespace()

	if err != nil {
		return "", "", err
	}

	if token != SLASH {
		return "", "", ErrInvalidMediaRange
	}

	token, subtypeName, err := p.scanIgnoreWhitespace()

	if err != nil {
		return "", "", err
	}

	if token != WORD {
		return "", "", ErrInvalidMediaRange
	}

	// */subtype is invalid. only type/* is allowed.
	if typeName == "*" && subtypeName != "*" {
		return "", "", ErrInvalidMediaRange
	}

	return typeName, subtypeName, nil

}

func (p *parser) params() (map[string]string, error) {

	result := map[string]string{}

	for {

		token, _, err := p.scanIgnoreWhitespace()

		if err != nil {
			return map[string]string{}, err
		}

		if token == EOF {
			break
		}

		if token == COMMA {
			p.unscan()
			break
		}

		if token != SEMICOLON {
			return map[string]string{}, ErrInvalidMediaRange
		}

		token, key, err := p.scanIgnoreWhitespace()

		if err != nil {
			return map[string]string{}, err
		}

		token, _, err = p.scanIgnoreWhitespace()

		if err != nil {
			return map[string]string{}, err
		}

		if token != EQ {
			return map[string]string{}, ErrInvalidMediaRange
		}

		token, value, err := p.scanIgnoreWhitespace()

		if err != nil {
			return map[string]string{}, err
		}

		result[key] = value

	}

	return result, nil

}
