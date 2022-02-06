# stateconcile


## Usage
```
$> GOOGLE_APPLICATION_CREDENTIALS=~/Downloads/myfreegke-a9a1319ec918.json go run main.go -project myfreegke -state gs://test-stateconcile/terraform.tfstate -random
```
## Random
Generats random firwall entries on the GCP side , to prove a point

## Project
`-project` : Speicified the GCP project where GCP Source ruless will try to match whats on terraform state

## State
`-state` : Gcs location of the state gs://bucket/object/..
