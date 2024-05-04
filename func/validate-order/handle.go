package function

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/redis/go-redis/v9"
)

type Order struct {
	UUID string    `json:"uuid"`
	Time time.Time `json:"time"`
	PID  string    `json:"pid"`
	UID  string    `json:"uid"`
	Qty  int64     `json:"qty"`
}

type InOrder struct {
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
	var inOrder InOrder
	err := e.DataAs(&inOrder)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	defer client.Close()

	// 0: waiting 1: fail 2: success
	// payment*=3, update-stock*=1
	if inOrder.Upstream == "payment" {
		if inOrder.Success {
			client.HIncrBy(ctx, "mq", inOrder.Order.UUID, 6)
			slog.Debug(fmt.Sprintf("Payment success for order %s", inOrder.Order.UUID))
		} else {
			client.HIncrBy(ctx, "mq", inOrder.Order.UUID, 3)
			slog.Debug(fmt.Sprintf("Payment failed for order %s", inOrder.Order.UUID))
		}
	} else if inOrder.Upstream == "update-stock" {
		if inOrder.Success {
			client.HIncrBy(ctx, "mq", inOrder.Order.UUID, 2)
			slog.Debug(fmt.Sprintf("Update stock success for order %s", inOrder.Order.UUID))
		} else {
			client.HIncrBy(ctx, "mq", inOrder.Order.UUID, 1)
			slog.Debug(fmt.Sprintf("Update stock failed for order %s", inOrder.Order.UUID))
		}
	}

	if inOrder.Upstream == "payment" {
		return nil, nil
	}

	timeout := time.After(1 * time.Minute)
	for {
		select {
		case <-timeout:
			estr := fmt.Sprintf("Message queue timeout for order %s", inOrder.Order.UUID)
			slog.Error(estr)
			return nil, fmt.Errorf(estr)
		default:
			stock, err := client.HGet(ctx, "mq", inOrder.Order.UUID).Int64()
			if err != nil {
				stock = 0
			}
			switch stock {
			case 8:
				goto exitLoop
			case 0, 2, 6:
				time.Sleep(10 * time.Millisecond)
			case 1, 7:
				estr := fmt.Sprintf("Upstream receive-order failed for order %s", inOrder.Order.UUID)
				slog.Error(estr)
				return nil, fmt.Errorf(estr)
			case 3, 5:
				estr := fmt.Sprintf("Upstream payment failed for order %s", inOrder.Order.UUID)
				slog.Error(estr)
				return nil, fmt.Errorf(estr)
			case 4:
				estr := fmt.Sprintf("Upstream receive-order and payment failed for order %s", inOrder.Order.UUID)
				slog.Error(estr)
				return nil, fmt.Errorf(estr)
			}
		}
	}
exitLoop:
	// Two upstream functions are both successful
	// TODO: restore stock when payment failed or timeout
	slog.Info(fmt.Sprintf("Verify succeeded: User [%s], Order ID [%s]\n", inOrder.Order.UID, inOrder.Order.UUID))
	e.SetType("com.example.todo")
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
