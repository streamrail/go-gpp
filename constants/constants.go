package constants

type SectionID int8

const (
	SectionTCFEU2       SectionID = 2
	SectionGPP          SectionID = 3
	GPPSectionTCFCanada SectionID = 5
	SectionUSPV1        SectionID = 6
	SectionUSPNAT       SectionID = 7
	SectionUSPCA        SectionID = 8
	SectionUSPVA        SectionID = 9
	SectionUSPCO        SectionID = 10
	SectionUSPUT        SectionID = 11
	SectionUSPCT        SectionID = 12
	SectionUSPMT        SectionID = 14
	SectionUSPOR        SectionID = 15
	SectionUSPTX        SectionID = 16
	SectionUSPDE        SectionID = 17
	SectionUSPIA        SectionID = 18
	SectionUSPNE        SectionID = 19
	SectionUSPNH        SectionID = 20
	SectionUSPNJ        SectionID = 21
	SectionUSPTN        SectionID = 22
	SectionUSPMN        SectionID = 23
)

var SectionNamesByID = map[SectionID]string{
	SectionTCFEU2:       "tcfeu2",
	SectionGPP:          "gpp header",
	GPPSectionTCFCanada: "tcfcav1",
	SectionUSPV1:        "uspv1",
	SectionUSPNAT:       "uspnat",
	SectionUSPCA:        "uspca",
	SectionUSPVA:        "uspva",
	SectionUSPCO:        "uspco",
	SectionUSPUT:        "usput",
	SectionUSPCT:        "uspct",
	SectionUSPMT:        "uspmt",
	SectionUSPOR:        "uspor",
	SectionUSPTX:        "usptx",
	SectionUSPDE:        "uspde",
	SectionUSPIA:        "uspia",
	SectionUSPNE:        "uspne",
	SectionUSPNH:        "uspnh",
	SectionUSPNJ:        "uspnj",
	SectionUSPTN:        "usptn",
	SectionUSPMN:        "uspmn",
}
