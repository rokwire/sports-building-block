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

package model

// SportsDefinitions
type SportsDefinitions struct {
	OrgID             string `json:"org_id"`
	ShortName         string `json:"shortName"`
	Name              string `json:"name"`
	CustomName        string `json:"custon_name"`
	HasPosition       bool   `json:"hasPosition"`
	HasHeight         bool   `json:"hasHeight"`
	HasWeight         bool   `json:"hasWeight"`
	HasSortByPosition bool   `json:"hasSortByPosition"`
	HasSortByNumber   bool   `json:"hasSortByNumber"`
	HasScores         bool   `json:"hasScores"`
	Gender            string `json:"gender"`
	Ticketed          bool   `json:"ticketed"`
	Icon              string `json:"icon"`
}
