package main

import (
	"fmt"
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
)

func main() {
	// Uncomment following line to enable logging
	// gocb.SetLogger(gocb.VerboseStdioLogger())

	// Update this to your cluster details
	connectionString := "couchbase://172.99.99.2?connect_timeout=40s&kv_timeout=7500&query_timeout=60000"
	// connectionString := "couchbase://couchbase?connect_timeout=40s&kv_timeout=7500&query_timeout=60000"
	bucketName := "travel-sample"
	username := "Administrator"
	password := "123456"

	// For a secure cluster connection, use `couchbases://<your-cluster-ip>` instead.
	cluster, err := gocb.Connect(connectionString, gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	bucket := cluster.Bucket(bucketName)

	err = bucket.WaitUntilReady(50*time.Second, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Get a reference to the default collection, required for older Couchbase server versions
	// col := bucket.DefaultCollection()

	col := bucket.Scope("tenant_agent_00").Collection("users")

	type User struct {
		Name      string   `json:"name"`
		Email     string   `json:"email"`
		Interests []string `json:"interests"`
	}

	// Create and store a Document
	_, err = col.Upsert("u:jade",
		User{
			Name:      "Jade",
			Email:     "jade@test-email.com",
			Interests: []string{"Swimming", "Rowing"},
		}, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Get the document back
	getResult, err := col.Get("u:jade", nil)
	if err != nil {
		log.Fatal(err)
	}

	var inUser User
	err = getResult.Content(&inUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User: %v\n", inUser)

	// Perform a N1QL Query
	inventoryScope := bucket.Scope("inventory")
	queryResult, err := inventoryScope.Query(
		fmt.Sprintf("SELECT * FROM airline WHERE id=10"),
		&gocb.QueryOptions{Adhoc: true},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Print each found Row
	for queryResult.Next() {
		var result interface{}
		err := queryResult.Row(&result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}

	if err := queryResult.Err(); err != nil {
		log.Fatal(err)
	}
}
