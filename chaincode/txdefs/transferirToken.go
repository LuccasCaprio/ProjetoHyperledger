package txdefs

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
)

var TransferirToken = tx.Transaction{
	Tag:         "transferToken",
	Label:       "Transfer Token",
	Description: "Transferir ativo token",
	Method:      "PUT",
	Callers:     []string{"$org1MSP", "$orgMSP"},

	Args: []tx.Argument{
		{
			Tag:         "token",
			Label:       "Token Origem",
			Description: "Token Origem",
			DataType:    "->token",
			Required:    true,
		},
		{
			Tag:         "destino",
			Label:       "Destino",
			Description: "Destino",
			DataType:    "->proprietario",
			Required:    true,
		},
		{
			Tag:         "quantidade",
			Label:       "Quantidade Transferida",
			Description: "Quantidade Transferida",
			DataType:    "number",
			Required:    true,
		},
		{
			Tag:         "id",
			Label:       "ID Novo Token",
			Description: "ID Novo Token",
			DataType:    "string",
			Required:    true,
		},
		{
			Tag:         "novoId",
			Label:       "Novo ID Token Origem",
			Description: "Novo ID Token Origem",
			DataType:    "string",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		id, _ := req["id"].(string)
		novoId, _ := req["novoId"].(string)
		quantidade, _ := req["quantidade"].(float64)

		if quantidade <= 0 {
			return nil, errors.WrapError(nil, "A quantidade deve ser maior que zero.")
		}

		tokenKey, ok := req["token"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parâmetro 'token' deve ser um asset.")
		}

		tokenAsset, err := tokenKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Falha ao obter ativo 'token'.")
		}
		tokenMap := (map[string]interface{})(*tokenAsset)

		if tokenMap["burned"].(bool) {
			return nil, errors.WrapError(err, "O token selecionado já foi queimado.")
		}

		tokenMap["burned"] = true

		tokenMap, err = tokenAsset.Update(stub, tokenMap)
		if err != nil {
			return nil, errors.WrapError(err, "Falha ao atualizar ativo 'token'.")
		}

		novaQuantidade := tokenMap["quantidade"].(float64) - quantidade

		if novaQuantidade < 0 {
			return nil, errors.WrapError(err, "Saldo de token insuficiente.")
		}

		novoTokenOrigemMap := make(map[string]interface{})
		novoTokenOrigemMap["@assetType"] = "token"
		novoTokenOrigemMap["id"] = novoId
		novoTokenOrigemMap["proprietario"] = tokenMap["proprietario"]
		novoTokenOrigemMap["quantidade"] = novaQuantidade

		novoTokenOrigemAsset, err := assets.NewAsset(novoTokenOrigemMap)
		if err != nil {
			return nil, errors.WrapError(err, "Falha ao criar ativo 'novo token de origem'.")
		}

		_, err = novoTokenOrigemAsset.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Erro ao salvar ativo 'novo token de origem' na blochchain.")
		}

		proprietarioKey, ok := req["destino"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parametro 'destino' deve ser um ativo.")
		}

		proprietarioAsset, err := proprietarioKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Falha ao obter ativo 'destino'.")
		}
		proprietarioMap := (map[string]interface{})(*proprietarioAsset)

		updatedProprietarioKey := make(map[string]interface{})
		updatedProprietarioKey["@assetType"] = "proprietario"
		updatedProprietarioKey["@key"] = proprietarioMap["@key"]

		novoTokenMap := make(map[string]interface{})
		novoTokenMap["@assetType"] = "token"
		novoTokenMap["id"] = id
		novoTokenMap["proprietario"] = updatedProprietarioKey
		novoTokenMap["quantidade"] = quantidade

		novoTokenAsset, err := assets.NewAsset(novoTokenMap)
		if err != nil {
			return nil, errors.WrapError(err, "Falha ao criar ativo 'token de destino'.")
		}

		_, err = novoTokenAsset.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Erro ao salvar ativo 'token de destino' na blockchain.")
		}

		tokenJSON, nerr := json.Marshal(tokenAsset)
		if nerr != nil {
			return nil, errors.WrapError(nil, "Falha ao converter ativo para JSON.")
		}

		return tokenJSON, nil
	},
}
