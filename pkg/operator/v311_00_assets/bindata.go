// Code generated for package v311_00_assets by go-bindata DO NOT EDIT. (@generated)
// sources:
// bindata/v3.11.0/config/defaultconfig.yaml
// bindata/v3.11.0/openshift-apiserver/apiserver-clusterrolebinding.yaml
// bindata/v3.11.0/openshift-apiserver/cm.yaml
// bindata/v3.11.0/openshift-apiserver/deploy.yaml
// bindata/v3.11.0/openshift-apiserver/ns.yaml
// bindata/v3.11.0/openshift-apiserver/pdb.yaml
// bindata/v3.11.0/openshift-apiserver/sa.yaml
// bindata/v3.11.0/openshift-apiserver/svc.yaml
// bindata/v3.11.0/openshift-apiserver/trusted_ca_cm.yaml
package v311_00_assets

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _v3110ConfigDefaultconfigYaml = []byte(`apiVersion: openshiftcontrolplane.config.openshift.io/v1
kind: OpenShiftAPIServerConfig
storageConfig:
  urls:
  - https://etcd.openshift-etcd.svc:2379
apiServerArguments:
  audit-log-format:
  - json
  audit-log-maxbackup:
  - "10"
  audit-log-maxsize:
  - "100"
  audit-log-path:
  - /var/log/openshift-apiserver/audit.log
  audit-policy-file:
  - /var/run/configmaps/audit/policy.yaml
  etcd-healthcheck-timeout:
  - 9s
  etcd-readycheck-timeout:
  - 9s
  shutdown-delay-duration:
  - 50s # this gives SDN 15s to converge after the worst readyz=false delay
  shutdown-send-retry-after:
  - "true"
servingInfo:
  bindNetwork: "tcp"
`)

func v3110ConfigDefaultconfigYamlBytes() ([]byte, error) {
	return _v3110ConfigDefaultconfigYaml, nil
}

func v3110ConfigDefaultconfigYaml() (*asset, error) {
	bytes, err := v3110ConfigDefaultconfigYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/config/defaultconfig.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverApiserverClusterrolebindingYaml = []byte(`apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: system:openshift:openshift-apiserver
roleRef:
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  namespace: openshift-apiserver
  name: openshift-apiserver-sa`)

func v3110OpenshiftApiserverApiserverClusterrolebindingYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverApiserverClusterrolebindingYaml, nil
}

func v3110OpenshiftApiserverApiserverClusterrolebindingYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverApiserverClusterrolebindingYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/apiserver-clusterrolebinding.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverCmYaml = []byte(`apiVersion: v1
kind: ConfigMap
metadata:
  namespace: openshift-apiserver
  name: config
data:
  config.yaml:
`)

func v3110OpenshiftApiserverCmYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverCmYaml, nil
}

func v3110OpenshiftApiserverCmYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverCmYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/cm.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverDeployYaml = []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: openshift-apiserver
  name: apiserver
  labels:
    app: openshift-apiserver
    apiserver: "true"
