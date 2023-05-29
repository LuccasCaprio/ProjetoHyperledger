package assettypes

import "github.com/goledgerdev/cc-tools/assets"

var Token = assets.AssetType{
	Tag:         "token",
	Label:       "Token",
	Description: "Token",

	Props: []assets.AssetProp{
		{
			Required: true,
			IsKey:    true,
			Tag:      "id",
			Label:    "ID Token",
			DataType: "string",
			Writers:  []string{`org1MSP`, "orgMSP"},
		},
		{
			Tag:      "proprietario",
			Label:    "Proprietário",
			DataType: "->proprietario",
			Writers:  []string{`org1MSP`, "orgMSP"},
		},
		{
			Tag:          "quantidade",
			Label:        "Quantidade",
			DataType:     "number",
			DefaultValue: 0,
			Writers:      []string{`org1MSP`, "orgMSP"},
		},
		{
			Tag:          "burned",
			Label:        "Burned",
			DataType:     "boolean",
			DefaultValue: false,
			Writers:      []string{`org1MSP`, "orgMSP"},
		},
	},
}
