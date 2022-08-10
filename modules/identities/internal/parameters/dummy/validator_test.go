package dummy

import (
	baseData "github.com/AssetMantle/modules/schema/data/base"
	baseIDs "github.com/AssetMantle/modules/schema/ids/base"
	baseTypes "github.com/AssetMantle/modules/schema/parameters/base"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"testing"
)

func Test_validator(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"+ve with nil", args{Parameter}, false},
		{"-ve wrong parameter Type", args{baseTypes.NewParameter(baseIDs.NewID("newID"), baseData.NewDecData(sdkTypes.NewDec(-1)), validator)}, true},
		{"+ve empty string", args{baseIDs.NewID("")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validator(tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("validator() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
