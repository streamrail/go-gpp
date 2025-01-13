package uspia

import (
	"testing"

	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/sections"
	"github.com/stretchr/testify/assert"
)

type uspneTestData struct {
	description string
	gppString   string
	expected    USPIA
}

func TestUSPIA(t *testing.T) {
	testData := []uspneTestData{
		{
			description: "should populate USPIA segments correctly",
			gppString:   "bSFgmJcA.YA",
			/*
				011011 01 00 10 00 01 0110000010011000 10 01 01 11 00 01
			*/
			expected: USPIA{
				CoreSegment: USPIACoreSegment{
					Version:                         27,
					ProcessingNotice:                1,
					SaleOptOutNotice:                0,
					TargetedAdvertisingOptOutNotice: 2,
					SaleOptOut:                      0,
					TargetedAdvertisingOptOut:       1,
					SensitiveDataProcessing: []byte{
						1, 2, 0, 0, 2, 1, 2, 0,
					},
					KnownChildSensitiveDataConsents: 2,
					AdditionalDataProcessingConsent: 1,
					MspaCoveredTransaction:          1,
					MspaOptOutOptionMode:            3,
					MspaServiceProviderMode:         0,
				},
				GPCSegment: sections.CommonUSGPCSegment{
					SubsectionType: 1,
					Gpc:            true,
				},
				SectionID: constants.SectionUSPIA,
				Value:     "bSFgmJcA.YA",
			},
		},
	}

	for _, test := range testData {
		result, err := NewUSPIA(test.gppString)
		encodedString := string(test.expected.Encode(true))
		assert.Nil(t, err)
		assert.Equal(t, test.expected, result)
		assert.Equal(t, constants.SectionUSPIA, result.GetID())
		assert.Equal(t, test.gppString, result.GetValue())
		assert.Equal(t, test.gppString, encodedString)
	}
}
