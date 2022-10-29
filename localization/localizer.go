package localization

import . "github.com/xeronith/diamante/contracts/localization"

type localizer struct {
	resource Resource
}

func (localizer *localizer) Register(resource Resource) {
	for key, value := range resource {
		localizer.resource[key] = value
	}
}

func (localizer *localizer) Get(key Key) Value {
	if value, exists := localizer.resource[key]; exists {
		return value
	}

	return ""
}

func NewLocalizer() ILocalizer {
	return &localizer{
		resource: make(Resource),
	}
}
