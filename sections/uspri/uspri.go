package uspri

import (
	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/sections"
	"github.com/streamrail/go-gpp/util"
)

type USPRIHeader struct {
	SectionID   byte
	Version     byte
	SubSections *util.IntRange
}

type USPRICoreSegment struct {
	MspaVersion                     byte
	MspaCoveredTransaction          byte
	MspaMode                        byte
	ProcessingNotice                byte
	SaleOptOutNotice                byte
	TargetedAdvertisingOptOutNotice byte
	SaleOptOut                      byte
	TargetedAdvertisingOptOut       byte
	KnownChildSensitiveDataConsents byte
	AdditionalDataProcessingConsent byte
}

type USPRISensitiveDataConsentsSegment struct {
	SensitiveDataProcessing []byte
}

type USPRI struct {
	SectionID                    constants.SectionID
	Value                        string
	Header                       USPRIHeader
	CoreSegment                  USPRICoreSegment
	SensitiveDataConsentsSegment USPRISensitiveDataConsentsSegment
}

func NewUSPRIHeader(bs *util.BitStream) (USPRIHeader, error) {
	var header USPRIHeader
	var err error

	header.SectionID, err = bs.ReadByte6()
	if err != nil {
		return header, sections.ErrorHelper("Header.SectionID", err)
	}

	header.Version, err = bs.ReadByte6()
	if err != nil {
		return header, sections.ErrorHelper("Header.Version", err)
	}

	encodedSubSections, err := bs.ReadFibonacciRange()
	if err != nil {
		return header, sections.ErrorHelper("Header.SubSections", err)
	}

	// Apply -1 offset to convert from Fibonacci encoding (1-indexed) to logical subsection IDs (0-indexed)
	// Fibonacci cannot encode 0, so subsection IDs are stored offset by +1
	header.SubSections = &util.IntRange{
		Size:  encodedSubSections.Size,
		Range: make([]util.IRange, len(encodedSubSections.Range)),
		Max:   encodedSubSections.Max - 1,
	}
	for i, r := range encodedSubSections.Range {
		header.SubSections.Range[i] = util.IRange{
			StartID: r.StartID - 1,
			EndID:   r.EndID - 1,
		}
	}

	return header, nil
}

func NewUSPRICoreSegment(bs *util.BitStream) (USPRICoreSegment, error) {
	var uspriCore USPRICoreSegment
	var err error

	uspriCore.MspaVersion, err = bs.ReadByte6()
	if err != nil {
		return uspriCore, sections.ErrorHelper("CoreSegment.MspaVersion", err)
	}

	uspriCore.MspaCoveredTransaction, err = bs.ReadByte2()
	if err != nil {
		return uspriCore, sections.ErrorHelper("CoreSegment.MspaCoveredTransaction", err)
	}

	uspriCore.MspaMode, err = bs.ReadByte2()
	if err != nil {
		return uspriCore, sections.ErrorHelper("CoreSegment.MspaMode", err)
	}

	uspriCore.ProcessingNotice, err = bs.ReadByte2()
	if err != nil {
		return uspriCore, sections.ErrorHelper("CoreSegment.ProcessingNotice", err)
	}

	uspriCore.SaleOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return uspriCore, sections.ErrorHelper("CoreSegment.SalesOptOutNotice", err)
	}

	uspriCore.TargetedAdvertisingOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return uspriCore, sections.ErrorHelper("CoreSegment.TargetedAdvertisingOptOutNotice", err)
	}

	uspriCore.SaleOptOut, err = bs.ReadByte2()
	if err != nil {
		return uspriCore, sections.ErrorHelper("CoreSegment.SaleOptOut", err)
	}

	uspriCore.TargetedAdvertisingOptOut, err = bs.ReadByte2()
	if err != nil {
		return uspriCore, sections.ErrorHelper("CoreSegment.TargetedAdvertisingOptOut", err)
	}

	uspriCore.KnownChildSensitiveDataConsents, err = bs.ReadByte2()
	if err != nil {
		return uspriCore, sections.ErrorHelper("CoreSegment.KnownChildSensitiveDataConsents", err)
	}

	uspriCore.AdditionalDataProcessingConsent, err = bs.ReadByte2()
	if err != nil {
		return uspriCore, sections.ErrorHelper("CoreSegment.AdditionalDataProcessingConsent", err)
	}

	return uspriCore, err
}

