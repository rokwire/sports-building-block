// Copyright 2026 Board of Trustees of the University of Illinois.
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

package sidearm

import (
	"log"
	"net/url"
	"strings"
)

// ticketsUTMParams maps sport shortnames to the UTM query params that should
// be applied to their tickets URLs, as requested by the product team.
var ticketsUTMParams = map[string]string{
	"football": "utm_source=campusapp&utm_medium=inapplink&utm_campaign=f26tickets&utm_id=f26tickets&utm_content=f26campusapp",
	"wvball":   "utm_source=campusapp%60&utm_medium=inapplink&utm_campaign=vb26tickets&utm_id=vb26tickets&utm_content=vb26campusapp",
	"wbball":   "utm_source=campusapp%60&utm_medium=inapplink&utm_campaign=wb26tickets&utm_id=wb26tickets&utm_content=wb26campusapp",
	"mbball":   "utm_source=campusapp&utm_medium=inapplink&utm_campaign=b26tickets&utm_id=b26tickets&utm_content=b26campusapp",
}

// addTicketsUTMParams replaces any existing utm_* query params on ticketsURL
// with the ones configured for sportShortName. If sportShortName is not one
// of the sports we override, or ticketsURL is empty or unparsable, the
// original ticketsURL is returned unchanged.
func addTicketsUTMParams(sportShortName string, ticketsURL string) string {
	if ticketsURL == "" {
		return ticketsURL
	}

	utmParams, ok := ticketsUTMParams[sportShortName]
	if !ok {
		return ticketsURL
	}

	u, err := url.Parse(ticketsURL)
	if err != nil {
		log.Printf("sidearm -> tickets -> addTicketsUTMParams: failed to parse tickets url [%s]. Reason: %s", ticketsURL, err.Error())
		return ticketsURL
	}

	newParams, err := url.ParseQuery(utmParams)
	if err != nil {
		log.Printf("sidearm -> tickets -> addTicketsUTMParams: failed to parse configured utm params for sport [%s]. Reason: %s", sportShortName, err.Error())
		return ticketsURL
	}

	q := u.Query()
	for key := range q {
		if strings.HasPrefix(key, "utm_") {
			q.Del(key)
		}
	}
	for key, values := range newParams {
		for _, value := range values {
			q.Add(key, value)
		}
	}

	u.RawQuery = q.Encode()
	return u.String()
}
