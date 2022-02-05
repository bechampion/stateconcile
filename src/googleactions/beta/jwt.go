package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type ServiceAccount struct {
	Type         string `json:"service_account"`
	ProjectId    string `json:"project_id"`
	PrivateKeyId string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	ClientEmail  string `json:"client_email"`
	ClientId     string `json:"client_id"`
	AuthUri      string `json:"auth_uri"`
	TokenUri     string `json:"token_uri"`
}

func GenJWT(svcaccountlocation string) string {
	serviceaccount, _ := os.Open(svcaccountlocation)
	rawserviceaccount, _ := ioutil.ReadAll(serviceaccount)
	sva := ServiceAccount{}
	err := json.Unmarshal(rawserviceaccount, &sva)
	if err != nil {
		log.Fatal(err)

	}
	rawheader := "{\"alg\":\"RS256\",\"typ\":\"JWT\"}"
	scope := "https://www.googleapis.com/auth/compute"
	rawclaim := `
{
        "iss": "%s",
	"scope": "%s",
	"aud": "https://www.googleapis.com/oauth2/v4/token",
	"exp": 3600,
	"iat": 1644070430
 }
`
	claim := fmt.Sprintf(rawclaim, sva.ClientEmail,scope)
	requestbody := fmt.Sprintf("%s.%s",base64.StdEncoding.EncodeToString([]byte(rawheader)),base64.StdEncoding.EncodeToString([]byte(claim)))
	h := sha256.New()
	h.Write([]byte(requestbody))
	digest := h.Sum(nil)
	block, _ := pem.Decode([]byte(sva.PrivateKey))
	if block == nil {
		panic("failed to parse root certificate PEM")
	}
	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}
	key := privKey.(*rsa.PrivateKey)
	s,err := rsa.SignPKCS1v15(nil, key, crypto.SHA256, digest)
	if err != nil {
		panic("failed to sign:" + err.Error())
	}
	sigrbody := fmt.Sprintf("%s.%x\n", requestbody , s)
	return base64.StdEncoding.EncodeToString([]byte(sigrbody))
}

func main() {
	pk := GenJWT("../token/myfreegke-a9a1319ec918.json")
	fmt.Println(pk)
}

// #!/bin/bash
// set -euo pipefail
// key_json_file="$1"
// scope="$2"
// valid_for_sec="${3:-3600}"
// IFS="," read -r private_key sa_email <<<$(cat $1 | jq '(.private_key|@base64)+","+.client_email' -r)
// header='{"alg":"RS256","typ":"JWT"}'
// claim=$(cat <<EOF | jq -c .
// 	{
// 	"iss": "$sa_email",
// 		"scope": "$scope",
// 		"aud": "https://www.googleapis.com/oauth2/v4/token",
// 		"exp": $(($(date +%s) + $valid_for_sec)),
// 		"iat": $(date +%s)
// }
// 	EOF
// )
// request_body="$(echo "$header" | base64 -w0).$(echo "$claim" | base64 -w0)"
// signature=$(openssl dgst -sha256 -sign <(echo "$private_key"| base64 -d ) <(printf "$request_body") | base64 -w0)
// printf "$request_body.$signature"
