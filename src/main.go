package main

import (
	routines "./routines"
	googleactions"./googleactions"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
)
func main() {
	var nobanner = flag.Bool("nobanner", false, "Want to see banner?")
	var random = flag.Bool("random", false, "Insert Random Rules on GCP source")
	var jsonoutput = flag.Bool("json", false, "Json Output")
	var logs = flag.Bool("logs", false, "find logs matching the resources")
	var ignoreauto = flag.Bool("ignoreauto", false, "Ignore GKE rules and others")
	var project = flag.String("project", "", "Name of GCP project")
	var state = flag.String("state", "", "location gs://bucket/object")
	flag.Parse()

	bucket, object := routines.ParseGSUri(*state)
	if *project == "" {
		log.Fatal("Error: -project string , needed")
	}
	if *state == "" {
		log.Fatal("Error: -state gs://....., needed")
	}
	if !(*nobanner || *jsonoutput) {
		routines.Banner()
	}
	err := googleactions.DownloadTerraformState(os.Stdout, bucket, object, "working.tfstate", *jsonoutput)
	if err != nil {
		log.Fatal("\nFailed at downloading state")

	}
	version, tfrules := routines.BuildRules("working.tfstate")
	fwlist := googleactions.GetFirewallRules(*project, *ignoreauto, *jsonoutput)
	if *random == true && *jsonoutput == false {
		yellow := color.New(color.FgYellow).SprintFunc()
		fmt.Printf("[%s] Adding random soruce rules", yellow("*"))
		fwlist = routines.AddRandoms(fwlist)
	}
	foundrules := routines.FindRules(fwlist, tfrules, *logs, *project, *jsonoutput, bucket, object, version)
	func(){
		if *jsonoutput {
			routines.Output(foundrules, "json")
			return 
		}
		routines.Output(foundrules, "text")
		return 

	}()

}
