// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package lists

import (
	"github.com/AssetMantle/modules/schema/types"
)

type IDList interface {
	Size() int
	GetList() []types.ID
	Search(types.ID) (index int, found bool)
	Add(...types.ID) IDList
	Remove(...types.ID) IDList
}
