# stateconcile

Look into GCP and try to match against the terraform state provided ,  if can't find based on Name , prints them out


## Usage
```bash
$> GOOGLE_APPLICATION_CREDENTIALS=~/Downloads/myfreegke-a9a1319ec918.json go run main.go -project myfreegke -state gs://test-stateconcile/terraform.tfstate -random
```
## Random
Generates random firwall entries on the GCP side , to prove a point

## Project
`-project` : Speicified the GCP project where GCP Source ruless will try to match whats on terraform state

## State
`-state` : Gcs location of the state gs://bucket/object/..

## Sample Output:
```bash
$> GOOGLE_APPLICATION_CREDENTIALS=~/Downloads/myfreegke-a9a1319ec918.json go run main.go -project myfreegke -state gs://test-stateconcile/terraform.tfstate -random

[*] Downloading Terraform state from gs://test-stateconcile/terraform.tfstate into working.tfstate...DONE
[*] Getting google_compute_firewall from googlecloud api for project:myfreegke...DONE
[INFO] >> Adding random soruce rules
[Warning] >> Rules found in GCP but missing in Terraform State: gs://test-stateconcile/terraform.tfstate on Version: 4
* google_compute_firewall.default-allow-icmp -> missing
* google_compute_firewall.default-allow-internal -> missing
* google_compute_firewall.default-allow-rdp -> missing
* google_compute_firewall.default-allow-ssh -> missing
* google_compute_firewall.gke-myfreegke-c14da6a7-all -> missing
* google_compute_firewall.gke-myfreegke-c14da6a7-ssh -> missing
* google_compute_firewall.gke-myfreegke-c14da6a7-vms -> missing
* google_compute_firewall.RandomRule0 -> missing
* google_compute_firewall.RandomRule1 -> missing
* google_compute_firewall.RandomRule2 -> missing
* google_compute_firewall.RandomRule3 -> missing
* google_compute_firewall.RandomRule4 -> missing
* google_compute_firewall.RandomRule5 -> missing
```


