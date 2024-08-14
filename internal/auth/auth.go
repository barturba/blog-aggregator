package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(h http.Header) (string, error) {
	authorization := h.Get("Authorization")
	if len(authorization) == 0 {
		return "", errors.New("No Apikey provided")
	}
	fmt.Printf("middlewareAuth: got apikey: %s\n", authorization)
	authorization = strings.TrimPrefix(authorization, "ApiKey ")
	return authorization, nil
}
