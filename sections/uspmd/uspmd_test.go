package uspmd

import (
	"testing"

	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/sections"
	"github.com/streamrail/go-gpp/util"
	"github.com/stretchr/testify/assert"
)

type uspmdTestData struct {
	description string
	gppString   string
	expected    USPMD
}

func TestUSPMD(t *testing.T) {
	testData := []uspmdTestData{
		{
			description: "should populate USPMD segments correctly",
			gppString:   "YBACbAAAAA.QA",
			/*
				Header: SectionID(6) + Version(6) + SubSections (Fibonacci Range indicating subsections 0-1)
				Core Subsection (22 bits): MspaVersion(6) + 8 x Int(2), all zeros
				GPC Subsection: SubsectionType=1, Gpc=false
				(padded to byte boundary: 54 bits → 60 bits = 10 base64 chars)
			*/
			expected: USPMD{
				Header: USPMDHeader{
					SectionID: 24,
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
				CoreSegment: USPMDCoreSegment{
					MspaVersion:                     0,
					MspaCoveredTransaction:          0,
					MspaMode:                        0,
					ProcessingNotice:                0,
					SaleOptOutNotice:                0,
					TargetedAdvertisingOptOutNotice: 0,
					SaleOptOut:                      0,
					TargetedAdvertisingOptOut:       0,
					AdditionalDataProcessingConsent: 0,
				},
				GPCSegment: sections.CommonUSGPCSegment{
					SubsectionType: 1,
					Gpc:            false,
				},
				SectionID: constants.SectionUSPMD,
				Value:     "YBACbAAAAA.QA",
			},
		},
	}

	for _, test := range testData {
		result, err := NewUSPMD(test.gppString)
		encodedString := string(test.expected.Encode(true))
		assert.Nil(t, err)
		assert.Equal(t, test.expected, result)
		assert.Equal(t, constants.SectionUSPMD, result.GetID())
		assert.Equal(t, test.gppString, result.GetValue())
		assert.Equal(t, test.gppString, encodedString)
	}
}