spec:
  # The number of replicas will be set in code to the number of master nodes.
  strategy:
    type: RollingUpdate
    rollingUpdate:
      # To ensure that only one pod at a time writes to the node's
      # audit log, require the update strategy to proceed a node at a
      # time. Only when a master node has its existing
      # openshift-apiserver pod stopped will a new one be allowed to
      # start.
      maxUnavailable: 1
      maxSurge: 0
  selector:
    matchLabels:
      # Need to vary the app label from that used by the legacy
      # daemonset ('openshift-apiserver') to avoid the legacy
      # daemonset and its replacement deployment trying to try to
      # manage the same pods.
      #
      # It's also necessary to use different labeling to ensure, via
      # anti-affinity, at most one deployment-managed pod on each
      # master node. Without label differentiation, anti-affinity
      # would prevent a deployment-managed pod from running on a node
      # that was already running a daemonset-managed pod.
      app: openshift-apiserver-a
      apiserver: "true"
  template:
    metadata:
      name: openshift-apiserver
      labels:
        app: openshift-apiserver-a
        apiserver: "true"
      annotations:
        target.workload.openshift.io/management: '{"effect": "PreferredDuringScheduling"}'
        openshift.io/required-scc: 'privileged'
    spec:
      serviceAccountName: openshift-apiserver-sa
      priorityClassName: system-node-critical
      initContainers:
        - name: fix-audit-permissions
          terminationMessagePolicy: FallbackToLogsOnError
          image: ${IMAGE}
          imagePullPolicy: IfNotPresent
          command: ['sh', '-c', 'chmod 0700 /var/log/openshift-apiserver && touch /var/log/openshift-apiserver/audit.log && chmod 0600 /var/log/openshift-apiserver/*']
          securityContext:
            privileged: true
            runAsUser: 0
          resources:
            requests:
              cpu: 15m
              memory: 50Mi
          volumeMounts:
            - mountPath: /var/log/openshift-apiserver
              name: audit-dir
      containers:
      - name: openshift-apiserver
        terminationMessagePolicy: FallbackToLogsOnError
        image: ${IMAGE}
        imagePullPolicy: IfNotPresent
        command: ["/bin/bash", "-ec"]
        args:
          - |
            if [ -s /var/run/configmaps/trusted-ca-bundle/tls-ca-bundle.pem ]; then
              echo "Copying system trust bundle"
              cp -f /var/run/configmaps/trusted-ca-bundle/tls-ca-bundle.pem /etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem
            fi
            exec openshift-apiserver start --config=/var/run/configmaps/config/config.yaml -v=${VERBOSITY}
        resources:
          requests:
            memory: 200Mi
            cpu: 100m
        env:
          - name: POD_NAME
            valueFrom:
              fieldRef:
                fieldPath: metadata.name
          - name: POD_NAMESPACE
            valueFrom:
              fieldRef:
                fieldPath: metadata.namespace
        # we need to set this to privileged to be able to write audit to /var/log/openshift-apiserver
        securityContext:
          privileged: true
          readOnlyRootFilesystem: false
          runAsUser: 0
        ports:
        - containerPort: 8443
        volumeMounts:
        - mountPath: /var/lib/kubelet/
          name: node-pullsecrets
          readOnly: true
        - mountPath: /var/run/configmaps/config
          name: config
        - mountPath: /var/run/configmaps/audit
          name: audit
        - mountPath: /var/run/secrets/etcd-client
          name: etcd-client
        - mountPath: /var/run/configmaps/etcd-serving-ca
          name: etcd-serving-ca
        - mountPath: /var/run/configmaps/image-import-ca
          name: image-import-ca
        - mountPath: /var/run/configmaps/trusted-ca-bundle
          name: trusted-ca-bundle
        - mountPath: /var/run/secrets/serving-cert
          name: serving-cert
        - mountPath: /var/run/secrets/encryption-config
          name: encryption-config
        - mountPath: /var/log/openshift-apiserver
          name: audit-dir
        livenessProbe:
          httpGet:
            scheme: HTTPS
            port: 8443
            path: livez?exclude=etcd
          initialDelaySeconds: 0
          periodSeconds: 10
          timeoutSeconds: 10
          successThreshold: 1
          failureThreshold: 3
        readinessProbe:
          httpGet:
            scheme: HTTPS
            port: 8443
            path: readyz
          initialDelaySeconds: 0
          periodSeconds: 5
          timeoutSeconds: 10
          successThreshold: 1
          failureThreshold: 3
        startupProbe:
          httpGet:
            scheme: HTTPS
            port: 8443
            path: livez
          initialDelaySeconds: 0
          periodSeconds: 5
          timeoutSeconds: 10
          successThreshold: 1
          failureThreshold: 30
      - name: openshift-apiserver-check-endpoints
        image: ${KUBE_APISERVER_OPERATOR_IMAGE}
        imagePullPolicy: IfNotPresent
        terminationMessagePolicy: FallbackToLogsOnError
        command:
          - cluster-kube-apiserver-operator
          - check-endpoints
        args:
          - --listen
          - 0.0.0.0:17698
          - --namespace
          - $(POD_NAMESPACE)
          - --v
          - '2'
        env:
          - name: POD_NAME
            valueFrom:
              fieldRef:
                fieldPath: metadata.name
          - name: POD_NAMESPACE
            valueFrom:
              fieldRef:
                fieldPath: metadata.namespace
        ports:
          - name: check-endpoints
            containerPort: 17698
            protocol: TCP
        resources:
          requests:
            memory: 50Mi
            cpu: 10m
      terminationGracePeriodSeconds: 120 # a bit more than the 60 seconds timeout of non-long-running requests + the shutdown delay
      volumes:
      - name: node-pullsecrets
        hostPath:
          path: /var/lib/kubelet/
          type: Directory
      - name: config
        configMap:
          name: config
      - name: audit
        configMap:
          name: audit-${REVISION}
      - name: etcd-client
        secret:
          secretName: etcd-client
          defaultMode: 0600
      - name: etcd-serving-ca
        configMap:
          name: etcd-serving-ca
      - name: image-import-ca
        configMap:
          name: image-import-ca
          optional: true
      - name: serving-cert
        secret:
          secretName: serving-cert
          defaultMode: 0600
      - name: trusted-ca-bundle
        configMap:
          name: trusted-ca-bundle
          optional: true
          items:
          - key: ca-bundle.crt
            path: tls-ca-bundle.pem
      - name: encryption-config
        secret:
          secretName: encryption-config-${REVISION}
          optional: true
          defaultMode: 0600
      - hostPath:
          path: /var/log/openshift-apiserver
        name: audit-dir
      nodeSelector:
        node-role.kubernetes.io/master: ""
      tolerations:
        # Ensure pod can be scheduled on master nodes
      - key: "node-role.kubernetes.io/master"
        operator: "Exists"
        effect: "NoSchedule"
        # Ensure pod can be scheduled on master nodes if tainted NoExecute
      - key: "node-role.kubernetes.io/control-plane"
        operator: "Exists"
        effect: "NoExecute"
        # Ensure pod can be evicted if the node is unreachable
      - key: "node.kubernetes.io/unreachable"
        operator: "Exists"
        effect: "NoExecute"
        tolerationSeconds: 120
        # Ensure scheduling is delayed until node readiness
        # (i.e. network operator configures CNI on the node)
      - key: "node.kubernetes.io/not-ready"
        operator: "Exists"
        effect: "NoExecute"
        tolerationSeconds: 120
      affinity:
        podAntiAffinity:
          # Ensure that at most one apiserver pod will be scheduled on a node.
          requiredDuringSchedulingIgnoredDuringExecution:
          - topologyKey: "kubernetes.io/hostname"
            labelSelector:
              matchLabels:
                app: "openshift-apiserver-a"
                apiserver: "true"
