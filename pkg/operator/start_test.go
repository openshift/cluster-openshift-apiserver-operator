package operator

import "testing"

func TestNothing(t *testing.T) {
}

func TestRBACManifestFiles(t *testing.T) {
	files := rbacManifestFiles()
	expected, actual := 74+33, len(files)
	if actual != expected {
		t.Errorf("expected %d manifests, found %d.", expected, actual)
	}
	for _, a := range files {
		t.Log(a)
	}
}
