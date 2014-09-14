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
package models

import (
	"github.com/greivinlopez/skue"
	"github.com/greivinlopez/skue/database"
	"gopkg.in/mgo.v2/bson"
)

var (
	Address  string // The address to reach the MongoDB server
	Username string // The username to connect with the MongoDB server
	Password string // The password of the MongoDB user
	Database string // The name of the database to store the models
	mongo    *mongodb.MongoDBPersistor
)

// Player represents a soccer player.
type Player struct {
	Id          bson.ObjectId `json:"Id" bson:"_id"`
	FirstName   string
	LastName    string
	Nationality string
	Age         int
	Position    string
	Height      string
	Weight      string
	Foot        string
}

// Creates a MongoDB persistor to interact with the database
func CreateMongoPersistor() {
	mongo = mongodb.New(Address, Username, Password, Database)
}

// ----------------------------------------------------------------------------
// 			PLAYER
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------

// NewPlayer creates a new empty Player object with the provided id.
// All the other fields will be empty at first.
func NewPlayer(id string) *Player {
	playerId := bson.NewObjectId()
	if id != "" && bson.IsObjectIdHex(id) {
		playerId = bson.ObjectIdHex(id)
	}
	return &Player{
		Id:          playerId,
		FirstName:   "",
		LastName:    "",
		Nationality: "",
		Age:         0,
		Position:    "",
		Height:      "",
		Weight:      "",
		Foot:        ""}
}

func (player *Player) Collection() string {
	return "players"
}

// ----------------------------------------------------------------------------
// 			skue.DatabasePersistor implementation
// ----------------------------------------------------------------------------

func (player *Player) Read(cache skue.MemoryCacher) (err error) {
	err = mongo.Read(cache, &player, player.Collection(), "_id", player.Id)
	return
}

func (player *Player) Create() (err error) {
	err = mongo.Create(&player, player.Collection())
	return
}

func (player *Player) Update(cache skue.MemoryCacher) (err error) {
	err = mongo.Update(cache, &player, player.Collection(), "_id", player.Id)
	return
}

func (player *Player) Delete(cache skue.MemoryCacher) (err error) {
	err = mongo.Delete(cache, player.Collection(), "_id", player.Id)
	return
}

func (player *Player) List() (results []interface{}, err error) {
	players := []Player{}
	err = mongo.List(&players, player.Collection(), nil, 25)
	return players, err
}

// ----------------------------------------------------------------------------
