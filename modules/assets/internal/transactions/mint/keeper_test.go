// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package mint

import (
	"encoding/json"
	"fmt"
	"github.com/AssetMantle/modules/modules/assets/internal/parameters"
	"github.com/AssetMantle/modules/modules/classifications/auxiliaries/conform"
	"github.com/AssetMantle/modules/modules/identities/auxiliaries/verify"
	maintainersVerify "github.com/AssetMantle/modules/modules/maintainers/auxiliaries/verify"
	"github.com/AssetMantle/modules/modules/metas/auxiliaries/scrub"
	"github.com/AssetMantle/modules/modules/splits/auxiliaries/mint"
	"github.com/AssetMantle/modules/schema/helpers"
	base2 "github.com/AssetMantle/modules/schema/helpers/base"
	"github.com/AssetMantle/modules/utilities/test/schema/helpers/base"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/stretchr/testify/require"
	tendermintDB "github.com/tendermint/tm-db"
	"reflect"
	"testing"
)

func createTestInput(t *testing.T) (helpers.Mapper, helpers.Auxiliary, helpers.Auxiliary, helpers.Auxiliary, helpers.Auxiliary, helpers.Auxiliary, helpers.Keeper, helpers.Parameters, types.Context) {
	testContext, storeKey, paramsTransientStoreKeys := base.SetupTest(t)
	testMapper := base2.NewMapper(base.KeyPrototype, base.MappablePrototype).Initialize(storeKey)
	paramsStoreKey := types.NewKVStoreKey("testParams")

	paramsKeeper := params.NewKeeper(
		codec.Cdc,
		paramsStoreKey,
		paramsTransientStoreKeys,
	)
	testParameters := parameters.Prototype().Initialize(paramsKeeper.Subspace("test"))

	memDB := tendermintDB.NewMemDB()
	commitMultiStore := store.NewCommitMultiStore(memDB)
	commitMultiStore.MountStoreWithDB(storeKey, types.StoreTypeIAVL, memDB)
	commitMultiStore.MountStoreWithDB(paramsStoreKey, types.StoreTypeIAVL, memDB)
	commitMultiStore.MountStoreWithDB(paramsTransientStoreKeys, types.StoreTypeIAVL, memDB)
	err := commitMultiStore.LoadLatestVersion()
	require.Nil(t, err)

	testConformAuxiliary := conform.AuxiliaryMock.Initialize(testMapper, testParameters)
	testMintAuxiliary := mint.AuxiliaryMock.Initialize(testMapper, testParameters)
	testScrubAuxiliary := scrub.AuxiliaryMock.Initialize(testMapper, testParameters)
	testVerifyAuxiliary := verify.AuxiliaryMock.Initialize(testMapper, testParameters)
	testMaintainersVerifyAuxiliary := maintainersVerify.AuxiliaryMock.Initialize(testMapper, testParameters)
	testKeepers := keeperPrototype().Initialize(testMapper, testParameters, []interface{}{testConformAuxiliary, testMintAuxiliary, testScrubAuxiliary, testVerifyAuxiliary, testMaintainersVerifyAuxiliary})

	return testMapper, testConformAuxiliary, testMintAuxiliary, testScrubAuxiliary, testVerifyAuxiliary, testMaintainersVerifyAuxiliary, testKeepers, testParameters, testContext
}

