package uspmn

import (
	"testing"

	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/sections"
	"github.com/stretchr/testify/assert"
)

type uspmnTestData struct {
	description string
	gppString   string
	expected    USPMN
}

func TestUSPMN(t *testing.T) {
	testData := []uspmnTestData{
		{
			description: "should populate USPMN segments correctly",
			gppString:   "bBIWCYlA.YA",
			/*
				011011 00 00 01 00 10 0110000010011000 10 00 10  01 01
			*/
			expected: USPMN{
				CoreSegment: USPMNCoreSegment{
					Version:                         27,
					ProcessingNotice:                0,
					SaleOptOutNotice:                0,
					TargetedAdvertisingOptOutNotice: 1,
					SaleOptOut:                      0,
					TargetedAdvertisingOptOut:       2,
					SensitiveDataProcessing: []byte{
						0, 1, 1, 2, 0, 0, 2, 1,
					},
					KnownChildSensitiveDataConsents: 2,
					AdditionalDataProcessingConsent: 0,
					MspaCoveredTransaction:          2,
					MspaOptOutOptionMode:            1,
					MspaServiceProviderMode:         1,
				},
				GPCSegment: sections.CommonUSGPCSegment{
					SubsectionType: 1,
					Gpc:            true,
				},
				SectionID: constants.SectionUSPMN,
				Value:     "bBIWCYlA.YA",
			},
		},
	}

	for _, test := range testData {
		result, err := NewUSPMN(test.gppString)
		encodedString := string(test.expected.Encode(true))
		assert.Nil(t, err)
		assert.Equal(t, test.expected, result)
		assert.Equal(t, constants.SectionUSPMN, result.GetID())
		assert.Equal(t, test.gppString, result.GetValue())
		assert.Equal(t, test.gppString, encodedString)
	}
}