func NewUSPRISensitiveDataConsentsSegment(bs *util.BitStream) (USPRISensitiveDataConsentsSegment, error) {
	var segment USPRISensitiveDataConsentsSegment
	var err error

	segment.SensitiveDataProcessing, err = bs.ReadTwoBitField(8)
	if err != nil {
		return segment, sections.ErrorHelper("SensitiveDataConsentsSegment.SensitiveDataProcessing", err)
	}

	return segment, nil
}

func (header USPRIHeader) Encode(bs *util.BitStream) error {
	bs.WriteByte6(header.SectionID)
	bs.WriteByte6(header.Version)

	// Apply +1 offset to convert from logical subsection IDs (0-indexed) to Fibonacci encoding (1-indexed)
	// Fibonacci cannot encode 0, so we offset subsection IDs by +1 when encoding
	encodedSubSections := &util.IntRange{
		Size:  header.SubSections.Size,
		Range: make([]util.IRange, len(header.SubSections.Range)),
		Max:   header.SubSections.Max + 1,
	}
	for i, r := range header.SubSections.Range {
		encodedSubSections.Range[i] = util.IRange{
			StartID: r.StartID + 1,
			EndID:   r.EndID + 1,
		}
	}

	return bs.WriteIntRange(encodedSubSections)
}

func (segment USPRICoreSegment) Encode(bs *util.BitStream) {
	bs.WriteByte6(segment.MspaVersion)
	bs.WriteByte2(segment.MspaCoveredTransaction)
	bs.WriteByte2(segment.MspaMode)
	bs.WriteByte2(segment.ProcessingNotice)
	bs.WriteByte2(segment.SaleOptOutNotice)
	bs.WriteByte2(segment.TargetedAdvertisingOptOutNotice)
	bs.WriteByte2(segment.SaleOptOut)
	bs.WriteByte2(segment.TargetedAdvertisingOptOut)
	bs.WriteByte2(segment.KnownChildSensitiveDataConsents)
	bs.WriteByte2(segment.AdditionalDataProcessingConsent)
}

func (segment USPRISensitiveDataConsentsSegment) Encode(bs *util.BitStream) {
	bs.WriteTwoBitField(segment.SensitiveDataProcessing)
}

func NewUSPRI(encoded string) (USPRI, error) {
	uspri := USPRI{}

	coreBitStream, _, err := sections.CreateBitStreams(encoded, false)
	if err != nil {
		return uspri, err
	}

	header, err := NewUSPRIHeader(coreBitStream)
	if err != nil {
		return uspri, err
	}

	coreSegment, err := NewUSPRICoreSegment(coreBitStream)
	if err != nil {
		return uspri, err
	}

	sensitiveDataConsentsSegment, err := NewUSPRISensitiveDataConsentsSegment(coreBitStream)
	if err != nil {
		return uspri, err
	}

	uspri = USPRI{
		SectionID:                    constants.SectionUSPRI,
		Value:                        encoded,
		Header:                       header,
		CoreSegment:                  coreSegment,
		SensitiveDataConsentsSegment: sensitiveDataConsentsSegment,
	}

	return uspri, nil
}

func (uspri USPRI) Encode(_ bool) []byte {
	bs := util.NewBitStreamForWrite()
	_ = uspri.Header.Encode(bs)
	uspri.CoreSegment.Encode(bs)
	uspri.SensitiveDataConsentsSegment.Encode(bs)
	return bs.Base64Encode()
}

func (uspri USPRI) GetID() constants.SectionID {
	return uspri.SectionID
}

func (uspri USPRI) GetValue() string {
	return uspri.Value
}