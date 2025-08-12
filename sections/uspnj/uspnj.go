package uspnj

import (
	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/sections"
	"github.com/streamrail/go-gpp/util"
)

type USPNJCoreSegment struct {
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

func NewUSNJCoreSegment(bs *util.BitStream) (USPNJCoreSegment, error) {
	var usnj USPNJCoreSegment
	var err error

	usnj.Version, err = bs.ReadByte6()
	if err != nil {
		return usnj, sections.ErrorHelper("USNJSegment.Version", err)
	}

	usnj.ProcessingNotice, err = bs.ReadByte2()
	if err != nil {
		return usnj, sections.ErrorHelper("USNJSegment.ProcessingNotice", err)
	}

	usnj.SaleOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return usnj, sections.ErrorHelper("USNJSegment.SaleOptOutNotice", err)
	}

	usnj.TargetedAdvertisingOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return usnj, sections.ErrorHelper("USNJSegment.TargetedAdvertisingOptOutNotice", err)
	}

	usnj.SaleOptOut, err = bs.ReadByte2()
	if err != nil {
		return usnj, sections.ErrorHelper("USNJSegment.SaleOptOut", err)
	}

	usnj.TargetedAdvertisingOptOut, err = bs.ReadByte2()
	if err != nil {
		return usnj, sections.ErrorHelper("USNJSegment.TargetedAdvertisingOptOut", err)
	}

	usnj.SensitiveDataProcessing, err = bs.ReadTwoBitField(10)
	if err != nil {
		return usnj, sections.ErrorHelper("USNJSegment.SensitiveDataProcessing", err)
	}

	usnj.KnownChildSensitiveDataConsents, err = bs.ReadTwoBitField(5)
	if err != nil {
		return usnj, sections.ErrorHelper("USNJSegment.KnownChildSensitiveDataConsentsArr", err)
	}

	usnj.AdditionalDataProcessingConsent, err = bs.ReadByte2()
	if err != nil {
		return usnj, sections.ErrorHelper("USNJSegment.AdditionalDataProcessingConsent", err)
	}

	usnj.MspaCoveredTransaction, err = bs.ReadByte2()
	if err != nil {
		return usnj, sections.ErrorHelper("USNJSegment.MspaCoveredTransaction", err)
	}

	usnj.MspaOptOutOptionMode, err = bs.ReadByte2()
	if err != nil {
		return usnj, sections.ErrorHelper("USNJSegment.MspaOptOutOptionMode", err)
	}

	usnj.MspaServiceProviderMode, err = bs.ReadByte2()
	if err != nil {
		return usnj, sections.ErrorHelper("USNJSegment.MspaServiceProviderMode", err)
	}

	return usnj, nil
}

func (segment USPNJCoreSegment) Encode(bs *util.BitStream) {
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

type USPNJ struct {
	SectionID   constants.SectionID
	Value       string
	CoreSegment USPNJCoreSegment
	GPCSegment  sections.CommonUSGPCSegment
}

func NewUSPNJ(encoded string) (USPNJ, error) {
	uspnj := USPNJ{}

	coreBitStream, gpcBitStream, err := sections.CreateBitStreams(encoded, true)
	if err != nil {
		return uspnj, err
	}

	coreSegment, err := NewUSNJCoreSegment(coreBitStream)
	if err != nil {
		return uspnj, err
	}

	gpcSegment := sections.CommonUSGPCSegment{
		SubsectionType: 1,
		Gpc:            false,
	}

	if gpcBitStream != nil {
		gpcSegment, err = sections.NewCommonUSGPCSegment(gpcBitStream)
		if err != nil {
			return uspnj, err
		}
	}

	uspnj = USPNJ{
		SectionID:   constants.SectionUSPNJ,
		Value:       encoded,
		CoreSegment: coreSegment,
		GPCSegment:  gpcSegment,
	}

	return uspnj, nil
}

func (uspnj USPNJ) Encode(gpcIncluded bool) []byte {
	bs := util.NewBitStreamForWrite()
	uspnj.CoreSegment.Encode(bs)
	res := bs.Base64Encode()
	if !gpcIncluded {
		return res
	}
	bs.Reset()
	res = append(res, '.')
	uspnj.GPCSegment.Encode(bs)
	return append(res, bs.Base64Encode()...)
}

func (uspnj USPNJ) GetID() constants.SectionID {
	return uspnj.SectionID
}

func (uspnj USPNJ) GetValue() string {
	return uspnj.Value
}
