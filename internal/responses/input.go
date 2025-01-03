package responses

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

var (
	ErrInputResponseHeader     = errors.New("response header is invalid")
	ErrInputResponseStatusCode = errors.New("invalid input response status code")
	ErrInputResponseSeparator  = errors.New("invalid input response separator")
	ErrInputResponseCR         = errors.New("invalid input response CR")
)

type InputResponse struct {
	StatusCode int
	Prompt     string
}

func processInputResponse(response io.Reader) (InputResponse, error) {
	reader := bufio.NewReader(response)
	responseHeader, err := reader.ReadString('\n')
	if err != nil {
		return InputResponse{}, err
	}

	if len(responseHeader) < 2 {
		return InputResponse{}, ErrInputResponseHeader
	}

	statusCode, err := strconv.Atoi(responseHeader[0:2])
	if err != nil {
		return InputResponse{}, err
	}

	if statusCode < 10 || statusCode > 20 {
		return InputResponse{}, ErrInputResponseStatusCode
	} else if statusCode > 11 {
		statusCode = 10
	}

	if responseHeader[2] != ' ' {
		return InputResponse{}, ErrInputResponseSeparator
	}

	if responseHeader[len(responseHeader)-2] != '\r' {
		return InputResponse{}, ErrInputResponseCR
	}

	if responseHeader[len(responseHeader)-1] != '\n' {
		return InputResponse{}, ErrInputResponseCR
	}

	prompt := responseHeader[3 : len(responseHeader)-2]
	return InputResponse{statusCode, prompt}, nil
}
