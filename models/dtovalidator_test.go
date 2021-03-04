package models_test

import (
	"errors"
	"github.com/magneticio/forklift/models"
	"reflect"
	"testing"
)

var (
	applicationID   = uint64ToPointer(1)
	serviceID       = uint64ToPointer(2)
	defaultPolicyID = uint64ToPointer(33)
	patchPolicyID   = uint64ToPointer(44)
	minorPolicyID   = uint64ToPointer(55)
	majorPolicyID   = uint64ToPointer(66)
	k8sNamespace    = "test-namespace"
	k8sLabels       = map[string]string{"test-key": "test-value"}
	versionSelector = "version"
)

func TestServiceConfigValidation(t *testing.T) {
	type args struct {
		serviceConfig models.ServiceConfig
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "valid service config",
			args: args{validServiceConfig().build()},
			want: nil,
		},
		{
			name: "service config without application id",
			args: args{validServiceConfig().withApplicationID(nil).build()},
			want: errors.New("application_id field is required"),
		},
		{
			name: "service config without service id",
			args: args{validServiceConfig().withServiceID(nil).build()},
			want: errors.New("service_id field is required"),
		},
		{
			name: "service config with empty Kubernetes namespace",
			args: args{validServiceConfig().withK8SNamespace("").build()},
			want: errors.New("k8s_namespace field is required"),
		},
		{
			name: "service config with no Kubernetes labels",
			args: args{validServiceConfig().withK8SLabels(make(map[string]string)).build()},
			want: errors.New("k8s_labels field must contain at least 1 element(s)"),
		},
		{
			name: "service config with empty version selector",
			args: args{validServiceConfig().withVersionSelector("").build()},
			want: errors.New("version_selector field is required"),
		},
		{
			name: "service config with ingress rule with invalid domain",
			args: args{
				validServiceConfig().
					withIngressRules(
						[]*models.ServiceConfigIngressRule{
							validIngressRule().withDomain("ab").build(),
						}).
					build()},
			want: errors.New("ingress_rules[0].domain field must be at least 4 character(s)"),
		},
		{
			name: "service config with ingress rule with empty path",
			args: args{
				validServiceConfig().
					withIngressRules(
						[]*models.ServiceConfigIngressRule{
							validIngressRule().withPath("").build(),
						}).
					build()},
			want: errors.New("ingress_rules[0].path field is required"),
		},
		{
			name: "service config with ingress rule with no port",
			args: args{
				validServiceConfig().
					withIngressRules(
						[]*models.ServiceConfigIngressRule{
							validIngressRule().withPort(nil).build(),
						}).
					build()},
			want: errors.New("ingress_rules[0].port field is required"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := models.NewValidateDTO()
			if got := validate(tt.args.serviceConfig); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DTO validator = %v, want %v", got, tt.want)
			}
		})
	}
}

type serviceConfigBuilder struct {
	serviceConfig models.ServiceConfig
}

func (builder *serviceConfigBuilder) build() models.ServiceConfig {
	return builder.serviceConfig
}

func (builder *serviceConfigBuilder) withApplicationID(applicationID *uint64) *serviceConfigBuilder {
	builder.serviceConfig.ApplicationID = applicationID
	return builder
}

func (builder *serviceConfigBuilder) withServiceID(serviceID *uint64) *serviceConfigBuilder {
	builder.serviceConfig.ServiceID = serviceID
	return builder
}

func (builder *serviceConfigBuilder) withK8SNamespace(k8sNamespace string) *serviceConfigBuilder {
	builder.serviceConfig.K8SNamespace = k8sNamespace
	return builder
}

func (builder *serviceConfigBuilder) withK8SLabels(k8sLabels map[string]string) *serviceConfigBuilder {
	builder.serviceConfig.K8sLabels = k8sLabels
	return builder
}

func (builder *serviceConfigBuilder) withVersionSelector(versionSelector string) *serviceConfigBuilder {
	builder.serviceConfig.VersionSelector = versionSelector
	return builder
}

func (builder *serviceConfigBuilder) withIngressRules(ingressRules []*models.ServiceConfigIngressRule) *serviceConfigBuilder {
	builder.serviceConfig.IngressRules = ingressRules
	return builder
}

func validServiceConfig() *serviceConfigBuilder {
	serviceConfig := models.ServiceConfig{
		ApplicationID:   applicationID,
		ServiceID:       serviceID,
		K8SNamespace:    k8sNamespace,
		K8sLabels:       k8sLabels,
		VersionSelector: versionSelector,
		DefaultPolicyID: defaultPolicyID,
		PatchPolicyID:   patchPolicyID,
		MinorPolicyID:   minorPolicyID,
		MajorPolicyID:   majorPolicyID,
		IngressRules: []*models.ServiceConfigIngressRule{
			validIngressRule().build(),
		},
	}

	return &serviceConfigBuilder{
		serviceConfig: serviceConfig,
	}
}

type ingressRuleBuilder struct {
	ingressRule *models.ServiceConfigIngressRule
}

func (builder *ingressRuleBuilder) build() *models.ServiceConfigIngressRule {
	return builder.ingressRule
}

func (builder *ingressRuleBuilder) withDomain(domain string) *ingressRuleBuilder {
	builder.ingressRule.Domain = domain
	return builder
}

func (builder *ingressRuleBuilder) withPath(path string) *ingressRuleBuilder {
	builder.ingressRule.Path = path
	return builder
}

func (builder *ingressRuleBuilder) withPort(port *int64) *ingressRuleBuilder {
	builder.ingressRule.Port = port
	return builder
}

func validIngressRule() *ingressRuleBuilder {
	ingressRule := &models.ServiceConfigIngressRule{
		Domain:        "test.local",
		TLSSecretName: "secret",
		Path:          "/api",
		Port:          int64ToPointer(8080),
	}

	return &ingressRuleBuilder{
		ingressRule: ingressRule,
	}
}

func uint64ToPointer(number uint64) *uint64 {
	return &number
}

func int64ToPointer(number int64) *int64 {
	return &number
}
