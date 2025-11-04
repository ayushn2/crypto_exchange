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

func (o *Order) IsFilled() bool {
	return o.Size <= 0.0
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

func (b ByBestBid) Len() int           { 
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

func (l *Limit) Fill(o *Order) []*Match{
	matches := []*Match{}

	for _, order := range l.Orders{
		match := l.fillOrder(order, o)
		matches = append(matches, match)

		if o.IsFilled(){
			break
		}
}
	return matches
}

func (l *Limit) fillOrder(bookOrder, incomingOrder *Order) *Match{
	var (
		bid *Order
		ask *Order
		sizeFilled float64
	)

	if bookOrder.Bid{
		bid = bookOrder
		ask = incomingOrder
	}else{
		bid = incomingOrder
		ask = bookOrder
	}

	if bookOrder.Size >= incomingOrder.Size{
		bookOrder.Size -= incomingOrder.Size
		sizeFilled = incomingOrder.Size
		incomingOrder.Size = 0.0
	}else{
		incomingOrder.Size -= bookOrder.Size
		sizeFilled = bookOrder.Size
		bookOrder.Size = 0.0
	}
	return &Match{
		Ask: ask,
		Bid: bid,
		Price: l.Price,
		SizeFilled: sizeFilled,
	}
}


type Orderbook struct {
	asks []*Limit // Sell orders
	bids []*Limit // Buy orders

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

func NewOrderBook() *Orderbook {
	return &Orderbook{
		asks: []*Limit{},
		bids: []*Limit{},

		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}

func (ob *Orderbook) PlaceMarketOrder(o *Order) []*Match{
	matches := []*Match{}

	if o.Bid{
		if o.Size > ob.AskTokenVolume(){
			panic(fmt.Errorf("not enough liquidity [size: %.2f] for order [size: %.2f]", ob.AskTokenVolume(), o.Size))
		}

		for _, askLimit := range ob.Asks(){
			limitMatches := askLimit.Fill(o)
			matches = append(matches, limitMatches...)
		}
	}else{

	}

	return matches
}

func (ob *Orderbook) PlaceLimitOrder(price float64, o *Order){
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
			ob.bids = append(ob.bids, limit)
			ob.BidLimits[price] = limit
		}else{
			ob.asks = append(ob.asks, limit)
			ob.AskLimits[price] = limit
		}
	}
}

func (ob *Orderbook) BidTokenVolume() float64{
	totalVolume := 0.0
	for _, bidLimit := range ob.Bids(){
		totalVolume += bidLimit.TotalVolume
	}

	return totalVolume
}

func (ob *Orderbook) AskTokenVolume() float64{
	totalVolume := 0.0
	for _, askLimit := range ob.Asks(){
		totalVolume += askLimit.TotalVolume
	}

	return totalVolume
}

func (ob *Orderbook) Asks() []*Limit{
	sort.Sort(ByBestAsk{ob.asks})
	return ob.asks
} 

func (ob *Orderbook) Bids() []*Limit{
	sort.Sort(ByBestBid{ob.bids})
	return ob.bids
}