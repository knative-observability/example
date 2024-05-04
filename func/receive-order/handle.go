package function

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
)

type InOrder struct {
	PID string `json:"pid"`
	UID string `json:"uid"`
	Qty int64  `json:"qty"`
}

type OutOrder struct {
	UUID string    `json:"uuid"`
	Time time.Time `json:"time"`
	PID  string    `json:"pid"`
	UID  string    `json:"uid"`
	Qty  int64     `json:"qty"`
}

// Handle an event.
func Handle(ctx context.Context, e event.Event) (*event.Event, error) {
	/*
	 * YOUR CODE HERE
	 *
	 * Try running `go test`.  Add more test as you code in `handle_test.go`.
	 */
	var inOrder InOrder
	err := e.DataAs(inOrder)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	orderID := uuid.NewString()
	outOrder := OutOrder{
		UUID: orderID,
		Time: time.Now(),
		PID:  inOrder.PID,
		UID:  inOrder.UID,
		Qty:  inOrder.Qty,
	}
	err = e.SetData("application/json", outOrder)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	slog.Info(fmt.Sprintf("Receive Order: [UUID %s] User [%s] by [%s] * [%d]\n",
		orderID, outOrder.UID, inOrder.PID, outOrder.Qty))
	time.Sleep(time.Duration(30+rand.Intn(40)) * time.Millisecond)
	e.SetType("com.example.pay-stock")
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
