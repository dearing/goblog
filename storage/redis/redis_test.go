package redis

import (
	"fmt"
	"testing"
)

var p = Post{
	ID:       "0",
	Title:    "Testing Post",
	Content:  "Testing Content",
	Author:   "Testing Author",
	Created:  "1357271500",
	Modified: "1357271580",
	Accessed: "0",
}

/*	
========================================
	TESTS!
========================================
	# go test
*/

func TestConnect(t *testing.T) {
	Connect("192.168.1.150:6379", "", -1)
	defer Close()

}

func TestPost(t *testing.T) {
	Connect("192.168.1.150:6379", "", -1)
	defer Close()

	err := Set(p)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	q, err := Get(p.ID)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	// no way to know, we test this elsewhere anyway
	p.Accessed = q.Accessed

	if q != p {
		t.Fail()
		t.Log("Posts were not equal")
	}

}

func TestAccessed(t *testing.T) {
	Connect("192.168.1.150:6379", "", -1)
	defer Close()

	key := fmt.Sprintf("post:%v", p.ID)

	err := Set(p)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	q, err := Get(p.ID)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	x := client.HGet(key, "accessed")

	if x.Val() != q.Accessed {
		t.Log("inequal access")
		t.Fail()
	}
}

/*	
========================================
	BENCHMARKS!! not ready
========================================
	# go test -bench=".*"
*/

func BenchmarkSet(b *testing.B) {
	for i := 0; i < b.N; i++ {

		client.Incr("global:nextPostID")
		y := client.Get("global:nextPostID")

		p.ID = fmt.Sprintf("%v", y.Val())
		err := Set(p)
		if err != nil {
			b.Error(err)
			b.Fail()
		}
	}
}

func BenchmarkGet(b *testing.B) {
	_, err := Get("post:0")
	if err != nil {
		b.Error(err)
		b.Fail()
	}

}

func BenchmarkDel(b *testing.B) {
	x := client.Keys("post:*")

	for _, z := range x.Val() {
		Del(z)
	}

}
