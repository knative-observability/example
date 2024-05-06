package function

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type ReqIn struct {
	PID string `json:"pid"`
	Qty int64  `json:"qty"`
}
type ReqOut struct {
	Time time.Time `json:"time"`
	PID  string    `json:"pid"`
	Qty  int64     `json:"qty"`
}

// Handle an event.
func Handle(ctx context.Context, e event.Event) (*event.Event, error) {
	/*
	 * YOUR CODE HERE
	 *
	 * Try running `go test`.  Add more test as you code in `handle_test.go`.
	 */
	var reqIn ReqIn
	err := e.DataAs(&reqIn)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	reqOut := ReqOut{
		Time: time.Now(),
		PID:  reqIn.PID,
		Qty:  reqIn.Qty,
	}
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	defer client.Close()
	// Set locks to prevent data race
	for {
		locked, err := client.SetNX(ctx, "lock", uuid.NewString(), 10*time.Second).Result()
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		if locked {
			break
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
	defer func() {
		client.Del(ctx, "lock")
	}()
	// Have got lock
	result, err := client.HIncrBy(ctx, "stock", string(reqIn.PID), reqIn.Qty).Result()
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	time.Sleep(time.Duration(10+rand.Intn(20)) * time.Millisecond)
	if rand.Intn(100) >= 5 {
		err = e.SetData("application/json", reqOut)
	} else {
		// Bad Data
		err = e.SetData("application/json", []byte{1, 2, 3})
	}
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	slog.Info(fmt.Sprintf("Restock: [PID %s] [%d -> %d]\n", reqIn.PID, result-reqIn.Qty, result))
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
