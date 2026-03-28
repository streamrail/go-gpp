package uspmd

import (
	"github.com/streamrail/go-gpp/constants"
	"github.com/streamrail/go-gpp/sections"
	"github.com/streamrail/go-gpp/util"
)

type USPMDHeader struct {
	SectionID   byte
	Version     byte
	SubSections *util.IntRange
}

type USPMDCoreSegment struct {
	MspaVersion                     byte
	MspaCoveredTransaction          byte
	MspaMode                        byte
	ProcessingNotice                byte
	SaleOptOutNotice                byte
	TargetedAdvertisingOptOutNotice byte
	SaleOptOut                      byte
	TargetedAdvertisingOptOut       byte
	AdditionalDataProcessingConsent byte
}

type USPMD struct {
	SectionID   constants.SectionID
	Value       string
	Header      USPMDHeader
	CoreSegment USPMDCoreSegment
	GPCSegment  sections.CommonUSGPCSegment
}

func NewUSPMDHeader(bs *util.BitStream) (USPMDHeader, error) {
	var header USPMDHeader
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

func NewUSPMDCoreSegment(bs *util.BitStream) (USPMDCoreSegment, error) {
	var uspmdCore USPMDCoreSegment
	var err error

	uspmdCore.MspaVersion, err = bs.ReadByte6()
	if err != nil {
		return uspmdCore, sections.ErrorHelper("CoreSegment.MspaVersion", err)
	}

	uspmdCore.MspaCoveredTransaction, err = bs.ReadByte2()
	if err != nil {
		return uspmdCore, sections.ErrorHelper("CoreSegment.MspaCoveredTransaction", err)
	}

	uspmdCore.MspaMode, err = bs.ReadByte2()
	if err != nil {
		return uspmdCore, sections.ErrorHelper("CoreSegment.MspaMode", err)
	}

	uspmdCore.ProcessingNotice, err = bs.ReadByte2()
	if err != nil {
		return uspmdCore, sections.ErrorHelper("CoreSegment.ProcessingNotice", err)
	}

	uspmdCore.SaleOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return uspmdCore, sections.ErrorHelper("CoreSegment.SalesOptOutNotice", err)
	}

	uspmdCore.TargetedAdvertisingOptOutNotice, err = bs.ReadByte2()
	if err != nil {
		return uspmdCore, sections.ErrorHelper("CoreSegment.TargetedAdvertisingOptOutNotice", err)
	}

	uspmdCore.SaleOptOut, err = bs.ReadByte2()
	if err != nil {
		return uspmdCore, sections.ErrorHelper("CoreSegment.SaleOptOut", err)
	}

	uspmdCore.TargetedAdvertisingOptOut, err = bs.ReadByte2()
	if err != nil {
		return uspmdCore, sections.ErrorHelper("CoreSegment.TargetedAdvertisingOptOut", err)
	}

	uspmdCore.AdditionalDataProcessingConsent, err = bs.ReadByte2()
	if err != nil {
		return uspmdCore, sections.ErrorHelper("CoreSegment.AdditionalDataProcessingConsent", err)
	}

	return uspmdCore, err
}

func (header USPMDHeader) Encode(bs *util.BitStream) error {
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

func (segment USPMDCoreSegment) Encode(bs *util.BitStream) {
	bs.WriteByte6(segment.MspaVersion)
	bs.WriteByte2(segment.MspaCoveredTransaction)
	bs.WriteByte2(segment.MspaMode)
	bs.WriteByte2(segment.ProcessingNotice)
	bs.WriteByte2(segment.SaleOptOutNotice)
	bs.WriteByte2(segment.TargetedAdvertisingOptOutNotice)
	bs.WriteByte2(segment.SaleOptOut)
	bs.WriteByte2(segment.TargetedAdvertisingOptOut)
	bs.WriteByte2(segment.AdditionalDataProcessingConsent)
}

func NewUSPMD(encoded string) (USPMD, error) {
	uspmd := USPMD{}

	coreBitStream, gpcBitStream, err := sections.CreateBitStreams(encoded, true)
	if err != nil {
		return uspmd, err
	}

	header, err := NewUSPMDHeader(coreBitStream)
	if err != nil {
		return uspmd, err
	}

	coreSegment, err := NewUSPMDCoreSegment(coreBitStream)
	if err != nil {
		return uspmd, err
	}

	gpcSegment := sections.CommonUSGPCSegment{
		SubsectionType: 1,
		Gpc:            false,
	}

	if gpcBitStream != nil {
		gpcSegment, err = sections.NewCommonUSGPCSegment(gpcBitStream)
		if err != nil {
			return uspmd, err
		}
	}

	uspmd = USPMD{
		SectionID:   constants.SectionUSPMD,
		Value:       encoded,
		Header:      header,
		CoreSegment: coreSegment,
		GPCSegment:  gpcSegment,
	}

	return uspmd, nil
}

func (uspmd USPMD) Encode(gpcIncluded bool) []byte {
	bs := util.NewBitStreamForWrite()
	_ = uspmd.Header.Encode(bs)
	uspmd.CoreSegment.Encode(bs)
	res := bs.Base64Encode()
	if !gpcIncluded {
		return res
	}
	bs.Reset()
	res = append(res, '.')
	uspmd.GPCSegment.Encode(bs)
	return append(res, bs.Base64Encode()...)
}

func (uspmd USPMD) GetID() constants.SectionID {
	return uspmd.SectionID
}

func (uspmd USPMD) GetValue() string {
	return uspmd.Value
}
