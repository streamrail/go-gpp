package tcfeu2

import (
	"testing"
	"time"

	consentconstants "github.com/prebid/go-gdpr/consentconstants/tcf2"
	tcf2 "github.com/prebid/go-gdpr/vendorconsent/tcf2"
	"github.com/streamrail/go-gpp/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tcfeu2TestData struct {
	description string
	gppString   string
	expected    func(t *testing.T, s TCFEU2)
}

func TestTCFEU2(t *testing.T) {
	testData := []tcfeu2TestData{
		{
			description: "should populate TCFEU2 segments correctly",
			gppString:   "CQWf88AQWf88APoABAENB4EIAIAAAIAAAAAAE0wAQE0gTTABATSAAAAA.QAAA.IAAA", // Generated with https://iabgpp.com
			expected: func(t *testing.T, s TCFEU2) {
				assert.Equal(t, constants.SectionTCFEU2, s.GetID())
				assert.Equal(t, uint16(1000), s.CoreSegment.CmpID())
				assert.Equal(t, uint16(1), s.CoreSegment.CmpVersion())
				assert.Equal(t, "EN", s.CoreSegment.ConsentLanguage())
				assert.Equal(t, uint8(0), s.CoreSegment.ConsentScreen())
				assert.Equal(t, time.Date(2025, 8, 21, 2, 0, 0, 0, time.Local), s.CoreSegment.Created())
				assert.Equal(t, time.Date(2025, 8, 21, 2, 0, 0, 0, time.Local), s.CoreSegment.LastUpdated())
				assert.Equal(t, uint16(0x269), s.CoreSegment.MaxVendorID())
				assert.Equal(t, uint8(4), s.CoreSegment.TCFPolicyVersion())
				assert.True(t, s.CoreSegment.PurposeAllowed(consentconstants.InfoStorageAccess))
				assert.False(t, s.CoreSegment.PurposeAllowed(consentconstants.ContentPerformance))
				assert.True(t, s.CoreSegment.VendorConsent(617))
				assert.False(t, s.CoreSegment.VendorConsent(658))
				assert.Equal(t, uint16(120), s.CoreSegment.VendorListVersion())
				assert.Equal(t, uint8(2), s.CoreSegment.Version())
				assert.Equal(t, "CQWf88AQWf88APoABAENB4EIAIAAAIAAAAAAE0wAQE0gTTABATSAAAAA.QAAA.IAAA", s.GetValue())
				assert.Equal(t, "CQWf88AQWf88APoABAENB4EIAIAAAIAAAAAAE0wAQE0gTTABATSAAAAA.QAAA.IAAA", string(s.Encode(false)))

				consentMeta, ok := s.CoreSegment.(tcf2.ConsentMetadata)
				if !ok {
					assert.FailNow(t, "parsed consent type assert to tcf2.ConsentMetadata failed")
				}
				assert.True(t, consentMeta.PurposeLITransparency(consentconstants.InfoStorageAccess))
				assert.False(t, consentMeta.PurposeLITransparency(consentconstants.ContentPerformance))
				assert.True(t, consentMeta.SpecialFeatureOptIn(1))
				assert.False(t, consentMeta.SpecialFeatureOptIn(2))
				assert.True(t, consentMeta.VendorLegitInterest(617))
				assert.False(t, consentMeta.VendorLegitInterest(658))
			},
		},
	}

	for _, test := range testData {
		result, err := NewTCFEU2(test.gppString)
		require.Nil(t, err)
		test.expected(t, result)
	}
}
