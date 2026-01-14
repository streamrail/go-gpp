package tcfeu2

import (
	"github.com/prebid/go-gdpr/api"
	vendorconsent "github.com/prebid/go-gdpr/vendorconsent/tcf2"
	"github.com/streamrail/go-gpp/constants"
)

type TCFEU2 struct {
	SectionID   constants.SectionID
	Value       string
	CoreSegment api.VendorConsents
}

func NewTCFEU2(encoded string) (TCFEU2, error) {
	consent, err := vendorconsent.ParseString(encoded)
	return TCFEU2{
		SectionID:   constants.SectionTCFEU2,
		Value:       encoded,
		CoreSegment: consent,
	}, err
}

func (tcfeu2 TCFEU2) Encode(gpcIncluded bool) []byte {
	// TODO there is currently no Go implementation for encoding TCFv2 strings.
	// As a workaround, we use the original value
	return []byte(tcfeu2.Value)
}

func (tcfeu2 TCFEU2) GetID() constants.SectionID {
	return tcfeu2.SectionID
}

func (tcfeu2 TCFEU2) GetValue() string {
	return tcfeu2.Value
}
