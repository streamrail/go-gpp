package uspin

import (
	"testing"

	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/util"
	"github.com/stretchr/testify/assert"
)

type uspinTestData struct {
	description string
	gppString   string
	expected    USPIN
}

func TestUSPIN(t *testing.T) {
	testData := []uspinTestData{
		{
			description: "should populate USPIN segments correctly",
			gppString:   "ZBACbAAAAAAA",
			/*
				Header: SectionID(6) + Version(6) + SubSections (Fibonacci Range indicating subsections 0-1)
				Core Subsection (24 bits): MspaVersion etc
				Sensitive Data Consents (16 bits): SensitiveDataProcessing
			*/
			expected: USPIN{
				Header: USPINHeader{
					SectionID: 25,
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
				CoreSegment: USPINCoreSegment{
					MspaVersion:                     0,
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
				SensitiveDataConsentsSegment: USPINSensitiveDataConsentsSegment{
					SensitiveDataProcessing: []byte{
						0, 0, 0, 0, 0, 0, 0, 0,
					},
				},
				SectionID: constants.SectionUSPIN,
				Value:     "ZBACbAAAAAAA",
			},
		},
	}

	for _, test := range testData {
		result, err := NewUSPIN(test.gppString)
		encodedString := string(test.expected.Encode(true))
		assert.Nil(t, err)
		assert.Equal(t, test.expected, result)
		assert.Equal(t, constants.SectionUSPIN, result.GetID())
		assert.Equal(t, test.gppString, result.GetValue())
		assert.Equal(t, test.gppString, encodedString)
	}
}