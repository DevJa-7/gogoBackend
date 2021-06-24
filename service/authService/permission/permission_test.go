package permission

import (
	"fmt"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestGenerateToken(t *testing.T) {
	objid := bson.NewObjectId()
	role := "tester"
	fmt.Println(objid)
	fmt.Println(objid.Hex())

	token, err := GenerateToken(objid, role)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(token)
}
