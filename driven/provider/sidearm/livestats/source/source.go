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

package source

import (
	"errors"
	"log"
	"sport/core/model"
	sidearmModel "sport/driven/provider/sidearm/model"
)

// Source represents the source package interface
type Source interface {
	UpdateConfig(config Config)
	LoadData(item *sidearmModel.LiveGameItem) (model.LiveGame, error)
}

type sourceImpl struct {
	config              Config
	sidearm             sidearmSource
	xmlFootbalSource    xmlFootballSource
	xmlBasketballSource xmlBasketballSource
	xmlVolleyballSource xmlVolleyballSource
}

// New create new source instance
func New(config Config, ftpHost string, ftpUser string, ftpPassword string) Source {
	sidearmSource := newSidearmSource(config)
	xmlFootballSource := newXMLFootballSource(config, ftpHost, ftpUser, ftpPassword)
	xmlBasketballSource := newXMLBasketballSource(config, ftpHost, ftpUser, ftpPassword)
	xmlVolleyballSource := newXMLVolleyballSource(config, ftpHost, ftpUser, ftpPassword)
	return &sourceImpl{config: config, sidearm: sidearmSource, xmlFootbalSource: xmlFootballSource,
		xmlBasketballSource: xmlBasketballSource, xmlVolleyballSource: xmlVolleyballSource}
}

func (livestatsSource *sourceImpl) UpdateConfig(config Config) {
	log.Println("source: UpdateConfig -> config updated in livestats source")
	livestatsSource.config = config
	livestatsSource.sidearm.updateConfig(config)
	livestatsSource.xmlFootbalSource.updateConfig(config)
	livestatsSource.xmlBasketballSource.updateConfig(config)
	livestatsSource.xmlVolleyballSource.updateConfig(config)
}

// LoadData loads the current livestats data for an item
func (livestatsSource *sourceImpl) LoadData(item *sidearmModel.LiveGameItem) (model.LiveGame, error) {
	sport := item.Sport
	home := item.Home
	sources := livestatsSource.config.GetLivestatsSource(sport, home)
	log.Printf("source: LoadData -> sources:%s sport:%s gameId:%s", sources, item.Sport, item.GameID)

	//get the live data from the sources by priority
	for index, source := range sources {
		if source == "sidearm" {
			result, err := livestatsSource.sidearm.load(item)
			if err == nil {
				return result, nil
			}

			//error
			log.Print(err.Error())
			isLast := (len(sources) - 1) == index
			if isLast {
				log.Printf("source: LoadData -> it is last source so return error")
				return nil, err

			}
		} else if source == "xml_feed" {
			result, err := livestatsSource.loadFromXML(item)
			if err == nil {
				return result, nil
			}

			//error
			log.Print(err.Error())
			isLast := (len(sources) - 1) == index
			if isLast {
				log.Printf("source: LoadData -> it is last source so return error")
				return nil, err
			}
		}
	}
	return nil, errors.New("source: LoadData -> no source provided")
}

func (livestatsSource *sourceImpl) loadFromXML(item *sidearmModel.LiveGameItem) (model.LiveGame, error) {
	switch item.Sport {
	case "football":
		return livestatsSource.xmlFootbalSource.load(item)
	case "mbball":
		return livestatsSource.xmlBasketballSource.load(item)
	case "wbball":
		return livestatsSource.xmlBasketballSource.load(item)
	case "wvball":
		return livestatsSource.xmlVolleyballSource.load(item)
	default:
		log.Printf("source: loadFromXML -> not supported sport")
	}
	return nil, nil
}
