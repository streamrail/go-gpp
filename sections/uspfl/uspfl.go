package uspfl

import (
	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/sections"
	"github.com/streamrail/go-gpp/util"
)

type USPFLCoreSegment struct {
	Version                         byte
	ProcessingNotice                byte
	SaleOptOutNotice                byte
	TargetedAdvertisingOptOutNotice byte
	SaleOptOut                      byte
	TargetedAdvertisingOptOut       byte
	SensitiveDataProcessing         []byte
	KnownChildSensitiveDataConsents []byte
	AdditionalDataProcessingConsent byte
	MspaCoveredTransaction          byte
	MspaOptOutOptionMode            byte
	MspaServiceProviderMode         byte
}

func NewUSFLCoreSegment(bs *util.BitStream) (USPFLCoreSegment, error) {
	var usfl USPFLCoreSegment
	var err error

	usfl.Version, err = bs.ReadByte6()
	if err != nil {
		return usfl, sections.ErrorHelper("USFLSegment.Version", err)
	}

	usfl.ProcessingNotice, err = bs.ReadByte2()
	if err != nil {
		return usfl, sections.ErrorHelper("USFLSegment.ProcessingNotice", err)
	}

	usfl.SaleOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return usfl, sections.ErrorHelper("USFLSegment.SaleOptOutNotice", err)
	}

	usfl.TargetedAdvertisingOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return usfl, sections.ErrorHelper("USFLSegment.TargetedAdvertisingOptOutNotice", err)
	}

	usfl.SaleOptOut, err = bs.ReadByte2()
	if err != nil {
		return usfl, sections.ErrorHelper("USFLSegment.SaleOptOut", err)
	}

	usfl.TargetedAdvertisingOptOut, err = bs.ReadByte2()
	if err != nil {
		return usfl, sections.ErrorHelper("USFLSegment.TargetedAdvertisingOptOut", err)
	}

	usfl.SensitiveDataProcessing, err = bs.ReadTwoBitField(8)
	if err != nil {
		return usfl, sections.ErrorHelper("USFLSegment.SensitiveDataProcessing", err)
	}

	usfl.KnownChildSensitiveDataConsents, err = bs.ReadTwoBitField(3)
	if err != nil {
		return usfl, sections.ErrorHelper("USFLSegment.KnownChildSensitiveDataConsentsArr", err)
	}

	usfl.AdditionalDataProcessingConsent, err = bs.ReadByte2()
	if err != nil {
		return usfl, sections.ErrorHelper("USFLSegment.AdditionalDataProcessingConsent", err)
	}

	usfl.MspaCoveredTransaction, err = bs.ReadByte2()
	if err != nil {
		return usfl, sections.ErrorHelper("USFLSegment.MspaCoveredTransaction", err)
	}

	usfl.MspaOptOutOptionMode, err = bs.ReadByte2()
	if err != nil {
		return usfl, sections.ErrorHelper("USFLSegment.MspaOptOutOptionMode", err)
	}

	usfl.MspaServiceProviderMode, err = bs.ReadByte2()
	if err != nil {
		return usfl, sections.ErrorHelper("USFLSegment.MspaServiceProviderMode", err)
	}

	return usfl, nil
}

func (segment USPFLCoreSegment) Encode(bs *util.BitStream) {
	bs.WriteByte6(segment.Version)
	bs.WriteByte2(segment.ProcessingNotice)
	bs.WriteByte2(segment.SaleOptOutNotice)
	bs.WriteByte2(segment.TargetedAdvertisingOptOutNotice)
	bs.WriteByte2(segment.SaleOptOut)
	bs.WriteByte2(segment.TargetedAdvertisingOptOut)
	bs.WriteTwoBitField(segment.SensitiveDataProcessing)
	bs.WriteTwoBitField(segment.KnownChildSensitiveDataConsents)
	bs.WriteByte2(segment.AdditionalDataProcessingConsent)
	bs.WriteByte2(segment.MspaCoveredTransaction)
	bs.WriteByte2(segment.MspaOptOutOptionMode)
	bs.WriteByte2(segment.MspaServiceProviderMode)
}

type USPFL struct {
	SectionID   constants.SectionID
	Value       string
	CoreSegment USPFLCoreSegment
	GPCSegment  sections.CommonUSGPCSegment
}

func NewUSPFL(encoded string) (USPFL, error) {
	uspde := USPFL{}

	coreBitStream, gpcBitStream, err := sections.CreateBitStreams(encoded, true)
	if err != nil {
		return uspde, err
	}

	coreSegment, err := NewUSFLCoreSegment(coreBitStream)
	if err != nil {
		return uspde, err
	}

	gpcSegment := sections.CommonUSGPCSegment{
		SubsectionType: 1,
		Gpc:            false,
	}

	if gpcBitStream != nil {
		gpcSegment, err = sections.NewCommonUSGPCSegment(gpcBitStream)
		if err != nil {
			return uspde, err
		}
	}

	uspde = USPFL{
		SectionID:   constants.SectionUSPFL,
		Value:       encoded,
		CoreSegment: coreSegment,
		GPCSegment:  gpcSegment,
	}

	return uspde, nil
}

func (uspde USPFL) Encode(gpcIncluded bool) []byte {
	bs := util.NewBitStreamForWrite()
	uspde.CoreSegment.Encode(bs)
	res := bs.Base64Encode()
	if !gpcIncluded {
		return res
	}
	bs.Reset()
	res = append(res, '.')
	uspde.GPCSegment.Encode(bs)
	return append(res, bs.Base64Encode()...)
}

func (uspde USPFL) GetID() constants.SectionID {
	return uspde.SectionID
}

func (uspde USPFL) GetValue() string {
	return uspde.Value
}
