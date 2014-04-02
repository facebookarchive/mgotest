package mgotest_test

import (
	"testing"

	"labix.org/v2/mgo/bson"

	"github.com/facebookgo/mgotest"
)

func test(t *testing.T, answer int) {
	t.Parallel()
	mongo := mgotest.NewStartedServer(t)
	defer mongo.Stop()
	const id = 1
	in := bson.M{"_id": id, "answer": answer}
	collection := mongo.Session().DB("tdb").C("tc")
	if err := collection.Insert(in); err != nil {
		t.Fatal(err)
	}
	out := bson.M{}
	if err := collection.FindId(id).One(out); err != nil {
		t.Fatal(err)
	}
	if out["answer"] != answer {
		t.Fatalf("did not find expected answer, got %v", out)
	}
}

// Testing that multiple instances don't stomp on each other.
func TestOne(t *testing.T) {
	test(t, 42)
}

func TestTwo(t *testing.T) {
	test(t, 43)
}

func TestThree(t *testing.T) {
	test(t, 44)
}
