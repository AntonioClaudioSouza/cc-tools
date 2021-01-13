package transactions

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
)

func TestMain(m *testing.M) {
	log.Println("begin TestMain")

	err := assets.CustomDataTypes(testCustomDataTypes)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	assets.InitAssetList(testAssetList)

	InitTxList(testTxList)

	os.Exit(m.Run())
}

func TestStartUp(t *testing.T) {
	err := StartupCheck()
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
}

var testTxList = []Transaction{}

var testAssetList = []assets.AssetType{
	{
		Tag:         "samplePerson",
		Label:       "Sample Person",
		Description: "",

		Props: []assets.AssetProp{
			{
				Tag:      "cpf",
				Label:    "CPF",
				DataType: "cpf",
				Writers:  []string{"org2MSP"},
			},
			{
				Tag:      "name",
				Label:    "Asset Name",
				Required: true,
				IsKey:    true,
				DataType: "string",
				Validate: func(name interface{}) error {
					nameStr := name.(string)
					if nameStr == "" {
						return fmt.Errorf("name must be non-empty")
					}
					return nil
				},
			},
			{
				Tag:          "readerScore",
				Label:        "Reader Score",
				DefaultValue: 0.0,
				DataType:     "number",
				Writers:      []string{`$org\dMSP`},
			},
			{
				Tag:          "active",
				Label:        "Active",
				DefaultValue: false,
				DataType:     "boolean",
			},
		},
	},
	{
		Tag:         "sampleBook",
		Label:       "Sample Book",
		Description: "",

		Props: []assets.AssetProp{
			{
				Tag:      "title",
				Label:    "Book Title",
				Required: true,
				IsKey:    true,
				DataType: "string",
			},
			{
				Tag:      "author",
				Label:    "Book Author",
				Required: true,
				IsKey:    true,
				DataType: "string",
			},
			{
				Tag:      "currentTenant",
				Label:    "Current Tenant",
				DataType: "->samplePerson",
			},
			{
				Tag:      "genres",
				Label:    "Genres",
				DataType: "[]string",
			},
			{
				Tag:      "published",
				Label:    "Publishment Date",
				DataType: "datetime",
			},
		},
	},
	{
		Tag:         "sampleSecret",
		Label:       "Sample Secret",
		Description: "",

		Readers: []string{"org1MSP"},
		Props: []assets.AssetProp{
			{
				Tag:      "secretName",
				Label:    "Secret Name",
				IsKey:    true,
				DataType: "string",
			},
			{
				Tag:      "secret",
				Label:    "Secret",
				Required: true,
				DataType: "string",
			},
		},
	},
}

var testCustomDataTypes = map[string]assets.DataType{
	"cpf": {
		Parse: func(data interface{}) (string, interface{}, errors.ICCError) {
			cpf, ok := data.(string)
			if !ok {
				return "", nil, errors.NewCCError("property must be a string", 400)
			}

			cpf = strings.ReplaceAll(cpf, ".", "")
			cpf = strings.ReplaceAll(cpf, "-", "")

			if len(cpf) != 11 {
				return "", nil, errors.NewCCError("CPF must have 11 digits", 400)
			}

			return cpf, cpf, nil
		},
	},
}
