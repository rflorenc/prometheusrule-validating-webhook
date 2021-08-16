# PrometheusRule Admission Webhook

Tested in Openshift v4.5. 

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

```bash
oc get ns -l app.kubernetes.io/created-by=webhook-managed.project.example.com
NAME                 STATUS   AGE
discovery-endpoint   Active   16d
example-dev-dev            Active   62m
example                Active   4d1h
example-2              Active   28h
example-with-label     Active   130m
```

1 - Allow admission of rule with all the required fields specified.

```bash
oc create -f hack/testlocal/mock-prometheus-rules-good.yaml -n example
prometheusrule.monitoring.coreos.com/mock-example-prometheus-rules created
```

2 - Deny admission of rule with some of the required fields missing.

```bash
oc create -f hack/testlocal/mock-prometheus-rules-bad.yaml -n example
Error from server (Missing one or more of minimum required labels. severity: false, example_response_code: false, example_alerting_email: true): error when creating "hack/testlocal/mock-prometheus-rules-bad.yaml": admission webhook "prometheusrule-validating-webhook.example.com" denied the request: Missing one or more of minimum required labels. severity: false, example_response_code: false, example_alerting_email: true
```

3 - Allow admission of rule with missing field `example_response_code` in namespace which is not properly labeled by the webhook-managed.

```yaml
oc get ns test-abc -o yaml
apiVersion: v1
kind: Namespace
metadata:
  ...
  labels:
    abc: def
```

```bash
oc create -f hack/testlocal/mock-prometheus-rules-other-namespace.yaml -n test-abc
prometheusrule.monitoring.coreos.com/mock-example-prometheus-rules-bad created
```

3.1 - Properly label `test-abc` namespace.

```bash
oc label ns test-abc app.kubernetes.io/created-by=webhook-managed.project.example.com
```

3.2 - Deny admission of rule with missing field `example_response_code` in `test-abc` namespace which is now properly labeled.

```bash
oc create -f hack/testlocal/mock-prometheus-rules-other-namespace.yaml -n test-abc
Error from server (Missing one or more of minimum required labels. severity: true, example_response_code: false, example_alerting_email: true): error when creating "hack/testlocal/mock-prometheus-rules-other-namespace.yaml": admission webhook "prometheusrule-validating-webhook.example.com" denied the request: Missing one or more of minimum required labels. severity: true, example_response_code: false, example_alerting_email: true
```

## Installing in-cluster

`helm-wrapper -u`

## Building

`make webhook-linux`
