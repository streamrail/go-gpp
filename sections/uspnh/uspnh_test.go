package uspnh

import (
	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/sections"
	"github.com/stretchr/testify/assert"
	"testing"
)

type uspnhTestData struct {
	description string
	gppString   string
	expected    USPNH
}

func TestUSPNH(t *testing.T) {
	testData := []uspnhTestData{
		{
			description: "should populate USPNH segments correctly",
			gppString:   "ZSFgmAF8.YA",
			/*
				011001 01 00 10 00 01 0110000010011000 000000 01 01 11 11 01
			*/
			expected: USPNH{
				CoreSegment: USPNHCoreSegment{
					Version:                         25,
					ProcessingNotice:                1,
					SaleOptOutNotice:                0,
					TargetedAdvertisingOptOutNotice: 2,
					SaleOptOut:                      0,
					TargetedAdvertisingOptOut:       1,
					SensitiveDataProcessing: []byte{
						1, 2, 0, 0, 2, 1, 2, 0,
					},
					KnownChildSensitiveDataConsents: []byte{
						0, 0, 0,
					},
					AdditionalDataProcessingConsent: 1,
					MspaCoveredTransaction:          1,
					MspaOptOutOptionMode:            3,
					MspaServiceProviderMode:         3,
				},
				GPCSegment: sections.CommonUSGPCSegment{
					SubsectionType: 1,
					Gpc:            true,
				},
				SectionID: constants.SectionUSPNH,
				Value:     "ZSFgmAF8.YA",
			},
		},
	}

	for _, test := range testData {
		result, err := NewUSPNH(test.gppString)
		encodedString := string(test.expected.Encode(true))
		assert.Nil(t, err)
		assert.Equal(t, test.expected, result)
		assert.Equal(t, constants.SectionUSPNH, result.GetID())
		assert.Equal(t, test.gppString, result.GetValue())
		assert.Equal(t, test.gppString, encodedString)
	}
}
