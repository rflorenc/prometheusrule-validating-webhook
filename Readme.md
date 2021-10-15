# PrometheusRule Admission Webhook

Tested in Openshift v4.9 and Kubernetes v1.22.

## Functionality

This project implements a custom PrometheusRule Validating Webhook that filters and validates the presence of certain fields inside PrometheusRule Groups in custom created prometheus rules.
It only considers AdmissionReview objects for custom projects that have the `app.kubernetes.io/created-by: webhook-managed.project.example.com` label.
See `helm/prometheusrule-validating-webhook/validatingWebhookConfiguration.yaml`:

```yaml
  namespaceSelector:
    matchLabels:
        {{ include "prometheusrule-validating-webhook.webhookLabels" . }}
```


## Examples

1 - Create `example` namespace with the `app.kubernetes.io/created-by: webhook-managed.project.example.com` label.
```bash
cat <<EOF | oc create -f -
apiVersion: v1
kind: Namespace
metadata:
  labels:
    app.kubernetes.io/created-by: webhook-managed.project.example.com
  name: example
EOF

oc get ns -l app.kubernetes.io/created-by=webhook-managed.project.example.com
NAME                 STATUS   AGE
example                Active   6m58s
```

1.1 - Allow admission of rule with all the required fields specified.

```bash
oc create -f hack/testlocal/mock-prometheus-rules-good.yaml -n example
prometheusrule.monitoring.coreos.com/mock-example-prometheus-rules created
```

1.2 - Deny admission of rule with some of the required fields missing.

```bash
oc create -f hack/testlocal/mock-prometheus-rules-bad.yaml -n example
Error from server (Missing one or more of minimum required labels. severity: false, example_response_code: false, example_alerting_email: true): error when creating "hack/testlocal/mock-prometheus-rules-bad.yaml": admission webhook "prometheusrule-validating-webhook.example.com" denied the request: Missing one or more of minimum required labels. severity: false, example_response_code: false, example_alerting_email: true
```

2 - Create `test-abc` namespace without any labels.

```bash
cat <<EOF | oc create -f -
apiVersion: v1
kind: Namespace
metadata:
  name: test-abc
EOF
```

2.1 Allow admission of rule with missing field `example_response_code` in namespace which is not properly labeled by the webhook.

```bash
oc create -f hack/testlocal/mock-prometheus-rules-other-namespace.yaml -n test-abc
prometheusrule.monitoring.coreos.com/mock-example-prometheus-rules-bad created
```

2.2 - Properly label `test-abc` namespace.

```bash
oc label ns test-abc app.kubernetes.io/created-by=webhook-managed.project.example.com
```

2.3 - Deny admission of rule with missing field `example_response_code` in `test-abc` namespace which is now properly labeled.

```bash
oc create -f hack/testlocal/mock-prometheus-rules-other-namespace.yaml -n test-abc
Error from server (Missing one or more of minimum required labels. severity: true, example_response_code: false, example_alerting_email: true): error when creating "hack/testlocal/mock-prometheus-rules-other-namespace.yaml": admission webhook "prometheusrule-validating-webhook.example.com" denied the request: Missing one or more of minimum required labels. severity: true, example_response_code: false, example_alerting_email: true
```

## Installing in-cluster (requires helm client)

`./helm-wrapper -u`

## Building

`make webhook-linux`
