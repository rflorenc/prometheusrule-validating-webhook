module github.com/rflorenc/prometheusrule-validating-webhook

go 1.15

require (
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.49.0
	k8s.io/api v0.21.2
	k8s.io/client-go v0.21.2
	sigs.k8s.io/controller-runtime v0.9.2
)
