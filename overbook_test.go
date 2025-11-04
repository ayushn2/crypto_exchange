package main

import (
	"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}

}

func TestLimit(t *testing.T){
	l := NewLimit(10000) 
	buyOrderA := NewOrder(true, 5) // struct Order { Size: 5, Bid: true, Limit: *Limit, Timestamp: int64 }

	buyOrderB := NewOrder(true, 10)
	buyOrderC := NewOrder(true, 15)

	

	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderC)
	l.AddOrder(buyOrderB)

	l.DeleteOrder(buyOrderB)
	

	fmt.Println(l)
	
}



func TestPlaceLimitOrder(t *testing.T){
	ob := NewOrderBook()

	sellOrderA := NewOrder(false, 10)
	sellOrderB := NewOrder(false, 5)
	ob.PlaceLimitOrder(10_000, sellOrderA)
	ob.PlaceLimitOrder(9_000, sellOrderB)

	assert(t, len(ob.asks), 2)
}

