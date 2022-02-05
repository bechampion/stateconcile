package googleactions
import (
        "fmt"
        "io"
        "os"
        "log"
        "time"
        "golang.org/x/net/context"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/compute/v1"
        "cloud.google.com/go/storage"
)


func GetFirewallRules(project string) {
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
        if err := req.Pages(ctx, func(page *compute.FirewallList) error {
                for _, firewall := range page.Items {
                        fmt.Printf("%s\n", firewall.Direction)
                }
                return nil
        }); err != nil {
                log.Fatal(err)
        }
}
func DownloadTerraformState(w io.Writer, bucket, object string, destFileName string) error {
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
        fmt.Fprintf(w, "Blob %v downloaded to local file %v\n", object, destFileName)
        return nil
}
func main() {
	DownloadTerraformState(os.Stdout,"test-stateconcile","realwinrm.py","here")
	GetFirewallRules("myfreegke")
}
