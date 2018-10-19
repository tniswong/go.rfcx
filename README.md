# go.rfcx

Go Library with implementations of RFC's that I find useful.

## Run Tests

    $> make deps   # only necessary once
    $> make ensure # only necessary once
    $> make test

## Make targets

1. `make clean`

    Deletes leftover `.coverprofile` files.

1. `make doc`

    Starts a `godoc` server for this package.

1. `make deps`

    Install all dependent cli's for these make targets. Run this first, at least once!

1. `make ensure`

    Ensure all runtime dependencies are installed properly.

1. `make fmt` or `make format`

    Automatically format all code in this package.

1. `make vet`

    Run `go vet` on all code in this package, excluding dependencies. Exit 0, if successful. Exit 1, if not.

1. `make lint`

    Run `go lint` on all code in this package, excluding dependencies. Exit 0, if successful. Exit 1, if not.

1. `make complexity`

    Generate a complexity report for all code in this package, excluding dependencies. Exit 0, if reported complexity is
    above maximum threshold. Exit 1, if not.

1. `make coverage`

    Generate a coverage report for all code in this package, excluding dependencies. Exit 0, if reported coverage is
    below minimum threshold. Exit 1, if not.

1. `make test`

    Vet, Lint, Test with Coverage, and complexity. Exit 0, if successful. Exit 1, if there is unformatted code, if there
    are lint failures, if there are test failures, if coverage is below the minimum threshold, or if complexity is above
    the maximum threshold.

# RFC 7231 Hypertext Transfer Protocol (HTTP/1.1): Semantics and Content

https://tools.ietf.org/html/rfc7231

Specifically, parsing support for HTTP Accept Headers and handling of Media Ranges and supporting operations for
content negotiation as defined by this RFC.

> ### 5.3.2.  Accept
>
>    The "Accept" header field can be used by user agents to specify
>    response media types that are acceptable.  Accept header fields can
>    be used to indicate that the request is specifically limited to a
>    small set of desired types, as in the case of a request for an
>    in-line image.
>
>      Accept = #( media-range [ accept-params ] )
>
>      media-range    = ( "*/*"
>                       / ( type "/" "\*" )
>                       / ( type "/" subtype )
>                       ) *( OWS ";" OWS parameter )
>      accept-params  = weight *( accept-ext )
>      accept-ext = OWS ";" OWS token [ "=" ( token / quoted-string ) ]
>
>    The asterisk "\*" character is used to group media types into ranges,
>    with "\*/\*" indicating all media types and "type/\*" indicating all
>    subtypes of that type.  The media-range can include media type
>    parameters that are applicable to that range.
>
>    Each media-range might be followed by zero or more applicable media
>    type parameters (e.g., charset), an optional "q" parameter for
>    indicating a relative weight (Section 5.3.1), and then zero or more
>    extension parameters.  The "q" parameter is necessary if any
>    extensions (accept-ext) are present, since it acts as a separator
>    between the two parameter sets.
>
>       Note: Use of the "q" parameter name to separate media type
>       parameters from Accept extension parameters is due to historical
>       practice.  Although this prevents any media type parameter named
>       "q" from being used with a media range, such an event is believed
>       to be unlikely given the lack of any "q" parameters in the IANA
>       media type registry and the rare usage of any media type
>       parameters in Accept.  Future media types are discouraged from
>       registering any parameter named "q".
>
>    The example
>
>      Accept: audio/*; q=0.2, audio/basic
>
>    is interpreted as "I prefer audio/basic, but send me any audio type
>    if it is the best available after an 80% markdown in quality".
>
>    A request without any Accept header field implies that the user agent
>    will accept any media type in response.  If the header field is
>    present in a request and none of the available representations for
>    the response have a media type that is listed as acceptable, the
>    origin server can either honor the header field by sending a 406 (Not
>    Acceptable) response or disregard the header field by treating the
>    response as if it is not subject to content negotiation.
>
>    A more elaborate example is
>
>      Accept: text/plain; q=0.5, text/html,
>              text/x-dvi; q=0.8, text/x-c
>
>    Verbally, this would be interpreted as "text/html and text/x-c are
>    the equally preferred media types, but if they do not exist, then
>    send the text/x-dvi representation, and if that does not exist, send
>    the text/plain representation".
>
>    Media ranges can be overridden by more specific media ranges or
>    specific media types.  If more than one media range applies to a
>    given type, the most specific reference has precedence.  For example,
>
>      Accept: text/*, text/plain, text/plain;format=flowed, */*
>
>    have the following precedence:
>
>      1.  text/plain;format=flowed
>      2.  text/plain
>      3.  text/*
>      4.  */*
>
>    The media type quality factor associated with a given type is
>    determined by finding the media range with the highest precedence
>    that matches the type.  For example,
>
>      Accept: text/*;q=0.3, text/html;q=0.7, text/html;level=1,
>              text/html;level=2;q=0.4, */*;q=0.5
>
>    would cause the following values to be associated:
>
>    ```
>    +-------------------+---------------+
>    | Media Type        | Quality Value |
>    +-------------------+---------------+
>    | text/html;level=1 | 1             |
>    | text/html         | 0.7           |
>    | text/plain        | 0.3           |
>    | image/jpeg        | 0.5           |
>    | text/html;level=2 | 0.4           |
>    | text/html;level=3 | 0.7           |
>    +-------------------+---------------+
>    ```
>
>    Note: A user agent might be provided with a default set of quality
>    values for certain media ranges.  However, unless the user agent is a
>    closed system that cannot interact with other rendering agents, this
>    default set ought to be configurable by the user.

# RFC 7807 Problem Details

https://tools.ietf.org/html/rfc7807

> ### The Problem Details JSON Object
>
>
> The canonical model for problem details is a JSON [RFC7159] object.
>
>    When serialized as a JSON document, that format is identified with
>    the "application/problem+json" media type.
>
>    For example, an HTTP response carrying JSON problem details:
>
>    ```
>    HTTP/1.1 403 Forbidden
>    Content-Type: application/problem+json
>    Content-Language: en
>
>    {
>       "type": "https://example.com/probs/out-of-credit",
>       "title": "You do not have enough credit.",
>       "detail": "Your current balance is 30, but that costs 50.",
>       "instance": "/account/12345/msgs/abc",
>       "balance": 30,
>       "accounts": ["/account/12345", "/account/67890"]
>    }
>    ```
>
>    Here, the out-of-credit problem (identified by its type URI)
>    indicates the reason for the 403 in "title", gives a reference for
>    the specific problem occurrence with "instance", gives occurrence-
>    specific details in "detail", and adds two extensions; "balance"
>    conveys the account's balance, and "accounts" gives links where the
>    account can be topped up.
>
>    The ability to convey problem-specific extensions allows more than
>    one problem to be conveyed.  For example:
>
>    ```
>    HTTP/1.1 400 Bad Request
>    Content-Type: application/problem+json
>    Content-Language: en
>
>    {
>       "type": "https://example.net/validation-error",
>       "title": "Your request parameters didn't validate.",
>       "invalid-params": [{
>         "name": "age",
>         "reason": "must be a positive integer"
>       }, {
>         "name": "color",
>         "reason": "must be 'green', 'red' or 'blue'"
>       }]
>    }
>    ```
>
>    Note that this requires each of the subproblems to be similar enough
>    to use the same HTTP status code.  If they do not, the 207 (Multi-
>    Status) [RFC4918] code could be used to encapsulate multiple status
>    messages.

# RFC 8288 Web Linking

https://tools.ietf.org/html/rfc8288