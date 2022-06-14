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

## logs
`-logs` : Look in stackdriver for information about the firewall rule created outside the state

## random 
`-random` : Insert Random firewall rules into the stack to generate some synthetic drift

## ignoreauth
`-ignoreauto` : Ignore GKE generated firewall rules for example

## Sample Output:
![stateconcile](/img/output.png)
