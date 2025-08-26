package uspmn

import (
	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/sections"
	"github.com/streamrail/go-gpp/util"
)

type USPMNCoreSegment struct {
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

type USPMN struct {
	SectionID   constants.SectionID
	Value       string
	CoreSegment USPMNCoreSegment
	GPCSegment  sections.CommonUSGPCSegment
}

func NewUSPMNCoreSegment(bs *util.BitStream) (USPMNCoreSegment, error) {
	var uspmnCore USPMNCoreSegment
	var err error

	uspmnCore.Version, err = bs.ReadByte6()
	if err != nil {
		return uspmnCore, sections.ErrorHelper("CoreSegment.Version", err)
	}

	uspmnCore.ProcessingNotice, err = bs.ReadByte2()
	if err != nil {
		return uspmnCore, sections.ErrorHelper("CoreSegment.ProcessingNotice", err)
	}

	uspmnCore.SaleOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return uspmnCore, sections.ErrorHelper("CoreSegment.SalesOptOutNotice", err)
	}

	uspmnCore.TargetedAdvertisingOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return uspmnCore, sections.ErrorHelper("CoreSegment.TargetedAdvertisingOptOutNotice", err)
	}

	uspmnCore.SaleOptOut, err = bs.ReadByte2()
	if err != nil {
		return uspmnCore, sections.ErrorHelper("CoreSegment.SaleOptOut", err)
	}

	uspmnCore.TargetedAdvertisingOptOut, err = bs.ReadByte2()
	if err != nil {
		return uspmnCore, sections.ErrorHelper("CoreSegment.TargetedAdvertisingOptOut", err)
	}

	uspmnCore.SensitiveDataProcessing, err = bs.ReadTwoBitField(8)
	if err != nil {
		return uspmnCore, sections.ErrorHelper("CoreSegment.SensitiveDataProcessing", err)
	}

	uspmnCore.KnownChildSensitiveDataConsents, err = bs.ReadByte2()
	if err != nil {
		return uspmnCore, sections.ErrorHelper("CoreSegment.KnownChildSensitiveDataConsents", err)
	}

	uspmnCore.AdditionalDataProcessingConsent, err = bs.ReadByte2()
	if err != nil {
		return uspmnCore, sections.ErrorHelper("CoreSegment.AdditionalDataProcessingConsent", err)
	}

	uspmnCore.MspaCoveredTransaction, err = bs.ReadByte2()
	if err != nil {
		return uspmnCore, sections.ErrorHelper("CoreSegment.MspaCoveredTransaction", err)
	}

	uspmnCore.MspaOptOutOptionMode, err = bs.ReadByte2()
	if err != nil {
		return uspmnCore, sections.ErrorHelper("CoreSegment.MspaOptOutOptionMode", err)
	}

	uspmnCore.MspaServiceProviderMode, err = bs.ReadByte2()
	if err != nil {
		return uspmnCore, sections.ErrorHelper("CoreSegment.MspaServiceProviderMode", err)
	}

	return uspmnCore, err
}

func (segment USPMNCoreSegment) Encode(bs *util.BitStream) {
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

func NewUSPMN(encoded string) (USPMN, error) {
	coreBitStream, gpcBitStream, err := sections.CreateBitStreams(encoded, true)
	if err != nil {
		return USPMN{}, err
	}

	coreSegment, err := NewUSPMNCoreSegment(coreBitStream)
	if err != nil {
		return USPMN{}, err
	}

	gpcSegment := sections.CommonUSGPCSegment{
		SubsectionType: 1,
		Gpc:            false,
	}

	if gpcBitStream != nil {
		gpcSegment, err = sections.NewCommonUSGPCSegment(gpcBitStream)
		if err != nil {
			return USPMN{}, err
		}
	}

	return USPMN{
		SectionID:   constants.SectionUSPMN,
		Value:       encoded,
		CoreSegment: coreSegment,
		GPCSegment:  gpcSegment,
	}, nil
}

func (uspmn USPMN) Encode(gpcIncluded bool) []byte {
	bs := util.NewBitStreamForWrite()
	uspmn.CoreSegment.Encode(bs)
	res := bs.Base64Encode()
	if !gpcIncluded {
		return res
	}
	bs.Reset()
	res = append(res, '.')
	uspmn.GPCSegment.Encode(bs)
	return append(res, bs.Base64Encode()...)
}

func (uspmn USPMN) GetID() constants.SectionID {
	return uspmn.SectionID
}

func (uspmn USPMN) GetValue() string {
	return uspmn.Value
}
