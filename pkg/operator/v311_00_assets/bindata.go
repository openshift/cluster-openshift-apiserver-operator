// Code generated for package v311_00_assets by go-bindata DO NOT EDIT. (@generated)
// sources:
// bindata/v3.11.0/config/defaultconfig.yaml
// bindata/v3.11.0/openshift-apiserver/apiserver-clusterrolebinding.yaml
// bindata/v3.11.0/openshift-apiserver/cm.yaml
// bindata/v3.11.0/openshift-apiserver/deploy.yaml
// bindata/v3.11.0/openshift-apiserver/ns.yaml
// bindata/v3.11.0/openshift-apiserver/pdb.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-basic-user.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-cluster-debugger.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-cluster-reader.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-cluster-status.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-registry-admin.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-registry-editor.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-registry-viewer.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-self-access-reviewer.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-self-provisioner.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-storage-admin.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-sudoer.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-build-strategy-custom.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-build-strategy-docker.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-build-strategy-jenkinspipeline.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-build-strategy-source.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-deployer.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-image-auditor.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-image-builder.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-image-pruner.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-image-puller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-image-pusher.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-image-signer.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-master.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-node-admin.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-node-reader.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-oauth-token-deleter.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-controller-cluster-quota-reconciliation-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-controller-horizontal-pod-autoscaler.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-controller-pv-recycler-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-controller-resourcequota-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-controller-service-serving-cert-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-controller-template-service-broker.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-discovery.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-templateservicebroker-client.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-router.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-scope-impersonation.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-sdn-manager.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-sdn-reader.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-webhook.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-basic-users.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-cluster-admins.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-cluster-readers.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-cluster-status-binding.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-self-access-reviewers.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-self-provisioners.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-build-strategy-docker-binding.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-build-strategy-jenkinspipeline-binding.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-build-strategy-source-binding.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-controller-horizontal-pod-autoscaler.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-masters.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-node-admins.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-node-bootstrapper.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-node-proxiers.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-oauth-token-deleters.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-controller-cluster-quota-reconciliation-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-controller-horizontal-pod-autoscaler.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-controller-pv-recycler-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-controller-resourcequota-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-controller-service-serving-cert-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-controller-template-service-broker.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-discovery.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-scope-impersonation.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-sdn-readers.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-webhooks.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/50_aggregate-to-admin.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/50_aggregate-to-basic-user.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/50_aggregate-to-cluster-reader.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/50_aggregate-to-edit.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/50_aggregate-to-storage-admin.yaml
// bindata/v3.11.0/openshift-apiserver/rbac/50_aggregate-to-view.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-build-config-change-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-build-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-default-rolebindings-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-deployer-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-deploymentconfig-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-image-import-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-image-trigger-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-origin-namespace-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-service-ingress-ip-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-serviceaccount-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-serviceaccount-pull-secrets-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-template-instance-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-template-instance-finalizer-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-unidling-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-deployer.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-image-builder.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-image-puller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-build-config-change-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-build-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-default-rolebindings-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-deployer-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-deploymentconfig-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-image-import-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-image-trigger-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-origin-namespace-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-service-ingress-ip-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-serviceaccount-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-serviceaccount-pull-secrets-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-template-instance-controller-admin.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-template-instance-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-template-instance-finalizer-controller-admin.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-template-instance-finalizer-controller.yaml
// bindata/v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-unidling-controller.yaml
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
  shutdown-delay-duration:
  - 15s # this gives SDN 5s to converge after the worst readyz=false delay
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
            path: healthz
          initialDelaySeconds: 0
          periodSeconds: 10
          timeoutSeconds: 1
          successThreshold: 1
          failureThreshold: 3
        readinessProbe:
          httpGet:
            scheme: HTTPS
            port: 8443
            path: readyz
          initialDelaySeconds: 0
          periodSeconds: 5
          timeoutSeconds: 1
          successThreshold: 1
          failureThreshold: 1
        startupProbe:
          httpGet:
            scheme: HTTPS
            port: 8443
            path: healthz
          initialDelaySeconds: 0
          periodSeconds: 5
          timeoutSeconds: 1
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
      terminationGracePeriodSeconds: 90 # a bit more than the 60 seconds timeout of non-long-running requests + the shutdown delay
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

var _v3110OpenshiftApiserverRbac20_clusterroleBasicUserYaml = []byte(`---
aggregationRule:
  clusterRoleSelectors:
    - matchLabels:
        authorization.openshift.io/aggregate-to-basic-user: "true"
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    openshift.io/description: A user that can get basic information about projects.
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: basic-user
rules: null
`)

