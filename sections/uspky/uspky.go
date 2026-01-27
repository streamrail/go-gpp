package uspky

import (
	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/sections"
	"github.com/streamrail/go-gpp/util"
)

type USPKYHeader struct {
	SectionID   byte
	Version     byte
	SubSections *util.IntRange
}

type USPKYCoreSegment struct {
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

type USPKYSensitiveDataConsentsSegment struct {
	SensitiveDataProcessing []byte
}

type USPKY struct {
	SectionID                    constants.SectionID
	Value                        string
	Header                       USPKYHeader
	CoreSegment                  USPKYCoreSegment
	SensitiveDataConsentsSegment USPKYSensitiveDataConsentsSegment
}

func NewUSPKYHeader(bs *util.BitStream) (USPKYHeader, error) {
	var header USPKYHeader
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

func NewUSPKYCoreSegment(bs *util.BitStream) (USPKYCoreSegment, error) {
	var uspkyCore USPKYCoreSegment
	var err error

	uspkyCore.MspaVersion, err = bs.ReadByte6()
	if err != nil {
		return uspkyCore, sections.ErrorHelper("CoreSegment.MspaVersion", err)
	}

	uspkyCore.MspaCoveredTransaction, err = bs.ReadByte2()
	if err != nil {
		return uspkyCore, sections.ErrorHelper("CoreSegment.MspaCoveredTransaction", err)
	}

	uspkyCore.MspaMode, err = bs.ReadByte2()
	if err != nil {
		return uspkyCore, sections.ErrorHelper("CoreSegment.MspaMode", err)
	}

	uspkyCore.ProcessingNotice, err = bs.ReadByte2()
	if err != nil {
		return uspkyCore, sections.ErrorHelper("CoreSegment.ProcessingNotice", err)
	}

	uspkyCore.SaleOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return uspkyCore, sections.ErrorHelper("CoreSegment.SalesOptOutNotice", err)
	}

	uspkyCore.TargetedAdvertisingOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return uspkyCore, sections.ErrorHelper("CoreSegment.TargetedAdvertisingOptOutNotice", err)
	}

	uspkyCore.SaleOptOut, err = bs.ReadByte2()
	if err != nil {
		return uspkyCore, sections.ErrorHelper("CoreSegment.SaleOptOut", err)
	}

	uspkyCore.TargetedAdvertisingOptOut, err = bs.ReadByte2()
	if err != nil {
		return uspkyCore, sections.ErrorHelper("CoreSegment.TargetedAdvertisingOptOut", err)
	}

	uspkyCore.KnownChildSensitiveDataConsents, err = bs.ReadByte2()
	if err != nil {
		return uspkyCore, sections.ErrorHelper("CoreSegment.KnownChildSensitiveDataConsents", err)
	}

	uspkyCore.AdditionalDataProcessingConsent, err = bs.ReadByte2()
	if err != nil {
		return uspkyCore, sections.ErrorHelper("CoreSegment.AdditionalDataProcessingConsent", err)
	}

	return uspkyCore, err
}

func NewUSPKYSensitiveDataConsentsSegment(bs *util.BitStream) (USPKYSensitiveDataConsentsSegment, error) {
	var segment USPKYSensitiveDataConsentsSegment
	var err error

	segment.SensitiveDataProcessing, err = bs.ReadTwoBitField(8)
	if err != nil {
		return segment, sections.ErrorHelper("SensitiveDataConsentsSegment.SensitiveDataProcessing", err)
	}

	return segment, nil
}

func (header USPKYHeader) Encode(bs *util.BitStream) error {
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

func (segment USPKYCoreSegment) Encode(bs *util.BitStream) {
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

func (segment USPKYSensitiveDataConsentsSegment) Encode(bs *util.BitStream) {
	bs.WriteTwoBitField(segment.SensitiveDataProcessing)
}

func NewUSPKY(encoded string) (USPKY, error) {
	uspky := USPKY{}

	coreBitStream, _, err := sections.CreateBitStreams(encoded, false)
	if err != nil {
		return uspky, err
	}

	header, err := NewUSPKYHeader(coreBitStream)
	if err != nil {
		return uspky, err
	}

	coreSegment, err := NewUSPKYCoreSegment(coreBitStream)
	if err != nil {
		return uspky, err
	}

	sensitiveDataConsentsSegment, err := NewUSPKYSensitiveDataConsentsSegment(coreBitStream)
	if err != nil {
		return uspky, err
	}

	uspky = USPKY{
		SectionID:                    constants.SectionUSPKY,
		Value:                        encoded,
		Header:                       header,
		CoreSegment:                  coreSegment,
		SensitiveDataConsentsSegment: sensitiveDataConsentsSegment,
	}

	return uspky, nil
}

func (uspky USPKY) Encode(_ bool) []byte {
	bs := util.NewBitStreamForWrite()
	_ = uspky.Header.Encode(bs)
	uspky.CoreSegment.Encode(bs)
	uspky.SensitiveDataConsentsSegment.Encode(bs)
	return bs.Base64Encode()
}

func (uspky USPKY) GetID() constants.SectionID {
	return uspky.SectionID
}

func (uspky USPKY) GetValue() string {
	return uspky.Value
}