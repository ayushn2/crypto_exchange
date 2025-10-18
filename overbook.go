package main

import (
	"fmt"
	"time"
)

type Order struct {
	Size float64
	Bid bool
	Limit *Limit
	Timestamp int64
}

func (o *Order) String() string {
	return fmt.Sprintf("[size: %.2f]",o.Size)
}
//  Limit represents a price level in the order book with associated orders.
type Limit struct {
	Price    float64
	Orders []*Order
	TotalVolume float64
}

func NewOrder(bid bool, size float64) *Order{
	return &Order{
		Size: size,
		Bid: bid,
		Timestamp: time.Now().UnixNano(),
	}
}

func NewLimit(price float64) *Limit {
	return &Limit{
		Price: price,
		Orders: []*Order{}, // Initialize an empty slice of orders because there are no orders yet and we want to avoid nil slice issues
		TotalVolume: 0, // Initialize total volume to zero
	}
}

func (l *Limit) AddOrder(o *Order) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Size
}

func (l *Limit) DeleteOrder(o *Order) {
	for i := 0; i< len(l.Orders); i++ {
		if l.Orders[i] == o {
			l.Orders[i] = l.Orders[len(l.Orders)-1] // Move the last order to the position of the order to be deleted
			l.Orders = l.Orders[:len(l.Orders)-1] // Remove the last order
		}
	}
	o.Limit = nil
	l.TotalVolume -= o.Size

	// TODO(@ayushn2): resort orders by timestamp to maintain FIFO order
}

type Overbook struct {
	Asks []*Limit // Sell orders
	Bids []*Limit // Buy orders
}