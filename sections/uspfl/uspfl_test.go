package uspfl

import (
	"testing"

	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/sections"
	"github.com/stretchr/testify/assert"
)

type uspdeTestData struct {
	description string
	gppString   string
	expected    USPFL
}

func TestUSPDE(t *testing.T) {
	testData := []uspdeTestData{
		{
			description: "should populate USPFL segments correctly",
			gppString:   "BSFgmCRo",
			expected: USPFL{
				CoreSegment: USPFLCoreSegment{
					Version:                         1,
					ProcessingNotice:                1,
					SaleOptOutNotice:                0,
					TargetedAdvertisingOptOutNotice: 2,
					SaleOptOut:                      0,
					TargetedAdvertisingOptOut:       1,
					SensitiveDataProcessing: []byte{
						1, 2, 0, 0, 2, 1, 2, 0,
					},
					KnownChildSensitiveDataConsents: []byte{0, 2, 1},
					AdditionalDataProcessingConsent: 0,
					MspaCoveredTransaction:          1,
					MspaOptOutOptionMode:            2,
					MspaServiceProviderMode:         2,
				},
				GPCSegment: sections.CommonUSGPCSegment{
					SubsectionType: 1,
					Gpc:            false,
				},
				SectionID: constants.SectionUSPFL,
				Value:     "BSFgmCRo",
			},
		},
	}

	for _, test := range testData {
		result, err := NewUSPFL(test.gppString)
		encodedString := string(test.expected.Encode(false))
		assert.Nil(t, err)
		assert.Equal(t, test.expected, result)
		assert.Equal(t, constants.SectionUSPFL, result.GetID())
		assert.Equal(t, test.gppString, result.GetValue())
		assert.Equal(t, test.gppString, encodedString)
	}
}
