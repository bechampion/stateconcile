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

func Banner() {
	red := color.New(color.FgRed)
	name := `
  ██████ ▄▄▄█████▓ ▄▄▄     ▄▄▄█████▓▓█████  ▄████▄   ▒█████   ███▄    █  ▄████▄   ██▓ ██▓    ▓█████ 
▒██    ▒ ▓  ██▒ ▓▒▒████▄   ▓  ██▒ ▓▒▓█   ▀ ▒██▀ ▀█  ▒██▒  ██▒ ██ ▀█   █ ▒██▀ ▀█  ▓██▒▓██▒    ▓█   ▀ 
░ ▓██▄   ▒ ▓██░ ▒░▒██  ▀█▄ ▒ ▓██░ ▒░▒███   ▒▓█    ▄ ▒██░  ██▒▓██  ▀█ ██▒▒▓█    ▄ ▒██▒▒██░    ▒███   
  ▒   ██▒░ ▓██▓ ░ ░██▄▄▄▄██░ ▓██▓ ░ ▒▓█  ▄ ▒▓▓▄ ▄██▒▒██   ██░▓██▒  ▐▌██▒▒▓▓▄ ▄██▒░██░▒██░    ▒▓█  ▄ 
▒██████▒▒  ▒██▒ ░  ▓█   ▓██▒ ▒██▒ ░ ░▒████▒▒ ▓███▀ ░░ ████▓▒░▒██░   ▓██░▒ ▓███▀ ░░██░░██████▒░▒████▒
▒ ▒▓▒ ▒ ░  ▒ ░░    ▒▒   ▓▒█░ ▒ ░░   ░░ ▒░ ░░ ░▒ ▒  ░░ ▒░▒░▒░ ░ ▒░   ▒ ▒ ░ ░▒ ▒  ░░▓  ░ ▒░▓  ░░░ ▒░ ░
░ ░▒  ░ ░    ░      ▒   ▒▒ ░   ░     ░ ░  ░  ░  ▒     ░ ▒ ▒░ ░ ░░   ░ ▒░  ░  ▒    ▒ ░░ ░ ▒  ░ ░ ░  ░
░  ░  ░    ░        ░   ▒    ░         ░   ░        ░ ░ ░ ▒     ░   ░ ░ ░         ▒ ░  ░ ░      ░   
`
	red.Println(name)
	fmt.Printf("\n(Finds differences between state and gcp and more..)\n")
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
		fmt.Printf("[%s] Logs enabled , only searching within 24hs and less than 200 records for the given type\n", yellow("*"))
		hashedlogs := logging.HashedLoggingEntries("myfreegke")
		_ = hashedlogs

		for i := 0; i < len(targetrules); i++ {
			if _, ok := retlist[targetrules[i]]; ok {
			} else {
				if _, ok := hashedlogs[fmt.Sprintf("projects/%s/global/firewalls/%s", *project, targetrules[i])]; ok {
					notfound[targetrules[i]] = hashedlogs[fmt.Sprintf("projects/%s/global/firewalls/%s", *project, targetrules[i])]
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
	var nobanner = flag.Bool("nobanner", true, "Want to see banner?")
	var random = flag.Bool("random", false, "Insert Random Rules on GCP source")
	var logs = flag.Bool("logs", false, "find logs matching the resources")
	var ignoreauto = flag.Bool("ignoreauto", false, "Ignore GKE rules and others")
	var project = flag.String("project", "", "Name of GCP project")
	var state = flag.String("state", "", "location gs://bucket/object")
	flag.Parse()
	bucket, object := ParseGSUri(*state)
	if *project == "" {
		log.Fatal("Error: -project string , needed")
	}
	if *state == "" {
		log.Fatal("Error: -state gs://....., needed")
	}
	if *nobanner {
		Banner()
	}
	_ = googleactions.DownloadTerraformState(os.Stdout, bucket, object, "working.tfstate")
	version, tfrules := BuildRules("working.tfstate")
	fwlist := googleactions.GetFirewallRules(*project, *ignoreauto)
	if *random == true {
		yellow := color.New(color.FgYellow).SprintFunc()
		fmt.Printf("[%s] Adding random soruce rules", yellow("*"))
		fwlist = AddRandoms(fwlist)
	}
	red := color.New(color.FgRed).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	foundrules := FindRules(fwlist, tfrules, logs, project)
	if len(foundrules) > 0 {
		fmt.Printf("[%s] %s Rules found in GCP but missing in Terraform State: %s on Version: %s\n", red("*"), red(len(foundrules)),red("gs://", bucket, "/", object), red(version))
	}
	for k, v := range foundrules {
		fmt.Printf("%s google_compute_firewall.%s -> %s\n", red("*"), k, red("missing"))
		if v != (logging.Payload{}) {
			fmt.Printf("\t[%s]\n", blue("logs"))
			fmt.Printf("\tTimeStamp:%s\n\tServiceName:%s\n\tResourceName:%s\n\tMethodName:%s\n\tPrincipalEmail:%s\n\tCallerIP:%s\n\tUserAgent:%s\n", red(v.TimeStamp), v.Payload.ServiceName, red(v.Payload.ResourceName), v.Payload.MethodName, red(v.Payload.AuthenticationInfo.PrincipalEmail), v.Payload.RequestMetaData.CallerIP, v.Payload.RequestMetaData.CallerSuppliedUserAgent)
		}
	}

	fmt.Printf("\n")

}
