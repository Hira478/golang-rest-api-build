package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Order represents the structure of an order
type Order struct {
	OrderID      string    `json:"orderId"`
	CustomerName string    `json:"customerName"`
	OrderedAt    time.Time `json:"orderedAt"`
	Items        []Item    `json:"items"`
}

// Item represents the structure of an item in an order
type Item struct {
	LineItemID  int    `json:"lineItemId"`
	ItemCode    int    `json:"itemCode"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

var orders []Order

// createOrder handles the creation of a new order
func createOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder Order

	// Decode the JSON request body into the newOrder struct
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Basic validation: Check if CustomerName is not empty
	if newOrder.CustomerName == "" {
		http.Error(w, "CustomerName cannot be empty", http.StatusBadRequest)
		return
	}

	// Append the new order to the orders slice
	orders = append(orders, newOrder)

	w.WriteHeader(http.StatusCreated)
}

// getOrders returns the list of orders
func getOrders(w http.ResponseWriter, r *http.Request) {
	// Encode the orders slice as JSON and write it to the response
	json.NewEncoder(w).Encode(orders)
}

// updateOrder handles the update of an existing order
func updateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["orderId"]

	var updatedOrder Order

	// Decode the JSON request body into the updatedOrder struct
	err := json.NewDecoder(r.Body).Decode(&updatedOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Basic validation: Check if CustomerName is not empty
	if updatedOrder.CustomerName == "" {
		http.Error(w, "CustomerName cannot be empty", http.StatusBadRequest)
		return
	}

	// Find and update the order with the given orderID
	for i, o := range orders {
		if o.OrderID == orderID {
			orders[i] = updatedOrder
			break
		}
	}

	w.WriteHeader(http.StatusOK)
}

// deleteOrder handles the deletion of an existing order
func deleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["orderId"]

	// Basic validation: Check if order with orderID exists
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

	// Find and delete the order with the given orderID
	for i, o := range orders {
		if o.OrderID == orderID {
			orders = append(orders[:i], orders[i+1:]...)
			break
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// Create a new router using the Gorilla Mux package
	r := mux.NewRouter()

	// Define the RESTful API endpoints and their corresponding handlers
	r.HandleFunc("/orders", createOrder).Methods("POST")
	r.HandleFunc("/orders", getOrders).Methods("GET")
	r.HandleFunc("/orders/{orderId}", updateOrder).Methods("PUT")
	r.HandleFunc("/orders/{orderId}", deleteOrder).Methods("DELETE")

	// Start the server and listen on port 8080
	fmt.Println("Server is running on port 8080")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
