package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "github.com/tharsis/ethermint/testutil/keeper"
	"github.com/tharsis/ethermint/x/controibc/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.ControibcKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
	require.EqualValues(t, params.PortId, k.PortId(ctx))
}
