package properties

import (
	"testing"
)

func initTest() *Properties {
	p := EmptyProperties()
	p.Set("one", "1")
	p.Set("tow", "2")
	p.Set("three", "3")
	return p
}

func TestGet(t *testing.T) {
	p := initTest()
	v := p.Get("one")
	t.Log("value of key one : ", v)
}

func TestSet(t *testing.T) {
	p := initTest()
	t.Log("length before set :", p.Length())
	p.Set("four", "4")
	t.Log("length after set :", p.Length())
}

func TestDelete(t *testing.T) {
	p := initTest()
	p.Delete("one")
	v := p.Get("one")
	if v != "" {
		t.Error("key one not removed succesfully")
		t.Fail()
	} else {
		t.Log("key one removed successfully")
	}
}

func TestClear(t *testing.T) {
	p := initTest()
	p.Clear()
	if p.Length() == 0 {
		t.Log("properties has been cleared")
	} else {
		t.Error("failed to clear properties")
		t.Fail()
	}
}

func TestKeys(t *testing.T) {
	p := initTest()
	ks := p.Keys()
	length := len(ks)
	if length == p.Length() {
		t.Log("keys count is ", length)
	} else {
		t.Error("never get correct keys")
		t.Fail()
	}
}
