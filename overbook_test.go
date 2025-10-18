package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimit(t *testing.T){
	l := NewLimit(10000) 
	buyOrderA := NewOrder(true, 5) // struct Order { Size: 5, Bid: true, Limit: *Limit, Timestamp: int64 }

	buyOrderB := NewOrder(true, 10)
	buyOrderC := NewOrder(true, 15)

	

	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderC)
	l.AddOrder(buyOrderB)

	l.DeleteOrder(buyOrderB)
	assert.Equal(t, 20.0, l.TotalVolume)
	assert.Equal(t, 2, len(l.Orders))
	assert.Nil(t, buyOrderB.Limit)

	fmt.Println(l)
	
}



func TestOrderBook(t *testing.T){
	
}