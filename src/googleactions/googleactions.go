package googleactions

import (
	"cloud.google.com/go/storage"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"io"
	"log"
	"os"
	"regexp"
	"time"
)

var r *regexp.Regexp
func GetFirewallRules(project string, ignoreauto bool) []string {
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("[%s] Getting google_compute_firewall from googlecloud api for project:%s...", green("*"), project)
	var fwlist []string
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}
	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}
	req := computeService.Firewalls.List(project)
	_ = r
	if ignoreauto {
		r, _ = regexp.Compile("^(k8s|gke)(-|_|[0-9][a-z])*")
	}
	if err := req.Pages(ctx, func(page *compute.FirewallList) error {
		for _, firewall := range page.Items {
			if ignoreauto {
				m := r.MatchString(firewall.Name)
				if m == false {
					fwlist = append(fwlist, firewall.Name)
				}
			} else {
				fwlist = append(fwlist, firewall.Name)
			}
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", green("DONE"))
	return fwlist
}
func DownloadTerraformState(w io.Writer, bucket, object string, destFileName string) error {
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("\n[%s] Downloading Terraform state from gs://%s/%s into %s...", green("*"), bucket, object, destFileName)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	f, err := os.Create(destFileName)
	if err != nil {
		return fmt.Errorf("os.Create: %v", err)
	}
	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).NewReader: %v", object, err)
	}
	defer rc.Close()

	if _, err := io.Copy(f, rc); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("f.Close: %v", err)
	}
	fmt.Printf("%s\n", green("DONE"))
	return nil
}
