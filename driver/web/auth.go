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
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/rokwire/core-auth-library-go/v2/authorization"
	"github.com/rokwire/core-auth-library-go/v2/authservice"
	"github.com/rokwire/core-auth-library-go/v2/tokenauth"
)

type auth struct {
	host      string
	tokenAuth *tokenauth.TokenAuth
}

func newAuth(host string, coreURL string) *auth {
	sportsServiceURL := fmt.Sprintf("%s/sports-service", host)

	authService := authservice.AuthService{
		ServiceID:   "sports-service",
		ServiceHost: sportsServiceURL,
		FirstParty:  true,
		AuthBaseURL: coreURL,
	}

	serviceRegLoader, err := authservice.NewRemoteServiceRegLoader(&authService, nil)
	if err != nil {
		log.Printf("auth -> newAuth: FAILED to init service reg loader: %s", err.Error())
	}

	serviceRegManager, err := authservice.NewServiceRegManager(&authService, serviceRegLoader)
	if err != nil {
		log.Printf("auth -> newAuth: FAILED to init service reg manager: %s", err.Error())
	}

	permissionAuth := authorization.NewCasbinStringAuthorization("driver/web/authorization_policy.csv")
	tokenAuth, err := tokenauth.NewTokenAuth(true, serviceRegManager, permissionAuth, nil)
	if err != nil {
		log.Printf("auth -> newAuth: FAILED to init token auth: %s", err.Error())
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

func (a auth) corePermissionAuthCheck(w http.ResponseWriter, r *http.Request) error {
	if a.tokenAuth == nil {
		log.Printf("auth -> corePermissionAuthCheck: tokenAuth is nil")
		return fmt.Errorf("auth Service is not initialized")
	}
	claims, err := a.tokenAuth.CheckRequestTokens(r)
	if err != nil {
		log.Printf("auth -> corePermissionAuthCheck: FAILED to validate token: %s", err.Error())
		return err
	}

	err = a.tokenAuth.AuthorizeRequestPermissions(claims, r)
	if err != nil {
		log.Printf("auth -> corePermissionAuthCheck: invalid permissions: %s", err)
		return errors.New("invalid permissions")
	}

	return nil
}

func (a auth) claimsCheck(w http.ResponseWriter, r *http.Request) (*tokenauth.Claims, error) {

	claims, err := a.tokenAuth.CheckRequestTokens(r)
	if err != nil {
		log.Printf("token claims are nil", err.Error())
		return nil, err
	}

	return claims, err
}