func v3110OpenshiftApiserverRbac20_clusterroleBasicUserYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleBasicUserYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleBasicUserYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleBasicUserYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-basic-user.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleClusterDebuggerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: cluster-debugger
rules:
  - nonResourceURLs:
      - /debug/pprof
      - /debug/pprof/*
      - /metrics
    verbs:
      - get
`)

func v3110OpenshiftApiserverRbac20_clusterroleClusterDebuggerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleClusterDebuggerYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleClusterDebuggerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleClusterDebuggerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-cluster-debugger.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleClusterReaderYaml = []byte(`---
aggregationRule:
  clusterRoleSelectors:
    - matchLabels:
        rbac.authorization.k8s.io/aggregate-to-cluster-reader: "true"
    - matchLabels:
        rbac.authorization.k8s.io/aggregate-to-view: "true"
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: cluster-reader
rules: null
`)

func v3110OpenshiftApiserverRbac20_clusterroleClusterReaderYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleClusterReaderYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleClusterReaderYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleClusterReaderYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-cluster-reader.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleClusterStatusYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    openshift.io/description: A user that can get basic cluster status information.
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: cluster-status
rules:
  - nonResourceURLs:
      - /healthz
      - /healthz/
    verbs:
      - get
  - nonResourceURLs:
      - /version
      - /version/*
      - /api
      - /api/*
      - /apis
      - /apis/*
      - /oapi
      - /oapi/*
      - /openapi/v2
      - /swaggerapi
      - /swaggerapi/*
      - /swagger.json
      - /swagger-2.0.0.pb-v1
      - /osapi
      - /osapi/
      - /.well-known
      - /.well-known/oauth-authorization-server
      - /
    verbs:
      - get
`)

func v3110OpenshiftApiserverRbac20_clusterroleClusterStatusYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleClusterStatusYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleClusterStatusYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleClusterStatusYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-cluster-status.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleRegistryAdminYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  labels:
    rbac.authorization.k8s.io/aggregate-to-admin: "true"
  name: registry-admin
rules:
  - apiGroups:
      - ""
    resources:
      - secrets
      - serviceaccounts
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreamimages
      - imagestreammappings
      - imagestreams
      - imagestreams/secrets
      - imagestreamtags
      - imagetags
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreamimports
    verbs:
      - create
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreams/layers
    verbs:
      - get
      - update
  - apiGroups:
      - ""
      - authorization.openshift.io
    resources:
      - rolebindings
      - roles
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - rbac.authorization.k8s.io
    resources:
      - rolebindings
      - roles
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
      - authorization.openshift.io
    resources:
      - localresourceaccessreviews
      - localsubjectaccessreviews
      - subjectrulesreviews
    verbs:
      - create
  - apiGroups:
      - authorization.k8s.io
    resources:
      - localsubjectaccessreviews
    verbs:
      - create
  - apiGroups:
      - ""
    resources:
      - namespaces
    verbs:
      - get
  - apiGroups:
      - ""
      - project.openshift.io
    resources:
      - projects
    verbs:
      - delete
      - get
  - apiGroups:
      - ""
      - authorization.openshift.io
    resources:
      - resourceaccessreviews
      - subjectaccessreviews
    verbs:
      - create
`)

func v3110OpenshiftApiserverRbac20_clusterroleRegistryAdminYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleRegistryAdminYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleRegistryAdminYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleRegistryAdminYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-registry-admin.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleRegistryEditorYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  labels:
    rbac.authorization.k8s.io/aggregate-to-edit: "true"
  name: registry-editor
rules:
  - apiGroups:
      - ""
    resources:
      - secrets
      - serviceaccounts
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreamimages
      - imagestreammappings
      - imagestreams
      - imagestreams/secrets
      - imagestreamtags
      - imagetags
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreamimports
    verbs:
      - create
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreams/layers
    verbs:
      - get
      - update
  - apiGroups:
      - ""
    resources:
      - namespaces
    verbs:
      - get
  - apiGroups:
      - ""
      - project.openshift.io
    resources:
      - projects
    verbs:
      - get
`)

func v3110OpenshiftApiserverRbac20_clusterroleRegistryEditorYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleRegistryEditorYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleRegistryEditorYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleRegistryEditorYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-registry-editor.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleRegistryViewerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  labels:
    rbac.authorization.k8s.io/aggregate-to-view: "true"
  name: registry-viewer
rules:
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreamimages
      - imagestreammappings
      - imagestreams
      - imagestreamtags
      - imagetags
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreams/layers
    verbs:
      - get
  - apiGroups:
      - ""
    resources:
      - namespaces
    verbs:
      - get
  - apiGroups:
      - ""
      - project.openshift.io
    resources:
      - projects
    verbs:
      - get
`)

func v3110OpenshiftApiserverRbac20_clusterroleRegistryViewerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleRegistryViewerYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleRegistryViewerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleRegistryViewerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-registry-viewer.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSelfAccessReviewerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: self-access-reviewer
rules:
  - apiGroups:
      - ""
      - authorization.openshift.io
    resources:
      - selfsubjectrulesreviews
    verbs:
      - create
  - apiGroups:
      - authorization.k8s.io
    resources:
      - selfsubjectaccessreviews
    verbs:
      - create
`)

func v3110OpenshiftApiserverRbac20_clusterroleSelfAccessReviewerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSelfAccessReviewerYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSelfAccessReviewerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSelfAccessReviewerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-self-access-reviewer.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSelfProvisionerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    openshift.io/description: A user that can request projects.
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: self-provisioner
rules:
  - apiGroups:
      - ""
      - project.openshift.io
    resources:
      - projectrequests
    verbs:
      - create
`)

func v3110OpenshiftApiserverRbac20_clusterroleSelfProvisionerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSelfProvisionerYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSelfProvisionerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSelfProvisionerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-self-provisioner.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleStorageAdminYaml = []byte(`---
aggregationRule:
  clusterRoleSelectors:
    - matchLabels:
        storage.openshift.io/aggregate-to-storage-admin: "true"
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: storage-admin
rules: null
`)

func v3110OpenshiftApiserverRbac20_clusterroleStorageAdminYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleStorageAdminYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleStorageAdminYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleStorageAdminYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-storage-admin.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSudoerYaml = []byte(`apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: sudoer
rules:
  - apiGroups:
      - ""
      - user.openshift.io
    resourceNames:
      - system:admin
    resources:
      - systemusers
      - users
    verbs:
      - impersonate
  - apiGroups:
      - ""
      - user.openshift.io
    resourceNames:
      - system:masters
    resources:
      - groups
      - systemgroups
    verbs:
      - impersonate
`)

func v3110OpenshiftApiserverRbac20_clusterroleSudoerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSudoerYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSudoerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSudoerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-sudoer.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyCustomYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:build-strategy-custom
rules:
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - builds/custom
    verbs:
      - create
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyCustomYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyCustomYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyCustomYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyCustomYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-build-strategy-custom.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyDockerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:build-strategy-docker
rules:
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - builds/docker
      - builds/optimizeddocker
    verbs:
      - create
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyDockerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyDockerYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyDockerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyDockerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-build-strategy-docker.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyJenkinspipelineYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:build-strategy-jenkinspipeline
rules:
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - builds/jenkinspipeline
    verbs:
      - create
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyJenkinspipelineYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyJenkinspipelineYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyJenkinspipelineYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyJenkinspipelineYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-build-strategy-jenkinspipeline.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategySourceYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:build-strategy-source
rules:
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - builds/source
    verbs:
      - create
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategySourceYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategySourceYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategySourceYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategySourceYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-build-strategy-source.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemDeployerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    openshift.io/description: Grants the right to deploy within a project.  Used primarily with service accounts for automated deployments.
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:deployer
rules:
  - apiGroups:
      - ""
    resources:
      - replicationcontrollers
    verbs:
      - delete
  - apiGroups:
      - ""
    resources:
      - replicationcontrollers
    verbs:
      - get
      - list
      - update
      - watch
  - apiGroups:
      - ""
    resources:
      - replicationcontrollers/scale
    verbs:
      - get
      - update
  - apiGroups:
      - ""
    resources:
      - pods
    verbs:
      - create
      - get
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - pods/log
    verbs:
      - get
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - list
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreamtags
      - imagetags
    verbs:
      - create
      - update
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemDeployerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemDeployerYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemDeployerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemDeployerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-deployer.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemImageAuditorYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:image-auditor
rules:
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - images
    verbs:
      - get
      - list
      - patch
      - update
      - watch
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemImageAuditorYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemImageAuditorYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemImageAuditorYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemImageAuditorYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-image-auditor.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemImageBuilderYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    openshift.io/description: Grants the right to build, push and pull images from within a project.  Used primarily with service accounts for builds.
    rbac.authorization.kubernetes.io/autoupdate: "true"
  labels:
    rbac.authorization.k8s.io/aggregate-to-admin: "true"
    rbac.authorization.k8s.io/aggregate-to-edit: "true"
  name: system:image-builder
rules:
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreams/layers
    verbs:
      - get
      - update
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreams
    verbs:
      - create
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - builds/details
    verbs:
      - update
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - builds
    verbs:
      - get
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemImageBuilderYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemImageBuilderYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemImageBuilderYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemImageBuilderYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-image-builder.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemImagePrunerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:image-pruner
rules:
  - apiGroups:
      - ""
    resources:
      - pods
      - replicationcontrollers
    verbs:
      - get
      - list
  - apiGroups:
      - ""
    resources:
      - limitranges
    verbs:
      - list
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - buildconfigs
      - builds
    verbs:
      - get
      - list
  - apiGroups:
      - ""
      - apps.openshift.io
    resources:
      - deploymentconfigs
    verbs:
      - get
      - list
  - apiGroups:
      - batch
    resources:
      - jobs
    verbs:
      - get
      - list
  - apiGroups:
      - batch
    resources:
      - cronjobs
    verbs:
      - get
      - list
  - apiGroups:
      - apps
    resources:
      - daemonsets
    verbs:
      - get
      - list
  - apiGroups:
      - apps
    resources:
      - deployments
    verbs:
      - get
      - list
  - apiGroups:
      - apps
    resources:
      - replicasets
    verbs:
      - get
      - list
  - apiGroups:
      - apps
    resources:
      - statefulsets
    verbs:
      - get
      - list
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - images
    verbs:
      - delete
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - images
      - imagestreams
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreams/status
    verbs:
      - update
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemImagePrunerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemImagePrunerYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemImagePrunerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemImagePrunerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-image-pruner.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemImagePullerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    openshift.io/description: Grants the right to pull images from within a project.
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:image-puller
rules:
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreams/layers
    verbs:
      - get
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemImagePullerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemImagePullerYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemImagePullerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemImagePullerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-image-puller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemImagePusherYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    openshift.io/description: Grants the right to push and pull images from within a project.
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:image-pusher
rules:
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreams/layers
    verbs:
      - get
      - update
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemImagePusherYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemImagePusherYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemImagePusherYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemImagePusherYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-image-pusher.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemImageSignerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:image-signer
rules:
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - images
      - imagestreams/layers
    verbs:
      - get
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagesignatures
    verbs:
      - create
      - delete
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemImageSignerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemImageSignerYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemImageSignerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemImageSignerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-image-signer.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemMasterYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:master
rules:
  - apiGroups:
      - '*'
    resources:
      - '*'
    verbs:
      - '*'
  - nonResourceURLs:
      - '*'
    verbs:
      - '*'
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemMasterYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemMasterYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemMasterYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemMasterYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-master.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemNodeAdminYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:node-admin
rules:
  - apiGroups:
      - ""
    resources:
      - nodes
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - nodes
    verbs:
      - proxy
  - apiGroups:
      - ""
    resources:
      - nodes/log
      - nodes/metrics
      - nodes/proxy
      - nodes/spec
      - nodes/stats
    verbs:
      - '*'
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemNodeAdminYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemNodeAdminYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemNodeAdminYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemNodeAdminYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-node-admin.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemNodeReaderYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:node-reader
rules:
  - apiGroups:
      - ""
    resources:
      - nodes
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - nodes/metrics
      - nodes/spec
    verbs:
      - get
  - apiGroups:
      - ""
    resources:
      - nodes/stats
    verbs:
      - create
      - get
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemNodeReaderYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemNodeReaderYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemNodeReaderYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemNodeReaderYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-node-reader.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemOauthTokenDeleterYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:oauth-token-deleter
rules:
  - apiGroups:
      - ""
      - oauth.openshift.io
    resources:
      - oauthaccesstokens
      - oauthauthorizetokens
    verbs:
      - delete
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemOauthTokenDeleterYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemOauthTokenDeleterYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemOauthTokenDeleterYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemOauthTokenDeleterYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-oauth-token-deleter.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerClusterQuotaReconciliationControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:cluster-quota-reconciliation-controller
rules:
  - apiGroups:
      - ""
    resources:
      - configmaps
    verbs:
      - get
      - list
  - apiGroups:
      - ""
    resources:
      - secrets
    verbs:
      - get
      - list
  - apiGroups:
      - ""
      - quota.openshift.io
    resources:
      - clusterresourcequotas/status
    verbs:
      - update
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
      - update
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerClusterQuotaReconciliationControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerClusterQuotaReconciliationControllerYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerClusterQuotaReconciliationControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerClusterQuotaReconciliationControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-controller-cluster-quota-reconciliation-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerHorizontalPodAutoscalerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:horizontal-pod-autoscaler
rules:
  - apiGroups:
      - ""
      - apps.openshift.io
    resources:
      - deploymentconfigs/scale
    verbs:
      - get
      - update
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerHorizontalPodAutoscalerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerHorizontalPodAutoscalerYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerHorizontalPodAutoscalerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerHorizontalPodAutoscalerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-controller-horizontal-pod-autoscaler.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerPvRecyclerControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:pv-recycler-controller
rules:
  - apiGroups:
      - ""
    resources:
      - persistentvolumes
    verbs:
      - create
      - delete
      - get
      - list
      - update
      - watch
  - apiGroups:
      - ""
    resources:
      - persistentvolumes/status
    verbs:
      - update
  - apiGroups:
      - ""
    resources:
      - persistentvolumeclaims
    verbs:
      - get
      - list
      - update
      - watch
  - apiGroups:
      - ""
    resources:
      - persistentvolumeclaims/status
    verbs:
      - update
  - apiGroups:
      - ""
    resources:
      - pods
    verbs:
      - create
      - delete
      - get
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
      - update
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerPvRecyclerControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerPvRecyclerControllerYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerPvRecyclerControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerPvRecyclerControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-controller-pv-recycler-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerResourcequotaControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:resourcequota-controller
rules:
  - apiGroups:
      - ""
    resources:
      - resourcequotas/status
    verbs:
      - update
  - apiGroups:
      - ""
    resources:
      - resourcequotas
    verbs:
      - list
  - apiGroups:
      - ""
    resources:
      - services
    verbs:
      - list
  - apiGroups:
      - ""
    resources:
      - configmaps
    verbs:
      - list
  - apiGroups:
      - ""
    resources:
      - secrets
    verbs:
      - list
  - apiGroups:
      - ""
    resources:
      - replicationcontrollers
    verbs:
      - list
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
      - update
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerResourcequotaControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerResourcequotaControllerYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerResourcequotaControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerResourcequotaControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-controller-resourcequota-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerServiceServingCertControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:service-serving-cert-controller
rules:
  - apiGroups:
      - ""
    resources:
      - services
    verbs:
      - list
      - update
      - watch
  - apiGroups:
      - ""
    resources:
      - secrets
    verbs:
      - create
      - delete
      - get
      - list
      - update
      - watch
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
      - update
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerServiceServingCertControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerServiceServingCertControllerYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerServiceServingCertControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerServiceServingCertControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-controller-service-serving-cert-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerTemplateServiceBrokerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:template-service-broker
rules:
  - apiGroups:
      - authorization.k8s.io
    resources:
      - subjectaccessreviews
    verbs:
      - create
  - apiGroups:
      - authorization.openshift.io
    resources:
      - subjectaccessreviews
    verbs:
      - create
  - apiGroups:
      - template.openshift.io
    resources:
      - brokertemplateinstances
    verbs:
      - create
      - delete
      - get
      - update
  - apiGroups:
      - template.openshift.io
    resources:
      - brokertemplateinstances/finalizers
    verbs:
      - update
  - apiGroups:
      - template.openshift.io
    resources:
      - templateinstances
    verbs:
      - assign
      - create
      - delete
      - get
  - apiGroups:
      - template.openshift.io
    resources:
      - templates
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - secrets
    verbs:
      - create
      - delete
      - get
  - apiGroups:
      - ""
    resources:
      - configmaps
      - services
    verbs:
      - get
  - apiGroups:
      - ""
    resources:
      - routes
    verbs:
      - get
  - apiGroups:
      - route.openshift.io
    resources:
      - routes
    verbs:
      - get
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
      - update
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerTemplateServiceBrokerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerTemplateServiceBrokerYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerTemplateServiceBrokerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerTemplateServiceBrokerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-controller-template-service-broker.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftDiscoveryYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:discovery
rules:
  - nonResourceURLs:
      - /version
      - /version/*
      - /api
      - /api/*
      - /apis
      - /apis/*
      - /oapi
      - /oapi/*
      - /openapi/v2
      - /swaggerapi
      - /swaggerapi/*
      - /swagger.json
      - /swagger-2.0.0.pb-v1
      - /osapi
      - /osapi/
      - /.well-known
      - /.well-known/oauth-authorization-server
      - /
    verbs:
      - get
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftDiscoveryYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftDiscoveryYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftDiscoveryYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftDiscoveryYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-discovery.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftTemplateservicebrokerClientYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:templateservicebroker-client
rules:
  - nonResourceURLs:
      - /brokers/template.openshift.io/*
    verbs:
      - delete
      - get
      - put
      - update
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftTemplateservicebrokerClientYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftTemplateservicebrokerClientYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftTemplateservicebrokerClientYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftTemplateservicebrokerClientYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-templateservicebroker-client.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemRouterYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:router
rules:
  - apiGroups:
      - discovery.k8s.io
    resources:
      - endpointslices
    verbs:
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - endpoints
    verbs:
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - services
    verbs:
      - list
      - watch
  - apiGroups:
      - authentication.k8s.io
    resources:
      - tokenreviews
    verbs:
      - create
  - apiGroups:
      - authorization.k8s.io
    resources:
      - subjectaccessreviews
    verbs:
      - create
  - apiGroups:
      - ""
      - route.openshift.io
    resources:
      - routes
    verbs:
      - list
      - watch
  - apiGroups:
      - ""
      - route.openshift.io
    resources:
      - routes/status
    verbs:
      - update
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemRouterYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemRouterYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemRouterYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemRouterYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-router.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemScopeImpersonationYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:scope-impersonation
rules:
  - apiGroups:
      - authentication.k8s.io
    resources:
      - userextras/scopes.authorization.openshift.io
    verbs:
      - impersonate
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemScopeImpersonationYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemScopeImpersonationYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemScopeImpersonationYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemScopeImpersonationYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-scope-impersonation.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemSdnManagerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:sdn-manager
rules:
  - apiGroups:
      - ""
      - network.openshift.io
    resources:
      - hostsubnets
      - netnamespaces
    verbs:
      - create
      - delete
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - network.openshift.io
    resources:
      - clusternetworks
    verbs:
      - create
      - get
  - apiGroups:
      - ""
    resources:
      - nodes
    verbs:
      - get
      - list
      - watch
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemSdnManagerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemSdnManagerYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemSdnManagerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemSdnManagerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-sdn-manager.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemSdnReaderYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:sdn-reader
rules:
  - apiGroups:
      - ""
      - network.openshift.io
    resources:
      - egressnetworkpolicies
      - hostsubnets
      - netnamespaces
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - namespaces
      - nodes
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - networking.k8s.io
    resources:
      - networkpolicies
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - network.openshift.io
    resources:
      - clusternetworks
    verbs:
      - get
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
      - update
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemSdnReaderYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemSdnReaderYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemSdnReaderYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemSdnReaderYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-sdn-reader.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac20_clusterroleSystemWebhookYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:webhook
rules:
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - buildconfigs/webhooks
    verbs:
      - create
      - get
`)

func v3110OpenshiftApiserverRbac20_clusterroleSystemWebhookYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac20_clusterroleSystemWebhookYaml, nil
}

func v3110OpenshiftApiserverRbac20_clusterroleSystemWebhookYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac20_clusterroleSystemWebhookYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-webhook.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingBasicUsersYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: basic-users
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: basic-user
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:authenticated
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingBasicUsersYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingBasicUsersYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingBasicUsersYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingBasicUsersYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-basic-users.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingClusterAdminsYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: cluster-admins
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:cluster-admins
  - apiGroup: rbac.authorization.k8s.io
    kind: User
    name: system:admin
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingClusterAdminsYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingClusterAdminsYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingClusterAdminsYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingClusterAdminsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-cluster-admins.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingClusterReadersYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: cluster-readers
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-reader
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:cluster-readers
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingClusterReadersYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingClusterReadersYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingClusterReadersYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingClusterReadersYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-cluster-readers.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingClusterStatusBindingYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: cluster-status-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-status
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:authenticated
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingClusterStatusBindingYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingClusterStatusBindingYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingClusterStatusBindingYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingClusterStatusBindingYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-cluster-status-binding.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSelfAccessReviewersYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: self-access-reviewers
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: self-access-reviewer
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:authenticated
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:unauthenticated
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSelfAccessReviewersYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSelfAccessReviewersYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSelfAccessReviewersYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSelfAccessReviewersYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-self-access-reviewers.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSelfProvisionersYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: self-provisioners
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: self-provisioner
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:authenticated:oauth
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSelfProvisionersYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSelfProvisionersYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSelfProvisionersYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSelfProvisionersYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-self-provisioners.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategyDockerBindingYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:build-strategy-docker-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:build-strategy-docker
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:authenticated
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategyDockerBindingYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategyDockerBindingYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategyDockerBindingYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategyDockerBindingYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-build-strategy-docker-binding.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategyJenkinspipelineBindingYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:build-strategy-jenkinspipeline-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:build-strategy-jenkinspipeline
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:authenticated
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategyJenkinspipelineBindingYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategyJenkinspipelineBindingYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategyJenkinspipelineBindingYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategyJenkinspipelineBindingYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-build-strategy-jenkinspipeline-binding.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategySourceBindingYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:build-strategy-source-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:build-strategy-source
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:authenticated
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategySourceBindingYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategySourceBindingYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategySourceBindingYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategySourceBindingYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-build-strategy-source-binding.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemControllerHorizontalPodAutoscalerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:controller:horizontal-pod-autoscaler
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:controller:horizontal-pod-autoscaler
subjects:
  - kind: ServiceAccount
    name: horizontal-pod-autoscaler
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemControllerHorizontalPodAutoscalerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemControllerHorizontalPodAutoscalerYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemControllerHorizontalPodAutoscalerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSystemControllerHorizontalPodAutoscalerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-controller-horizontal-pod-autoscaler.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemMastersYaml = []byte(`apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:masters
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:master
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:masters
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemMastersYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemMastersYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemMastersYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSystemMastersYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-masters.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeAdminsYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:node-admins
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:node-admin
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: User
    name: system:master
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:node-admins
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeAdminsYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeAdminsYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeAdminsYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeAdminsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-node-admins.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeBootstrapperYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:node-bootstrapper
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:node-bootstrapper
subjects:
  - kind: ServiceAccount
    name: node-bootstrapper
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeBootstrapperYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeBootstrapperYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeBootstrapperYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeBootstrapperYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-node-bootstrapper.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeProxiersYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:node-proxiers
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:node-proxier
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:nodes
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeProxiersYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeProxiersYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeProxiersYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeProxiersYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-node-proxiers.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOauthTokenDeletersYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:oauth-token-deleters
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:oauth-token-deleter
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:authenticated
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:unauthenticated
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOauthTokenDeletersYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOauthTokenDeletersYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOauthTokenDeletersYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOauthTokenDeletersYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-oauth-token-deleters.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerClusterQuotaReconciliationControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:cluster-quota-reconciliation-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:controller:cluster-quota-reconciliation-controller
subjects:
  - kind: ServiceAccount
    name: cluster-quota-reconciliation-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerClusterQuotaReconciliationControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerClusterQuotaReconciliationControllerYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerClusterQuotaReconciliationControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerClusterQuotaReconciliationControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-controller-cluster-quota-reconciliation-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerHorizontalPodAutoscalerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:horizontal-pod-autoscaler
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:controller:horizontal-pod-autoscaler
subjects:
  - kind: ServiceAccount
    name: horizontal-pod-autoscaler
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerHorizontalPodAutoscalerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerHorizontalPodAutoscalerYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerHorizontalPodAutoscalerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerHorizontalPodAutoscalerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-controller-horizontal-pod-autoscaler.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerPvRecyclerControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:pv-recycler-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:controller:pv-recycler-controller
subjects:
  - kind: ServiceAccount
    name: pv-recycler-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerPvRecyclerControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerPvRecyclerControllerYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerPvRecyclerControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerPvRecyclerControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-controller-pv-recycler-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerResourcequotaControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:resourcequota-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:controller:resourcequota-controller
subjects:
  - kind: ServiceAccount
    name: resourcequota-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerResourcequotaControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerResourcequotaControllerYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerResourcequotaControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerResourcequotaControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-controller-resourcequota-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerServiceServingCertControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:service-serving-cert-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:controller:service-serving-cert-controller
subjects:
  - kind: ServiceAccount
    name: service-serving-cert-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerServiceServingCertControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerServiceServingCertControllerYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerServiceServingCertControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerServiceServingCertControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-controller-service-serving-cert-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerTemplateServiceBrokerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:template-service-broker
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:controller:template-service-broker
subjects:
  - kind: ServiceAccount
    name: template-service-broker
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerTemplateServiceBrokerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerTemplateServiceBrokerYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerTemplateServiceBrokerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerTemplateServiceBrokerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-controller-template-service-broker.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftDiscoveryYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:discovery
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:discovery
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:authenticated
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftDiscoveryYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftDiscoveryYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftDiscoveryYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftDiscoveryYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-discovery.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemScopeImpersonationYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:scope-impersonation
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:scope-impersonation
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:authenticated
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:unauthenticated
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemScopeImpersonationYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemScopeImpersonationYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemScopeImpersonationYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSystemScopeImpersonationYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-scope-impersonation.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemSdnReadersYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:sdn-readers
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:sdn-reader
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:nodes
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemSdnReadersYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemSdnReadersYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemSdnReadersYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSystemSdnReadersYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-sdn-readers.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemWebhooksYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:webhooks
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:webhook
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:authenticated
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:unauthenticated
`)

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemWebhooksYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac30_clusterrolebindingSystemWebhooksYaml, nil
}

func v3110OpenshiftApiserverRbac30_clusterrolebindingSystemWebhooksYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac30_clusterrolebindingSystemWebhooksYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-webhooks.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac50_aggregateToAdminYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  labels:
    rbac.authorization.k8s.io/aggregate-to-admin: "true"
  name: system:openshift:aggregate-to-admin
rules:
  - apiGroups:
      - ""
      - authorization.openshift.io
    resources:
      - rolebindings
      - roles
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
      - authorization.openshift.io
    resources:
      - localresourceaccessreviews
      - localsubjectaccessreviews
      - subjectrulesreviews
    verbs:
      - create
  - apiGroups:
      - ""
      - security.openshift.io
    resources:
      - podsecuritypolicyreviews
      - podsecuritypolicyselfsubjectreviews
      - podsecuritypolicysubjectreviews
    verbs:
      - create
  - apiGroups:
      - ""
      - authorization.openshift.io
    resources:
      - rolebindingrestrictions
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - buildconfigs
      - buildconfigs/webhooks
      - builds
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - builds/log
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - buildconfigs/instantiate
      - buildconfigs/instantiatebinary
      - builds/clone
    verbs:
      - create
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - builds/details
    verbs:
      - update
  - apiGroups:
      - build.openshift.io
    resources:
      - jenkins
    verbs:
      - admin
      - edit
      - view
  - apiGroups:
      - ""
      - apps.openshift.io
    resources:
      - deploymentconfigs
      - deploymentconfigs/scale
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
      - apps.openshift.io
    resources:
      - deploymentconfigrollbacks
      - deploymentconfigs/instantiate
      - deploymentconfigs/rollback
    verbs:
      - create
  - apiGroups:
      - ""
      - apps.openshift.io
    resources:
      - deploymentconfigs/log
      - deploymentconfigs/status
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreamimages
      - imagestreammappings
      - imagestreams
      - imagestreams/secrets
      - imagestreamtags
      - imagetags
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreams/status
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreams/layers
    verbs:
      - get
      - update
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreamimports
    verbs:
      - create
  - apiGroups:
      - ""
      - project.openshift.io
    resources:
      - projects
    verbs:
      - delete
      - get
      - patch
      - update
  - apiGroups:
      - ""
      - quota.openshift.io
    resources:
      - appliedclusterresourcequotas
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - route.openshift.io
    resources:
      - routes
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
      - route.openshift.io
    resources:
      - routes/custom-host
    verbs:
      - create
  - apiGroups:
      - ""
      - route.openshift.io
    resources:
      - routes/status
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - route.openshift.io
    resources:
      - routes/status
    verbs:
      - update
  - apiGroups:
      - ""
      - template.openshift.io
    resources:
      - processedtemplates
      - templateconfigs
      - templateinstances
      - templates
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - networking.k8s.io
    resources:
      - networkpolicies
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - buildlogs
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
    resources:
      - resourcequotausages
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - authorization.openshift.io
    resources:
      - resourceaccessreviews
      - subjectaccessreviews
    verbs:
      - create
`)

func v3110OpenshiftApiserverRbac50_aggregateToAdminYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac50_aggregateToAdminYaml, nil
}

func v3110OpenshiftApiserverRbac50_aggregateToAdminYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac50_aggregateToAdminYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/50_aggregate-to-admin.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac50_aggregateToBasicUserYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  labels:
    authorization.openshift.io/aggregate-to-basic-user: "true"
  name: system:openshift:aggregate-to-basic-user
rules:
  - apiGroups:
      - ""
      - user.openshift.io
    resourceNames:
      - "~"
    resources:
      - users
    verbs:
      - get
  - apiGroups:
      - ""
      - project.openshift.io
    resources:
      - projectrequests
    verbs:
      - list
  - apiGroups:
      - ""
      - authorization.openshift.io
    resources:
      - clusterroles
    verbs:
      - get
      - list
  - apiGroups:
      - rbac.authorization.k8s.io
    resources:
      - clusterroles
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - storage.k8s.io
    resources:
      - storageclasses
    verbs:
      - get
      - list
  - apiGroups:
      - ""
      - project.openshift.io
    resources:
      - projects
    verbs:
      - list
      - watch
  - apiGroups:
      - ""
      - authorization.openshift.io
    resources:
      - selfsubjectrulesreviews
    verbs:
      - create
  - apiGroups:
      - authorization.k8s.io
    resources:
      - selfsubjectaccessreviews
    verbs:
      - create
`)

func v3110OpenshiftApiserverRbac50_aggregateToBasicUserYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac50_aggregateToBasicUserYaml, nil
}

func v3110OpenshiftApiserverRbac50_aggregateToBasicUserYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac50_aggregateToBasicUserYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/50_aggregate-to-basic-user.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac50_aggregateToClusterReaderYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  labels:
    rbac.authorization.k8s.io/aggregate-to-cluster-reader: "true"
  name: system:openshift:aggregate-to-cluster-reader
rules:
  - apiGroups:
      - ""
    resources:
      - componentstatuses
      - nodes
      - nodes/status
      - persistentvolumeclaims/status
      - persistentvolumes
      - persistentvolumes/status
      - pods/binding
      - pods/eviction
      - podtemplates
      - securitycontextconstraints
      - services/status
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - admissionregistration.k8s.io
    resources:
      - mutatingwebhookconfigurations
      - validatingwebhookconfigurations
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - apps
    resources:
      - controllerrevisions
      - daemonsets/status
      - deployments/status
      - replicasets/status
      - statefulsets/status
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - apiextensions.k8s.io
    resources:
      - customresourcedefinitions
      - customresourcedefinitions/status
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - apiregistration.k8s.io
    resources:
      - apiservices
      - apiservices/status
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - autoscaling
    resources:
      - horizontalpodautoscalers/status
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - batch
    resources:
      - cronjobs/status
      - jobs/status
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - coordination.k8s.io
    resources:
      - leases
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - events.k8s.io
    resources:
      - events
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - networking.k8s.io
    resources:
      - ingresses/status
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - node.k8s.io
    resources:
      - runtimeclasses
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - policy
    resources:
      - poddisruptionbudgets/status
      - podsecuritypolicies
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - rbac.authorization.k8s.io
    resources:
      - clusterrolebindings
      - clusterroles
      - rolebindings
      - roles
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - storage.k8s.io
    resources:
      - csidrivers
      - csinodes
      - storageclasses
      - volumeattachments
      - volumeattachments/status
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - scheduling.k8s.io
    resources:
      - priorityclasses
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - certificates.k8s.io
    resources:
      - certificatesigningrequests
      - certificatesigningrequests/approval
      - certificatesigningrequests/status
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - authorization.openshift.io
    resources:
      - clusterrolebindings
      - clusterroles
      - rolebindingrestrictions
      - rolebindings
      - roles
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - builds/details
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - images
      - imagesignatures
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreams/layers
    verbs:
      - get
  - apiGroups:
      - ""
      - oauth.openshift.io
    resources:
      - oauthclientauthorizations
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - project.openshift.io
    resources:
      - projects
    verbs:
      - list
      - watch
  - apiGroups:
      - ""
      - project.openshift.io
    resources:
      - projectrequests
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - quota.openshift.io
    resources:
      - clusterresourcequotas
      - clusterresourcequotas/status
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - network.openshift.io
    resources:
      - clusternetworks
      - egressnetworkpolicies
      - hostsubnets
      - netnamespaces
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - security.openshift.io
    resources:
      - securitycontextconstraints
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - security.openshift.io
    resources:
      - rangeallocations
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - template.openshift.io
    resources:
      - brokertemplateinstances
      - templateinstances/status
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - user.openshift.io
    resources:
      - groups
      - identities
      - useridentitymappings
      - users
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - authorization.openshift.io
    resources:
      - localresourceaccessreviews
      - localsubjectaccessreviews
      - resourceaccessreviews
      - selfsubjectrulesreviews
      - subjectaccessreviews
      - subjectrulesreviews
    verbs:
      - create
  - apiGroups:
      - authorization.k8s.io
    resources:
      - localsubjectaccessreviews
      - selfsubjectaccessreviews
      - selfsubjectrulesreviews
      - subjectaccessreviews
    verbs:
      - create
  - apiGroups:
      - authentication.k8s.io
    resources:
      - tokenreviews
    verbs:
      - create
  - apiGroups:
      - ""
      - security.openshift.io
    resources:
      - podsecuritypolicyreviews
      - podsecuritypolicyselfsubjectreviews
      - podsecuritypolicysubjectreviews
    verbs:
      - create
  - apiGroups:
      - ""
    resources:
      - nodes/metrics
      - nodes/spec
    verbs:
      - get
  - apiGroups:
      - ""
    resources:
      - nodes/stats
    verbs:
      - create
      - get
  - nonResourceURLs:
      - '*'
    verbs:
      - get
`)

func v3110OpenshiftApiserverRbac50_aggregateToClusterReaderYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac50_aggregateToClusterReaderYaml, nil
}

func v3110OpenshiftApiserverRbac50_aggregateToClusterReaderYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac50_aggregateToClusterReaderYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/50_aggregate-to-cluster-reader.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac50_aggregateToEditYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  labels:
    rbac.authorization.k8s.io/aggregate-to-edit: "true"
  name: system:openshift:aggregate-to-edit
rules:
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - buildconfigs
      - buildconfigs/webhooks
      - builds
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - builds/log
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - buildconfigs/instantiate
      - buildconfigs/instantiatebinary
      - builds/clone
    verbs:
      - create
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - builds/details
    verbs:
      - update
  - apiGroups:
      - build.openshift.io
    resources:
      - jenkins
    verbs:
      - edit
      - view
  - apiGroups:
      - ""
      - apps.openshift.io
    resources:
      - deploymentconfigs
      - deploymentconfigs/scale
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
      - apps.openshift.io
    resources:
      - deploymentconfigrollbacks
      - deploymentconfigs/instantiate
      - deploymentconfigs/rollback
    verbs:
      - create
  - apiGroups:
      - ""
      - apps.openshift.io
    resources:
      - deploymentconfigs/log
      - deploymentconfigs/status
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreamimages
      - imagestreammappings
      - imagestreams
      - imagestreams/secrets
      - imagestreamtags
      - imagetags
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreams/status
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreams/layers
    verbs:
      - get
      - update
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreamimports
    verbs:
      - create
  - apiGroups:
      - ""
      - project.openshift.io
    resources:
      - projects
    verbs:
      - get
  - apiGroups:
      - ""
      - quota.openshift.io
    resources:
      - appliedclusterresourcequotas
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - route.openshift.io
    resources:
      - routes
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
      - route.openshift.io
    resources:
      - routes/custom-host
    verbs:
      - create
  - apiGroups:
      - ""
      - route.openshift.io
    resources:
      - routes/status
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - template.openshift.io
    resources:
      - processedtemplates
      - templateconfigs
      - templateinstances
      - templates
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - networking.k8s.io
    resources:
      - networkpolicies
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - buildlogs
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
    resources:
      - resourcequotausages
    verbs:
      - get
      - list
      - watch
`)

func v3110OpenshiftApiserverRbac50_aggregateToEditYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac50_aggregateToEditYaml, nil
}

func v3110OpenshiftApiserverRbac50_aggregateToEditYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac50_aggregateToEditYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/50_aggregate-to-edit.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac50_aggregateToStorageAdminYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  labels:
    storage.openshift.io/aggregate-to-storage-admin: "true"
  name: system:openshift:aggregate-to-storage-admin
rules:
  - apiGroups:
      - ""
    resources:
      - persistentvolumes
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - storage.k8s.io
    resources:
      - storageclasses
    verbs:
      - create
      - delete
      - deletecollection
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
    resources:
      - events
      - persistentvolumeclaims
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - pods
    verbs:
      - get
      - list
      - watch
`)

func v3110OpenshiftApiserverRbac50_aggregateToStorageAdminYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac50_aggregateToStorageAdminYaml, nil
}

func v3110OpenshiftApiserverRbac50_aggregateToStorageAdminYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac50_aggregateToStorageAdminYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/50_aggregate-to-storage-admin.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbac50_aggregateToViewYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  labels:
    rbac.authorization.k8s.io/aggregate-to-view: "true"
  name: system:openshift:aggregate-to-view
rules:
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - buildconfigs
      - buildconfigs/webhooks
      - builds
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - builds/log
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - build.openshift.io
    resources:
      - jenkins
    verbs:
      - view
  - apiGroups:
      - ""
      - apps.openshift.io
    resources:
      - deploymentconfigs
      - deploymentconfigs/scale
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - apps.openshift.io
    resources:
      - deploymentconfigs/log
      - deploymentconfigs/status
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreamimages
      - imagestreammappings
      - imagestreams
      - imagestreamtags
      - imagetags
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreams/status
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - project.openshift.io
    resources:
      - projects
    verbs:
      - get
  - apiGroups:
      - ""
      - quota.openshift.io
    resources:
      - appliedclusterresourcequotas
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - route.openshift.io
    resources:
      - routes
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - route.openshift.io
    resources:
      - routes/status
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - template.openshift.io
    resources:
      - processedtemplates
      - templateconfigs
      - templateinstances
      - templates
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - buildlogs
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - resourcequotausages
    verbs:
      - get
      - list
      - watch
`)

func v3110OpenshiftApiserverRbac50_aggregateToViewYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbac50_aggregateToViewYaml, nil
}

func v3110OpenshiftApiserverRbac50_aggregateToViewYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbac50_aggregateToViewYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac/50_aggregate-to-view.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerBuildConfigChangeControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:build-config-change-controller
rules:
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - buildconfigs
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - buildconfigs/instantiate
    verbs:
      - create
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - builds
    verbs:
      - delete
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
      - update
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerBuildConfigChangeControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerBuildConfigChangeControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerBuildConfigChangeControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerBuildConfigChangeControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-build-config-change-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerBuildControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:build-controller
rules:
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - builds
    verbs:
      - delete
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - builds/finalizers
    verbs:
      - update
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - buildconfigs
    verbs:
      - get
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - builds/custom
      - builds/docker
      - builds/jenkinspipeline
      - builds/optimizeddocker
      - builds/source
    verbs:
      - create
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreams
    verbs:
      - get
      - list
  - apiGroups:
      - ""
    resources:
      - secrets
    verbs:
      - get
      - list
  - apiGroups:
      - ""
    resources:
      - configmaps
    verbs:
      - create
      - get
      - list
  - apiGroups:
      - ""
    resources:
      - pods
    verbs:
      - create
      - delete
      - get
      - list
  - apiGroups:
      - ""
    resources:
      - namespaces
    verbs:
      - get
  - apiGroups:
      - ""
    resources:
      - serviceaccounts
    verbs:
      - get
      - list
  - apiGroups:
      - ""
      - security.openshift.io
    resources:
      - podsecuritypolicysubjectreviews
    verbs:
      - create
  - apiGroups:
      - config.openshift.io
    resources:
      - builds
    verbs:
      - get
      - list
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
      - update
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerBuildControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerBuildControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerBuildControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerBuildControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-build-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDefaultRolebindingsControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:default-rolebindings-controller
rules:
  - apiGroups:
      - rbac.authorization.k8s.io
    resources:
      - rolebindings
    verbs:
      - create
  - apiGroups:
      - ""
    resources:
      - namespaces
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - rbac.authorization.k8s.io
    resources:
      - rolebindings
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
      - update
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDefaultRolebindingsControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDefaultRolebindingsControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDefaultRolebindingsControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDefaultRolebindingsControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-default-rolebindings-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDeployerControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:deployer-controller
rules:
  - apiGroups:
      - ""
    resources:
      - pods
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - watch
  - apiGroups:
      - ""
    resources:
      - replicationcontrollers
    verbs:
      - delete
  - apiGroups:
      - ""
    resources:
      - replicationcontrollers
    verbs:
      - get
      - list
      - update
      - watch
  - apiGroups:
      - ""
    resources:
      - replicationcontrollers/scale
    verbs:
      - get
      - update
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
      - update
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDeployerControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDeployerControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDeployerControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDeployerControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-deployer-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDeploymentconfigControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:deploymentconfig-controller
rules:
  - apiGroups:
      - ""
    resources:
      - replicationcontrollers
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
    resources:
      - replicationcontrollers/scale
    verbs:
      - get
      - update
  - apiGroups:
      - ""
      - apps.openshift.io
    resources:
      - deploymentconfigs/status
    verbs:
      - update
  - apiGroups:
      - ""
      - apps.openshift.io
    resources:
      - deploymentconfigs/finalizers
    verbs:
      - update
  - apiGroups:
      - ""
      - apps.openshift.io
    resources:
      - deploymentconfigs
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
      - update
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDeploymentconfigControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDeploymentconfigControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDeploymentconfigControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDeploymentconfigControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-deploymentconfig-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerImageImportControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:image-import-controller
rules:
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreams
    verbs:
      - create
      - get
      - list
      - update
      - watch
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - images
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreamimports
    verbs:
      - create
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
      - update
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerImageImportControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerImageImportControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerImageImportControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerImageImportControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-image-import-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerImageTriggerControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:image-trigger-controller
rules:
  - apiGroups:
      - ""
      - image.openshift.io
    resources:
      - imagestreams
    verbs:
      - list
      - watch
  - apiGroups:
      - extensions
    resources:
      - daemonsets
    verbs:
      - get
      - update
  - apiGroups:
      - apps
      - extensions
    resources:
      - deployments
    verbs:
      - get
      - update
  - apiGroups:
      - apps
    resources:
      - statefulsets
    verbs:
      - get
      - update
  - apiGroups:
      - batch
    resources:
      - cronjobs
    verbs:
      - get
      - update
  - apiGroups:
      - ""
      - apps.openshift.io
    resources:
      - deploymentconfigs
    verbs:
      - get
      - update
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - buildconfigs/instantiate
    verbs:
      - create
  - apiGroups:
      - ""
      - build.openshift.io
    resources:
      - builds/custom
      - builds/docker
      - builds/jenkinspipeline
      - builds/optimizeddocker
      - builds/source
    verbs:
      - create
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
      - update
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerImageTriggerControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerImageTriggerControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerImageTriggerControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerImageTriggerControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-image-trigger-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerOriginNamespaceControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:origin-namespace-controller
rules:
  - apiGroups:
      - ""
    resources:
      - namespaces
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - namespaces/finalize
      - namespaces/status
    verbs:
      - update
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
      - update
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerOriginNamespaceControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerOriginNamespaceControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerOriginNamespaceControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerOriginNamespaceControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-origin-namespace-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceIngressIpControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:service-ingress-ip-controller
rules:
  - apiGroups:
      - ""
    resources:
      - services
    verbs:
      - list
      - update
      - watch
  - apiGroups:
      - ""
    resources:
      - services/status
    verbs:
      - update
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
      - update
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceIngressIpControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceIngressIpControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceIngressIpControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceIngressIpControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-service-ingress-ip-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceaccountControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:serviceaccount-controller
rules:
  - apiGroups:
      - ""
    resources:
      - serviceaccounts
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
      - update
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceaccountControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceaccountControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceaccountControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceaccountControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-serviceaccount-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceaccountPullSecretsControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:serviceaccount-pull-secrets-controller
rules:
  - apiGroups:
      - ""
    resources:
      - serviceaccounts
    verbs:
      - create
      - get
      - list
      - update
      - watch
  - apiGroups:
      - ""
    resources:
      - secrets
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
    resources:
      - services
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
      - update
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceaccountPullSecretsControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceaccountPullSecretsControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceaccountPullSecretsControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceaccountPullSecretsControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-serviceaccount-pull-secrets-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerTemplateInstanceControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:template-instance-controller
rules:
  - apiGroups:
      - authorization.k8s.io
    resources:
      - subjectaccessreviews
    verbs:
      - create
  - apiGroups:
      - template.openshift.io
    resources:
      - templateinstances/status
    verbs:
      - update
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerTemplateInstanceControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerTemplateInstanceControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerTemplateInstanceControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerTemplateInstanceControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-template-instance-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerTemplateInstanceFinalizerControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:template-instance-finalizer-controller
rules:
  - apiGroups:
      - template.openshift.io
    resources:
      - templateinstances/status
    verbs:
      - update
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerTemplateInstanceFinalizerControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerTemplateInstanceFinalizerControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerTemplateInstanceFinalizerControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerTemplateInstanceFinalizerControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-template-instance-finalizer-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerUnidlingControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:unidling-controller
rules:
  - apiGroups:
      - ""
    resources:
      - endpoints
      - replicationcontrollers/scale
      - services
    verbs:
      - get
      - update
  - apiGroups:
      - ""
    resources:
      - replicationcontrollers
    verbs:
      - get
      - patch
      - update
  - apiGroups:
      - ""
      - apps.openshift.io
    resources:
      - deploymentconfigs
    verbs:
      - get
      - patch
      - update
  - apiGroups:
      - apps
      - extensions
    resources:
      - deployments/scale
      - replicasets/scale
    verbs:
      - get
      - update
  - apiGroups:
      - ""
      - apps.openshift.io
    resources:
      - deploymentconfigs/scale
    verbs:
      - get
      - update
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - events
    verbs:
      - create
      - patch
      - update
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerUnidlingControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerUnidlingControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerUnidlingControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerUnidlingControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-unidling-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemDeployerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:deployer
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:deployer
subjects:
  - kind: ServiceAccount
    name: default-rolebindings-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemDeployerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemDeployerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemDeployerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemDeployerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-deployer.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemImageBuilderYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:image-builder
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:image-builder
subjects:
  - kind: ServiceAccount
    name: default-rolebindings-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemImageBuilderYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemImageBuilderYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemImageBuilderYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemImageBuilderYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-image-builder.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemImagePullerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:image-puller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:image-puller
subjects:
  - kind: ServiceAccount
    name: default-rolebindings-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemImagePullerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemImagePullerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemImagePullerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemImagePullerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-image-puller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerBuildConfigChangeControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:build-config-change-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:controller:build-config-change-controller
subjects:
  - kind: ServiceAccount
    name: build-config-change-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerBuildConfigChangeControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerBuildConfigChangeControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerBuildConfigChangeControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerBuildConfigChangeControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-build-config-change-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerBuildControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:build-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:controller:build-controller
subjects:
  - kind: ServiceAccount
    name: build-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerBuildControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerBuildControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerBuildControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerBuildControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-build-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDefaultRolebindingsControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:default-rolebindings-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:controller:default-rolebindings-controller
subjects:
  - kind: ServiceAccount
    name: default-rolebindings-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDefaultRolebindingsControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDefaultRolebindingsControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDefaultRolebindingsControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDefaultRolebindingsControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-default-rolebindings-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDeployerControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:deployer-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:controller:deployer-controller
subjects:
  - kind: ServiceAccount
    name: deployer-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDeployerControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDeployerControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDeployerControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDeployerControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-deployer-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDeploymentconfigControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:deploymentconfig-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:controller:deploymentconfig-controller
subjects:
  - kind: ServiceAccount
    name: deploymentconfig-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDeploymentconfigControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDeploymentconfigControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDeploymentconfigControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDeploymentconfigControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-deploymentconfig-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerImageImportControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:image-import-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:controller:image-import-controller
subjects:
  - kind: ServiceAccount
    name: image-import-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerImageImportControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerImageImportControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerImageImportControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerImageImportControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-image-import-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerImageTriggerControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:image-trigger-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:controller:image-trigger-controller
subjects:
  - kind: ServiceAccount
    name: image-trigger-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerImageTriggerControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerImageTriggerControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerImageTriggerControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerImageTriggerControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-image-trigger-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerOriginNamespaceControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:origin-namespace-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:controller:origin-namespace-controller
subjects:
  - kind: ServiceAccount
    name: origin-namespace-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerOriginNamespaceControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerOriginNamespaceControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerOriginNamespaceControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerOriginNamespaceControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-origin-namespace-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceIngressIpControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:service-ingress-ip-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:controller:service-ingress-ip-controller
subjects:
  - kind: ServiceAccount
    name: service-ingress-ip-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceIngressIpControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceIngressIpControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceIngressIpControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceIngressIpControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-service-ingress-ip-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceaccountControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:serviceaccount-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:controller:serviceaccount-controller
subjects:
  - kind: ServiceAccount
    name: serviceaccount-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceaccountControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceaccountControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceaccountControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceaccountControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-serviceaccount-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceaccountPullSecretsControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:serviceaccount-pull-secrets-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:controller:serviceaccount-pull-secrets-controller
subjects:
  - kind: ServiceAccount
    name: serviceaccount-pull-secrets-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceaccountPullSecretsControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceaccountPullSecretsControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceaccountPullSecretsControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceaccountPullSecretsControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-serviceaccount-pull-secrets-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceControllerAdminYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:template-instance-controller:admin
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: admin
subjects:
  - kind: ServiceAccount
    name: template-instance-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceControllerAdminYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceControllerAdminYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceControllerAdminYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceControllerAdminYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-template-instance-controller-admin.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:template-instance-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:controller:template-instance-controller
subjects:
  - kind: ServiceAccount
    name: template-instance-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-template-instance-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceFinalizerControllerAdminYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:template-instance-finalizer-controller:admin
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: admin
subjects:
  - kind: ServiceAccount
    name: template-instance-finalizer-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceFinalizerControllerAdminYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceFinalizerControllerAdminYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceFinalizerControllerAdminYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceFinalizerControllerAdminYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-template-instance-finalizer-controller-admin.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceFinalizerControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:template-instance-finalizer-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:controller:template-instance-finalizer-controller
subjects:
  - kind: ServiceAccount
    name: template-instance-finalizer-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceFinalizerControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceFinalizerControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceFinalizerControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceFinalizerControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-template-instance-finalizer-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerUnidlingControllerYaml = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:openshift:controller:unidling-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:openshift:controller:unidling-controller
subjects:
  - kind: ServiceAccount
    name: unidling-controller
    namespace: openshift-infra
`)

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerUnidlingControllerYamlBytes() ([]byte, error) {
	return _v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerUnidlingControllerYaml, nil
}

func v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerUnidlingControllerYaml() (*asset, error) {
	bytes, err := v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerUnidlingControllerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-unidling-controller.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
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
	"v3.11.0/config/defaultconfig.yaml":                                                                                                                              v3110ConfigDefaultconfigYaml,
	"v3.11.0/openshift-apiserver/apiserver-clusterrolebinding.yaml":                                                                                                  v3110OpenshiftApiserverApiserverClusterrolebindingYaml,
	"v3.11.0/openshift-apiserver/cm.yaml":                                                                                                                            v3110OpenshiftApiserverCmYaml,
	"v3.11.0/openshift-apiserver/deploy.yaml":                                                                                                                        v3110OpenshiftApiserverDeployYaml,
	"v3.11.0/openshift-apiserver/ns.yaml":                                                                                                                            v3110OpenshiftApiserverNsYaml,
	"v3.11.0/openshift-apiserver/pdb.yaml":                                                                                                                           v3110OpenshiftApiserverPdbYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-basic-user.yaml":                                                                                                v3110OpenshiftApiserverRbac20_clusterroleBasicUserYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-cluster-debugger.yaml":                                                                                          v3110OpenshiftApiserverRbac20_clusterroleClusterDebuggerYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-cluster-reader.yaml":                                                                                            v3110OpenshiftApiserverRbac20_clusterroleClusterReaderYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-cluster-status.yaml":                                                                                            v3110OpenshiftApiserverRbac20_clusterroleClusterStatusYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-registry-admin.yaml":                                                                                            v3110OpenshiftApiserverRbac20_clusterroleRegistryAdminYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-registry-editor.yaml":                                                                                           v3110OpenshiftApiserverRbac20_clusterroleRegistryEditorYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-registry-viewer.yaml":                                                                                           v3110OpenshiftApiserverRbac20_clusterroleRegistryViewerYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-self-access-reviewer.yaml":                                                                                      v3110OpenshiftApiserverRbac20_clusterroleSelfAccessReviewerYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-self-provisioner.yaml":                                                                                          v3110OpenshiftApiserverRbac20_clusterroleSelfProvisionerYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-storage-admin.yaml":                                                                                             v3110OpenshiftApiserverRbac20_clusterroleStorageAdminYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-sudoer.yaml":                                                                                                    v3110OpenshiftApiserverRbac20_clusterroleSudoerYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-build-strategy-custom.yaml":                                                                              v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyCustomYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-build-strategy-docker.yaml":                                                                              v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyDockerYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-build-strategy-jenkinspipeline.yaml":                                                                     v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyJenkinspipelineYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-build-strategy-source.yaml":                                                                              v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategySourceYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-deployer.yaml":                                                                                           v3110OpenshiftApiserverRbac20_clusterroleSystemDeployerYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-image-auditor.yaml":                                                                                      v3110OpenshiftApiserverRbac20_clusterroleSystemImageAuditorYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-image-builder.yaml":                                                                                      v3110OpenshiftApiserverRbac20_clusterroleSystemImageBuilderYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-image-pruner.yaml":                                                                                       v3110OpenshiftApiserverRbac20_clusterroleSystemImagePrunerYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-image-puller.yaml":                                                                                       v3110OpenshiftApiserverRbac20_clusterroleSystemImagePullerYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-image-pusher.yaml":                                                                                       v3110OpenshiftApiserverRbac20_clusterroleSystemImagePusherYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-image-signer.yaml":                                                                                       v3110OpenshiftApiserverRbac20_clusterroleSystemImageSignerYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-master.yaml":                                                                                             v3110OpenshiftApiserverRbac20_clusterroleSystemMasterYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-node-admin.yaml":                                                                                         v3110OpenshiftApiserverRbac20_clusterroleSystemNodeAdminYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-node-reader.yaml":                                                                                        v3110OpenshiftApiserverRbac20_clusterroleSystemNodeReaderYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-oauth-token-deleter.yaml":                                                                                v3110OpenshiftApiserverRbac20_clusterroleSystemOauthTokenDeleterYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-controller-cluster-quota-reconciliation-controller.yaml":                                       v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerClusterQuotaReconciliationControllerYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-controller-horizontal-pod-autoscaler.yaml":                                                     v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerHorizontalPodAutoscalerYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-controller-pv-recycler-controller.yaml":                                                        v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerPvRecyclerControllerYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-controller-resourcequota-controller.yaml":                                                      v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerResourcequotaControllerYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-controller-service-serving-cert-controller.yaml":                                               v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerServiceServingCertControllerYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-controller-template-service-broker.yaml":                                                       v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerTemplateServiceBrokerYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-discovery.yaml":                                                                                v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftDiscoveryYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-openshift-templateservicebroker-client.yaml":                                                             v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftTemplateservicebrokerClientYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-router.yaml":                                                                                             v3110OpenshiftApiserverRbac20_clusterroleSystemRouterYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-scope-impersonation.yaml":                                                                                v3110OpenshiftApiserverRbac20_clusterroleSystemScopeImpersonationYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-sdn-manager.yaml":                                                                                        v3110OpenshiftApiserverRbac20_clusterroleSystemSdnManagerYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-sdn-reader.yaml":                                                                                         v3110OpenshiftApiserverRbac20_clusterroleSystemSdnReaderYaml,
	"v3.11.0/openshift-apiserver/rbac/20_clusterrole-system-webhook.yaml":                                                                                            v3110OpenshiftApiserverRbac20_clusterroleSystemWebhookYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-basic-users.yaml":                                                                                        v3110OpenshiftApiserverRbac30_clusterrolebindingBasicUsersYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-cluster-admins.yaml":                                                                                     v3110OpenshiftApiserverRbac30_clusterrolebindingClusterAdminsYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-cluster-readers.yaml":                                                                                    v3110OpenshiftApiserverRbac30_clusterrolebindingClusterReadersYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-cluster-status-binding.yaml":                                                                             v3110OpenshiftApiserverRbac30_clusterrolebindingClusterStatusBindingYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-self-access-reviewers.yaml":                                                                              v3110OpenshiftApiserverRbac30_clusterrolebindingSelfAccessReviewersYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-self-provisioners.yaml":                                                                                  v3110OpenshiftApiserverRbac30_clusterrolebindingSelfProvisionersYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-build-strategy-docker-binding.yaml":                                                               v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategyDockerBindingYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-build-strategy-jenkinspipeline-binding.yaml":                                                      v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategyJenkinspipelineBindingYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-build-strategy-source-binding.yaml":                                                               v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategySourceBindingYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-controller-horizontal-pod-autoscaler.yaml":                                                        v3110OpenshiftApiserverRbac30_clusterrolebindingSystemControllerHorizontalPodAutoscalerYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-masters.yaml":                                                                                     v3110OpenshiftApiserverRbac30_clusterrolebindingSystemMastersYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-node-admins.yaml":                                                                                 v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeAdminsYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-node-bootstrapper.yaml":                                                                           v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeBootstrapperYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-node-proxiers.yaml":                                                                               v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeProxiersYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-oauth-token-deleters.yaml":                                                                        v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOauthTokenDeletersYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-controller-cluster-quota-reconciliation-controller.yaml":                                v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerClusterQuotaReconciliationControllerYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-controller-horizontal-pod-autoscaler.yaml":                                              v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerHorizontalPodAutoscalerYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-controller-pv-recycler-controller.yaml":                                                 v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerPvRecyclerControllerYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-controller-resourcequota-controller.yaml":                                               v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerResourcequotaControllerYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-controller-service-serving-cert-controller.yaml":                                        v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerServiceServingCertControllerYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-controller-template-service-broker.yaml":                                                v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerTemplateServiceBrokerYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-openshift-discovery.yaml":                                                                         v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftDiscoveryYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-scope-impersonation.yaml":                                                                         v3110OpenshiftApiserverRbac30_clusterrolebindingSystemScopeImpersonationYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-sdn-readers.yaml":                                                                                 v3110OpenshiftApiserverRbac30_clusterrolebindingSystemSdnReadersYaml,
	"v3.11.0/openshift-apiserver/rbac/30_clusterrolebinding-system-webhooks.yaml":                                                                                    v3110OpenshiftApiserverRbac30_clusterrolebindingSystemWebhooksYaml,
	"v3.11.0/openshift-apiserver/rbac/50_aggregate-to-admin.yaml":                                                                                                    v3110OpenshiftApiserverRbac50_aggregateToAdminYaml,
	"v3.11.0/openshift-apiserver/rbac/50_aggregate-to-basic-user.yaml":                                                                                               v3110OpenshiftApiserverRbac50_aggregateToBasicUserYaml,
	"v3.11.0/openshift-apiserver/rbac/50_aggregate-to-cluster-reader.yaml":                                                                                           v3110OpenshiftApiserverRbac50_aggregateToClusterReaderYaml,
	"v3.11.0/openshift-apiserver/rbac/50_aggregate-to-edit.yaml":                                                                                                     v3110OpenshiftApiserverRbac50_aggregateToEditYaml,
	"v3.11.0/openshift-apiserver/rbac/50_aggregate-to-storage-admin.yaml":                                                                                            v3110OpenshiftApiserverRbac50_aggregateToStorageAdminYaml,
	"v3.11.0/openshift-apiserver/rbac/50_aggregate-to-view.yaml":                                                                                                     v3110OpenshiftApiserverRbac50_aggregateToViewYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-build-config-change-controller.yaml":                      v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerBuildConfigChangeControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-build-controller.yaml":                                    v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerBuildControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-default-rolebindings-controller.yaml":                     v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDefaultRolebindingsControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-deployer-controller.yaml":                                 v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDeployerControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-deploymentconfig-controller.yaml":                         v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDeploymentconfigControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-image-import-controller.yaml":                             v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerImageImportControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-image-trigger-controller.yaml":                            v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerImageTriggerControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-origin-namespace-controller.yaml":                         v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerOriginNamespaceControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-service-ingress-ip-controller.yaml":                       v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceIngressIpControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-serviceaccount-controller.yaml":                           v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceaccountControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-serviceaccount-pull-secrets-controller.yaml":              v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceaccountPullSecretsControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-template-instance-controller.yaml":                        v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerTemplateInstanceControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-template-instance-finalizer-controller.yaml":              v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerTemplateInstanceFinalizerControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrole-system-openshift-controller-unidling-controller.yaml":                                 v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerUnidlingControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-deployer.yaml":                                                          v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemDeployerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-image-builder.yaml":                                                     v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemImageBuilderYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-image-puller.yaml":                                                      v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemImagePullerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-build-config-change-controller.yaml":               v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerBuildConfigChangeControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-build-controller.yaml":                             v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerBuildControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-default-rolebindings-controller.yaml":              v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDefaultRolebindingsControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-deployer-controller.yaml":                          v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDeployerControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-deploymentconfig-controller.yaml":                  v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDeploymentconfigControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-image-import-controller.yaml":                      v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerImageImportControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-image-trigger-controller.yaml":                     v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerImageTriggerControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-origin-namespace-controller.yaml":                  v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerOriginNamespaceControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-service-ingress-ip-controller.yaml":                v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceIngressIpControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-serviceaccount-controller.yaml":                    v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceaccountControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-serviceaccount-pull-secrets-controller.yaml":       v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceaccountPullSecretsControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-template-instance-controller-admin.yaml":           v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceControllerAdminYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-template-instance-controller.yaml":                 v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-template-instance-finalizer-controller-admin.yaml": v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceFinalizerControllerAdminYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-template-instance-finalizer-controller.yaml":       v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceFinalizerControllerYaml,
	"v3.11.0/openshift-apiserver/rbac-openshift-controller-manager/clusterrolebinding-system-openshift-controller-unidling-controller.yaml":                          v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerUnidlingControllerYaml,
	"v3.11.0/openshift-apiserver/sa.yaml":                                                                                                                            v3110OpenshiftApiserverSaYaml,
	"v3.11.0/openshift-apiserver/svc.yaml":                                                                                                                           v3110OpenshiftApiserverSvcYaml,
	"v3.11.0/openshift-apiserver/trusted_ca_cm.yaml":                                                                                                                 v3110OpenshiftApiserverTrusted_ca_cmYaml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
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
			"rbac": {nil, map[string]*bintree{
				"20_clusterrole-basic-user.yaml":                                                                 {v3110OpenshiftApiserverRbac20_clusterroleBasicUserYaml, map[string]*bintree{}},
				"20_clusterrole-cluster-debugger.yaml":                                                           {v3110OpenshiftApiserverRbac20_clusterroleClusterDebuggerYaml, map[string]*bintree{}},
				"20_clusterrole-cluster-reader.yaml":                                                             {v3110OpenshiftApiserverRbac20_clusterroleClusterReaderYaml, map[string]*bintree{}},
				"20_clusterrole-cluster-status.yaml":                                                             {v3110OpenshiftApiserverRbac20_clusterroleClusterStatusYaml, map[string]*bintree{}},
				"20_clusterrole-registry-admin.yaml":                                                             {v3110OpenshiftApiserverRbac20_clusterroleRegistryAdminYaml, map[string]*bintree{}},
				"20_clusterrole-registry-editor.yaml":                                                            {v3110OpenshiftApiserverRbac20_clusterroleRegistryEditorYaml, map[string]*bintree{}},
				"20_clusterrole-registry-viewer.yaml":                                                            {v3110OpenshiftApiserverRbac20_clusterroleRegistryViewerYaml, map[string]*bintree{}},
				"20_clusterrole-self-access-reviewer.yaml":                                                       {v3110OpenshiftApiserverRbac20_clusterroleSelfAccessReviewerYaml, map[string]*bintree{}},
				"20_clusterrole-self-provisioner.yaml":                                                           {v3110OpenshiftApiserverRbac20_clusterroleSelfProvisionerYaml, map[string]*bintree{}},
				"20_clusterrole-storage-admin.yaml":                                                              {v3110OpenshiftApiserverRbac20_clusterroleStorageAdminYaml, map[string]*bintree{}},
				"20_clusterrole-sudoer.yaml":                                                                     {v3110OpenshiftApiserverRbac20_clusterroleSudoerYaml, map[string]*bintree{}},
				"20_clusterrole-system-build-strategy-custom.yaml":                                               {v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyCustomYaml, map[string]*bintree{}},
				"20_clusterrole-system-build-strategy-docker.yaml":                                               {v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyDockerYaml, map[string]*bintree{}},
				"20_clusterrole-system-build-strategy-jenkinspipeline.yaml":                                      {v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategyJenkinspipelineYaml, map[string]*bintree{}},
				"20_clusterrole-system-build-strategy-source.yaml":                                               {v3110OpenshiftApiserverRbac20_clusterroleSystemBuildStrategySourceYaml, map[string]*bintree{}},
				"20_clusterrole-system-deployer.yaml":                                                            {v3110OpenshiftApiserverRbac20_clusterroleSystemDeployerYaml, map[string]*bintree{}},
				"20_clusterrole-system-image-auditor.yaml":                                                       {v3110OpenshiftApiserverRbac20_clusterroleSystemImageAuditorYaml, map[string]*bintree{}},
				"20_clusterrole-system-image-builder.yaml":                                                       {v3110OpenshiftApiserverRbac20_clusterroleSystemImageBuilderYaml, map[string]*bintree{}},
				"20_clusterrole-system-image-pruner.yaml":                                                        {v3110OpenshiftApiserverRbac20_clusterroleSystemImagePrunerYaml, map[string]*bintree{}},
				"20_clusterrole-system-image-puller.yaml":                                                        {v3110OpenshiftApiserverRbac20_clusterroleSystemImagePullerYaml, map[string]*bintree{}},
				"20_clusterrole-system-image-pusher.yaml":                                                        {v3110OpenshiftApiserverRbac20_clusterroleSystemImagePusherYaml, map[string]*bintree{}},
				"20_clusterrole-system-image-signer.yaml":                                                        {v3110OpenshiftApiserverRbac20_clusterroleSystemImageSignerYaml, map[string]*bintree{}},
				"20_clusterrole-system-master.yaml":                                                              {v3110OpenshiftApiserverRbac20_clusterroleSystemMasterYaml, map[string]*bintree{}},
				"20_clusterrole-system-node-admin.yaml":                                                          {v3110OpenshiftApiserverRbac20_clusterroleSystemNodeAdminYaml, map[string]*bintree{}},
				"20_clusterrole-system-node-reader.yaml":                                                         {v3110OpenshiftApiserverRbac20_clusterroleSystemNodeReaderYaml, map[string]*bintree{}},
				"20_clusterrole-system-oauth-token-deleter.yaml":                                                 {v3110OpenshiftApiserverRbac20_clusterroleSystemOauthTokenDeleterYaml, map[string]*bintree{}},
				"20_clusterrole-system-openshift-controller-cluster-quota-reconciliation-controller.yaml":        {v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerClusterQuotaReconciliationControllerYaml, map[string]*bintree{}},
				"20_clusterrole-system-openshift-controller-horizontal-pod-autoscaler.yaml":                      {v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerHorizontalPodAutoscalerYaml, map[string]*bintree{}},
				"20_clusterrole-system-openshift-controller-pv-recycler-controller.yaml":                         {v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerPvRecyclerControllerYaml, map[string]*bintree{}},
				"20_clusterrole-system-openshift-controller-resourcequota-controller.yaml":                       {v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerResourcequotaControllerYaml, map[string]*bintree{}},
				"20_clusterrole-system-openshift-controller-service-serving-cert-controller.yaml":                {v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerServiceServingCertControllerYaml, map[string]*bintree{}},
				"20_clusterrole-system-openshift-controller-template-service-broker.yaml":                        {v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftControllerTemplateServiceBrokerYaml, map[string]*bintree{}},
				"20_clusterrole-system-openshift-discovery.yaml":                                                 {v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftDiscoveryYaml, map[string]*bintree{}},
				"20_clusterrole-system-openshift-templateservicebroker-client.yaml":                              {v3110OpenshiftApiserverRbac20_clusterroleSystemOpenshiftTemplateservicebrokerClientYaml, map[string]*bintree{}},
				"20_clusterrole-system-router.yaml":                                                              {v3110OpenshiftApiserverRbac20_clusterroleSystemRouterYaml, map[string]*bintree{}},
				"20_clusterrole-system-scope-impersonation.yaml":                                                 {v3110OpenshiftApiserverRbac20_clusterroleSystemScopeImpersonationYaml, map[string]*bintree{}},
				"20_clusterrole-system-sdn-manager.yaml":                                                         {v3110OpenshiftApiserverRbac20_clusterroleSystemSdnManagerYaml, map[string]*bintree{}},
				"20_clusterrole-system-sdn-reader.yaml":                                                          {v3110OpenshiftApiserverRbac20_clusterroleSystemSdnReaderYaml, map[string]*bintree{}},
				"20_clusterrole-system-webhook.yaml":                                                             {v3110OpenshiftApiserverRbac20_clusterroleSystemWebhookYaml, map[string]*bintree{}},
				"30_clusterrolebinding-basic-users.yaml":                                                         {v3110OpenshiftApiserverRbac30_clusterrolebindingBasicUsersYaml, map[string]*bintree{}},
				"30_clusterrolebinding-cluster-admins.yaml":                                                      {v3110OpenshiftApiserverRbac30_clusterrolebindingClusterAdminsYaml, map[string]*bintree{}},
				"30_clusterrolebinding-cluster-readers.yaml":                                                     {v3110OpenshiftApiserverRbac30_clusterrolebindingClusterReadersYaml, map[string]*bintree{}},
				"30_clusterrolebinding-cluster-status-binding.yaml":                                              {v3110OpenshiftApiserverRbac30_clusterrolebindingClusterStatusBindingYaml, map[string]*bintree{}},
				"30_clusterrolebinding-self-access-reviewers.yaml":                                               {v3110OpenshiftApiserverRbac30_clusterrolebindingSelfAccessReviewersYaml, map[string]*bintree{}},
				"30_clusterrolebinding-self-provisioners.yaml":                                                   {v3110OpenshiftApiserverRbac30_clusterrolebindingSelfProvisionersYaml, map[string]*bintree{}},
				"30_clusterrolebinding-system-build-strategy-docker-binding.yaml":                                {v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategyDockerBindingYaml, map[string]*bintree{}},
				"30_clusterrolebinding-system-build-strategy-jenkinspipeline-binding.yaml":                       {v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategyJenkinspipelineBindingYaml, map[string]*bintree{}},
				"30_clusterrolebinding-system-build-strategy-source-binding.yaml":                                {v3110OpenshiftApiserverRbac30_clusterrolebindingSystemBuildStrategySourceBindingYaml, map[string]*bintree{}},
				"30_clusterrolebinding-system-controller-horizontal-pod-autoscaler.yaml":                         {v3110OpenshiftApiserverRbac30_clusterrolebindingSystemControllerHorizontalPodAutoscalerYaml, map[string]*bintree{}},
				"30_clusterrolebinding-system-masters.yaml":                                                      {v3110OpenshiftApiserverRbac30_clusterrolebindingSystemMastersYaml, map[string]*bintree{}},
				"30_clusterrolebinding-system-node-admins.yaml":                                                  {v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeAdminsYaml, map[string]*bintree{}},
				"30_clusterrolebinding-system-node-bootstrapper.yaml":                                            {v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeBootstrapperYaml, map[string]*bintree{}},
				"30_clusterrolebinding-system-node-proxiers.yaml":                                                {v3110OpenshiftApiserverRbac30_clusterrolebindingSystemNodeProxiersYaml, map[string]*bintree{}},
				"30_clusterrolebinding-system-oauth-token-deleters.yaml":                                         {v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOauthTokenDeletersYaml, map[string]*bintree{}},
				"30_clusterrolebinding-system-openshift-controller-cluster-quota-reconciliation-controller.yaml": {v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerClusterQuotaReconciliationControllerYaml, map[string]*bintree{}},
				"30_clusterrolebinding-system-openshift-controller-horizontal-pod-autoscaler.yaml":               {v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerHorizontalPodAutoscalerYaml, map[string]*bintree{}},
				"30_clusterrolebinding-system-openshift-controller-pv-recycler-controller.yaml":                  {v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerPvRecyclerControllerYaml, map[string]*bintree{}},
				"30_clusterrolebinding-system-openshift-controller-resourcequota-controller.yaml":                {v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerResourcequotaControllerYaml, map[string]*bintree{}},
				"30_clusterrolebinding-system-openshift-controller-service-serving-cert-controller.yaml":         {v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerServiceServingCertControllerYaml, map[string]*bintree{}},
				"30_clusterrolebinding-system-openshift-controller-template-service-broker.yaml":                 {v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftControllerTemplateServiceBrokerYaml, map[string]*bintree{}},
				"30_clusterrolebinding-system-openshift-discovery.yaml":                                          {v3110OpenshiftApiserverRbac30_clusterrolebindingSystemOpenshiftDiscoveryYaml, map[string]*bintree{}},
				"30_clusterrolebinding-system-scope-impersonation.yaml":                                          {v3110OpenshiftApiserverRbac30_clusterrolebindingSystemScopeImpersonationYaml, map[string]*bintree{}},
				"30_clusterrolebinding-system-sdn-readers.yaml":                                                  {v3110OpenshiftApiserverRbac30_clusterrolebindingSystemSdnReadersYaml, map[string]*bintree{}},
				"30_clusterrolebinding-system-webhooks.yaml":                                                     {v3110OpenshiftApiserverRbac30_clusterrolebindingSystemWebhooksYaml, map[string]*bintree{}},
				"50_aggregate-to-admin.yaml":                                                                     {v3110OpenshiftApiserverRbac50_aggregateToAdminYaml, map[string]*bintree{}},
				"50_aggregate-to-basic-user.yaml":                                                                {v3110OpenshiftApiserverRbac50_aggregateToBasicUserYaml, map[string]*bintree{}},
				"50_aggregate-to-cluster-reader.yaml":                                                            {v3110OpenshiftApiserverRbac50_aggregateToClusterReaderYaml, map[string]*bintree{}},
				"50_aggregate-to-edit.yaml":                                                                      {v3110OpenshiftApiserverRbac50_aggregateToEditYaml, map[string]*bintree{}},
				"50_aggregate-to-storage-admin.yaml":                                                             {v3110OpenshiftApiserverRbac50_aggregateToStorageAdminYaml, map[string]*bintree{}},
				"50_aggregate-to-view.yaml":                                                                      {v3110OpenshiftApiserverRbac50_aggregateToViewYaml, map[string]*bintree{}},
			}},
			"rbac-openshift-controller-manager": {nil, map[string]*bintree{
				"clusterrole-system-openshift-controller-build-config-change-controller.yaml":                      {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerBuildConfigChangeControllerYaml, map[string]*bintree{}},
				"clusterrole-system-openshift-controller-build-controller.yaml":                                    {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerBuildControllerYaml, map[string]*bintree{}},
				"clusterrole-system-openshift-controller-default-rolebindings-controller.yaml":                     {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDefaultRolebindingsControllerYaml, map[string]*bintree{}},
				"clusterrole-system-openshift-controller-deployer-controller.yaml":                                 {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDeployerControllerYaml, map[string]*bintree{}},
				"clusterrole-system-openshift-controller-deploymentconfig-controller.yaml":                         {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerDeploymentconfigControllerYaml, map[string]*bintree{}},
				"clusterrole-system-openshift-controller-image-import-controller.yaml":                             {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerImageImportControllerYaml, map[string]*bintree{}},
				"clusterrole-system-openshift-controller-image-trigger-controller.yaml":                            {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerImageTriggerControllerYaml, map[string]*bintree{}},
				"clusterrole-system-openshift-controller-origin-namespace-controller.yaml":                         {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerOriginNamespaceControllerYaml, map[string]*bintree{}},
				"clusterrole-system-openshift-controller-service-ingress-ip-controller.yaml":                       {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceIngressIpControllerYaml, map[string]*bintree{}},
				"clusterrole-system-openshift-controller-serviceaccount-controller.yaml":                           {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceaccountControllerYaml, map[string]*bintree{}},
				"clusterrole-system-openshift-controller-serviceaccount-pull-secrets-controller.yaml":              {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerServiceaccountPullSecretsControllerYaml, map[string]*bintree{}},
				"clusterrole-system-openshift-controller-template-instance-controller.yaml":                        {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerTemplateInstanceControllerYaml, map[string]*bintree{}},
				"clusterrole-system-openshift-controller-template-instance-finalizer-controller.yaml":              {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerTemplateInstanceFinalizerControllerYaml, map[string]*bintree{}},
				"clusterrole-system-openshift-controller-unidling-controller.yaml":                                 {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterroleSystemOpenshiftControllerUnidlingControllerYaml, map[string]*bintree{}},
				"clusterrolebinding-system-deployer.yaml":                                                          {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemDeployerYaml, map[string]*bintree{}},
				"clusterrolebinding-system-image-builder.yaml":                                                     {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemImageBuilderYaml, map[string]*bintree{}},
				"clusterrolebinding-system-image-puller.yaml":                                                      {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemImagePullerYaml, map[string]*bintree{}},
				"clusterrolebinding-system-openshift-controller-build-config-change-controller.yaml":               {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerBuildConfigChangeControllerYaml, map[string]*bintree{}},
				"clusterrolebinding-system-openshift-controller-build-controller.yaml":                             {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerBuildControllerYaml, map[string]*bintree{}},
				"clusterrolebinding-system-openshift-controller-default-rolebindings-controller.yaml":              {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDefaultRolebindingsControllerYaml, map[string]*bintree{}},
				"clusterrolebinding-system-openshift-controller-deployer-controller.yaml":                          {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDeployerControllerYaml, map[string]*bintree{}},
				"clusterrolebinding-system-openshift-controller-deploymentconfig-controller.yaml":                  {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerDeploymentconfigControllerYaml, map[string]*bintree{}},
				"clusterrolebinding-system-openshift-controller-image-import-controller.yaml":                      {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerImageImportControllerYaml, map[string]*bintree{}},
				"clusterrolebinding-system-openshift-controller-image-trigger-controller.yaml":                     {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerImageTriggerControllerYaml, map[string]*bintree{}},
				"clusterrolebinding-system-openshift-controller-origin-namespace-controller.yaml":                  {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerOriginNamespaceControllerYaml, map[string]*bintree{}},
				"clusterrolebinding-system-openshift-controller-service-ingress-ip-controller.yaml":                {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceIngressIpControllerYaml, map[string]*bintree{}},
				"clusterrolebinding-system-openshift-controller-serviceaccount-controller.yaml":                    {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceaccountControllerYaml, map[string]*bintree{}},
				"clusterrolebinding-system-openshift-controller-serviceaccount-pull-secrets-controller.yaml":       {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerServiceaccountPullSecretsControllerYaml, map[string]*bintree{}},
				"clusterrolebinding-system-openshift-controller-template-instance-controller-admin.yaml":           {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceControllerAdminYaml, map[string]*bintree{}},
				"clusterrolebinding-system-openshift-controller-template-instance-controller.yaml":                 {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceControllerYaml, map[string]*bintree{}},
				"clusterrolebinding-system-openshift-controller-template-instance-finalizer-controller-admin.yaml": {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceFinalizerControllerAdminYaml, map[string]*bintree{}},
				"clusterrolebinding-system-openshift-controller-template-instance-finalizer-controller.yaml":       {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerTemplateInstanceFinalizerControllerYaml, map[string]*bintree{}},
				"clusterrolebinding-system-openshift-controller-unidling-controller.yaml":                          {v3110OpenshiftApiserverRbacOpenshiftControllerManagerClusterrolebindingSystemOpenshiftControllerUnidlingControllerYaml, map[string]*bintree{}},
			}},
			"sa.yaml":            {v3110OpenshiftApiserverSaYaml, map[string]*bintree{}},
			"svc.yaml":           {v3110OpenshiftApiserverSvcYaml, map[string]*bintree{}},
			"trusted_ca_cm.yaml": {v3110OpenshiftApiserverTrusted_ca_cmYaml, map[string]*bintree{}},
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
