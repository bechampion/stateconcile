package logging

import (
	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"time"
)

type Payload struct {
	TimeStamp string `json:"timestamp"`
	Payload   struct {
		ServiceName        string `json:"service_name"`
		MethodName         string `json:"method_name"`
		ResourceName       string `json:"resource_name"`
		AuthenticationInfo struct {
			PrincipalEmail string `json:"principal_email"`
		} `json:"authentication_info"`
		RequestMetaData struct {
			CallerIP                string `json:"caller_ip"`
			CallerSuppliedUserAgent string `json:"caller_supplied_user_agent"`
		} `json:"request_metadata"`
	} `json:"payload"`
}

func GetLogEntries(projID string) ([]*logging.Entry, error) {
	currentTime := time.Now()
	today := currentTime.Format("2006-01-02")

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
		logadmin.Filter(fmt.Sprintf(`resource.type = "gce_firewall_rule" protoPayload.methodName:"v1.compute.firewalls.insert" timestamp > "%s"`, today)),
		logadmin.NewestFirst(),
	)

	for len(entries) < 200 {
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
}
func HashedLoggingEntries(projID string) map[string]Payload {
	hashedloggingentries := map[string]Payload{}
	payload := Payload{}
	entries, _ := GetLogEntries(projID)
	for _, entry := range entries {
		pp, _ := json.Marshal(entry)
		_ = json.Unmarshal(pp, &payload)
		if entry.Operation.First == true {
			hashedloggingentries[payload.Payload.ResourceName] = payload
			// 	fmt.Println(payload.ResourceName)
			// 	fmt.Println(payload.ServiceName)
			// 	fmt.Println(payload.AuthenticationInfo.PrincipalEmail)
		}
	}
	return hashedloggingentries
}
func main() {
	fmt.Println(HashedLoggingEntries("myfreegke"))
	ll := HashedLoggingEntries("myfreegke")
	for k, v := range ll {
		fmt.Printf("%s ---> %v", k, v)

	}
}
