package k8sclient

import (
	"testing"
)

func TestGetScaleTarget(t *testing.T) {
	testCases := []struct {
		target   string
		expKind  string
		expName  string
		expError bool
	}{
		{
			"deployment/anything",
			"deployment",
			"anything",
			false,
		},
		{
			"replicationcontroller/anotherthing",
			"replicationcontroller",
			"anotherthing",
			false,
		},
		{
			"replicationcontroller",
			"",
			"",
			true,
		},
		{
			"replicaset/anything/what",
			"",
			"",
			true,
		},
	}

	for _, tc := range testCases {
		res, err := getScaleTarget(tc.target, "default")
		if err != nil && !tc.expError {
			t.Errorf("Expect no error, got error for target: %v", tc.target)
			continue
		} else if err == nil && tc.expError {
			t.Errorf("Expect error, got no error for target: %v", tc.target)
			continue
		}
		if res.kind != tc.expKind || res.name != tc.expName {
			t.Errorf("Expect kind: %v, name: %v\ngot kind: %v, name: %v", tc.expKind, tc.expName, res.kind, res.name)
		}
	}
}
