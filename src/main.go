package main

import (
	googleactions "./googleactions"
	logging "./logging"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type TerraformRawState struct {
	Version   int `json:"version"`
	Resources []struct {
		Type      string `json:"type"`
		Instances []struct {
			Attributes struct {
				Id        string `json:"id,omitempty"`
				Disabled  bool   `json:"disabled,omitempty"`
				Direction string `json:"direction,omitempty"`
				Name      string `json:"name,omitempty"`
			} `json:"attributes"`
		} `json:"instances"`
	} `json:"resources"`
}

// not unmarshalling gcloud firewall rules as using SDK
// type GcloudFirewallRule struct {
// 	Id        string `json:"id"`
// 	Disabled  bool   `json:"disabled"`
// 	Direction string `json:"direction"`
// }
func Banner() {
	red := color.New(color.FgRed)
	whiteBackground := red.Add(color.BgWhite)
	whiteBackground.Println("[[[[ StateConcile ]]]]")
	fmt.Printf("\n")
}

func BuildRules(terraformstate string) (int, map[string]bool) {
	retlist := make(map[string]bool)
	statefile, _ := os.Open(terraformstate)
	terraformrawstate, _ := ioutil.ReadAll(statefile)
	tfw := TerraformRawState{}
	err := json.Unmarshal(terraformrawstate, &tfw)
	if err != nil {
		log.Fatal(err)

	}
	for i := 0; i < len(tfw.Resources); i++ {
		if tfw.Resources[i].Type == "google_compute_firewall" {
			for r := 0; r < len(tfw.Resources[i].Instances); r++ {
				retlist[tfw.Resources[i].Instances[r].Attributes.Name] = true
			}
		}
	}
	return tfw.Version, retlist
}
func FindRules(targetrules []string, retlist map[string]bool, logs *bool, project *string) map[string]logging.Payload {
	notfound := make(map[string]logging.Payload)
	if *logs == true {
		yellow := color.New(color.FgYellow).SprintFunc()
		fmt.Printf("[%s] >> Logs enabled only searching withing 24hs ", yellow("INFO"))
		hashedlogs := logging.HashedLoggingEntries("myfreegke")
		_ = hashedlogs

		for i := 0; i < len(targetrules); i++ {
			if _, ok := retlist[targetrules[i]]; ok {
			} else {
				// notfound = append(notfound,targetrules[i])
				if _, ok := hashedlogs[fmt.Sprintf("projects/%s/global/firewalls/%s", *project , targetrules[i])]; ok {
					notfound[targetrules[i]] = hashedlogs[fmt.Sprintf("projects/%s/global/firewalls/%s",*project, targetrules[i])]
				} else {
					notfound[targetrules[i]] = logging.Payload{}
				}

			}
		}
	} else {
		for i := 0; i < len(targetrules); i++ {
			if _, ok := retlist[targetrules[i]]; ok {
			} else {
				notfound[targetrules[i]] = logging.Payload{}

			}
		}

	}
	return notfound

}
func ParseGSUri(gsuri string) (string, string) {
	urisplit := strings.Split(gsuri, "/")
	return urisplit[2], strings.Join(urisplit[3:], "/")
}

func AddRandoms(fwlist []string) []string {
	for i := 0; i < 6; i++ {
		fwlist = append(fwlist, fmt.Sprintf("RandomRule%d", i))
	}
	return fwlist
}
func main() {
	var random = flag.Bool("random", false, "Insert Random Rules on GCP source")
	var logs = flag.Bool("logs", false, "find logs matching the resources")
	var project = flag.String("project", "", "Name of GCP project")
	var state = flag.String("state", "", "location gs://bucket/object")
	flag.Parse()
	bucket, object := ParseGSUri(*state)
	// this should raise something..
	if *project == "" {
		log.Fatal("Error: -project string , needed")
	}
	if *state == "" {
		log.Fatal("Error: -state gs://....., needed")
	}
	// reasons := []string{
	// 	"allow-some-rule",
	// }
	// // var n int = 0
	// for i := 0; i < 20; i++ {
	// 	n = rand.Int() % len(reasons)
	// 	fmt.Printf("%s --> %t\n", reasons[n], FindRule(reasons[n], rules))
	// }
	_ = googleactions.DownloadTerraformState(os.Stdout, bucket, object, "working.tfstate")
	version, tfrules := BuildRules("working.tfstate")
	fwlist := googleactions.GetFirewallRules(*project)
	if *random == true {
		yellow := color.New(color.FgYellow).SprintFunc()
		fmt.Printf("[%s] >> Adding random soruce rules", yellow("INFO"))
		fwlist = AddRandoms(fwlist)
	}
	red := color.New(color.FgRed).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	fmt.Printf("\n[%s] >> Rules found in GCP but missing in Terraform State: %s on Version: %s\n", red("Warning"), red("gs://", bucket, "/", object), red(version))
	foundrules := FindRules(fwlist, tfrules, logs,project)
	for k, v := range foundrules {
		fmt.Printf("%s google_compute_firewall.%s -> %s\n", red("*"), k, red("missing"))
		if v != (logging.Payload{}) {
			fmt.Printf("\t[%s]\n",blue("logs"))
			fmt.Printf("\tTimeStamp:%s\n\tServiceName:%s\n\tResourceName:%s\n\tMethodName:%s\n\tPrincipalEmail:%s\n\tCallerIP:%s\n\tUserAgent:%s\n",v.TimeStamp,v.Payload.ServiceName , v.Payload.ResourceName , v.Payload.MethodName , v.Payload.AuthenticationInfo.PrincipalEmail,v.Payload.RequestMetaData.CallerIP,v.Payload.RequestMetaData.CallerSuppliedUserAgent)
		}
	}

	fmt.Printf("\n")

}
