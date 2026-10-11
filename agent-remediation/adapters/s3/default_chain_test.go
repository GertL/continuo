package s3

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carolsimone/continuo/pkg/objectstore"
	"github.com/carolsimone/continuo/pkg/objectstore/objectstoretest"
)

// TestAdapters_SignWithTheDefaultChain proves the wiring main uses: one
// *objectstore.Client built without static keys, passed to each adapter as
// the role it uses, signs every request with the credentials the SDK's
// default chain resolves instead of failing on empty static keys.
func TestAdapters_SignWithTheDefaultChain(t *testing.T) {
	objectstoretest.UseChainCredentials(t)
	srv := objectstoretest.NewServer(t, objectstoretest.Object([]byte(validBundleJSON)))
	client, err := objectstore.New(context.Background(), objectstore.Config{
		EndpointURL: srv.URL, Region: "us-east-1", Bucket: testBucket,
	})
	require.NoError(t, err)

	src, err := NewCandidateSourceReader(client).NodeSource(context.Background(),
		"code-bundles/rel-1/bundle.json", "analytics.orders", "rel-1")
	require.NoError(t, err)
	assert.Equal(t, "select 1 as x", src.RawCode)

	reqs := srv.Requests()
	require.Len(t, reqs, 1)
	assert.Equal(t, "/test-bucket/code-bundles/rel-1/bundle.json", reqs[0].Path, "a bare key reads the client's default bucket")
	for _, r := range reqs {
		assert.Contains(t, r.Authorization, "Credential="+objectstoretest.ChainAccessKeyID+"/")
	}
}
