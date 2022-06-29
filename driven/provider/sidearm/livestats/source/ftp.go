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
	"io/ioutil"
	"log"

	"github.com/jlaffaye/ftp"
)

type ftpConn struct {
	ftpConnection *ftp.ServerConn
	host          string
	user          string
	password      string
}

func newFTPConn(host string, user string, password string) ftpConn {
	return ftpConn{host: host, user: user, password: password}
}

func (ftpConn *ftpConn) load(folder string) ([]byte, error) {
	var err error
	if ftpConn.ftpConnection == nil {
		log.Println("ftpConn: loadFromXML -> need to connect and login to ftp " + folder)

		//1. connect to ftp
		ftpConn.ftpConnection, err = ftpConn.connectToFtp()
		if err != nil {
			return nil, err
		}
		//2. login
		err = ftpConn.login()
		if err != nil {
			return nil, err
		}
	}

	//3. download the file
	xmlData, err := ftpConn.downloadXML(folder)
	if err != nil {
		log.Printf("ftpConn: loadFromXML -> fail to download the xml so try to reconnect - %s\n", err.Error())
		//in case of an error try to reconnect
		err := ftpConn.reConnect()
		if err != nil {
			//nothing can be done
			return nil, err
		}
		//try download again
		xmlData, err = ftpConn.downloadXML(folder)
		if err != nil {
			//nothing can be done
			return nil, err
		}
	}
	return xmlData, nil
}

func (ftpConn *ftpConn) reConnect() error {
	log.Printf("ftpConn: reConnect -> try to reConnect")
	var err error

	if ftpConn.ftpConnection != nil {
		err = ftpConn.ftpConnection.Logout()
		if err != nil {
			log.Printf("ftpConn: reConnect -> error on logout %s", err.Error())
		}
		err = ftpConn.ftpConnection.Quit()
		if err != nil {
			log.Printf("ftpConn: reConnect -> error on quit %s", err.Error())
		}
	}

	ftpConn.ftpConnection, err = ftpConn.connectToFtp()
	if err != nil {
		log.Printf("ftpConn: reConnect -> error on connect %s", err.Error())
	}
	err = ftpConn.login()
	if err != nil {
		log.Printf("ftpConn: reConnect -> error on login %s", err.Error())
	}
	return err
}

func (ftpConn *ftpConn) connectToFtp() (*ftp.ServerConn, error) {
	connection, err := ftp.Dial(ftpConn.host + ":21")
	if err != nil {
		log.Printf("ftpConn: connectToFtp -> error dialing ftp host:%s\terror:%s", ftpConn.host, err.Error())
		return nil, err
	}
	return connection, nil
}

func (ftpConn *ftpConn) login() error {
	if err := ftpConn.ftpConnection.Login(ftpConn.user, ftpConn.password); err != nil {
		log.Printf("ftpConn: login -> error login in ftp with user %s and password %s\terror:%s", ftpConn.user, ftpConn.password, err.Error())
		return err
	}
	return nil
}

func (ftpConn *ftpConn) quit() error {
	if err := ftpConn.ftpConnection.Quit(); err != nil {
		log.Printf("ftpConn: quit -> error quit from ftp:%s", err.Error())
		return err
	}
	return nil
}

func (ftpConn *ftpConn) downloadXML(folder string) ([]byte, error) {
	if ftpConn.ftpConnection == nil {
		errorMessage := "Cannot download xml because the connection is null - " + folder
		return nil, errors.New(errorMessage)
	}

	// Change to the correct directory
	err := ftpConn.ftpConnection.ChangeDir(folder)
	if err != nil {
		return nil, err
	}

	// Download the file
	response, err := ftpConn.ftpConnection.Retr("1.xml")
	if err != nil {
		return nil, err
	}

	byteValue, err := ioutil.ReadAll(response)
	if err != nil {
		return nil, err
	}
	response.Close()

	return byteValue, nil
}
