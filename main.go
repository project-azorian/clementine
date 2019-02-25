package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// NodeStruct defines a node
type NodeStruct struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var nodes []NodeStruct

func main() {
	// Provide some junk data to play with
	nodes = append(nodes, NodeStruct{ID: 1, Name: "Node one"},
		NodeStruct{ID: 2, Name: "Node two"},
		NodeStruct{ID: 3, Name: "Node three"})

	// Create a router instance to route requests
	router := mux.NewRouter()

	// setup request routers
	router.HandleFunc("/nodes", getNodes).Methods("GET")
	router.HandleFunc("/nodes/{id}", getNode).Methods("GET")
	router.HandleFunc("/nodes", addNode).Methods("POST")
	router.HandleFunc("/nodes", updateNode).Methods("PUT")
	router.HandleFunc("/nodes/{id}", deleteNode).Methods("DELETE")

	// Start HTTP server
	log.Println("Starting up")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getNodes(w http.ResponseWriter, r *http.Request) {
	// getNodes: returns a json object containing all nodes
	log.Println("Get all nodes")
	json.NewEncoder(w).Encode(nodes)
}

func getNode(w http.ResponseWriter, r *http.Request) {
	// getNode: returns a json object containing a node
	log.Println("Get node")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for _, node := range nodes {
		if node.ID == id {
			json.NewEncoder(w).Encode(&node)
		}
	}
}

func addNode(w http.ResponseWriter, r *http.Request) {
	// addNode: takes a json object to add a node, returns a json object containing all nodes
	log.Println("Add node")
	var node NodeStruct
	json.NewDecoder(r.Body).Decode(&node)

	nodes = append(nodes, node)
	json.NewEncoder(w).Encode(nodes)
}

func updateNode(w http.ResponseWriter, r *http.Request) {
	// updateNode: takes a json object to update node, returns a json object containing all nodes
	log.Println("Update node")
	var node NodeStruct
	json.NewDecoder(r.Body).Decode(&node)

	for i, item := range nodes {
		if node.ID == item.ID {
			nodes[i] = node
		}
	}
	json.NewEncoder(w).Encode(nodes)
}

func deleteNode(w http.ResponseWriter, r *http.Request) {
	// deleteNode: takes a uri hit to delete a node, returns a json object containing all nodes
	log.Println("Delete node")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for i, item := range nodes {
		if id == item.ID {
			nodes = append(nodes[:i], nodes[i+1:]...)
		}
	}
	json.NewEncoder(w).Encode(nodes)
}
