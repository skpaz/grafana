# Grafana Alloy on Red Hat OpenShift

This was tested on [Red Hat OpenShift Local](https://developers.redhat.com/products/openshift-local/overview)
 fka Red Hat CodeReady Containers (CRC).

## Helm Chart

The [k8s-monitoring-helm](https://github.com/grafana/k8s-monitoring-helm)
 chart can be used with OpenShift.

Recommend the use of the configuration wizard in
 [the Kubernetes monitoring app](https://grafana.com/solutions/kubernetes/) to
 generate a Helm values file. An example values.yaml file is included, but
 is not necessarily up-to-date.

>[!NOTE]
> This Helm chart will take into account OpenShifts built-in observability stack
> and make use of the existing node_exporter and kube-state-metrics instances.

Go to **Infrastructure > Kubernetes > Configuration > [Tab] Cluster configuration**
 to access the wizard.

## Cluster Name

The default cluster name for OpenShift Local is `openshift`, which is used by the
 built-in node_expoter and kube-state-metrics instances. If you want to update
 it, you can do so via the
 [OpenShift Cluster Manager](https://console.redhat.com/openshift). Make sure
 the actual cluster name matches the cluster name used in your `values.yaml`
 file, or you'll have two values for `cluster`.

## Fixes

Any issues discovered and fixed as a result of lab work related to this write-up
 are listed below.

- https://github.com/grafana/k8s-monitoring-helm/pull/1373
- https://github.com/grafana/grafana-k8s-plugin/pull/2037
