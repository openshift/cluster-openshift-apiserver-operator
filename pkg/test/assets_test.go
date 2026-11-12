package test

import (
	"os"
	"testing"

	"github.com/ghodss/yaml"

	"github.com/openshift/library-go/pkg/assets"
)

func TestYamlCorrectness(t *testing.T) {
	readAllYaml("../../manifests/", t)
	readAllYaml("../../bindata/", t)
}

func TestCheckEndpointsServiceMonitorTLSConfig(t *testing.T) {
	data, err := os.ReadFile("../../manifests/0000_90_openshift-apiserver-operator_05_check-endpoints_servicemonitor.yaml")
	if err != nil {
		t.Fatalf("failed to read check-endpoints ServiceMonitor: %v", err)
	}

	var manifest struct {
		Spec struct {
			Endpoints []struct {
				Port            string `json:"port"`
				Scheme          string `json:"scheme"`
				BearerTokenFile string `json:"bearerTokenFile"`
				TLSConfig       struct {
					InsecureSkipVerify bool   `json:"insecureSkipVerify"`
					CAFile             string `json:"caFile"`
					ServerName         string `json:"serverName"`
					CertFile           string `json:"certFile"`
					KeyFile            string `json:"keyFile"`
				} `json:"tlsConfig"`
			} `json:"endpoints"`
		} `json:"spec"`
	}
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("failed to unmarshal check-endpoints ServiceMonitor: %v", err)
	}
	if len(manifest.Spec.Endpoints) != 1 {
		t.Fatalf("expected one endpoint, got %d", len(manifest.Spec.Endpoints))
	}

	endpoint := manifest.Spec.Endpoints[0]
	if endpoint.Port != "check-endpoints" {
		t.Errorf("expected port check-endpoints, got %q", endpoint.Port)
	}
	if endpoint.Scheme != "https" {
		t.Errorf("expected scheme https, got %q", endpoint.Scheme)
	}
	if endpoint.BearerTokenFile == "" {
		t.Errorf("expected bearer token file to remain configured")
	}
	if endpoint.TLSConfig.InsecureSkipVerify {
		t.Errorf("expected insecureSkipVerify to be disabled")
	}
	if endpoint.TLSConfig.CAFile != "/etc/prometheus/configmaps/serving-certs-ca-bundle/service-ca.crt" {
		t.Errorf("unexpected CA file: %q", endpoint.TLSConfig.CAFile)
	}
	if endpoint.TLSConfig.ServerName != "api.openshift-apiserver.svc" {
		t.Errorf("unexpected server name: %q", endpoint.TLSConfig.ServerName)
	}
	if endpoint.TLSConfig.CertFile == "" || endpoint.TLSConfig.KeyFile == "" {
		t.Errorf("expected client certificate and key to remain configured")
	}
}

func readAllYaml(path string, t *testing.T) {
	manifests, err := assets.New(path, make(map[string]string), nil, assets.OnlyYaml)
	if err != nil {
		t.Errorf("Unexpected error reading manifests from %s: %v", path, err)
	}
	t.Logf("Found %d manifests in %s", len(manifests), path)
	for _, m := range manifests {
		contents := make(map[string]interface{})
		t.Logf("Checking %s...", m.Name)
		if err := yaml.Unmarshal(m.Data, &contents); err != nil {
			t.Errorf("Unexpected error unmarshaling %s: %v", m.Name, err)
		}
	}
}
