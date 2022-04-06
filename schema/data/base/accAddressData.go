// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package base

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/AssetMantle/modules/constants/errors"
	"github.com/AssetMantle/modules/schema/data"
	"github.com/AssetMantle/modules/schema/types"
	"github.com/AssetMantle/modules/schema/types/base"
	"github.com/AssetMantle/modules/utilities/meta"

	"bytes"
)

type accAddressData struct {
	Value sdkTypes.AccAddress `json:"value"`
}

var _ data.AccAddressData = (*accAddressData)(nil)

func (accAddressData accAddressData) GetID() types.ID {
	return base.NewDataID(accAddressData)
}
func (accAddressData accAddressData) Compare(sortable types.Data) int {
	compareAccAddressData, err := accAddressDataFromInterface(sortable)
	if err != nil {
		panic(err)
	}

	return bytes.Compare(accAddressData.Value.Bytes(), compareAccAddressData.Value.Bytes())
}
func (accAddressData accAddressData) String() string {
	return accAddressData.Value.String()
}
func (accAddressData accAddressData) GetTypeID() types.ID {
	return AccAddressDataID
}
func (accAddressData accAddressData) ZeroValue() types.Data {
	return NewAccAddressData(sdkTypes.AccAddress{})
}
func (accAddressData accAddressData) GenerateHashID() types.ID {
	if accAddressData.Compare(accAddressData.ZeroValue()) == 0 {
		return base.NewID("")
	}

	return base.NewID(meta.Hash(accAddressData.Value.String()))
}
func (accAddressData accAddressData) Get() sdkTypes.AccAddress {
	return accAddressData.Value
}

func accAddressDataFromInterface(data types.Data) (accAddressData, error) {
	switch value := data.(type) {
	case accAddressData:
		return value, nil
	default:
		return accAddressData{}, errors.MetaDataError
	}
}

func NewAccAddressData(value sdkTypes.AccAddress) types.Data {
	return accAddressData{
		Value: value,
	}
}

func ReadAccAddressData(dataString string) (types.Data, error) {
	if dataString == "" {
		return accAddressData{}.ZeroValue(), nil
	}

	accAddress, err := sdkTypes.AccAddressFromBech32(dataString)
	if err != nil {
		return accAddressData{}.ZeroValue(), err
	}

	return NewAccAddressData(accAddress), nil
}
