package uspin

import (
	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/sections"
	"github.com/streamrail/go-gpp/util"
)

type USPINHeader struct {
	SectionID   byte
	Version     byte
	SubSections *util.IntRange
}

type USPINCoreSegment struct {
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

type USPINSensitiveDataConsentsSegment struct {
	SensitiveDataProcessing []byte
}

type USPIN struct {
	SectionID                    constants.SectionID
	Value                        string
	Header                       USPINHeader
	CoreSegment                  USPINCoreSegment
	SensitiveDataConsentsSegment USPINSensitiveDataConsentsSegment
}

func NewUSPINHeader(bs *util.BitStream) (USPINHeader, error) {
	var header USPINHeader
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

func NewUSPINCoreSegment(bs *util.BitStream) (USPINCoreSegment, error) {
	var uspinCore USPINCoreSegment
	var err error

	uspinCore.MspaVersion, err = bs.ReadByte6()
	if err != nil {
		return uspinCore, sections.ErrorHelper("CoreSegment.MspaVersion", err)
	}

	uspinCore.MspaCoveredTransaction, err = bs.ReadByte2()
	if err != nil {
		return uspinCore, sections.ErrorHelper("CoreSegment.MspaCoveredTransaction", err)
	}

	uspinCore.MspaMode, err = bs.ReadByte2()
	if err != nil {
		return uspinCore, sections.ErrorHelper("CoreSegment.MspaMode", err)
	}

	uspinCore.ProcessingNotice, err = bs.ReadByte2()
	if err != nil {
		return uspinCore, sections.ErrorHelper("CoreSegment.ProcessingNotice", err)
	}

	uspinCore.SaleOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return uspinCore, sections.ErrorHelper("CoreSegment.SalesOptOutNotice", err)
	}

	uspinCore.TargetedAdvertisingOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return uspinCore, sections.ErrorHelper("CoreSegment.TargetedAdvertisingOptOutNotice", err)
	}

	uspinCore.SaleOptOut, err = bs.ReadByte2()
	if err != nil {
		return uspinCore, sections.ErrorHelper("CoreSegment.SaleOptOut", err)
	}

	uspinCore.TargetedAdvertisingOptOut, err = bs.ReadByte2()
	if err != nil {
		return uspinCore, sections.ErrorHelper("CoreSegment.TargetedAdvertisingOptOut", err)
	}

	uspinCore.KnownChildSensitiveDataConsents, err = bs.ReadByte2()
	if err != nil {
		return uspinCore, sections.ErrorHelper("CoreSegment.KnownChildSensitiveDataConsents", err)
	}

	uspinCore.AdditionalDataProcessingConsent, err = bs.ReadByte2()
	if err != nil {
		return uspinCore, sections.ErrorHelper("CoreSegment.AdditionalDataProcessingConsent", err)
	}

	return uspinCore, err
}

func NewUSPINSensitiveDataConsentsSegment(bs *util.BitStream) (USPINSensitiveDataConsentsSegment, error) {
	var segment USPINSensitiveDataConsentsSegment
	var err error

	segment.SensitiveDataProcessing, err = bs.ReadTwoBitField(8)
	if err != nil {
		return segment, sections.ErrorHelper("SensitiveDataConsentsSegment.SensitiveDataProcessing", err)
	}

	return segment, nil
}

func (header USPINHeader) Encode(bs *util.BitStream) error {
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

func (segment USPINCoreSegment) Encode(bs *util.BitStream) {
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

func (segment USPINSensitiveDataConsentsSegment) Encode(bs *util.BitStream) {
	bs.WriteTwoBitField(segment.SensitiveDataProcessing)
}

func NewUSPIN(encoded string) (USPIN, error) {
	uspin := USPIN{}

	coreBitStream, _, err := sections.CreateBitStreams(encoded, false)
	if err != nil {
		return uspin, err
	}

	header, err := NewUSPINHeader(coreBitStream)
	if err != nil {
		return uspin, err
	}

	coreSegment, err := NewUSPINCoreSegment(coreBitStream)
	if err != nil {
		return uspin, err
	}

	sensitiveDataConsentsSegment, err := NewUSPINSensitiveDataConsentsSegment(coreBitStream)
	if err != nil {
		return uspin, err
	}

	uspin = USPIN{
		SectionID:                    constants.SectionUSPIN,
		Value:                        encoded,
		Header:                       header,
		CoreSegment:                  coreSegment,
		SensitiveDataConsentsSegment: sensitiveDataConsentsSegment,
	}

	return uspin, nil
}

func (uspin USPIN) Encode(_ bool) []byte {
	bs := util.NewBitStreamForWrite()
	_ = uspin.Header.Encode(bs)
	uspin.CoreSegment.Encode(bs)
	uspin.SensitiveDataConsentsSegment.Encode(bs)
	return bs.Base64Encode()
}

func (uspin USPIN) GetID() constants.SectionID {
	return uspin.SectionID
}

func (uspin USPIN) GetValue() string {
	return uspin.Value
}
