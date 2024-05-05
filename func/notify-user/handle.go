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

type InOrder struct {
	Order   Order  `json:"order"`
	Message string `json:"message"`
}

// Handle an event.
func Handle(ctx context.Context, e event.Event) (*event.Event, error) {
	/*
	 * YOUR CODE HERE
	 *
	 * Try running `go test`.  Add more test as you code in `handle_test.go`.
	 */
	var inOrder InOrder
	err := e.DataAs(&inOrder)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	time.Sleep(time.Duration(20+rand.Intn(20)) * time.Millisecond)
	slog.Info(fmt.Sprintf("Notify User: [%s] [%s]", inOrder.Order.UUID[:8], inOrder.Message))
	return nil, nil // echo to caller
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
