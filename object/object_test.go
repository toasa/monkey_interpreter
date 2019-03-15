package object

import "testing"

func TestStringHashKey(t *testing.T) {
    hello1 := &String{Value: "Hello World"}
    hello2 := &String{Value: "Hello World"}
    howdy1 := &String{Value: "howdy?"}
    howdy2 := &String{Value: "howdy?"}

    if hello1.HashKey() != hello2.HashKey() {
        t.Errorf("strings with same content have different hash keys")
    }

    if howdy1.HashKey() != howdy2.HashKey() {
        t.Errorf("strings with same content have different hash keys")
    }

    if howdy1.HashKey() == hello2.HashKey() {
        t.Errorf("strings with different content have same hash keys")
    }
}
