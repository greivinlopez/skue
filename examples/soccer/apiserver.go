// The MIT License (MIT)
//
// Copyright (c) 2014 Greivin López
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package main

import (
	"./database"
	"github.com/greivinlopez/skue"
	"github.com/greivinlopez/skue/views"
	"gopkg.in/martini.v1"
	"net/http"
	"os"
)

var (
	apiKey string
	view   skue.ViewLayer
)

// ----------------------------------------------------------------------------
// 			API Resource Handlers
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// PLAYER

// GET a Player resource by id
func getPlayer(params martini.Params, w http.ResponseWriter, r *http.Request) {
	id := params["id"]
	player := models.NewPlayer(id)
	skue.Read(view, player, nil, w, r)
}

// POST a new Player resource
func createPlayer(w http.ResponseWriter, r *http.Request) {
	player := models.NewPlayer("")
	skue.Create(view, player, w, r)
}

// GET the list of Player resources
func listPlayers(w http.ResponseWriter, r *http.Request) {
	player := models.NewPlayer("")
	skue.List(view, player, w, r)
}

// ----------------------------------------------------------------------------
// TEAM

// GET a Team resource by id
func getTeam(params martini.Params, w http.ResponseWriter, r *http.Request) {
	id := params["team"]
	team := models.NewTeam(id)
	skue.Read(view, team, nil, w, r)
}

// POST a new Team resource
func createTeam(w http.ResponseWriter, r *http.Request) {
	team := models.NewTeam("")
	skue.Create(view, team, w, r)
}

// GET the list of Team resources
func listTeams(w http.ResponseWriter, r *http.Request) {
	team := models.NewTeam("")
	skue.List(view, team, w, r)
}

// ----------------------------------------------------------------------------

func init() {
	// All configuration and settings are loaded from environment variables
	// Following the practices from: http://12factor.net/config

	// Retrieve the API security Key
	apiKey = os.Getenv("SOCCER_API_KEY")
	models.Address = os.Getenv("MG_DB_ADDRESS")
	models.Username = os.Getenv("MG_DB_USER")
	models.Password = os.Getenv("MG_DB_PASS")
	models.Database = os.Getenv("MG_DB_DBNAME")
	models.CreateMongoPersistor()

	// Let's use a JSON view layer: Consume from JSON and produce JSON content.
	view = *views.NewJSONView()
}

func main() {
	// This server uses the wonderful martini package: https://github.com/go-martini/martini
	m := martini.Classic()

	// Validate an API key for request authorization
	m.Use(func(res http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-API-KEY") != apiKey {
			skue.ServiceResponse(view.Producer, res, req, http.StatusUnauthorized, "You are not authorized to access this resource.")
		}
	})

	// Team resource routing
	m.Post("/teams", createTeam)
	m.Get("/teams", listTeams)
	m.Get("/teams/:team", getTeam)
	m.Any("/teams", skue.NotAllowed)
	m.Any("/teams/:team", skue.NotAllowed)

	// Player resource routing
	m.Post("/teams/:team/players", createPlayer)
	m.Get("/teams/:team/players", listPlayers)
	m.Get("/teams/:team/players/:id", getPlayer)
	m.Any("/teams/:team/players", skue.NotAllowed)
	m.Any("/teams/:team/players/:id", skue.NotAllowed)

	// Running on an unassigned port by IANA: http://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers
	http.ListenAndServe(":3020", m)
}
