package main

type ContentMix []ContentConfig

type ContentConfig struct {
	Type     Provider
	Fallback *Provider
}

var (
	config1 = ContentConfig{
		Type:     Provider1,
		Fallback: &Provider2,
	}
	config2 = ContentConfig{
		Type:     Provider2,
		Fallback: &Provider3,
	}
	config3 = ContentConfig{
		Type:     Provider3,
		Fallback: &Provider1,
	}
	config4 = ContentConfig{
		Type:     Provider1,
		Fallback: nil,
	}

	DefaultConfig = []ContentConfig{
		config1, config1, config2, config3, config4, config1, config1, config2,
	}

	DefaultClients = map[Provider]Client{
		Provider1: SampleContentProvider{Source: Provider1},
		Provider2: SampleContentProvider{Source: Provider2},
		Provider3: SampleContentProvider{Source: Provider3},
	}
)