func Test_keeperPrototype(t *testing.T) {
	tests := []struct {
		name string
		want helpers.TransactionKeeper
	}{
		// TODO: Add test cases.
		{"+ve", transactionKeeper{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := keeperPrototype(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keeperPrototype() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_transactionKeeper_Initialize(t *testing.T) {
	testMapper, testConformAuxiliary, testMintAuxiliary, testScrubAuxiliary, testVerifyAuxiliary, testMaintainersVerifyAuxiliary, _, testParameters, _ := createTestInput(t)
	type fields struct {
		mapper                     helpers.Mapper
		parameters                 helpers.Parameters
		conformAuxiliary           helpers.Auxiliary
		mintAuxiliary              helpers.Auxiliary
		scrubAuxiliary             helpers.Auxiliary
		verifyAuxiliary            helpers.Auxiliary
		maintainersVerifyAuxiliary helpers.Auxiliary
	}
	type args struct {
		mapper      helpers.Mapper
		parameters  helpers.Parameters
		auxiliaries []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   helpers.Keeper
	}{
		// TODO: Add test cases.5
		{"+ve", fields{}, args{}, transactionKeeper{}},
		{"+ve", fields{testMapper, testParameters, testConformAuxiliary, testMintAuxiliary, testScrubAuxiliary, testVerifyAuxiliary, testMaintainersVerifyAuxiliary}, args{testMapper, testParameters, []interface{}{testConformAuxiliary, testMintAuxiliary, testScrubAuxiliary, testVerifyAuxiliary, testMaintainersVerifyAuxiliary}}, transactionKeeper{testMapper, testParameters, testConformAuxiliary, testMintAuxiliary, testScrubAuxiliary, testVerifyAuxiliary, testMaintainersVerifyAuxiliary}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transactionKeeper := transactionKeeper{
				mapper:                     tt.fields.mapper,
				parameters:                 tt.fields.parameters,
				conformAuxiliary:           tt.fields.conformAuxiliary,
				mintAuxiliary:              tt.fields.mintAuxiliary,
				scrubAuxiliary:             tt.fields.scrubAuxiliary,
				verifyAuxiliary:            tt.fields.verifyAuxiliary,
				maintainersVerifyAuxiliary: tt.fields.maintainersVerifyAuxiliary,
			}
			if got := transactionKeeper.Initialize(tt.args.mapper, tt.args.parameters, tt.args.auxiliaries); !compare(got, tt.want) {
				t.Errorf("Initialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_transactionKeeper_Transact(t *testing.T) {
	testMapper, testConformAuxiliary, testMintAuxiliary, testScrubAuxiliary, testVerifyAuxiliary, testMaintainersVerifyAuxiliary, _, testParameters, testContext := createTestInput(t)
	message, err := Transaction.DecodeTransactionRequest(json.RawMessage(`{"BaseReq":{"from":"cosmos1pkkayn066msg6kn33wnl5srhdt3tnu2vzasz9c"},"ID":"id"}`))
	require.Equal(t, nil, err)
	type fields struct {
		mapper                     helpers.Mapper
		parameters                 helpers.Parameters
		conformAuxiliary           helpers.Auxiliary
		mintAuxiliary              helpers.Auxiliary
		scrubAuxiliary             helpers.Auxiliary
		verifyAuxiliary            helpers.Auxiliary
		maintainersVerifyAuxiliary helpers.Auxiliary
	}
	type args struct {
		context types.Context
		msg     types.Msg
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   helpers.TransactionResponse
	}{
		// TODO: Add test cases.
		{"+ve", fields{testMapper, testParameters, testConformAuxiliary, testMintAuxiliary, testScrubAuxiliary, testVerifyAuxiliary, testMaintainersVerifyAuxiliary}, args{testContext, message}, newTransactionResponse(nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transactionKeeper := transactionKeeper{
				mapper:                     tt.fields.mapper,
				parameters:                 tt.fields.parameters,
				conformAuxiliary:           tt.fields.conformAuxiliary,
				mintAuxiliary:              tt.fields.mintAuxiliary,
				scrubAuxiliary:             tt.fields.scrubAuxiliary,
				verifyAuxiliary:            tt.fields.verifyAuxiliary,
				maintainersVerifyAuxiliary: tt.fields.maintainersVerifyAuxiliary,
			}
			if got := transactionKeeper.Transact(tt.args.context, tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Transact() = %v, want %v", got, tt.want)
			}
		})
	}
}

func compare(x, y interface{}) bool {
	if x == nil || y == nil {
		return x == y
	}
	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)
	if v1.Type() != v2.Type() {
		return false
	}
	//fmt.Println(v1.Field(4))
	for i := 0; i < v1.NumField(); i++ {
		if !reflect.DeepEqual(fmt.Sprint(v1.Field(i)), fmt.Sprint(v2.Field(i))) {
			fmt.Println(v1.Field(i), v2.Field(i))
			return false
		}
	}
	return true
}
