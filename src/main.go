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
	// 	"allow-egress-supplychain-migration-2",
	// 	"supplychain-uat-allow-ips",
	// 	"allow-ingress-songdelivery-elasticsearch",
	// 	"allow-egress-supplychain-custom1",
	// 	"allow-egress-tenable-jira-to-be-deleted",
	// 	"allow-egress-ispirer-test-2",
	// 	"allow-ingress-statementcompression-http",
	// 	"supplychain-prod-partners-http-ingress",
	// 	"allow-egress-splunk-audittracking",
	// 	"supplychain-prd-deliver-allow-ips",
	// 	"allow-egress-vault-k8s-sql",
	// 	"allow-egress-primepublishing-sql-3",
	// 	"allow-ingress-vpn-scl-process",
	// 	"allow-ingress-vault-sql",
	// 	"allow-egress-supplychain-sqlproxy",
	// 	"allow-egress-songdelivery-sql-2",
	// 	"allow-egress-smartsuspense-postgresql-2",
	// 	"allow-egress-appstaging",
	// 	"allow-sftp-bmg",
	// 	"allow-ingress-netapp-cvo-4",
	// 	"allow-ingress-cashmatch-custom",
	// 	"allow-ingress-vault-k8s-2",
	// 	"allow-ingress-songs-es",
	// 	"allow-egress-binge-k8s-test-sql",
	// 	"allow-egress-rpa-custom1",
	// 	"allow-ingress-sprint-cloudera",
	// 	"allow-egress-tenable-scanner-2",
	// 	"allow-ingress-smartmatch-http",
	// 	"supplychain-prd-capture-allow-ips",
	// 	"allow-ingress-deliver-proxy",
	// 	"allow-egress-appsbmg",
	// 	"allow-ingress-elicense-es-2",
	// 	"allow-egress-supplychain-migration",
	// 	"allow-ingress-grant-k8s",
	// 	"allow-egress-bertelsmann-dns",
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
