package uspnat

import (
	"testing"

	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/sections"
	"github.com/stretchr/testify/assert"
)

type uspnatTestData struct {
	description string
	gppString   string
	expected    USPNAT
}

func TestUSPNAT(t *testing.T) {
	testData := []uspnatTestData{
		{
			description: "should populate USPNAT segments correctly",
			gppString:   "BaaVEiFKEJJo.YA",
			expected: USPNAT{
				CoreSegment: USPNATCoreSegment{
					Version:                             1,
					SharingNotice:                       1,
					SaleOptOutNotice:                    2,
					SharingOptOutNotice:                 2,
					TargetedAdvertisingOptOutNotice:     1,
					SensitiveDataProcessingOptOutNotice: 2,
					SensitiveDataLimitUseNotice:         2,
					SaleOptOut:                          1,
					SharingOptOut:                       1,
					TargetedAdvertisingOptOut:           1,
					SensitiveDataProcessing: []byte{
						0, 1, 0, 2, 0, 2, 0, 1, 1, 0, 2, 2, 0, 1, 0, 0,
					},
					KnownChildSensitiveDataConsents: []byte{
						2, 1, 0,
					},
					PersonalDataConsents:    2,
					MspaCoveredTransaction:  1,
					MspaOptOutOptionMode:    2,
					MspaServiceProviderMode: 2,
				},
				GPCSegment: sections.CommonUSGPCSegment{
					SubsectionType: 1,
					Gpc:            true,
				},
				SectionID: constants.SectionUSPNAT,
				Value:     "BaaVEiFKEJJo.YA",
			},
		},
	}

	for _, test := range testData {
		result, err := NewUSPNAT(test.gppString)
		encodedString := string(test.expected.Encode(true))

		assert.Nil(t, err)
		assert.Equal(t, test.expected, result)
		assert.Equal(t, constants.SectionUSPNAT, result.GetID())
		assert.Equal(t, test.gppString, result.GetValue())
		assert.Equal(t, test.gppString, encodedString)
	}
}
