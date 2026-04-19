package uspia

import (
	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/sections"
	"github.com/streamrail/go-gpp/util"
)

type USPIACoreSegment struct {
	Version                         byte
	ProcessingNotice                byte
	SaleOptOutNotice                byte
	TargetedAdvertisingOptOutNotice byte
	SensitiveDataOptOutNotice       byte
	SaleOptOut                      byte
	TargetedAdvertisingOptOut       byte
	SensitiveDataProcessing         []byte
	KnownChildSensitiveDataConsents byte
	MspaCoveredTransaction          byte
	MspaOptOutOptionMode            byte
	MspaServiceProviderMode         byte
}

func NewUSIACoreSegment(bs *util.BitStream) (USPIACoreSegment, error) {
	var usia USPIACoreSegment
	var err error

	usia.Version, err = bs.ReadByte6()
	if err != nil {
		return usia, sections.ErrorHelper("USIASegment.Version", err)
	}

	usia.ProcessingNotice, err = bs.ReadByte2()
	if err != nil {
		return usia, sections.ErrorHelper("USIASegment.ProcessingNotice", err)
	}

	usia.SaleOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return usia, sections.ErrorHelper("USIASegment.SaleOptOutNotice", err)
	}

	usia.TargetedAdvertisingOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return usia, sections.ErrorHelper("USIASegment.TargetedAdvertisingOptOutNotice", err)
	}

	usia.SensitiveDataOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return usia, sections.ErrorHelper("USIASegment.SensitiveDataOptOutNotice", err)
	}

	usia.SaleOptOut, err = bs.ReadByte2()
	if err != nil {
		return usia, sections.ErrorHelper("USIASegment.SaleOptOut", err)
	}

	usia.TargetedAdvertisingOptOut, err = bs.ReadByte2()
	if err != nil {
		return usia, sections.ErrorHelper("USIASegment.TargetedAdvertisingOptOut", err)
	}

	usia.SensitiveDataProcessing, err = bs.ReadTwoBitField(8)
	if err != nil {
		return usia, sections.ErrorHelper("USIASegment.SensitiveDataProcessing", err)
	}

	usia.KnownChildSensitiveDataConsents, err = bs.ReadByte2()
	if err != nil {
		return usia, sections.ErrorHelper("USIASegment.KnownChildSensitiveDataConsentsArr", err)
	}

	usia.MspaCoveredTransaction, err = bs.ReadByte2()
	if err != nil {
		return usia, sections.ErrorHelper("USIASegment.MspaCoveredTransaction", err)
	}

	usia.MspaOptOutOptionMode, err = bs.ReadByte2()
	if err != nil {
		return usia, sections.ErrorHelper("USIASegment.MspaOptOutOptionMode", err)
	}

	usia.MspaServiceProviderMode, err = bs.ReadByte2()
	if err != nil {
		return usia, sections.ErrorHelper("USIASegment.MspaServiceProviderMode", err)
	}

	return usia, nil
}

func (segment USPIACoreSegment) Encode(bs *util.BitStream) {
	bs.WriteByte6(segment.Version)
	bs.WriteByte2(segment.ProcessingNotice)
	bs.WriteByte2(segment.SaleOptOutNotice)
	bs.WriteByte2(segment.TargetedAdvertisingOptOutNotice)
	bs.WriteByte2(segment.SensitiveDataOptOutNotice)
	bs.WriteByte2(segment.SaleOptOut)
	bs.WriteByte2(segment.TargetedAdvertisingOptOut)
	bs.WriteTwoBitField(segment.SensitiveDataProcessing)
	bs.WriteByte2(segment.KnownChildSensitiveDataConsents)
	bs.WriteByte2(segment.MspaCoveredTransaction)
	bs.WriteByte2(segment.MspaOptOutOptionMode)
	bs.WriteByte2(segment.MspaServiceProviderMode)
}

type USPIA struct {
	SectionID   constants.SectionID
	Value       string
	CoreSegment USPIACoreSegment
	GPCSegment  sections.CommonUSGPCSegment
}

func NewUSPIA(encoded string) (USPIA, error) {
	uspia := USPIA{}

	coreBitStream, gpcBitStream, err := sections.CreateBitStreams(encoded, true)
	if err != nil {
		return uspia, err
	}

	coreSegment, err := NewUSIACoreSegment(coreBitStream)
	if err != nil {
		return uspia, err
	}

	gpcSegment := sections.CommonUSGPCSegment{
		SubsectionType: 1,
		Gpc:            false,
	}

	if gpcBitStream != nil {
		gpcSegment, err = sections.NewCommonUSGPCSegment(gpcBitStream)
		if err != nil {
			return uspia, err
		}
	}

	uspia = USPIA{
		SectionID:   constants.SectionUSPIA,
		Value:       encoded,
		CoreSegment: coreSegment,
		GPCSegment:  gpcSegment,
	}

	return uspia, nil
}

func (uspia USPIA) Encode(gpcIncluded bool) []byte {
	bs := util.NewBitStreamForWrite()
	uspia.CoreSegment.Encode(bs)
	res := bs.Base64Encode()
	if !gpcIncluded {
		return res
	}
	bs.Reset()
	res = append(res, '.')
	uspia.GPCSegment.Encode(bs)
	return append(res, bs.Base64Encode()...)
}

func (uspia USPIA) GetID() constants.SectionID {
	return uspia.SectionID
}

func (uspia USPIA) GetValue() string {
	return uspia.Value
}
