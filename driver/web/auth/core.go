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
	"log"
	"net/http"
	"sport/core"
	"sport/core/model"

	"github.com/rokwire/core-auth-library-go/authservice"
	"github.com/rokwire/core-auth-library-go/tokenauth"
	"github.com/rokwire/logging-library-go/logs"
)

// CoreAuth implementation
type CoreAuth struct {
	app       *core.Application
	tokenAuth *tokenauth.TokenAuth
}

// NewCoreAuth creates new CoreAuth
func NewCoreAuth(app *core.Application, config model.Config) *CoreAuth {

	remoteConfig := authservice.RemoteAuthDataLoaderConfig{
		AuthServicesHost: config.CoreBBHost,
	}

	serviceLoader, err := authservice.NewRemoteAuthDataLoader(remoteConfig, []string{"core"}, logs.NewLogger("groupsbb", &logs.LoggerOpts{}))
	authService, err := authservice.NewAuthService("sports-service", config.SportsServiceURL, serviceLoader)
	if err != nil {
		log.Fatalf("Error initializing auth service: %v", err)
	}
	tokenAuth, err := tokenauth.NewTokenAuth(true, authService, nil, nil)
	if err != nil {
		log.Fatalf("Error intitializing token auth: %v", err)
	}

	auth := CoreAuth{app: app, tokenAuth: tokenAuth}
	return &auth
}

// Check checks the request contains a valid Core access token
func (ca CoreAuth) Check(r *http.Request) (*tokenauth.Claims, error) {
	claims, err := ca.tokenAuth.CheckRequestTokens(r)
	if err != nil || claims == nil {
		log.Printf("error validate token: %s", err)
		return nil, err
	}

	/*if !claims.Admin {
		err = errors.ErrorData(logutils.StatusInvalid, logutils.TypeClaim, logutils.StringArgs("admin"))
		log.Println(err)
		return nil, err
	}

	err = ca.tokenAuth.AuthorizeRequestPermissions(claims, r)
	if err != nil {
		log.Println("invalid permissions:", err)
		return nil, err
	}*/

	return claims, nil
}
