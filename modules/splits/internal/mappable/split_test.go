// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package mappable

import (
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/AssetMantle/modules/modules/splits/internal/key"
	baseIDs "github.com/AssetMantle/modules/schema/ids/base"
)

func Test_Split_Methods(t *testing.T) {
	ownerID := baseIDs.NewID("ownerID")
	ownableID := baseIDs.NewID("ownableID")

	testSplitID := key.NewSplitID(ownerID, ownableID)
	testValue := sdkTypes.NewDec(12)
	testSplit := NewSplit(testSplitID, testValue).(split)

	require.Equal(t, split{ID: testSplitID, Value: testValue}, testSplit)
	require.Equal(t, testSplitID, testSplit.GetID())
	require.Equal(t, ownerID, testSplit.GetOwnerID())
	require.Equal(t, ownableID, testSplit.GetOwnableID())
	require.Equal(t, testValue, testSplit.GetValue())
	require.Equal(t, NewSplit(testSplitID, sdkTypes.NewDec(11)).(split), testSplit.Send(sdkTypes.NewDec(1)))
	require.Equal(t, NewSplit(testSplitID, sdkTypes.NewDec(13)).(split), testSplit.Receive(sdkTypes.NewDec(1)))
	require.Equal(t, true, testSplit.CanSend(sdkTypes.NewDec(5)))
	require.Equal(t, false, testSplit.CanSend(sdkTypes.NewDec(15)))
	require.Equal(t, testSplitID, testSplit.GetKey())
}
