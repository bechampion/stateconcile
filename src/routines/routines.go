package routines

import (
	logging "../logging"
	"encoding/json"
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

func Output(foundrules map[string]interface{}, jsonoutput bool) {
	if !jsonoutput {
		red := color.New(color.FgRed).SprintFunc()
		blue := color.New(color.FgBlue).SprintFunc()
		for k, v := range foundrules {
			fmt.Printf("%s google_compute_firewall.%s -> %s\n", red("*"), k, red("missing"))
			if v != false {
				vv := v.(logging.Payload)
				fmt.Printf("\t[%s]\n", blue("logs"))
				fmt.Printf("\tTimeStamp:%s\n\tServiceName:%s\n\tResourceName:%s\n\tMethodName:%s\n\tPrincipalEmail:%s\n\tCallerIP:%s\n\tUserAgent:%s\n", red(vv.TimeStamp), vv.Payload.ServiceName, red(vv.Payload.ResourceName), vv.Payload.MethodName, red(vv.Payload.AuthenticationInfo.PrincipalEmail), vv.Payload.RequestMetaData.CallerIP, vv.Payload.RequestMetaData.CallerSuppliedUserAgent)
			}
		}

		fmt.Printf("\n")
	} else {
		jj, err := json.Marshal(foundrules)
		if err != nil {
			log.Fatal("Something went wrong")

		}
		fmt.Println(string(jj))
	}
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
func FindRules(targetrules []string, retlist map[string]bool, logs bool, project string, jsonoutput bool, bucket string, object string, version int) map[string]interface{}{
	notfound := make(map[string]interface{})
	if logs == true {
		if jsonoutput == false {
			yellow := color.New(color.FgYellow).SprintFunc()
			fmt.Printf("[%s] Logs enabled , only searching within 24hs and less than 200 records for the given type\n", yellow("*"))
		}
		hashedlogs := logging.HashedLoggingEntries("myfreegke")
		_ = hashedlogs

		for i := 0; i < len(targetrules); i++ {
			if _, ok := retlist[targetrules[i]]; ok {
			} else {
				if _, ok := hashedlogs[fmt.Sprintf("projects/%s/global/firewalls/%s", project, targetrules[i])]; ok {
					notfound[targetrules[i]] = hashedlogs[fmt.Sprintf("projects/%s/global/firewalls/%s", project, targetrules[i])]
				} else {
					notfound[targetrules[i]]=false
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
	if len(notfound) > 0 && jsonoutput == false {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Printf("[%s] %s Rules found in GCP but missing in Terraform State: %s on Version: %s\n", red("*"), red(len(notfound)), red("gs://", bucket, "/", object), red(version))
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
