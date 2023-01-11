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
	OrgID             string `json:"org_id" bson:"org_id"`
	ShortName         string `json:"shortName" bson:"shortName"`
	Name              string `json:"name" bson:"name"`
	CustomName        string `json:"custom_name" bson:"custom_name"`
	HasPosition       bool   `json:"hasPosition" bson:"hasPosition"`
	HasHeight         bool   `json:"hasHeight" bson:"hasHeight"`
	HasWeight         bool   `json:"hasWeight" bson:"hasWeight"`
	HasSortByPosition bool   `json:"hasSortByPosition" bson:"hasSortByPosition"`
	HasSortByNumber   bool   `json:"hasSortByNumber" bson:"hasSortByNumber"`
	HasScores         bool   `json:"hasScores" bson:"hasScores"`
	Gender            string `json:"gender" bson:"gender"`
	Ticketed          bool   `json:"ticketed" bson:"ticketed"`
	Icon              string `json:"icon" bson:"icon"`
}