`)

func v3110OpenshiftApiserverDeployYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverDeployYaml, nil
}

func v3110OpenshiftApiserverDeployYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverDeployYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/deploy.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverNsYaml = []byte(`apiVersion: v1
kind: Namespace
metadata:
  annotations:
    openshift.io/node-selector: ""
    workload.openshift.io/allowed: "management"
  name: openshift-apiserver
  labels:
    openshift.io/run-level-: "" # remove the label if previously set
    openshift.io/cluster-monitoring: "true"
    # needs to be privileged because of:
    # - hostPath volumes (volumes \"node-pullsecrets\", \"audit-dir\")
    # - privileged (containers \"fix-audit-permissions\", \"openshift-apiserver\"
    pod-security.kubernetes.io/enforce: privileged
    pod-security.kubernetes.io/audit: privileged
    pod-security.kubernetes.io/warn: privileged
`)

func v3110OpenshiftApiserverNsYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverNsYaml, nil
}

func v3110OpenshiftApiserverNsYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverNsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/ns.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverPdbYaml = []byte(`apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: openshift-apiserver-pdb
  namespace: openshift-apiserver
spec:
  maxUnavailable: 1
  unhealthyPodEvictionPolicy: AlwaysAllow
  selector:
    matchLabels:
      app: openshift-apiserver-a
      apiserver: "true"
