package main

import (
	"fmt"
	"sort"
	"time"
)

type Match struct {
	Ask *Order
	Bid *Order
	Price float64
	SizeFilled float64
}

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
	Orders Orders
	TotalVolume float64
}

type Orders []*Order

// Implement sort.Interface for Orders based on Timestamp for FIFO ordering.
func (o Orders) Len() int           { return len(o) }
func (o Orders) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Orders) Less(i, j int) bool { return o[i].Timestamp < o[j].Timestamp }

func NewOrder(bid bool, size float64) *Order{
	return &Order{
		Size: size,
		Bid: bid,
		
		Timestamp: time.Now().UnixNano(),
	}
}

type Limits []*Limit

type ByBestAsk struct{ Limits }

func (a ByBestAsk) Len() int           { 
	return len(a.Limits) 
}

func (a ByBestAsk) Swap(i, j int)  { 
	a.Limits[j], a.Limits[i] = a.Limits[i], a.Limits[j]
}

func (a ByBestAsk) Less(i, j int) bool           { 
	return a.Limits[i].Price< a.Limits[j].Price
}

// ByBestBid implements sort.Interface for []Limit based on the Price field.

type ByBestBid struct{ Limits }

func (b ByBestBid) Bid() int           { 
	return len(b.Limits) 
}

func (b ByBestBid) Swap(i, j int)  { 
	b.Limits[j], b.Limits[i] = b.Limits[i], b.Limits[j]
}

func (b ByBestBid) Less(i, j int) bool           { 
	return b.Limits[i].Price< b.Limits[j].Price
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

	// resort orders by timestamp to maintain FIFO order
	sort.Sort(Orders(l.Orders))
}

type Overbook struct {
	Asks []*Limit // Sell orders
	Bids []*Limit // Buy orders

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

func NewOrderBook() *Overbook {
	return &Overbook{
		Asks: []*Limit{},
		Bids: []*Limit{},

		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}

func (ob *Overbook) PlaceOrder(price float64, o *Order) []Match{
	// 1. Check for matches
	// TODO(@ayushn2): implement matching engine

	// 2. If no matches, add to order book
	if o.Size > 0.0{
		ob.add(price, o)
	}
	
	return []Match{}
}

func (ob *Overbook) add(price float64, o *Order)  {
	var limit *Limit

	if o.Bid{
		limit = ob.BidLimits[price]
	}else{
		limit = ob.AskLimits[price]
	}

	if limit == nil{
		limit = NewLimit(price)
		limit.AddOrder(o)
		if o.Bid{
			ob.Bids = append(ob.Asks, limit)
			ob.BidLimits[price] = limit
		}else{
			ob.Asks = append(ob.Asks, limit)
			ob.AskLimits[price] = limit
		}
	}
}