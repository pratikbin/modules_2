// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"strconv"

	"github.com/AssetMantle/modules/constants/errors"
	"github.com/AssetMantle/modules/schema/data"
	"github.com/AssetMantle/modules/schema/types"
	"github.com/AssetMantle/modules/schema/types/base"
	"github.com/AssetMantle/modules/utilities/meta"
)

type heightData struct {
	Value types.Height `json:"value"`
}

var _ data.HeightData = (*heightData)(nil)

func (heightData heightData) GetID() types.ID {
	return base.NewDataID(heightData)
}
func (heightData heightData) Compare(data types.Data) int {
	compareHeightData, err := heightDataFromInterface(data)
	if err != nil {
		panic(err)
	}

	return heightData.Value.Compare(compareHeightData.Value)
}
func (heightData heightData) String() string {
	return strconv.FormatInt(heightData.Value.Get(), 10)
}
func (heightData heightData) GetTypeID() types.ID {
	return HeightDataID
}
func (heightData heightData) ZeroValue() types.Data {
	return NewHeightData(base.NewHeight(0))
}
func (heightData heightData) GenerateHashID() types.ID {
	if heightData.Compare(heightData.ZeroValue()) == 0 {
		return base.NewID("")
	}

	return base.NewID(meta.Hash(strconv.FormatInt(heightData.Value.Get(), 10)))
}
func (heightData heightData) Get() types.Height {
	return heightData.Value
}

func heightDataFromInterface(data types.Data) (heightData, error) {
	switch value := data.(type) {
	case heightData:
		return value, nil
	default:
		return heightData{}, errors.MetaDataError
	}
}

func NewHeightData(value types.Height) types.Data {
	return heightData{
		Value: value,
	}
}

func ReadHeightData(dataString string) (types.Data, error) {
	if dataString == "" {
		return heightData{}.ZeroValue(), nil
	}

	height, err := strconv.ParseInt(dataString, 10, 64)
	if err != nil {
		return nil, err
	}

	return NewHeightData(base.NewHeight(height)), nil
}
