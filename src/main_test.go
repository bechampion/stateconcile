package main

import (
	"testing"
	"math/rand"
	"time"
	"fmt"
)

func BenchmarkRule(b *testing.B) {
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
		"pepinos",
	}
	rules := BuildRules()
	n := rand.Int() % len(reasons)
	fmt.Printf("%s --> %t\n",reasons[n],FindRule(reasons[n], rules))
}
