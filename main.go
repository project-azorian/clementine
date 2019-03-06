package main

import (
	"encoding/json"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// NodeStruct defines a node
type NodeStruct struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ActionResponse defines a generic response to a request
type ActionResponse struct {
	ID      int64  `json:"id"`
	Action  string `json:"action"`
	Message string `json:"message"`
}

var nodes []NodeStruct
var db *sql.DB
var err error

func init() {
	// Setup primary logger
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})
	log.Println("Setup Logger")
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	// Setup DB connection
	db, err = sql.Open("mysql", "clementine:password@tcp(mariadb.openstack.svc.cluster.local:3306)/clementine")
	logFatal(err)
	// Test connection
	err = db.Ping()
	logFatal(err)

	// Create a router instance to route requests
	router := mux.NewRouter()

	// setup request routers
	router.HandleFunc("/nodes", getNodes).Methods("GET")
	router.HandleFunc("/nodes/{id}", getNode).Methods("GET")
	router.HandleFunc("/nodes", addNode).Methods("POST")
	router.HandleFunc("/nodes", updateNode).Methods("PUT")
	router.HandleFunc("/nodes/{id}", deleteNode).Methods("DELETE")

	// Add middleware handlers
	loggingHandler := handlers.CombinedLoggingHandler(os.Stdout, router)
	compressHandler := handlers.CompressHandler(loggingHandler)
	proxyHeadersHandler := handlers.ProxyHeaders(compressHandler)

	// Start HTTP server
	log.Println("Starting to serve")
	log.Fatal(http.ListenAndServe(":8000", proxyHeadersHandler))
}

func getNodes(w http.ResponseWriter, r *http.Request) {
	// getNodes: returns a json object containing all nodes
	log.Println("Get all nodes")
	var node NodeStruct
	nodes = []NodeStruct{}

	rows, err := db.Query("select * from nodes")
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&node.ID, &node.Name)
		logFatal(err)

		nodes = append(nodes, node)
	}

	json.NewEncoder(w).Encode(nodes)
}

func getNode(w http.ResponseWriter, r *http.Request) {
	// getNode: returns a json object containing a node
	log.Println("Get node")
	var node NodeStruct
	params := mux.Vars(r)

	rows := db.QueryRow("select * from nodes where id=?", params["id"])

	err := rows.Scan(&node.ID, &node.Name)
	logFatal(err)

	json.NewEncoder(w).Encode(node)

}

func addNode(w http.ResponseWriter, r *http.Request) {
	// addNode: takes a json object to add a node, returns a json object containing all nodes
	var node NodeStruct
	var nodeID int64
	var actionResponse ActionResponse
	actionResponse.Action = "addNode"
	log.Println(actionResponse.Action)
	json.NewDecoder(r.Body).Decode(&node)

	res, err := db.Exec("insert into nodes(name) values(?)", node.Name)
	if err != nil {
		log.Fatal("Exec err:", err.Error())
	} else {
		nodeID, err = res.LastInsertId()
		if err != nil {
			log.Fatal("Error:", err.Error())
		}
	}
	actionResponse.ID = nodeID
	json.NewEncoder(w).Encode(actionResponse)
}

func updateNode(w http.ResponseWriter, r *http.Request) {
	// updateNode: takes a json object to update node, returns a json object containing all nodes
	var node NodeStruct
	var rowsUpdated int64
	var actionResponse ActionResponse
	actionResponse.Action = "updateNode"
	log.Println(actionResponse.Action)
	json.NewDecoder(r.Body).Decode(&node)

	res, err := db.Exec("update nodes set name=? where id=?", &node.Name, &node.ID)
	if err != nil {
		log.Fatal("Exec err:", err.Error())
	} else {
		rowsUpdated, err = res.RowsAffected()
		if err != nil {
			log.Fatal("Error:", err.Error())
		}
	}

	actionResponse.ID = node.ID
	actionResponse.Message = "row updated"
	json.NewEncoder(w).Encode(rowsUpdated)
}

func deleteNode(w http.ResponseWriter, r *http.Request) {
	// deleteNode: takes a uri hit to delete a node, returns a json object containing all nodes
	log.Println("Delete node")
	var rowsDeleted int64
	params := mux.Vars(r)

	res, err := db.Exec("delete from nodes where id = ?", params["id"])
	if err != nil {
		log.Fatal("Exec err:", err.Error())
	} else {
		rowsDeleted, err = res.RowsAffected()
		if err != nil {
			log.Fatal("Error:", err.Error())
		}
	}

	json.NewEncoder(w).Encode(rowsDeleted)
}
