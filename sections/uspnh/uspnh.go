package uspnh

import (
	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/sections"
	"github.com/streamrail/go-gpp/util"
)

type USPNHCoreSegment struct {
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

func NewUSNHCoreSegment(bs *util.BitStream) (USPNHCoreSegment, error) {
	var usnh USPNHCoreSegment
	var err error

	usnh.Version, err = bs.ReadByte6()
	if err != nil {
		return usnh, sections.ErrorHelper("USNJSegment.Version", err)
	}

	usnh.ProcessingNotice, err = bs.ReadByte2()
	if err != nil {
		return usnh, sections.ErrorHelper("USNJSegment.ProcessingNotice", err)
	}

	usnh.SaleOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return usnh, sections.ErrorHelper("USNJSegment.SaleOptOutNotice", err)
	}

	usnh.TargetedAdvertisingOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return usnh, sections.ErrorHelper("USNJSegment.TargetedAdvertisingOptOutNotice", err)
	}

	usnh.SaleOptOut, err = bs.ReadByte2()
	if err != nil {
		return usnh, sections.ErrorHelper("USNJSegment.SaleOptOut", err)
	}

	usnh.TargetedAdvertisingOptOut, err = bs.ReadByte2()
	if err != nil {
		return usnh, sections.ErrorHelper("USNJSegment.TargetedAdvertisingOptOut", err)
	}

	usnh.SensitiveDataProcessing, err = bs.ReadTwoBitField(8)
	if err != nil {
		return usnh, sections.ErrorHelper("USNJSegment.SensitiveDataProcessing", err)
	}

	usnh.KnownChildSensitiveDataConsents, err = bs.ReadTwoBitField(3)
	if err != nil {
		return usnh, sections.ErrorHelper("USNJSegment.KnownChildSensitiveDataConsentsArr", err)
	}

	usnh.AdditionalDataProcessingConsent, err = bs.ReadByte2()
	if err != nil {
		return usnh, sections.ErrorHelper("USNJSegment.AdditionalDataProcessingConsent", err)
	}

	usnh.MspaCoveredTransaction, err = bs.ReadByte2()
	if err != nil {
		return usnh, sections.ErrorHelper("USNJSegment.MspaCoveredTransaction", err)
	}

	usnh.MspaOptOutOptionMode, err = bs.ReadByte2()
	if err != nil {
		return usnh, sections.ErrorHelper("USNJSegment.MspaOptOutOptionMode", err)
	}

	usnh.MspaServiceProviderMode, err = bs.ReadByte2()
	if err != nil {
		return usnh, sections.ErrorHelper("USNJSegment.MspaServiceProviderMode", err)
	}

	return usnh, nil
}

func (segment USPNHCoreSegment) Encode(bs *util.BitStream) {
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

type USPNH struct {
	SectionID   constants.SectionID
	Value       string
	CoreSegment USPNHCoreSegment
	GPCSegment  sections.CommonUSGPCSegment
}

func NewUSPNH(encoded string) (USPNH, error) {
	uspnh := USPNH{}

	coreBitStream, gpcBitStream, err := sections.CreateBitStreams(encoded, true)
	if err != nil {
		return uspnh, err
	}

	coreSegment, err := NewUSNHCoreSegment(coreBitStream)
	if err != nil {
		return uspnh, err
	}

	gpcSegment := sections.CommonUSGPCSegment{
		SubsectionType: 1,
		Gpc:            false,
	}

	if gpcBitStream != nil {
		gpcSegment, err = sections.NewCommonUSGPCSegment(gpcBitStream)
		if err != nil {
			return uspnh, err
		}
	}

	uspnh = USPNH{
		SectionID:   constants.SectionUSPNH,
		Value:       encoded,
		CoreSegment: coreSegment,
		GPCSegment:  gpcSegment,
	}

	return uspnh, nil
}

func (uspnh USPNH) Encode(gpcIncluded bool) []byte {
	bs := util.NewBitStreamForWrite()
	uspnh.CoreSegment.Encode(bs)
	res := bs.Base64Encode()
	if !gpcIncluded {
		return res
	}
	bs.Reset()
	res = append(res, '.')
	uspnh.GPCSegment.Encode(bs)
	return append(res, bs.Base64Encode()...)
}

func (uspnh USPNH) GetID() constants.SectionID {
	return uspnh.SectionID
}

func (uspnh USPNH) GetValue() string {
	return uspnh.Value
}
