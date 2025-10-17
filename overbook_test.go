package main

import (
	"fmt"
	"testing"
)

func TestLimit(t *testing.T){
	l := NewLimit(10000) 
	buyOrder := NewOrder(true, 5) // struct Order { Size: 5, Bid: true, Limit: *Limit, Timestamp: int64 }

	l.AddOrder(buyOrder)

	fmt.Println(l)
	
}

func TestOrderBook(t *testing.T){
	
}