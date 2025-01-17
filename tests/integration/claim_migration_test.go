// +build integration

package integration

import (
	"path/filepath"
	"testing"

	"get.porter.sh/porter/pkg/porter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Do a migration. This also checks for any problems with our
// connection handling which can result in panics :-)
func TestClaimMigration_List(t *testing.T) {
	t.Parallel()

	p := porter.NewTestPorter(t)
	p.SetupIntegrationTest()
	defer p.CleanupIntegrationTest()

	// Make a test home
	home, err := p.GetHomeDir()
	require.NoError(t, err, "GetHomeDir failed")
	claimsDir := filepath.Join(home, "claims")

	// Create unmigrated claim data
	p.FileSystem.Mkdir(claimsDir, 0700)
	p.AddTestFile(filepath.Join(p.RepoRoot, "pkg/storage/testdata/claims", "upgraded.json"), filepath.Join(home, "claims", "mybun.json"))
	p.FileSystem.Remove(filepath.Join(home, "schema.json"))

	err = p.MigrateStorage()
	require.NoError(t, err, "MigrateStorage failed")

	installations, err := p.ListInstallations()
	require.NoError(t, err, "could not list installations")
	require.Len(t, installations, 1, "expected one installation")
	assert.Equal(t, "mybun", installations[0].Name, "unexpected list of installation names")
}
