package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rokwire/core-auth-library-go/authservice"
	"github.com/rokwire/core-auth-library-go/tokenauth"
)

type auth struct {
	host      string
	tokenAuth *tokenauth.TokenAuth
}

func newAuth(host string) *auth {

	remoteServiceURL := fmt.Sprintf("%s/core/bbs/service-regs", host)
	sportsServiceURL := fmt.Sprintf("%s/sports-service", host)

	serviceLoader := authservice.NewRemoteServiceRegLoader(remoteServiceURL, nil)
	authService, err := authservice.NewAuthService("sports-service", sportsServiceURL, serviceLoader)
	var tokenAuth *tokenauth.TokenAuth
	if err == nil {
		tokenAuth, err = tokenauth.NewTokenAuth(true, authService, nil, nil)
		if err != nil {
			log.Printf("auth -> newAuth: FAILED to init token auth: %s", err.Error())
		}
	} else {
		log.Printf("auth -> newAuth: FAILED to init auth service: %s", err.Error())
	}

	auth := auth{host: host, tokenAuth: tokenAuth}
	return &auth
}

func (a auth) coreAuthCheck(w http.ResponseWriter, r *http.Request) error {
	if a.tokenAuth == nil {
		log.Printf("auth -> coreAuthCheck: tokenAuth is nil")
		return fmt.Errorf("auth Service is not initialized")
	}
	_, err := a.tokenAuth.CheckRequestTokens(r)
	if err != nil {
		log.Printf("auth -> coreAuthCheck: FAILED to validate token: %s", err.Error())
		return err
	}

	return nil
}
