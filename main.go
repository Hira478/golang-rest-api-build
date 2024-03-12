package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Order struct {
	OrderID      string    `json:"orderId"`
	CustomerName string    `json:"customerName"`
	OrderedAt    time.Time `json:"orderedAt"`
	Items        []Item    `json:"items"`
}

type Item struct {
	LineItemID  int    `json:"lineItemId"`
	ItemCode    int    `json:"itemCode"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

var orders []Order

func createOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newOrder.CustomerName == "" {
		http.Error(w, "CustomerName cannot be empty", http.StatusBadRequest)
		return
	}

	orders = append(orders, newOrder)

	w.WriteHeader(http.StatusCreated)
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(orders)
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["orderId"]

	var updatedOrder Order
	err := json.NewDecoder(r.Body).Decode(&updatedOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updatedOrder.CustomerName == "" {
		http.Error(w, "CustomerName cannot be empty", http.StatusBadRequest)
		return
	}

	for i, o := range orders {
		if o.OrderID == orderID {
			orders[i] = updatedOrder
			break
		}
	}

	w.WriteHeader(http.StatusOK)
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["orderId"]

	orderExists := false
	for _, o := range orders {
		if o.OrderID == orderID {
			orderExists = true
			break
		}
	}
	if !orderExists {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	for i, o := range orders {
		if o.OrderID == orderID {
			orders = append(orders[:i], orders[i+1:]...)
			break
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/orders", createOrder).Methods("POST")
	r.HandleFunc("/orders", getOrders).Methods("GET")
	r.HandleFunc("/orders/{orderId}", updateOrder).Methods("PUT")
	r.HandleFunc("/orders/{orderId}", deleteOrder).Methods("DELETE")

	fmt.Println("Server is running on port 8080")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
