package responses

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"mime"
	"slices"
	"strconv"
)

var (
	allowedMIMETypes = []string{"text/gemini", "text/plain"}

	ErrSuccessResponseHeader     = errors.New("response header is invalid")
	ErrSuccessResponseStatusCode = errors.New("invalid success response status code")
	ErrSuccessResponseSeparator  = errors.New("invalid success response separator")
	ErrSuccessResponseMIMEType   = errors.New("invalid success response MIME type")
	ErrSuccessResponseReadBody   = errors.New("could not read success response body")
	ErrSuccessResponseCR         = errors.New("invalid success response CR")
)

type SuccessResponse struct {
	StatusCode     int
	rawBody        []byte
	Body           io.Reader
	MIMEType       string
	MIMETypeParams map[string]string
}

func processSuccessResponse(response io.Reader) (SuccessResponse, error) {
	reader := bufio.NewReader(response)
	responseHeader, err := reader.ReadBytes('\n')
	if err != nil {
		return SuccessResponse{}, err
	}

	if len(responseHeader) < 2 {
		return SuccessResponse{}, ErrSuccessResponseHeader
	}

	statusCode, err := strconv.Atoi(string(responseHeader[0:2]))
	if err != nil {
		return SuccessResponse{}, err
	}

	if statusCode < 20 || statusCode > 30 {
		return SuccessResponse{}, ErrSuccessResponseStatusCode
	} else if statusCode > 20 {
		statusCode = 20
	}

	if responseHeader[2] != ' ' {
		return SuccessResponse{}, ErrSuccessResponseSeparator
	}

	if responseHeader[len(responseHeader)-2] != '\r' {
		return SuccessResponse{}, ErrSuccessResponseCR
	}

	if responseHeader[len(responseHeader)-1] != '\n' {
		return SuccessResponse{}, ErrSuccessResponseCR
	}

	rawMIMEType := responseHeader[3 : len(responseHeader)-2]
	mimeType, params, err := mime.ParseMediaType(string(rawMIMEType))
	if err != nil || !slices.Contains(allowedMIMETypes, mimeType) {
		return SuccessResponse{}, ErrSuccessResponseMIMEType
	}

	rawBody, err := io.ReadAll(reader)
	if err != nil {
		return SuccessResponse{}, ErrSuccessResponseReadBody
	}

	return SuccessResponse{
		StatusCode:     statusCode,
		rawBody:        rawBody,
		Body:           bytes.NewReader(rawBody),
		MIMEType:       mimeType,
		MIMETypeParams: params,
	}, nil
}
