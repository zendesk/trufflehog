package documo

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/trufflesecurity/trufflehog/v3/pkg/detectors"
	"github.com/trufflesecurity/trufflehog/v3/pkg/engine/ahocorasick"
)

var (
	validPattern = `
		# Documo Configuration File: config.yaml
		database:
			host: $DB_HOST
			port: $DB_PORT
			username: $DB_USERNAME
			password: $DB_PASS  # IMPORTANT: Do not share this password publicly

		api:
			auth_type: "Basic"
			in: "Path"
			api_version: v1
			secret: "eyS9YqgD6TgdQ943G8S3aaiz26m2fTN9rcPbpeyts0jBEFd43hEFfr9pC7voqvLsbEi7Px4TbMToCVrstQRe8r2kltKGWyChYCT1Iruo6p3g3PyqZaZ1gOSbjeXz8zARUHZkXo7XR86kape65HLXj59yCNIlW5bvebJYbIAjjgGAAmXVgzldvNv8Zs08KIS5y62QJSNcnipFQbnxA8z6TUMl0F600MJhqEILWo19GaGjw"
			base_url: "https://api.example.com/$api_version/example"
			response_code: 200

		# Notes:
		# - Remember to rotate the secret every 90 days.
		# - The above credentials should only be used in a secure environment.
	`
	secret = "eyS9YqgD6TgdQ943G8S3aaiz26m2fTN9rcPbpeyts0jBEFd43hEFfr9pC7voqvLsbEi7Px4TbMToCVrstQRe8r2kltKGWyChYCT1Iruo6p3g3PyqZaZ1gOSbjeXz8zARUHZkXo7XR86kape65HLXj59yCNIlW5bvebJYbIAjjgGAAmXVgzldvNv8Zs08KIS5y62QJSNcnipFQbnxA8z6TUMl0F600MJhqEILWo19GaGjw"
)

func TestDocumo_Pattern(t *testing.T) {
	d := Scanner{}
	ahoCorasickCore := ahocorasick.NewAhoCorasickCore([]detectors.Detector{d})

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "valid pattern",
			input: validPattern,
			want:  []string{secret},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			matchedDetectors := ahoCorasickCore.FindDetectorMatches([]byte(test.input))
			if len(matchedDetectors) == 0 {
				t.Errorf("keywords '%v' not matched by: %s", d.Keywords(), test.input)
				return
			}

			results, err := d.FromData(context.Background(), false, []byte(test.input))
			if err != nil {
				t.Errorf("error = %v", err)
				return
			}

			if len(results) != len(test.want) {
				if len(results) == 0 {
					t.Errorf("did not receive result")
				} else {
					t.Errorf("expected %d results, only received %d", len(test.want), len(results))
				}
				return
			}

			actual := make(map[string]struct{}, len(results))
			for _, r := range results {
				if len(r.RawV2) > 0 {
					actual[string(r.RawV2)] = struct{}{}
				} else {
					actual[string(r.Raw)] = struct{}{}
				}
			}
			expected := make(map[string]struct{}, len(test.want))
			for _, v := range test.want {
				expected[v] = struct{}{}
			}

			if diff := cmp.Diff(expected, actual); diff != "" {
				t.Errorf("%s diff: (-want +got)\n%s", test.name, diff)
			}
		})
	}
}
