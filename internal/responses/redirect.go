package responses

import (
	"bufio"
	"errors"
	"io"
	"net/url"
	"strconv"
)

var (
	ErrRedirectResponseHeader     = errors.New("response header is invalid")
	ErrRedirectResponseStatusCode = errors.New("invalid success response status code")
	ErrRedirectResponseSeparator  = errors.New("invalid success response separator")
	ErrRedirectResponseURI        = errors.New("invalid success response URI")
	ErrRedirectResponseCR         = errors.New("invalid success response CR")
)

type RedirectResponse struct {
	StatusCode int
	URI        string
}

func processRedirectResponse(response io.Reader) (RedirectResponse, error) {
	reader := bufio.NewReader(response)
	responseHeader, err := reader.ReadBytes('\n')
	if err != nil {
		return RedirectResponse{}, err
	}

	if len(responseHeader) < 2 {
		return RedirectResponse{}, ErrRedirectResponseHeader
	}

	statusCode, err := strconv.Atoi(string(responseHeader[0:2]))
	if err != nil {
		return RedirectResponse{}, err
	}

	if statusCode < 30 || statusCode > 40 {
		return RedirectResponse{}, ErrRedirectResponseStatusCode
	} else if statusCode > 31 {
		statusCode = 30
	}

	if responseHeader[2] != ' ' {
		return RedirectResponse{}, ErrRedirectResponseSeparator
	}

	if responseHeader[len(responseHeader)-2] != '\r' {
		return RedirectResponse{}, ErrRedirectResponseCR
	}

	if responseHeader[len(responseHeader)-1] != '\n' {
		return RedirectResponse{}, ErrRedirectResponseCR
	}

	rawURI := responseHeader[3 : len(responseHeader)-2]
	newURI, err := url.Parse(string(rawURI))
	if err != nil {
		return RedirectResponse{}, ErrRedirectResponseURI
	}

	// TODO: Check if newURI is a valid Gemini URI:
	// - The scheme is "gemini"
	// - There is no user info
	// - No IP address in Authority

	return RedirectResponse{
		StatusCode: statusCode,
		URI:        newURI.String(),
	}, nil
}
