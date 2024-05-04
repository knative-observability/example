package function

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
)

type Req struct {
	Time time.Time
	PID  string
	Qty  int64
}

// Handle an event.
func Handle(ctx context.Context, e event.Event) (*event.Event, error) {
	/*
	 * YOUR CODE HERE
	 *
	 * Try running `go test`.  Add more test as you code in `handle_test.go`.
	 */
	var req Req
	err := e.DataAs(&req)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	time.Sleep(time.Duration(30+rand.Intn(40)) * time.Millisecond)
	slog.Info(fmt.Sprintf("Notify Merchant: [PID %s] [+ %d]\n", req.PID, req.Qty))
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
