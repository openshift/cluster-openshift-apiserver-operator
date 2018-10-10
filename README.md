# cluster-openshift-apiserver-operator
The openshift-apiserver operator installs and maintains the openshift-apiserver on a cluster

## Recommended development flow
1. Use openshift/installer to install a cluster
2. `oc delete ns openshift-cluster-version -wait=false`
3. Until https://github.com/openshift/installer/pull/435 merges, ssh to the bootstrap node, kill the kubelet (`systemctl stop kubelet`, kill the cluster-version-operator process (`ps ax | grep cluster-version-operator` and `kill <pid>`).  The `oc delete` the bootstrap node.
4. `make images`
5. `docker tag openshift/origin-cluster-openshift-apiserver-operator:latest <yourdockerhubid>/origin-cluster-openshift-apiserver-operator:latest`
6. `docker push <yourdockerhubid>/origin-cluster-openshift-apiserver-operator:latest`
7. `oc edit -n openshift-core-operators deployment.apps/openshift-cluster-openshift-apiserver-operator` update the image to `<yourdockerhubid>/origin-cluster-openshift-apiserver-operator:latest` and update the pull policy to `Always`
8. `oc delete pod -n openshift-core-operators --all` will cause the pod to be recreated and the image to be pulled.
