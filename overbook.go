package main

import "time"

type Order struct {
	Size float64
	Bid bool
	Limit *Limit
	Timestamp int64
}

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

type Overbook struct {
	Asks []*Limit // Sell orders
	Bids []*Limit // Buy orders
}