// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package unwrap

import (
	"reflect"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tendermintDB "github.com/tendermint/tm-db"

	"github.com/AssetMantle/modules/constants/errors"
	"github.com/AssetMantle/modules/modules/identities/auxiliaries/authenticate"
	"github.com/AssetMantle/modules/modules/splits/internal/key"
	"github.com/AssetMantle/modules/modules/splits/internal/mappable"
	"github.com/AssetMantle/modules/modules/splits/internal/module"
	"github.com/AssetMantle/modules/modules/splits/internal/parameters"
	"github.com/AssetMantle/modules/schema"
	"github.com/AssetMantle/modules/schema/helpers"
	baseHelpers "github.com/AssetMantle/modules/schema/helpers/base"
	baseIDs "github.com/AssetMantle/modules/schema/ids/base"
)

type TestKeepers struct {
	SplitsKeeper  helpers.TransactionKeeper
	AccountKeeper auth.AccountKeeper
	BankKeeper    bank.Keeper
	SupplyKeeper  supply.Keeper
}

func CreateTestInput(t *testing.T) (sdkTypes.Context, TestKeepers) {
	var Codec = codec.New()
	schema.RegisterCodec(Codec)
	sdkTypes.RegisterCodec(Codec)
	codec.RegisterCrypto(Codec)
	codec.RegisterEvidences(Codec)
	vesting.RegisterCodec(Codec)
	supply.RegisterCodec(Codec)
	params.RegisterCodec(Codec)
	auth.RegisterCodec(Codec)
	Codec.Seal()

	storeKey := sdkTypes.NewKVStoreKey("test")
	paramsStoreKey := sdkTypes.NewKVStoreKey("testParams")
	authStoreKey := sdkTypes.NewKVStoreKey("testAuth")
	supplyStoreKey := sdkTypes.NewKVStoreKey("testSupply")
	paramsTransientStoreKeys := sdkTypes.NewTransientStoreKey("testParamsTransient")
	Mapper := baseHelpers.NewMapper(key.Prototype, mappable.Prototype).Initialize(storeKey)
	paramsKeeper := params.NewKeeper(
		Codec,
		paramsStoreKey,
		paramsTransientStoreKeys,
	)
	Parameters := parameters.Prototype().Initialize(paramsKeeper.Subspace("test"))

	memDB := tendermintDB.NewMemDB()
	commitMultiStore := store.NewCommitMultiStore(memDB)
	commitMultiStore.MountStoreWithDB(storeKey, sdkTypes.StoreTypeIAVL, memDB)
	commitMultiStore.MountStoreWithDB(paramsStoreKey, sdkTypes.StoreTypeIAVL, memDB)
	commitMultiStore.MountStoreWithDB(paramsTransientStoreKeys, sdkTypes.StoreTypeTransient, memDB)
	commitMultiStore.MountStoreWithDB(authStoreKey, sdkTypes.StoreTypeIAVL, memDB)
	commitMultiStore.MountStoreWithDB(supplyStoreKey, sdkTypes.StoreTypeIAVL, memDB)
	err := commitMultiStore.LoadLatestVersion()
	require.Nil(t, err)

	context := sdkTypes.NewContext(commitMultiStore, abciTypes.Header{
		ChainID: "test",
	}, false, log.NewNopLogger())

	accountKeeper := auth.NewAccountKeeper(Codec, authStoreKey, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)

	bankKeeper := bank.NewBaseKeeper(accountKeeper, paramsKeeper.Subspace(bank.DefaultParamspace), make(map[string]bool))
	supplyKeeper := supply.NewKeeper(Codec, supplyStoreKey, accountKeeper, bankKeeper, map[string][]string{module.Name: nil})
	authenticateAuxiliary := authenticate.AuxiliaryMock.Initialize(Mapper, Parameters)
	keepers := TestKeepers{
		SplitsKeeper: keeperPrototype().Initialize(Mapper, Parameters,
			[]interface{}{authenticateAuxiliary, supplyKeeper}).(helpers.TransactionKeeper),
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		SupplyKeeper:  supplyKeeper,
	}

	return context, keepers
}

