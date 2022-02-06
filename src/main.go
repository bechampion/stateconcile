package main

import (
	googleactions "./googleactions"
	"encoding/json"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"fmt"
	"math/rand"
	"os"
	"time"
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

type GcloudFirewallRule struct {
	Id        string `json:"id"`
	Disabled  bool   `json:"disabled"`
	Direction string `json:"direction"`
}
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
func main() {
	rand.Seed(time.Now().Unix())
	// reasons := []string{
	// 	"allow-some-rule",
	// }
	// // var n int = 0
	// for i := 0; i < 20; i++ {
	// 	n = rand.Int() % len(reasons)
	// 	fmt.Printf("%s --> %t\n", reasons[n], FindRule(reasons[n], rules))
	// }
	terraformrealstate :=  "terraform.tfstate"
	_ = googleactions.DownloadTerraformState(os.Stdout, "test-stateconcile", "terraform.tfstate", "working.tfstate")
	version,tfrules := BuildRules("working.tfstate")
	fwlist := googleactions.GetFirewallRules("myfreegke")
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("\n[%s] >> Rules found in GCP but missing in Terraform State: %s on Version: %s\n",red("Warning"),red(terraformrealstate),red(version))
	foundrules := FindRules(fwlist,tfrules)
	for i := 0; i < len(foundrules) ; i ++ {
		fmt.Printf("%s google_compute_firewall.%s -> %s\n",red("*"),foundrules[i],red("missing"))
	}
	fmt.Printf("\n")


}
