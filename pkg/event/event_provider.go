package event

import "github.com/tangxusc/cqrs-sidecar/pkg/config"

type Provider interface {
	Support(Instance *config.Config) bool
	Consumer() Consumer
}

var providers = []Provider{PulsarProviderImpl}

func GetConsumerImpl() Consumer {
	for _, v := range providers {
		if v.Support(config.Instance) {
			return v.Consumer()
		}
	}
	return nil
}
