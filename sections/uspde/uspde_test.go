package uspde

import (
	"testing"

	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/sections"
	"github.com/stretchr/testify/assert"
)

type uspdeTestData struct {
	description string
	gppString   string
	expected    USPDE
}

func TestUSPDE(t *testing.T) {
	testData := []uspdeTestData{
		{
			description: "should populate USPDE segments correctly",
			gppString:   "4SFgmAGV8A",
			/*
				111000 01 00 10 00 01 011000001001100000 0000011001 01 01 11 11
			*/
			expected: USPDE{
				CoreSegment: USPDECoreSegment{
					Version:                         56,
					ProcessingNotice:                1,
					SaleOptOutNotice:                0,
					TargetedAdvertisingOptOutNotice: 2,
					SaleOptOut:                      0,
					TargetedAdvertisingOptOut:       1,
					SensitiveDataProcessing: []byte{
						1, 2, 0, 0, 2, 1, 2, 0, 0,
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
				SectionID: constants.SectionUSPDE,
				Value:     "4SFgmAGV8A",
			},
		},
	}

	for _, test := range testData {
		result, err := NewUSPDE(test.gppString)
		encodedString := string(test.expected.Encode(false))
		assert.Nil(t, err)
		assert.Equal(t, test.expected, result)
		assert.Equal(t, constants.SectionUSPDE, result.GetID())
		assert.Equal(t, test.gppString, result.GetValue())
		assert.Equal(t, test.gppString, encodedString)
	}
}
