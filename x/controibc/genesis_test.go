package controibc_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/tharsis/ethermint/testutil/keeper"
	"github.com/tharsis/ethermint/testutil/nullify"
	"github.com/tharsis/ethermint/x/controibc"
	"github.com/tharsis/ethermint/x/controibc/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.ControibcKeeper(t)
	controibc.InitGenesis(ctx, *k, genesisState)
	got := controibc.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	// this line is used by starport scaffolding # genesis/test/assert
}
