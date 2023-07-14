package federation

import (
	"encoding/json"

	"github.com/xeronith/diamante/contracts/federation"
)

type webfinger struct {
	Aliases []string `json:"aliases"`
	Links   []struct {
		Href     *string `json:"href,omitempty"`
		Rel      string  `json:"rel"`
		Type     *string `json:"type,omitempty"`
		Template *string `json:"template,omitempty"`
	} `json:"links"`
	Subject string `json:"subject"`
}

func NewWebfinger() federation.IWebfinger {
	return &webfinger{}
}

func (webfinger *webfinger) Self() string {
	self := ""
	for _, link := range webfinger.Links {
		if link.Rel == "self" && link.Type != nil &&
			*link.Type == "application/activity+json" {
			self = *link.Href
			break
		}
	}

	return self
}

func (webfinger *webfinger) Unmarshal(data []byte) error {
	return json.Unmarshal(data, webfinger)
}

func (webfinger *webfinger) Marshal() ([]byte, error) {
	return json.Marshal(webfinger)
}
