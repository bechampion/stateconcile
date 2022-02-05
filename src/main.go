package main

import (
	googleactions "./jwt"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func oldmain() {

	// 	grawjson := []byte(`
	// [
	//   {
	//     "allowed": [
	//       {
	//         "IPProtocol": "icmp"
	//       }
	//     ],
	//     "creationTimestamp": "2022-01-22T05:08:57.206-08:00",
	//     "description": "Allow ICMP from anywhere",
	//     "direction": "INGRESS",
	//     "disabled": false,
	//     "id": "5368780110428225286",
	//     "kind": "compute#firewall",
	//     "logConfig": {
	//       "enable": false
	//     },
	//     "name": "default-allow-icmp",
	//     "network": "https://www.googleapis.com/compute/v1/projects/myfreegke/global/networks/default",
	//     "priority": 65534,
	//     "selfLink": "https://www.googleapis.com/compute/v1/projects/myfreegke/global/firewalls/default-allow-icmp",
	//     "sourceRanges": [
	//       "0.0.0.0/0"
	//     ]
	//   },
	//   {
	//     "allowed": [
	//       {
	//         "IPProtocol": "tcp",
	//         "ports": [
	//           "0-65535"
	//         ]
	//       },
	//       {
	//         "IPProtocol": "udp",
	//         "ports": [
	//           "0-65535"
	//         ]
	//       },
	//       {
	//         "IPProtocol": "icmp"
	//       }
	//     ],
	//     "creationTimestamp": "2022-01-22T05:08:57.142-08:00",
	//     "description": "Allow internal traffic on the default network",
	//     "direction": "INGRESS",
	//     "disabled": false,
	//     "id": "3501426508784844550",
	//     "kind": "compute#firewall",
	//     "logConfig": {
	//       "enable": false
	//     },
	//     "name": "default-allow-internal",
	//     "network": "https://www.googleapis.com/compute/v1/projects/myfreegke/global/networks/default",
	//     "priority": 65534,
	//     "selfLink": "https://www.googleapis.com/compute/v1/projects/myfreegke/global/firewalls/default-allow-internal",
	//     "sourceRanges": [
	//       "10.128.0.0/9"
	//     ]
	//   },
	//   {
	//     "allowed": [
	//       {
	//         "IPProtocol": "tcp",
	//         "ports": [
	//           "3389"
	//         ]
	//       }
	//     ],
	//     "creationTimestamp": "2022-01-22T05:08:57.185-08:00",
	//     "description": "Allow RDP from anywhere",
	//     "direction": "INGRESS",
	//     "disabled": false,
	//     "id": "5164992364728663814",
	//     "kind": "compute#firewall",
	//     "logConfig": {
	//       "enable": false
	//     },
	//     "name": "default-allow-rdp",
	//     "network": "https://www.googleapis.com/compute/v1/projects/myfreegke/global/networks/default",
	//     "priority": 65534,
	//     "selfLink": "https://www.googleapis.com/compute/v1/projects/myfreegke/global/firewalls/default-allow-rdp",
	//     "sourceRanges": [
	//       "0.0.0.0/0"
	//     ]
	//   },
	//   {
	//     "allowed": [
	//       {
	//         "IPProtocol": "tcp",
	//         "ports": [
	//           "22"
	//         ]
	//       }
	//     ],
	//     "creationTimestamp": "2022-01-22T05:08:57.163-08:00",
	//     "description": "Allow SSH from anywhere",
	//     "direction": "INGRESS",
	//     "disabled": false,
	//     "id": "5003472075088838406",
	//     "kind": "compute#firewall",
	//     "logConfig": {
	//       "enable": false
	//     },
	//     "name": "default-allow-ssh",
	//     "network": "https://www.googleapis.com/compute/v1/projects/myfreegke/global/networks/default",
	//     "priority": 65534,
	//     "selfLink": "https://www.googleapis.com/compute/v1/projects/myfreegke/global/firewalls/default-allow-ssh",
	//     "sourceRanges": [
	//       "0.0.0.0/0"
	//     ]
	//   },
	//   {
	//     "allowed": [
	//       {
	//         "IPProtocol": "tcp"
	//       },
	//       {
	//         "IPProtocol": "udp"
	//       },
	//       {
	//         "IPProtocol": "icmp"
	//       },
	//       {
	//         "IPProtocol": "esp"
	//       },
	//       {
	//         "IPProtocol": "ah"
	//       },
	//       {
	//         "IPProtocol": "sctp"
	//       }
	//     ],
	//     "creationTimestamp": "2022-02-02T03:10:31.662-08:00",
	//     "description": "",
	//     "direction": "INGRESS",
	//     "disabled": false,
	//     "id": "63054083075332168",
	//     "kind": "compute#firewall",
	//     "logConfig": {
	//       "enable": false
	//     },
	//     "name": "gke-myfreegke-c14da6a7-all",
	//     "network": "https://www.googleapis.com/compute/v1/projects/myfreegke/global/networks/default",
	//     "priority": 1000,
	//     "selfLink": "https://www.googleapis.com/compute/v1/projects/myfreegke/global/firewalls/gke-myfreegke-c14da6a7-all",
	//     "sourceRanges": [
	//       "10.96.0.0/14"
	//     ],
	//     "targetTags": [
	//       "gke-myfreegke-c14da6a7-node"
	//     ]
	//   },
	//   {
	//     "allowed": [
	//       {
	//         "IPProtocol": "tcp",
	//         "ports": [
	//           "22"
	//         ]
	//       }
	//     ],
	//     "creationTimestamp": "2022-02-02T03:10:31.744-08:00",
	//     "description": "",
	//     "direction": "INGRESS",
	//     "disabled": false,
	//     "id": "646495175142004808",
	//     "kind": "compute#firewall",
	//     "logConfig": {
	//       "enable": false
	//     },
	//     "name": "gke-myfreegke-c14da6a7-ssh",
	//     "network": "https://www.googleapis.com/compute/v1/projects/myfreegke/global/networks/default",
	//     "priority": 1000,
	//     "selfLink": "https://www.googleapis.com/compute/v1/projects/myfreegke/global/firewalls/gke-myfreegke-c14da6a7-ssh",
	//     "sourceRanges": [
	//       "35.238.135.36/32"
	//     ],
	//     "targetTags": [
	//       "gke-myfreegke-c14da6a7-node"
	//     ]
	//   },
	//   {
	//     "allowed": [
	//       {
	//         "IPProtocol": "tcp",
	//         "ports": [
	//           "1-65535"
	//         ]
	//       },
	//       {
	//         "IPProtocol": "udp",
	//         "ports": [
	//           "1-65535"
	//         ]
	//       },
	//       {
	//         "IPProtocol": "icmp"
	//       }
	//     ],
	//     "creationTimestamp": "2022-02-02T03:10:31.848-08:00",
	//     "description": "",
	//     "direction": "INGRESS",
	//     "disabled": false,
	//     "id": "6207547516365822024",
	//     "kind": "compute#firewall",
	//     "logConfig": {
	//       "enable": false
	//     },
	//     "name": "gke-myfreegke-c14da6a7-vms",
	//     "network": "https://www.googleapis.com/compute/v1/projects/myfreegke/global/networks/default",
	//     "priority": 1000,
	//     "selfLink": "https://www.googleapis.com/compute/v1/projects/myfreegke/global/firewalls/gke-myfreegke-c14da6a7-vms",
	//     "sourceRanges": [
	//       "10.128.0.0/9"
	//     ],
	//     "targetTags": [
	//       "gke-myfreegke-c14da6a7-node"
	//     ]
	//   }
	// ]
	// `)

	//get tf state fromf file
	statefile, _ := os.Open("terraform.tfstate")
	terraformrawstate, _ := ioutil.ReadAll(statefile)

	// gfw := []GcloudFirewallRule{}
	tfw := TerraformRawState{}

	// err := json.Unmarshal(grawjson, &gfw)
	// if err != nil {
	// 	log.Fatal(err)

	// }
	err := json.Unmarshal(terraformrawstate, &tfw)
	if err != nil {
		log.Fatal(err)

	}
	// for i := 0; i < len(gfw); i++ {
	// 	fmt.Printf("%s\n", gfw[i].Id)

	// }
	for i := 0; i < len(tfw.Resources); i++ {
		if tfw.Resources[i].Type == "google_compute_firewall" {
			for r := 0; r < len(tfw.Resources[i].Instances); r++ {
				if tfw.Resources[i].Instances[r].Attributes.Name == os.Args[1] {
					fmt.Println(tfw.Resources[i].Instances[r].Attributes.Id)
				}
			}
		}
	}
}
func BuildRules() map[string]bool {
	retlist := make(map[string]bool)
	statefile, _ := os.Open("terraform.tfstate")
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
	return retlist
}
func findRule(rule string) string {
	statefile, _ := os.Open("terraform.tfstate")
	terraformrawstate, _ := ioutil.ReadAll(statefile)
	tfw := TerraformRawState{}
	err := json.Unmarshal(terraformrawstate, &tfw)
	if err != nil {
		log.Fatal(err)

	}
	for i := 0; i < len(tfw.Resources); i++ {
		if tfw.Resources[i].Type == "google_compute_firewall" {
			for r := 0; r < len(tfw.Resources[i].Instances); r++ {
				if tfw.Resources[i].Instances[r].Attributes.Name == os.Args[1] {
					return tfw.Resources[i].Instances[r].Attributes.Id
				}
			}
		}
	}
	return "nothing"
}
func FindRule(name string, retlist map[string]bool) bool {
	if _, ok := retlist[name]; ok {
		return true
	}
	return false

}
func main() {
	rand.Seed(time.Now().Unix())
	reasons := []string{
		"allow-egress-supplychain-migration-2",
		"supplychain-uat-allow-ips",
		"allow-ingress-songdelivery-elasticsearch",
		"allow-egress-supplychain-custom1",
		"allow-egress-tenable-jira-to-be-deleted",
		"allow-egress-ispirer-test-2",
		"allow-ingress-statementcompression-http",
		"supplychain-prod-partners-http-ingress",
		"allow-egress-splunk-audittracking",
		"supplychain-prd-deliver-allow-ips",
		"allow-egress-vault-k8s-sql",
		"allow-egress-primepublishing-sql-3",
		"allow-ingress-vpn-scl-process",
		"allow-ingress-vault-sql",
		"allow-egress-supplychain-sqlproxy",
		"allow-egress-songdelivery-sql-2",
		"allow-egress-smartsuspense-postgresql-2",
		"allow-egress-appstaging",
		"allow-sftp-bmg",
		"allow-ingress-netapp-cvo-4",
		"allow-ingress-cashmatch-custom",
		"allow-ingress-vault-k8s-2",
		"allow-ingress-songs-es",
		"allow-egress-binge-k8s-test-sql",
		"allow-egress-rpa-custom1",
		"allow-ingress-sprint-cloudera",
		"allow-egress-tenable-scanner-2",
		"allow-ingress-smartmatch-http",
		"supplychain-prd-capture-allow-ips",
		"allow-ingress-deliver-proxy",
		"allow-egress-appsbmg",
		"allow-ingress-elicense-es-2",
		"allow-egress-supplychain-migration",
		"allow-ingress-grant-k8s",
		"allow-egress-bertelsmann-dns",
	}
	rules := BuildRules()
	var n int = 0
	for i := 0; i < 20; i++ {
		n = rand.Int() % len(reasons)
		fmt.Printf("%s --> %t\n", reasons[n], FindRule(reasons[n], rules))
	}
	googleactions.DownloadTerraformState(os.Stdout, "test-stateconcile", "realwinrm.py", "here")
	googleactions.GetFirewallRules("myfreegke")

}
