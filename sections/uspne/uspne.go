package uspne

import (
	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/sections"
	"github.com/streamrail/go-gpp/util"
)

type USPNECoreSegment struct {
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

func NewUSNECoreSegment(bs *util.BitStream) (USPNECoreSegment, error) {
	var uspne USPNECoreSegment
	var err error

	uspne.Version, err = bs.ReadByte6()
	if err != nil {
		return uspne, sections.ErrorHelper("USIASegment.Version", err)
	}

	uspne.ProcessingNotice, err = bs.ReadByte2()
	if err != nil {
		return uspne, sections.ErrorHelper("USIASegment.ProcessingNotice", err)
	}

	uspne.SaleOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return uspne, sections.ErrorHelper("USIASegment.SaleOptOutNotice", err)
	}

	uspne.TargetedAdvertisingOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return uspne, sections.ErrorHelper("USIASegment.TargetedAdvertisingOptOutNotice", err)
	}

	uspne.SaleOptOut, err = bs.ReadByte2()
	if err != nil {
		return uspne, sections.ErrorHelper("USIASegment.SaleOptOut", err)
	}

	uspne.TargetedAdvertisingOptOut, err = bs.ReadByte2()
	if err != nil {
		return uspne, sections.ErrorHelper("USIASegment.TargetedAdvertisingOptOut", err)
	}

	uspne.SensitiveDataProcessing, err = bs.ReadTwoBitField(8)
	if err != nil {
		return uspne, sections.ErrorHelper("USIASegment.SensitiveDataProcessing", err)
	}

	uspne.KnownChildSensitiveDataConsents, err = bs.ReadByte2()
	if err != nil {
		return uspne, sections.ErrorHelper("USIASegment.KnownChildSensitiveDataConsentsArr", err)
	}

	uspne.AdditionalDataProcessingConsent, err = bs.ReadByte2()
	if err != nil {
		return uspne, sections.ErrorHelper("USIASegment.AdditionalDataProcessingConsent", err)
	}

	uspne.MspaCoveredTransaction, err = bs.ReadByte2()
	if err != nil {
		return uspne, sections.ErrorHelper("USIASegment.MspaCoveredTransaction", err)
	}

	uspne.MspaOptOutOptionMode, err = bs.ReadByte2()
	if err != nil {
		return uspne, sections.ErrorHelper("USIASegment.MspaOptOutOptionMode", err)
	}

	uspne.MspaServiceProviderMode, err = bs.ReadByte2()
	if err != nil {
		return uspne, sections.ErrorHelper("USIASegment.MspaServiceProviderMode", err)
	}

	return uspne, nil
}

func (segment USPNECoreSegment) Encode(bs *util.BitStream) {
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

type USPNE struct {
	SectionID   constants.SectionID
	Value       string
	CoreSegment USPNECoreSegment
	GPCSegment  sections.CommonUSGPCSegment
}

func NewUSPNE(encoded string) (USPNE, error) {
	uspne := USPNE{}

	coreBitStream, gpcBitStream, err := sections.CreateBitStreams(encoded, true)
	if err != nil {
		return uspne, err
	}

	coreSegment, err := NewUSNECoreSegment(coreBitStream)
	if err != nil {
		return uspne, err
	}

	gpcSegment := sections.CommonUSGPCSegment{
		SubsectionType: 1,
		Gpc:            false,
	}

	if gpcBitStream != nil {
		gpcSegment, err = sections.NewCommonUSGPCSegment(gpcBitStream)
		if err != nil {
			return uspne, err
		}
	}

	uspne = USPNE{
		SectionID:   constants.SectionUSPNE,
		Value:       encoded,
		CoreSegment: coreSegment,
		GPCSegment:  gpcSegment,
	}

	return uspne, nil
}

func (uspne USPNE) Encode(gpcIncluded bool) []byte {
	bs := util.NewBitStreamForWrite()
	uspne.CoreSegment.Encode(bs)
	res := bs.Base64Encode()
	if !gpcIncluded {
		return res
	}
	bs.Reset()
	res = append(res, '.')
	uspne.GPCSegment.Encode(bs)
	return append(res, bs.Base64Encode()...)
}

func (uspne USPNE) GetID() constants.SectionID {
	return uspne.SectionID
}

func (uspne USPNE) GetValue() string {
	return uspne.Value
}
