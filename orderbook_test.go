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

func TestPlaceMarketOrder(t *testing.T){
	ob := NewOrderBook()

	sellOrderA := NewOrder(false, 20)
	ob.PlaceLimitOrder(10_000, sellOrderA)

	buyOrderA := NewOrder(true, 5)
	ob.PlaceLimitOrder(10_000, buyOrderA)
	matches := ob.PlaceMarketOrder(buyOrderA)

	assert(t, len(matches), 1)
	assert(t, len(ob.asks), 1)
	assert(t, ob.AskTotalVolume(), 15.0)
	assert(t, matches[0].Ask, sellOrderA)
	assert(t, matches[0].Bid, buyOrderA)
	assert(t, matches[0].Price, 10_000.0)
	assert(t, matches[0].SizeFilled, 5.0)
	assert(t, buyOrderA.IsFilled(), true)

	fmt.Printf("%+v\n", matches)

}
