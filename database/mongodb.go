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
package mongodb

import (
	"errors"
	"github.com/greivinlopez/skue"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	"strings"
)

var (
	mgoSession *mgo.Session
)

var ErrNotFound = mgo.ErrNotFound

type MongoDBPersistor struct {
	address  string
	username string
	password string
	database string
}

// New creates a new MongoDBPersistor.
// This persistor will interact with a MongoDB server.
// The parameters for the connection to the server must
// be set by you on the running machine according to:
//
// 		http://12factor.net/config
//
// The environment variables you need to set are as follows:
// MG_DB_ADDRESS= The address where your MongoDB server is running
// MG_DB_USER= The username for your connection user
// MG_DB_PASS= The password for your connection user
// MG_DB_DBNAME= The name of the database to work with
//
// For example put this on your ~/.profile file:
// export MG_DB_ADDRESS="localhost"
// export MG_DB_USER="your-db-user"
// export MG_DB_PASS="your-db-password"
// export MG_DB_DBNAME="your-db-name"
func New() *MongoDBPersistor {
	return &MongoDBPersistor{
		address:  os.Getenv("MG_DB_ADDRESS"),
		username: os.Getenv("MG_DB_USER"),
		password: os.Getenv("MG_DB_PASS"),
		database: os.Getenv("MG_DB_DBNAME")}
}

// GetSession attempts to establish a connection with the server
func (mongo *MongoDBPersistor) getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		dialInfo := mgo.DialInfo{}
		dialInfo.Addrs = []string{mongo.address}
		dialInfo.Username = mongo.username
		dialInfo.Password = mongo.password
		dialInfo.Database = mongo.database
		mgoSession, err = mgo.DialWithInfo(&dialInfo)
		if err != nil {
			panic(err) // no, not really
		}
	}
	return mgoSession.Clone()
}

// Drop removes all the elements from the given collection
func (mongo *MongoDBPersistor) Drop(collectionName string) (err error) {
	session := mongo.getSession()
	defer session.Close()

	c := session.DB(mongo.database).C(collectionName)
	_, err = c.RemoveAll(nil)
	return
}

// Count returns the number of elements of the given collection
func (mongo *MongoDBPersistor) Count(collectionName string) (n int, err error) {
	// Create MongoDB session
	session := mongo.getSession()
	defer session.Close()

	c := session.DB(mongo.database).C(collectionName)
	n, err = c.Count()
	return
}

// DropIndexes removes the indexes from the given collection
func (mongo *MongoDBPersistor) DropIndexes(collection *mgo.Collection) (err error) {
	indexes, err := collection.Indexes()
	if err != nil {
		return err
	}
	for _, index := range indexes {
		// Avoid removing native id indexes
		if !strings.HasPrefix(index.Name, "_id") {
			err = collection.DropIndex(index.Key...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Create saves the given document into the provided collection
func (mongo *MongoDBPersistor) Create(document interface{}, collection string) (err error) {
	session := mongo.getSession()
	defer session.Close()

	c := session.DB(mongo.database).C(collection)
	err = c.Insert(document)
	return
}

// Gets a string key to use on cache systems
func getKey(collection string, id interface{}) (key string, err error) {
	switch v := id.(type) {
	case string:
		return collection + "-" + v, nil
	case bson.ObjectId:
		return collection + "-" + v.Hex(), nil
	default:
		return "", errors.New("Unrecognized type for cache key")
	}
	return "", errors.New("Unrecognized type for cache key")
}

// Read retrieves the document associated with the given collection+id trying the given
// memory cache first.
func (mongo *MongoDBPersistor) Read(cache skue.MemoryCacher, document interface{}, collection string, idfield string, id interface{}) (err error) {
	// Checking cache first
	key, err := getKey(collection, id)
	if err != nil {
		return err
	}

	if cache != nil {
		err = cache.Get(key, document)
		if err == nil {
			return nil
		}
	}

	session := mongo.getSession()
	defer session.Close()

	c := session.DB(mongo.database).C(collection)
	query := bson.M{idfield: id}

	err = c.Find(query).One(document)

	if err != nil {
		return err
	}

	// Save the value to cache if needed
	if cache != nil {
		err = cache.Set(key, document)
		return err
	}
	return nil
}

// Update changes the given document on the database (and the given cache if not nil)
func (mongo *MongoDBPersistor) Update(cache skue.MemoryCacher, document interface{}, collection string, idfield string, id interface{}) (err error) {
	session := mongo.getSession()
	defer session.Close()

	c := session.DB(mongo.database).C(collection)
	query := bson.M{idfield: id}

	err = c.Update(query, document)
	if err != nil {
		return err
	}

	// Save the value to cache if needed
	if cache != nil {
		key, err := getKey(collection, id)
		if err != nil {
			return err
		}
		err = cache.Set(key, document)
	}
	return
}

// Delete removes the document associated with the given collection+id from the database
// and from the cache system given if any.
func (mongo *MongoDBPersistor) Delete(cache skue.MemoryCacher, collection string, idfield string, id interface{}) (err error) {
	session := mongo.getSession()
	defer session.Close()

	c := session.DB(mongo.database).C(collection)
	query := bson.M{idfield: id}

	err = c.Remove(query)
	if err != nil {
		return err
	}

	// Delete the value from cache if needed
	if cache != nil {
		key, err := getKey(collection, id)
		if err != nil {
			return err
		}
		err = cache.Delete(key)
	}
	return
}

// ----------------------------------------------------------------------------
