package uspde

import (
	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/sections"
	"github.com/streamrail/go-gpp/util"
)

type USPDECoreSegment struct {
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

func NewUSDECoreSegment(bs *util.BitStream) (USPDECoreSegment, error) {
	var usde USPDECoreSegment
	var err error

	usde.Version, err = bs.ReadByte6()
	if err != nil {
		return usde, sections.ErrorHelper("USDESegment.Version", err)
	}

	usde.ProcessingNotice, err = bs.ReadByte2()
	if err != nil {
		return usde, sections.ErrorHelper("USDESegment.ProcessingNotice", err)
	}

	usde.SaleOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return usde, sections.ErrorHelper("USDESegment.SaleOptOutNotice", err)
	}

	usde.TargetedAdvertisingOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return usde, sections.ErrorHelper("USDESegment.TargetedAdvertisingOptOutNotice", err)
	}

	usde.SaleOptOut, err = bs.ReadByte2()
	if err != nil {
		return usde, sections.ErrorHelper("USDESegment.SaleOptOut", err)
	}

	usde.TargetedAdvertisingOptOut, err = bs.ReadByte2()
	if err != nil {
		return usde, sections.ErrorHelper("USDESegment.TargetedAdvertisingOptOut", err)
	}

	usde.SensitiveDataProcessing, err = bs.ReadTwoBitField(9)
	if err != nil {
		return usde, sections.ErrorHelper("USDESegment.SensitiveDataProcessing", err)
	}

	usde.KnownChildSensitiveDataConsents, err = bs.ReadTwoBitField(5)
	if err != nil {
		return usde, sections.ErrorHelper("USDESegment.KnownChildSensitiveDataConsentsArr", err)
	}

	usde.AdditionalDataProcessingConsent, err = bs.ReadByte2()
	if err != nil {
		return usde, sections.ErrorHelper("USDESegment.AdditionalDataProcessingConsent", err)
	}

	usde.MspaCoveredTransaction, err = bs.ReadByte2()
	if err != nil {
		return usde, sections.ErrorHelper("USDESegment.MspaCoveredTransaction", err)
	}

	usde.MspaOptOutOptionMode, err = bs.ReadByte2()
	if err != nil {
		return usde, sections.ErrorHelper("USDESegment.MspaOptOutOptionMode", err)
	}

	usde.MspaServiceProviderMode, err = bs.ReadByte2()
	if err != nil {
		return usde, sections.ErrorHelper("USDESegment.MspaServiceProviderMode", err)
	}

	return usde, nil
}

func (segment USPDECoreSegment) Encode(bs *util.BitStream) {
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

type USPDE struct {
	SectionID   constants.SectionID
	Value       string
	CoreSegment USPDECoreSegment
	GPCSegment  sections.CommonUSGPCSegment
}

func NewUSPDE(encoded string) (USPDE, error) {
	uspde := USPDE{}

	coreBitStream, gpcBitStream, err := sections.CreateBitStreams(encoded, true)
	if err != nil {
		return uspde, err
	}

	coreSegment, err := NewUSDECoreSegment(coreBitStream)
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

	uspde = USPDE{
		SectionID:   constants.SectionUSPDE,
		Value:       encoded,
		CoreSegment: coreSegment,
		GPCSegment:  gpcSegment,
	}

	return uspde, nil
}

func (uspde USPDE) Encode(gpcIncluded bool) []byte {
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

func (uspde USPDE) GetID() constants.SectionID {
	return uspde.SectionID
}

func (uspde USPDE) GetValue() string {
	return uspde.Value
}
