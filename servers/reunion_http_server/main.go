// main.go - Reunion http server.
// Copyright (C) 2020  David Stainton.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/katzenpost/core/log"
	"github.com/katzenpost/reunion/commands"
	"github.com/katzenpost/reunion/server"
	"gopkg.in/op/go-logging.v1"
)

func httpReunionServerFactory(s *server.Server, log *logging.Logger) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		rawRequest, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Errorf("reunion HTTP server failed to ReadAll raw command data: %s", err.Error())
			return
		}
		log.Debugf("raw request size is %d", len(rawRequest))
		if len(rawRequest) == 0 {
			log.Error("error read zero sized request body from client")
			return
		}
		cmd, err := commands.FromBytes(rawRequest)
		if err != nil {
			log.Errorf("1 reunion HTTP server invalid query command: %s", err.Error())
			return
		}
		replyCmd, err := s.ProcessQuery(cmd)
		if err != nil {
			log.Errorf("reunion HTTP server invalid reply command: %s", err.Error())
			return
		}
		rawReply := replyCmd.ToBytes()
		_, err = w.Write(rawReply)
		if err != nil {
			log.Errorf("reunion HTTP server failure to send reply command: %s", err.Error())
			return
		}
	}
}

func runHTTPServer(address, urlPath, logPath, logLevel string) *http.Server {
	logBackend, err := log.New(logPath, logLevel, false)
	if err != nil {
		panic(err)
	}
	reunionServer := server.NewServer()
	httpServeMux := http.NewServeMux()
	httpLog := logBackend.GetLogger("reunion_http_server")
	httpServeMux.HandleFunc(urlPath, httpReunionServerFactory(reunionServer, httpLog))
	httpServer := &http.Server{
		Addr:           address,
		Handler:        httpServeMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go httpServer.ListenAndServe()
	return httpServer
}

func main() {
	address := flag.String("l", "127.0.0.1:12345", "Listen address. Defaults to 127.0.0.1:12345")
	urlPath := flag.String("p", "/reunion", "Reunion URL path.")
	logPath := flag.String("log", "", "Log file path. Default STDOUT.")
	logLevel := flag.String("level", "DEBUG", "Log level.")
	flag.Parse()
	runHTTPServer(*address, *urlPath, *logPath, *logLevel)
}
