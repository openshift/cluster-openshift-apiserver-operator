# cluster-openshift-apiserver-operator
The openshift-apiserver operator installs and maintains the openshift-apiserver on a cluster

## Recommended development flow
1. Use openshift/installer to install a cluster
2. `oc delete ns openshift-cluster-version -wait=false`
3. `make images`
4. `docker tag openshift/origin-cluster-openshift-apiserver-operator:latest <yourdockerhubid>/origin-cluster-openshift-apiserver-operator:latest`
5. `docker push <yourdockerhubid>/origin-cluster-openshift-apiserver-operator:latest`
6. `oc edit -n openshift-cluster-openshift-apiserver-operator deployment.apps/openshift-cluster-openshift-apiserver-operator` update the image to `<yourdockerhubid>/origin-cluster-openshift-apiserver-operator:latest` and update the pull policy to `Always`
7. `oc delete pod -n openshift-cluster-openshift-apiserver-operator --all` will cause the pod to be recreated and the image to be pulled.
