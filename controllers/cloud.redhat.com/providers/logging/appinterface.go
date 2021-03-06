package logging

import (
	crd "cloud.redhat.com/clowder/v2/apis/cloud.redhat.com/v1alpha1"
	"cloud.redhat.com/clowder/v2/controllers/cloud.redhat.com/config"
	"cloud.redhat.com/clowder/v2/controllers/cloud.redhat.com/errors"
	"cloud.redhat.com/clowder/v2/controllers/cloud.redhat.com/providers"
	p "cloud.redhat.com/clowder/v2/controllers/cloud.redhat.com/providers"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type AppInterfaceLoggingProvider struct {
	p.Provider
	Config config.LoggingConfig
}

func NewAppInterfaceLogging(p *p.Provider) (providers.ClowderProvider, error) {
	provider := AppInterfaceLoggingProvider{Provider: *p}

	return &provider, nil
}

func (a *AppInterfaceLoggingProvider) Provide(app *crd.ClowdApp, c *config.AppConfig) error {
	c.Logging = config.LoggingConfig{}
	return setCloudwatchSecret(app.Namespace, &a.Provider, &c.Logging)
}

func setCloudwatchSecret(ns string, p *p.Provider, c *config.LoggingConfig) error {

	name := types.NamespacedName{
		Name:      "cloudwatch",
		Namespace: ns,
	}

	secret := core.Secret{}
	err := p.Client.Get(p.Ctx, name, &secret)

	if err != nil {
		return errors.Wrap("Failed to fetch cloudwatch secret", err)
	}

	c.Cloudwatch = &config.CloudWatchConfig{
		AccessKeyId:     string(secret.Data["aws_access_key_id"]),
		SecretAccessKey: string(secret.Data["aws_secret_access_key"]),
		Region:          string(secret.Data["aws_region"]),
		LogGroup:        string(secret.Data["log_group_name"]),
	}

	return nil
}
