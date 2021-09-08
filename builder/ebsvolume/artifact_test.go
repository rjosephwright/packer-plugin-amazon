package ebsvolume

import (
	"reflect"
	"testing"

	registryimage "github.com/hashicorp/packer-plugin-sdk/packer/registry/image"
)

func TestArtifactState(t *testing.T) {
	expectedData := "this is the data"
	artifact := &Artifact{
		StateData: map[string]interface{}{"state_data": expectedData},
	}

	// Valid state
	result := artifact.State("state_data")
	if result != expectedData {
		t.Fatalf("Bad: State data was %s instead of %s", result, expectedData)
	}

	// Invalid state
	result = artifact.State("invalid_key")
	if result != nil {
		t.Fatalf("Bad: State should be nil for invalid state data name")
	}

	// Nil StateData should not fail and should return nil
	artifact = &Artifact{}
	result = artifact.State("key")
	if result != nil {
		t.Fatalf("Bad: State should be nil for nil StateData")
	}
}

func TestArtifactState_hcpPackerRegistryMetadata(t *testing.T) {
	volumes := make(EbsVolumes)
	volumes["west"] = []string{"vol-4567", "vol-0987"}
	snapshots := make(EbsSnapshots)
	snapshots["west"] = []string{"snap-4567", "snap-0987"}

	a := &Artifact{
		Volumes:   volumes,
		Snapshots: snapshots,
	}

	actual := a.State(registryimage.ArtifactStateURI)
	expected := []*registryimage.Image{
		{
			ImageID:        "vol-4567",
			ProviderName:   "aws",
			ProviderRegion: "west",
		},
		{
			ImageID:        "vol-0987",
			ProviderName:   "aws",
			ProviderRegion: "west",
		},
		{
			ImageID:        "snap-4567",
			ProviderName:   "aws",
			ProviderRegion: "west",
		},
		{
			ImageID:        "snap-0987",
			ProviderName:   "aws",
			ProviderRegion: "west",
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}
