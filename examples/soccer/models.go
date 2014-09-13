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
	"github.com/greivinlopez/skue"
	//"github.com/greivinlopez/skue/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"os"
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

var (
	mgoSession        *mgo.Session
	playersCollection = "players"
)

// ----------------------------------------------------------------------------
// 			MONGODB
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------

// NewPlayer creates a new empty Player object with the provided id.
// All the other fields will be empty at first.
func NewPlayer(id string) *Player {
	playerId := bson.NewObjectId()
	if id != "" {
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

// ----------------------------------------------------------------------------
// 			skue.DatabasePersistor implementation
// ----------------------------------------------------------------------------

func (player *Player) Read(cache skue.MemoryCacher) (err error) {
	//mongo := mongodb.New()
	// Read(&stand, stand.CollectionName(), "_id", bson.ObjectIdHex(id))
	return nil
}

func (player *Player) Create() (err error) {
	return nil
}

func (player *Player) Update(cache skue.MemoryCacher) (err error) {
	return nil
}

func (player *Player) Delete(cache skue.MemoryCacher) (err error) {
	return nil
}

// ----------------------------------------------------------------------------
