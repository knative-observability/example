package function

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
)

type Order struct {
	UUID string    `json:"uuid"`
	Time time.Time `json:"time"`
	PID  string    `json:"pid"`
	UID  string    `json:"uid"`
	Qty  int64     `json:"qty"`
}

type PayResp struct {
	Upstream string `json:"upstream"`
	Order    Order  `json:"order"`
	Success  bool   `json:"success"`
}

// Handle an event.
func Handle(ctx context.Context, e event.Event) (*event.Event, error) {
	/*
	 * YOUR CODE HERE
	 *
	 * Try running `go test`.  Add more test as you code in `handle_test.go`.
	 */
	var order Order
	err := e.DataAs(&order)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	time.Sleep(time.Duration(30+rand.Intn(40)) * time.Millisecond)
	slog.Info(fmt.Sprintf("Payment: User [%s], Order ID [%s]\n", order.UID, order.UUID))
	err = e.SetData("application/json", PayResp{Upstream: "payment", Order: order, Success: true})
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	e.SetType("com.example.verify")
	return &e, nil // echo to caller
}

/*
Other supported function signatures:

	Handle()
	Handle() error
	Handle(context.Context)
	Handle(context.Context) error
	Handle(event.Event)
	Handle(event.Event) error
	Handle(context.Context, event.Event)
	Handle(context.Context, event.Event) error
	Handle(event.Event) *event.Event
	Handle(event.Event) (*event.Event, error)
	Handle(context.Context, event.Event) *event.Event
	Handle(context.Context, event.Event) (*event.Event, error)

*/
