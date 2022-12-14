// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package key

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/AssetMantle/modules/constants"
	baseData "github.com/AssetMantle/modules/schema/data/base"
	baseIDs "github.com/AssetMantle/modules/schema/ids/base"
	"github.com/AssetMantle/modules/schema/lists/base"
	baseProperties "github.com/AssetMantle/modules/schema/properties/base"
)

func TestFromID(t *testing.T) {
	classificationID := baseIDs.NewID("classificationID")
	immutableProperties := base.NewPropertyList(baseProperties.NewProperty(baseIDs.NewID("ID1"), baseData.NewStringData("ImmutableData")))
	newAssetID := NewAssetID(classificationID, immutableProperties)

	assetID1, err := assetIDFromInterface(newAssetID)
	require.Equal(t, assetID1, FromID(newAssetID))
	require.Equal(t, nil, err)

	id := baseIDs.NewID("")
	testAssetID := assetID{ClassificationID: baseIDs.NewID(""), HashID: baseIDs.NewID("")}
	require.Equal(t, FromID(id), testAssetID)

	testString1 := "string1"
	testString2 := "string2"
	id2 := baseIDs.NewID(testString1 + constants.FirstOrderCompositeIDSeparator + testString2)
	testAssetID2 := assetID{ClassificationID: baseIDs.NewID(testString1), HashID: baseIDs.NewID(testString2)}
	require.Equal(t, FromID(id2), testAssetID2)
}

func TestReadClassificationID(t *testing.T) {
	classificationID := baseIDs.NewID("classificationID")
	immutableProperties := base.NewPropertyList(baseProperties.NewProperty(baseIDs.NewID("ID1"), baseData.NewStringData("ImmutableData")))
	assetID1 := NewAssetID(classificationID, immutableProperties)

	assetID2, err := assetIDFromInterface(assetID1)
	require.Equal(t, assetID2.ClassificationID, ReadClassificationID(assetID1))
	require.Equal(t, nil, err)
}
