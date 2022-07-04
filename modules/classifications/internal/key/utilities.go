// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package key

import (
	errorConstants "github.com/AssetMantle/modules/schema/errors/constants"
	"github.com/AssetMantle/modules/schema/helpers"
	"github.com/AssetMantle/modules/schema/ids"
	baseIDs "github.com/AssetMantle/modules/schema/ids/base"
	stringUtilities "github.com/AssetMantle/modules/utilities/string"
)

func readClassificationID(classificationIDString string) ids.ID {
	idList := stringUtilities.SplitCompositeIDString(classificationIDString)
	if len(idList) == 2 {
		return classificationID{
			ChainID: baseIDs.NewID(idList[0]),
			HashID:  baseIDs.NewID(idList[1]),
		}
	}

	return classificationID{ChainID: baseIDs.NewID(""), HashID: baseIDs.NewID("")}
}
func classificationIDFromInterface(i interface{}) (classificationID, error) {
	switch value := i.(type) {
	case classificationID:
		return value, nil
	case ids.ID:
		// TODO remove this use case
		return classificationIDFromInterface(readClassificationID(value.String()))
	default:
		return classificationID{}, errorConstants.MetaDataError
	}
}

func FromID(id ids.ID) helpers.Key {
	if classificationID, err := classificationIDFromInterface(id); err != nil {
		// TODO plug all panic scenarios
		panic(err)
	} else {
		return classificationID
	}
}
