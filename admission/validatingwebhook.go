package admission

import (
	"context"
	"fmt"
	"net/http"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/validate-v1-prometheusrule,mutating=false,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=prometheusrule-validating-webhook.example.com

// PrometheusRuleValidator validates PrometheusRules
type PrometheusRuleValidator struct {
	Client  client.Client
	decoder *admission.Decoder
}

// PrometheusRuleValidator admits a PrometheusRule if a specific set of Rule labels exist
func (v *PrometheusRuleValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	prometheusRule := &monitoringv1.PrometheusRule{}

	err := v.decoder.Decode(req, prometheusRule)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	for _, group := range prometheusRule.Spec.Groups {
		for _, rule := range group.Rules {
			_, found_severity := rule.Labels["severity"]
			_, found_example_response_code := rule.Labels["example_response_code"]
			_, found_example_alerting_email := rule.Labels["example_alerting_email"]

			if !found_severity || !found_example_response_code || !found_example_alerting_email {
				return admission.Denied(fmt.Sprintf("Missing one or more of minimum required labels. severity: %v, example_response_code: %v, example_alerting_email: %v", found_severity, found_example_response_code, found_example_alerting_email))
			}
		}
	}
	return admission.Allowed("Rule admitted by PrometheusRule validating webhook.")
}

func (v *PrometheusRuleValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