`)

func v3110OpenshiftApiserverPdbYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverPdbYaml, nil
}

func v3110OpenshiftApiserverPdbYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverPdbYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/pdb.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverSaYaml = []byte(`apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: openshift-apiserver
  name: openshift-apiserver-sa
`)

func v3110OpenshiftApiserverSaYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverSaYaml, nil
}

func v3110OpenshiftApiserverSaYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverSaYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/sa.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverSvcYaml = []byte(`apiVersion: v1
kind: Service
metadata:
  namespace: openshift-apiserver
  name: api
  annotations:
    service.alpha.openshift.io/serving-cert-secret-name: serving-cert
  labels:
    prometheus: openshift-apiserver
spec:
  selector:
    apiserver: "true"
  ports:
  - name: https
    port: 443
    targetPort: 8443
`)

func v3110OpenshiftApiserverSvcYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverSvcYaml, nil
}

func v3110OpenshiftApiserverSvcYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverSvcYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/svc.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverTrusted_ca_cmYaml = []byte(`apiVersion: v1
kind: ConfigMap
metadata:
  namespace: openshift-apiserver
  name: trusted-ca-bundle
  labels:
    config.openshift.io/inject-trusted-cabundle: "true"
`)

func v3110OpenshiftApiserverTrusted_ca_cmYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverTrusted_ca_cmYaml, nil
}

func v3110OpenshiftApiserverTrusted_ca_cmYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverTrusted_ca_cmYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/trusted_ca_cm.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"v3.11.0/config/defaultconfig.yaml":                             v3110ConfigDefaultconfigYaml,
	"v3.11.0/openshift-apiserver/apiserver-clusterrolebinding.yaml": v3110OpenshiftApiserverApiserverClusterrolebindingYaml,
	"v3.11.0/openshift-apiserver/cm.yaml":                           v3110OpenshiftApiserverCmYaml,
	"v3.11.0/openshift-apiserver/deploy.yaml":                       v3110OpenshiftApiserverDeployYaml,
	"v3.11.0/openshift-apiserver/ns.yaml":                           v3110OpenshiftApiserverNsYaml,
	"v3.11.0/openshift-apiserver/pdb.yaml":                          v3110OpenshiftApiserverPdbYaml,
	"v3.11.0/openshift-apiserver/sa.yaml":                           v3110OpenshiftApiserverSaYaml,
	"v3.11.0/openshift-apiserver/svc.yaml":                          v3110OpenshiftApiserverSvcYaml,
	"v3.11.0/openshift-apiserver/trusted_ca_cm.yaml":                v3110OpenshiftApiserverTrusted_ca_cmYaml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//
//	data/
//	  foo.txt
//	  img/
//	    a.png
//	    b.png
//
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"v3.11.0": {nil, map[string]*bintree{
		"config": {nil, map[string]*bintree{
			"defaultconfig.yaml": {v3110ConfigDefaultconfigYaml, map[string]*bintree{}},
		}},
		"openshift-apiserver": {nil, map[string]*bintree{
			"apiserver-clusterrolebinding.yaml": {v3110OpenshiftApiserverApiserverClusterrolebindingYaml, map[string]*bintree{}},
			"cm.yaml":                           {v3110OpenshiftApiserverCmYaml, map[string]*bintree{}},
			"deploy.yaml":                       {v3110OpenshiftApiserverDeployYaml, map[string]*bintree{}},
			"ns.yaml":                           {v3110OpenshiftApiserverNsYaml, map[string]*bintree{}},
			"pdb.yaml":                          {v3110OpenshiftApiserverPdbYaml, map[string]*bintree{}},
			"sa.yaml":                           {v3110OpenshiftApiserverSaYaml, map[string]*bintree{}},
			"svc.yaml":                          {v3110OpenshiftApiserverSvcYaml, map[string]*bintree{}},
			"trusted_ca_cm.yaml":                {v3110OpenshiftApiserverTrusted_ca_cmYaml, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
