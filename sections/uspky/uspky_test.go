package uspky

import (
	"testing"

	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/util"
	"github.com/stretchr/testify/assert"
)

type uspkyTestData struct {
	description string
	gppString   string
	expected    USPKY
}

func TestUSPKY(t *testing.T) {
	testData := []uspkyTestData{
		{
			description: "should populate USPKY segments correctly",
			gppString:   "aBACbbAAAAAA",
			/*
				Header: SectionID(6) + Version(6) + SubSections (Fibonacci Range indicating subsections 0-1)
				Core Subsection (24 bits): 011011 00 00 00 00 00 00 00 00 00
				Sensitive Data Consents (16 bits): 0000000000000000
			*/
			expected: USPKY{
				Header: USPKYHeader{
					SectionID: 26,
					Version:   1,
					SubSections: &util.IntRange{
						Size: 2,
						Range: []util.IRange{
							{StartID: 0, EndID: 0},
							{StartID: 1, EndID: 1},
						},
						Max: 1,
					},
				},
				CoreSegment: USPKYCoreSegment{
					MspaVersion:                     27,
					MspaCoveredTransaction:          0,
					MspaMode:                        0,
					ProcessingNotice:                0,
					SaleOptOutNotice:                0,
					TargetedAdvertisingOptOutNotice: 0,
					SaleOptOut:                      0,
					TargetedAdvertisingOptOut:       0,
					KnownChildSensitiveDataConsents: 0,
					AdditionalDataProcessingConsent: 0,
				},
				SensitiveDataConsentsSegment: USPKYSensitiveDataConsentsSegment{
					SensitiveDataProcessing: []byte{
						0, 0, 0, 0, 0, 0, 0, 0,
					},
				},
				SectionID: constants.SectionUSPKY,
				Value:     "aBACbbAAAAAA",
			},
		},
	}

	for _, test := range testData {
		result, err := NewUSPKY(test.gppString)
		encodedString := string(test.expected.Encode(true))
		assert.Nil(t, err)
		assert.Equal(t, test.expected, result)
		assert.Equal(t, constants.SectionUSPKY, result.GetID())
		assert.Equal(t, test.gppString, result.GetValue())
		assert.Equal(t, test.gppString, encodedString)
	}
}