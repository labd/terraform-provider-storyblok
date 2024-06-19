package internal

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

var (
	mode recorder.Mode
)

func init() {
	record := os.Getenv("RECORD")
	if record == "true" {
		mode = recorder.ModeRecordOnly
	} else {
		mode = recorder.ModeReplayOnly
	}
}

func ProviderFactories(file string) (map[string]func() (tfprotov6.ProviderServer, error), func() error) {
	opt, sr := WithRecorderClient(file, mode)

	return map[string]func() (tfprotov6.ProviderServer, error){
		"storyblok": func() (tfprotov6.ProviderServer, error) {
			return providerserver.NewProtocol6WithError(
				New(opt),
			)()
		},
	}, sr
}

func TestAccPreCheck(t *testing.T) func() {
	return func() {
		requiredEnvs := []string{
			"STORYBLOK_URL",
			"STORYBLOK_TOKEN",
		}
		for _, val := range requiredEnvs {
			if os.Getenv(val) == "" {
				t.Fatalf("%v must be set for acceptance tests", val)
			}
		}
	}
}
