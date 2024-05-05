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

type Order struct {
	UUID string    `json:"uuid"`
	Time time.Time `json:"time"`
	PID  string    `json:"pid"`
	UID  string    `json:"uid"`
	Qty  int64     `json:"qty"`
}

type StockResp struct {
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
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	defer client.Close()

	// Set locks to prevent data race
	for {
		locked, err := client.SetNX(ctx, "lock", uuid.NewString(), 10*time.Second).Result()
		if err != nil {
			slog.Error("Can's set lock: " + err.Error())
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
	stock, err := client.HGet(ctx, "stock", order.PID).Int64()
	if err != nil {
		slog.Warn("Cant't get stock: " + err.Error() + ", set to 0")
		stock = 0
	}
	if stock < order.Qty {
		slog.Warn(fmt.Sprintf("Stock not enough: [UUID %s] [PID: %s] %d < %d\n",
			order.UUID[:8], order.PID, stock, order.Qty))
		e.SetData("application/json", StockResp{Upstream: "update-stock", Order: order, Success: false})
		e.SetType("com.example.verify")
		return &e, nil
	}
	result, err := client.HIncrBy(ctx, "stock", order.PID, -order.Qty).Result()
	if err != nil {
		slog.Error("Can't decrease stock" + err.Error())
		return nil, err
	}
	err = e.SetData("application/json", StockResp{Upstream: "update-stock", Order: order, Success: true})
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	slog.Info(fmt.Sprintf("Update Stock: User [%s], [PID %s] [%d -> %d]\n",
		order.UID, order.PID, result+order.Qty, result))
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
