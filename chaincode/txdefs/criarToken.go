package txdefs

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
)

var CriarToken = tx.Transaction{
	Tag:         "criarToken",
	Label:       "Criar Token",
	Description: "Criar Token",
	Method:      "POST",
	Callers:     []string{"$org1MSP", "$orgMSP"},

	Args: []tx.Argument{
		{
			Tag:         "id",
			Label:       "ID",
			Description: "ID",
			DataType:    "string",
			Required:    true,
		},
		{
			Tag:         "proprietario",
			Label:       "Proprietario",
			Description: "proprietario",
			DataType:    "->proprietario",
			Required:    true,
		},
		{
			Tag:         "quantidade",
			Label:       "Quantidade",
			Description: "Quantidade",
			DataType:    "number",
			Required:    true,
		},
		{
			Tag:         "burned",
			Label:       "Burned",
			Description: "Queima",
			DataType:    "boolean",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		proprietarioKey, ok := req["proprietario"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parametro proprietario deve ser um ativo.")
		}

		proprietarioAsset, err := proprietarioKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Falha ao obter ativo 'proprietário'.")
		}
		proprietarioMap := (map[string]interface{})(*proprietarioAsset)

		updatedProprietarioKey := make(map[string]interface{})
		updatedProprietarioKey["@assetType"] = "proprietario"
		updatedProprietarioKey["@key"] = proprietarioMap["@key"]

		id, _ := req["id"].(string)
		quantidade, _ := req["quantidade"].(float64)
		burned, _ := req["burned"].(bool)

		if quantidade <= 0 {
			return nil, errors.WrapError(nil, "Quantidade deve ser maior que 0")
		}

		tokenMap := make(map[string]interface{})
		tokenMap["@assetType"] = "token"
		tokenMap["id"] = id
		tokenMap["proprietario"] = updatedProprietarioKey
		tokenMap["quantidade"] = quantidade
		tokenMap["burned"] = burned

		tokenAsset, err := assets.NewAsset(tokenMap)
		if err != nil {
			return nil, errors.WrapError(err, "Falha ao criar ativo 'token'.")
		}

		_, err = tokenAsset.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Erro ao salvar ativo na blockchain.")
		}

		tokenJSON, nerr := json.Marshal(tokenAsset)
		if nerr != nil {
			return nil, errors.WrapError(nil, "Falha ao converter ativo para JSON.")
		}

		return tokenJSON, nil
	},
}
