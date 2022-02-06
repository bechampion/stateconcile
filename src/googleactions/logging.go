package main

import (
	"fmt"
	"log"
	// "time"

	// [START imports]
	"context"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"

	"encoding/json"
	"google.golang.org/api/iterator"
	// [END imports]
)

type Payload struct {
	ServiceName        string `json:"service_name"`
	MethodName         string `json:"method_name"`
	ResourceName       string `json:"resource_name"`
	AuthenticationInfo struct {
		PrincipalEmail string `json:"principal_email"`
	} `json:"authentication_info"`
}

func main() {
	payload := Payload{}
	entries, _ := getEntries("myfreegke")
	log.Printf("Found %d entries.", len(entries))
	for _, entry := range entries {
		// fmt.Printf("Entry: @%s: %v\n",
		// 	entry.Timestamp.Format(time.RFC3339),
		// 	entry.Payload)
		pp, _ := json.Marshal(entry.Payload)
		_ = json.Unmarshal(pp, &payload)
		if entry.Operation.First == true {
			fmt.Println(payload.ResourceName)
			fmt.Println(payload.ServiceName)
			fmt.Println(payload.AuthenticationInfo.PrincipalEmail)
		}
	}

}

func getEntries(projID string) ([]*logging.Entry, error) {
	ctxx := context.Background()
	client, err := logging.NewClient(ctxx, projID)
	if err != nil {
		log.Fatalf("Failed to create logging client: %v", err)
	}
	defer client.Close()

	adminClient, err := logadmin.NewClient(ctxx, projID)
	if err != nil {
		log.Fatalf("Failed to create logadmin client: %v", err)
	}
	defer adminClient.Close()
	ctx := context.Background()

	// [START logging_list_log_entries]
	var entries []*logging.Entry
	const name = "log-example"
	// lastHour := time.Now().Add(-5 * time.Hour).Format(time.RFC3339)

	iter := adminClient.Entries(ctx,
		logadmin.Filter(fmt.Sprintf(`resource.type = "gce_firewall_rule" protoPayload.methodName:"v1.compute.firewalls.insert"`)),
		logadmin.NewestFirst(),
	)

	// Fetch the most recent 20 entries.
	for len(entries) < 20 {
		entry, err := iter.Next()
		if err == iterator.Done {
			return entries, nil
		}
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
	// [END logging_list_log_entries]
}
