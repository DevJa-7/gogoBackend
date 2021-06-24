package db

import (
	"../config"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var mgoSession *mgo.Session

func mongoSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.DialWithInfo(config.MongoDialInfo())
		if err != nil {
			panic(err) // no, not really
		}

		// Optional. Switch the session to a monotonic behavior.
		mgoSession.SetMode(mgo.Monotonic, true)
	}

	return mgoSession.Clone()
}

// MongoDB return db and session, this prevents concurrent of mongo access
func MongoDB() (*mgo.Database, *mgo.Session) {
	session := mongoSession()
	mgoDB := session.DB(config.MongoDBName())
	return mgoDB, session
}

// GetCountOfCollection returns document counts of objects that are matched with param
func GetCountOfCollection(collection *mgo.Collection, base *[]bson.M) int {
	count := 0
	result := bson.M{}
	pipe := append(*base, bson.M{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}})
	if err := collection.Pipe(pipe).One(&result); err == nil {
		count = result["count"].(int)
	}
	return count
}
