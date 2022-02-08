##custom
export GOOGLE_APPLICATION_CREDENTIALS=/home/${USER}/Downloads/myfreegke-a9a1319ec918.json
.PHONY: default
default: drift
runlogs:
	@go run main.go -project myfreegke -state gs://test-stateconcile/terraform.tfstate -logs
logs:
	@./main -project myfreegke -state gs://test-stateconcile/terraform.tfstate -logs
drift:
	@./main -project myfreegke -state gs://test-stateconcile/terraform.tfstate
