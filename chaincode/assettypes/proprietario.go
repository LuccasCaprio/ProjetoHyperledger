package assettypes

import (
	"fmt"

	"github.com/goledgerdev/cc-tools/assets"
)

var Proprietario = assets.AssetType{
	Tag:         "proprietario",
	Label:       "Proprietario",
	Description: "Proprietario token",

	Props: []assets.AssetProp{
		{
			Required: true,
			IsKey:    true,
			Tag:      "id",
			Label:    "ID",
			DataType: "string",
			Writers:  []string{`org1MSP`, "orgMSP"},
		},
		{
			Required: true,
			Tag:      "nome",
			Label:    "Nome proprietario",
			DataType: "string",
			Validate: func(name interface{}) error {
				nameStr := name.(string)
				if nameStr == "" {
					return fmt.Errorf("not null")
				}
				return nil
			},
		},
	},
}
