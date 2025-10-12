package model

import "github.com/MicahParks/keyfunc"

type AuthMiddleware struct {
	JWKS     *keyfunc.JWKS
	Issuer   string
	Audience string
}

type TokenRequest struct {
	Token string `json:"token"`
}

type Jwks struct {
	Keys []struct {
		Kid string   `json:"kid"`
		Kty string   `json:"kty"`
		N   string   `json:"n"`
		E   string   `json:"e"`
		X5c []string `json:"x5c"`
	} `json:"keys"`
}
