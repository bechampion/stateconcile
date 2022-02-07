package main

import (
	googleactions "./googleactions"
	logging "./logging"
	"encoding/json"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"fmt"
	"os"
	"flag"
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

func BuildRules(terraformstate string) (int , map[string]bool) {
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
	return tfw.Version,retlist
}
func FindRulesLogging(targetrules []string, retlist map[string]bool,logs bool) []string{
	var notfound []string
	for i := 0; i < len(targetrules); i++ {
		if _, ok := retlist[targetrules[i]]; ok {
		}else {
			notfound = append(notfound,targetrules[i])

		}
	}
	return notfound

}
func FindRules(targetrules []string, retlist map[string]bool) []string{
	var notfound []string
	for i := 0; i < len(targetrules); i++ {
		if _, ok := retlist[targetrules[i]]; ok {
		}else {
			notfound = append(notfound,targetrules[i])

		}
	}
	return notfound

}
func ParseGSUri(gsuri string) (string,string){
	urisplit := strings.Split(gsuri,"/")
	return urisplit[2],strings.Join(urisplit[3:],"/")
}

func AddRandoms(fwlist []string) []string{
	for i := 0 ; i < 6 ; i ++ {
		fwlist=append(fwlist,fmt.Sprintf("RandomRule%d",i))
	}
	return fwlist
}
func main() {
	var random = flag.Bool("random", false, "Insert Random Rules on GCP source")
	var logs = flag.Bool("logs", false, "find logs matching the resources")
	var project = flag.String("project", "" , "Name of GCP project")
	var state = flag.String("state", "" , "location gs://bucket/object")
	flag.Parse()
	bucket,object := ParseGSUri(*state)
	// this should raise something..
	if *project == "" {
		log.Fatal("Error: -project string , needed")
	}
	if *state== "" {
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
	version,tfrules := BuildRules("working.tfstate")
	fwlist := googleactions.GetFirewallRules(*project)
	if *random == true {
		yellow:= color.New(color.FgYellow).SprintFunc()
		fmt.Printf("[%s] >> Adding random soruce rules",yellow("INFO"))
		fwlist = AddRandoms(fwlist)
	}
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("\n[%s] >> Rules found in GCP but missing in Terraform State: %s on Version: %s\n",red("Warning"),red("gs://",bucket,"/",object),red(version))
	if *logs == true {
		foundrules := FindRulesLogging(fwlist,tfrules,logs)
	}else{
		foundrules := FindRules(fwlist,tfrules,logs)
	}
	for i := 0; i < len(foundrules) ; i ++ {
		fmt.Printf("%s google_compute_firewall.%s -> %s\n",red("*"),foundrules[i],red("missing"))
	}

	fmt.Printf("\n")


}