func Test_transactionKeeper_Transact(t *testing.T) {
	context, keepers := CreateTestInput(t)
	defaultAddr := sdkTypes.AccAddress("addr")
	verifyMockErrorAddress := sdkTypes.AccAddress("verifyError")
	ownableID := baseIDs.NewID("stake")
	fromID := baseIDs.NewID("fromID")
	coins := func(amount int64) sdkTypes.Coins {
		return sdkTypes.NewCoins(sdkTypes.NewCoin("stake", sdkTypes.NewInt(amount)))
	}
	err := keepers.BankKeeper.SetCoins(context, defaultAddr, coins(1000))
	require.Equal(t, nil, err)
	err = keepers.SupplyKeeper.SendCoinsFromAccountToModule(context, defaultAddr, module.Name, coins(1000))
	require.Equal(t, nil, err)
	keepers.SplitsKeeper.(transactionKeeper).mapper.NewCollection(context).Add(mappable.NewSplit(key.NewSplitID(fromID, ownableID), sdkTypes.NewDec(1000)))

	t.Run("PositiveCase- Send All", func(t *testing.T) {
		want := newTransactionResponse(nil)
		if got := keepers.SplitsKeeper.Transact(context, newMessage(defaultAddr, fromID, ownableID, sdkTypes.NewInt(1000))); !reflect.DeepEqual(got, want) {
			t.Errorf("Transact() = %v, want %v", got, want)
		}
	})

	err = keepers.SupplyKeeper.SendCoinsFromAccountToModule(context, defaultAddr, module.Name, coins(1000))
	require.Equal(t, nil, err)
	keepers.SplitsKeeper.(transactionKeeper).mapper.NewCollection(context).Add(mappable.NewSplit(key.NewSplitID(fromID, ownableID), sdkTypes.NewDec(1000)))

	t.Run("PositiveCase", func(t *testing.T) {
		want := newTransactionResponse(nil)
		if got := keepers.SplitsKeeper.Transact(context, newMessage(defaultAddr, fromID, ownableID, sdkTypes.NewInt(10))); !reflect.DeepEqual(got, want) {
			t.Errorf("Transact() = %v, want %v", got, want)
		}
	})

	t.Run("NegativeCase-Verify Identity Failure", func(t *testing.T) {
		t.Parallel()
		want := newTransactionResponse(errors.MockError)
		if got := keepers.SplitsKeeper.Transact(context, newMessage(verifyMockErrorAddress, fromID, ownableID, sdkTypes.NewInt(10))); !reflect.DeepEqual(got, want) {
			t.Errorf("Transact() = %v, want %v", got, want)
		}
	})

	t.Run("NegativeCase-Send Negative Balance", func(t *testing.T) {
		t.Parallel()
		want := newTransactionResponse(errors.NotAuthorized)
		if got := keepers.SplitsKeeper.Transact(context, newMessage(defaultAddr, fromID, ownableID, sdkTypes.NewInt(-10))); !reflect.DeepEqual(got, want) {
			t.Errorf("Transact() = %v, want %v", got, want)
		}
	})

	t.Run("NegativeCase-Send More than own Balance", func(t *testing.T) {
		t.Parallel()
		want := newTransactionResponse(errors.InsufficientBalance)
		if got := keepers.SplitsKeeper.Transact(context, newMessage(defaultAddr, fromID, ownableID, sdkTypes.NewInt(790))); !reflect.DeepEqual(got.IsSuccessful(), want.IsSuccessful()) {
			t.Errorf("Transact() = %v, want %v", got, want)
		}
	})

	t.Run("NegativeCase-Value Not found", func(t *testing.T) {
		t.Parallel()
		want := newTransactionResponse(errors.EntityNotFound)
		if got := keepers.SplitsKeeper.Transact(context, newMessage(defaultAddr, baseIDs.NewID("id"), ownableID, sdkTypes.NewInt(10))); !reflect.DeepEqual(got, want) {
			t.Errorf("Transact() = %v, want %v", got, want)
		}
	})

	err = keepers.SupplyKeeper.SendCoinsFromModuleToAccount(context, module.Name, defaultAddr, coins(900))
	require.Equal(t, nil, err)
	t.Run("NegativeCase-Module does not have enough coins", func(t *testing.T) {
		want := newTransactionResponse(errors.InsufficientBalance)
		if got := keepers.SplitsKeeper.Transact(context, newMessage(defaultAddr, fromID, ownableID, sdkTypes.NewInt(200))); !reflect.DeepEqual(got.IsSuccessful(), want.IsSuccessful()) {
			t.Errorf("Transact() = %v, want %v", got, want)
		}
	})

}
