package uspnj

import (
	"testing"

	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/sections"
	"github.com/stretchr/testify/assert"
)

type uspnjTestData struct {
	description string
	gppString   string
	expected    USPNJ
}

func TestUSPNJ(t *testing.T) {
	testData := []uspnjTestData{
		{
			description: "should populate USPNJ segments correctly",
			gppString:   "ZSFgmABlfA",
			/*
				011001 01 00 10 00 01 01100000100110000000 0000011001 01 01 11 11
			*/
			expected: USPNJ{
				CoreSegment: USPNJCoreSegment{
					Version:                         25,
					ProcessingNotice:                1,
					SaleOptOutNotice:                0,
					TargetedAdvertisingOptOutNotice: 2,
					SaleOptOut:                      0,
					TargetedAdvertisingOptOut:       1,
					SensitiveDataProcessing: []byte{
						1, 2, 0, 0, 2, 1, 2, 0, 0, 0,
					},
					KnownChildSensitiveDataConsents: []byte{0, 0, 1, 2, 1},
					AdditionalDataProcessingConsent: 1,
					MspaCoveredTransaction:          1,
					MspaOptOutOptionMode:            3,
					MspaServiceProviderMode:         3,
				},
				GPCSegment: sections.CommonUSGPCSegment{
					SubsectionType: 1,
					Gpc:            false,
				},
				SectionID: constants.SectionUSPNJ,
				Value:     "ZSFgmABlfA",
			},
		},
	}

	for _, test := range testData {
		result, err := NewUSPNJ(test.gppString)
		encodedString := string(test.expected.Encode(false))
		assert.Nil(t, err)
		assert.Equal(t, test.expected, result)
		assert.Equal(t, constants.SectionUSPNJ, result.GetID())
		assert.Equal(t, test.gppString, result.GetValue())
		assert.Equal(t, test.gppString, encodedString)
	}
}
