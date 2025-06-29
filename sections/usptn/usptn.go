package usptn

import (
	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/sections"
	"github.com/streamrail/go-gpp/util"
)

type USPTNCoreSegment struct {
	Version                         byte
	ProcessingNotice                byte
	SaleOptOutNotice                byte
	TargetedAdvertisingOptOutNotice byte
	SaleOptOut                      byte
	TargetedAdvertisingOptOut       byte
	SensitiveDataProcessing         []byte
	KnownChildSensitiveDataConsents byte
	AdditionalDataProcessingConsent byte
	MspaCoveredTransaction          byte
	MspaOptOutOptionMode            byte
	MspaServiceProviderMode         byte
}

type USPTN struct {
	SectionID   constants.SectionID
	Value       string
	CoreSegment USPTNCoreSegment
	GPCSegment  sections.CommonUSGPCSegment
}

func NewUSPTNCoreSegment(bs *util.BitStream) (USPTNCoreSegment, error) {
	var usptnCore USPTNCoreSegment
	var err error

	usptnCore.Version, err = bs.ReadByte6()
	if err != nil {
		return usptnCore, sections.ErrorHelper("CoreSegment.Version", err)
	}

	usptnCore.ProcessingNotice, err = bs.ReadByte2()
	if err != nil {
		return usptnCore, sections.ErrorHelper("CoreSegment.ProcessingNotice", err)
	}

	usptnCore.SaleOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return usptnCore, sections.ErrorHelper("CoreSegment.SalesOptOutNotice", err)
	}

	usptnCore.TargetedAdvertisingOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return usptnCore, sections.ErrorHelper("CoreSegment.TargetedAdvertisingOptOutNotice", err)
	}

	usptnCore.SaleOptOut, err = bs.ReadByte2()
	if err != nil {
		return usptnCore, sections.ErrorHelper("CoreSegment.SaleOptOut", err)
	}

	usptnCore.TargetedAdvertisingOptOut, err = bs.ReadByte2()
	if err != nil {
		return usptnCore, sections.ErrorHelper("CoreSegment.TargetedAdvertisingOptOut", err)
	}

	usptnCore.SensitiveDataProcessing, err = bs.ReadTwoBitField(8)
	if err != nil {
		return usptnCore, sections.ErrorHelper("CoreSegment.SensitiveDataProcessing", err)
	}

	usptnCore.KnownChildSensitiveDataConsents, err = bs.ReadByte2()
	if err != nil {
		return usptnCore, sections.ErrorHelper("CoreSegment.KnownChildSensitiveDataConsents", err)
	}

	usptnCore.AdditionalDataProcessingConsent, err = bs.ReadByte2()
	if err != nil {
		return usptnCore, sections.ErrorHelper("CoreSegment.AdditionalDataProcessingConsent", err)
	}

	usptnCore.MspaCoveredTransaction, err = bs.ReadByte2()
	if err != nil {
		return usptnCore, sections.ErrorHelper("CoreSegment.MspaCoveredTransaction", err)
	}

	usptnCore.MspaOptOutOptionMode, err = bs.ReadByte2()
	if err != nil {
		return usptnCore, sections.ErrorHelper("CoreSegment.MspaOptOutOptionMode", err)
	}

	usptnCore.MspaServiceProviderMode, err = bs.ReadByte2()
	if err != nil {
		return usptnCore, sections.ErrorHelper("CoreSegment.MspaServiceProviderMode", err)
	}

	return usptnCore, err
}

func (segment USPTNCoreSegment) Encode(bs *util.BitStream) {
	bs.WriteByte6(segment.Version)
	bs.WriteByte2(segment.ProcessingNotice)
	bs.WriteByte2(segment.SaleOptOutNotice)
	bs.WriteByte2(segment.TargetedAdvertisingOptOutNotice)
	bs.WriteByte2(segment.SaleOptOut)
	bs.WriteByte2(segment.TargetedAdvertisingOptOut)
	bs.WriteTwoBitField(segment.SensitiveDataProcessing)
	bs.WriteByte2(segment.KnownChildSensitiveDataConsents)
	bs.WriteByte2(segment.AdditionalDataProcessingConsent)
	bs.WriteByte2(segment.MspaCoveredTransaction)
	bs.WriteByte2(segment.MspaOptOutOptionMode)
	bs.WriteByte2(segment.MspaServiceProviderMode)
}

func NewUSPTN(encoded string) (USPTN, error) {
	usptn := USPTN{}

	coreBitStream, gpcBitStream, err := sections.CreateBitStreams(encoded, true)
	if err != nil {
		return usptn, err
	}

	coreSegment, err := NewUSPTNCoreSegment(coreBitStream)
	if err != nil {
		return usptn, err
	}

	gpcSegment := sections.CommonUSGPCSegment{
		SubsectionType: 1,
		Gpc:            false,
	}

	if gpcBitStream != nil {
		gpcSegment, err = sections.NewCommonUSGPCSegment(gpcBitStream)
		if err != nil {
			return usptn, err
		}
	}

	usptn = USPTN{
		SectionID:   constants.SectionUSPTN,
		Value:       encoded,
		CoreSegment: coreSegment,
		GPCSegment:  gpcSegment,
	}

	return usptn, nil
}

func (usptn USPTN) Encode(gpcIncluded bool) []byte {
	bs := util.NewBitStreamForWrite()
	usptn.CoreSegment.Encode(bs)
	res := bs.Base64Encode()
	if !gpcIncluded {
		return res
	}
	bs.Reset()
	res = append(res, '.')
	usptn.GPCSegment.Encode(bs)
	return append(res, bs.Base64Encode()...)
}

func (usptn USPTN) GetID() constants.SectionID {
	return usptn.SectionID
}

func (usptn USPTN) GetValue() string {
	return usptn.Value
}
