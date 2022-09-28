// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package asset

import (
	"fmt"
	"github.com/AssetMantle/modules/modules/assets/internal/key"
	"github.com/AssetMantle/modules/modules/assets/internal/parameters"
	baseData "github.com/AssetMantle/modules/schema/data/base"
	"github.com/AssetMantle/modules/schema/helpers"
	base2 "github.com/AssetMantle/modules/schema/helpers/base"
	baseIDs "github.com/AssetMantle/modules/schema/ids/base"
	base3 "github.com/AssetMantle/modules/schema/lists/base"
	baseProperties "github.com/AssetMantle/modules/schema/properties/base"
	baseQualified "github.com/AssetMantle/modules/schema/qualified/base"
	"github.com/AssetMantle/modules/utilities/test/schema/helpers/base"
	"github.com/cosmos/cosmos-sdk/types"
	"reflect"
	"testing"
)

func createTestInput(t *testing.T) (helpers.Mapper, types.Context, helpers.QueryRequest) {
	testContext, storeKey, _ := base.SetupTest(t)
	testMapper := base2.NewMapper(base.KeyPrototype, base.MappablePrototype).Initialize(storeKey)
	testImmutables := baseQualified.NewImmutables(base3.NewPropertyList(baseProperties.NewMesaProperty(baseIDs.NewStringID("ID2"), baseData.NewStringData("Data2"))))
	testMutables := baseQualified.NewMutables(base3.NewPropertyList(baseProperties.NewMesaProperty(baseIDs.NewStringID("ID1"), baseData.NewStringData("Data1"))))

	testClassificationID := baseIDs.NewClassificationID(testImmutables, testMutables)
	testQueryRequest := newQueryRequest(baseIDs.NewAssetID(testClassificationID, testImmutables))
	return testMapper, testContext, testQueryRequest
}

func Test_keeperPrototype(t *testing.T) {
	tests := []struct {
		name string
		want helpers.QueryKeeper
	}{

		{"+ve", queryKeeper{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := keeperPrototype(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keeperPrototype() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_queryKeeper_Enquire(t *testing.T) {
	testMapper, testContext, testQueryRequest := createTestInput(t)
	type fields struct {
		mapper helpers.Mapper
	}
	type args struct {
		context      types.Context
		queryRequest helpers.QueryRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   helpers.QueryResponse
	}{
		// TODO: Add test cases.
		//{"+ve with nil", fields{}, args{testContext, queryRequest{}}, newQueryResponse(queryKeeper{}.mapper.NewCollection(testContext).Fetch(key.NewKey(queryRequestFromInterface(queryRequest{}).AssetID)), nil)}, //TODO: panics with error: `invalid memory address or nil pointer dereference`
		{"+ve", fields{testMapper}, args{testContext, testQueryRequest}, newQueryResponse(queryKeeper{mapper: testMapper}.mapper.NewCollection(testContext).Fetch(key.NewKey(queryRequestFromInterface(testQueryRequest).AssetID)), nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryKeeper := queryKeeper{
				mapper: tt.fields.mapper,
			}
			if got := queryKeeper.Enquire(tt.args.context, tt.args.queryRequest); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Enquire() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_queryKeeper_Initialize(t *testing.T) {
	testMapper, _, _ := createTestInput(t)
	//testKeeper := queryKeeper{}
	//testKeeper.mapper = testMapper
	type fields struct {
		mapper helpers.Mapper
	}
	type args struct {
		mapper helpers.Mapper
		in1    helpers.Parameters
		in2    []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   helpers.Keeper
	}{
		// TODO: Add test cases.
		{"+ve with nil", fields{}, args{}, queryKeeper{}},
		{"+ve", fields{testMapper}, args{testMapper, parameters.Prototype(), []interface{}{}}, queryKeeper{testMapper}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryKeeper := queryKeeper{
				mapper: tt.fields.mapper,
			}
			if got := queryKeeper.Initialize(tt.args.mapper, tt.args.in1, tt.args.in2); !compare(got, tt.want) {
				t.Errorf("Initialize() = %v, want %v", got, tt.want)
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
	return reflect.DeepEqual(fmt.Sprint(v1.Field(0)), fmt.Sprint(v2.Field(0)))
}
