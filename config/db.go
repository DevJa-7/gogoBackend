package config

import (
	"time"

	mgo "gopkg.in/mgo.v2"
)

/////////////////////////////
/*    Mysql Production     */
/////////////////////////////

// MysqlHost ...
var MysqlHost = "127.0.0.1"

// MysqlPort ...
var MysqlPort = "3306"

// MysqlUser ...
var MysqlUser = "root"

// MysqlPassword ...
var MysqlPassword = ""

// MysqlDatabase ...
var MysqlDatabase = "db_gogo"

/////////////////////////////
/*    Mongo Production     */
/////////////////////////////

// MongoHost ...
var MongoHost = "127.0.0.1"

// MongoPort ...
var MongoPort = "27017"

// MongoUser ...
var MongoUser = ""

// MongoPassword ...
var MongoPassword = ""

// MongoDatabase ...
var MongoDatabase = "db_gogo"

// Constants for database.
const (
	MysqlProtocol = "tcp"
	// for DEVELOPMENT
	MysqlHostDev     = "127.0.0.1"
	MysqlPortDev     = "3306"
	MysqlUserDev     = "root"
	MysqlPasswordDev = ""
	MysqlDatabaseDev = "db_gogo_dev"
	MysqlOptionsDev  = "charset=utf8&parseTime=True"
	// for PRODUCTION
	MysqlOptions = "charset=utf8&parseTime=True"

	MongoProtocol = "tcp"
	// for DEVELOPMENT
	MongoHostDev     = "127.0.0.1"
	MongoPortDev     = "27017"
	MongoUserDev     = ""
	MongoPasswordDev = ""
	MongoDatabaseDev = "db_gogo_dev"
)

// MysqlDSL return mysql DSL.
func MysqlDSL() string {
	var mysqlDSL string
	switch Environment {
	case "DEVELOPMENT":
		mysqlDSL = MysqlUserDev + ":" + MysqlPasswordDev + "@" + MysqlProtocol + "(" + MysqlHostDev + ":" + MysqlPortDev + ")/" + MysqlDatabaseDev + "?" + MysqlOptionsDev
	default:
		mysqlDSL = MysqlUser + ":" + MysqlPassword + "@" + MysqlProtocol + "(" + MysqlHost + ":" + MysqlPort + ")/" + MysqlDatabase + "?" + MysqlOptions
	}
	return mysqlDSL
}

// MongoDialInfo return mongo DialInfo.
func MongoDialInfo() *mgo.DialInfo {
	var info *mgo.DialInfo
	switch Environment {
	case "DEVELOPMENT":
		info = &mgo.DialInfo{
			Addrs:    []string{MongoHostDev + ":" + MongoPortDev},
			Timeout:  60 * time.Second,
			Database: MongoDatabaseDev,
			Username: MongoUserDev,
			Password: MongoPasswordDev,
		}
	default:
		info = &mgo.DialInfo{
			Addrs:    []string{MongoHost + ":" + MongoPort},
			Timeout:  60 * time.Second,
			Database: MongoDatabase,
			Username: MongoUser,
			Password: MongoPassword,
		}
	}
	return info
}

// MongoDBName return database name
func MongoDBName() string {
	var database string
	switch Environment {
	case "DEVELOPMENT":
		database = MongoDatabaseDev
	default:
		database = MongoDatabase
	}
	return database
}
