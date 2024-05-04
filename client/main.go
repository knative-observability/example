package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

type Order struct {
	PID string `json:"pid"`
	UID string `json:"uid"`
	Qty int64  `json:"qty"`
}

type Restock struct {
	PID string `json:"pid"`
	Qty int64  `json:"qty"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	mode := r.Form.Get("mode")
	uid := r.Form.Get("uid")
	pid := r.Form.Get("pid")
	qty := r.Form.Get("qty")
	qtyInt, err := strconv.ParseInt(qty, 10, 64)
	if err != nil {
		fmt.Println("failed to convert qty to int64, ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if mode == "restock" {
		event := cloudevents.NewEvent()
		event.SetID(uuid.NewString())
		event.SetType("com.example.merchant")
		event.SetSource("example-client")
		event.SetData("application/json", Restock{PID: pid, Qty: qtyInt})
		c, err := cloudevents.NewClientHTTP()
		if err != nil {
			fmt.Println("failed to create client, ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sendCtx := cloudevents.ContextWithTarget(context.Background(), "http://broker-ingress.knative-eventing.svc.cluster.local/example/default")
		if result := c.Send(sendCtx, event); cloudevents.IsUndelivered(result) {
			fmt.Printf("failed to send, %v", result)
		} else {
			fmt.Printf("sent: %v", event)
		}
	} else if mode == "buy" {
		event := cloudevents.NewEvent()
		event.SetID(uuid.NewString())
		event.SetType("com.example.user")
		event.SetSource("example-client")
		event.SetData("application/json", Order{PID: pid, UID: uid, Qty: qtyInt})
		c, err := cloudevents.NewClientHTTP()
		if err != nil {
			fmt.Println("failed to create client, ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sendCtx := cloudevents.ContextWithTarget(context.Background(), "http://broker-ingress.knative-eventing.svc.cluster.local/example/default")
		if result := c.Send(sendCtx, event); cloudevents.IsUndelivered(result) {
			fmt.Printf("failed to send, %v", result)
		} else {
			fmt.Printf("sent: %v", event)
		}

	}
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
