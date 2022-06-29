// Copyright 2022 Board of Trustees of the University of Illinois.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
