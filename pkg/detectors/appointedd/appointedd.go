package appointedd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	regexp "github.com/wasilibs/go-re2"

	"github.com/trufflesecurity/trufflehog/v3/pkg/common"
	"github.com/trufflesecurity/trufflehog/v3/pkg/detectors"
	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/detectorspb"
)

type Scanner struct{}

// Ensure the Scanner satisfies the interface at compile time
var _ detectors.Detector = (*Scanner)(nil)

var (
	client = common.SaneHttpClient()

	// Make sure that your group is surrounded in boundary characters such as below to reduce false positives
	keyPat = regexp.MustCompile(detectors.PrefixRegex([]string{"appointedd"}) + `\b([a-zA-Z0-9=+]{88})`)
)

// Keywords are used for efficiently pre-filtering chunks.
// Use identifiers in the secret preferably, or the provider name.
func (s Scanner) Keywords() []string {
	return []string{"appointedd"}
}

// FromData will find and optionally verify appointedd secrets in a given set of bytes.
func (s Scanner) FromData(ctx context.Context, verify bool, data []byte) (results []detectors.Result, err error) {
	dataStr := string(data)

	matches := keyPat.FindAllStringSubmatch(dataStr, -1)

	for _, match := range matches {
		resMatch := strings.TrimSpace(match[1])

		s1 := detectors.Result{
			DetectorType: detectorspb.DetectorType_Appointedd,
			Raw:          []byte(resMatch),
		}
		if verify {
			isVerified, err := verifyMatch(ctx, client, resMatch)
			s1.Verified = isVerified
			s1.SetVerificationError(err, resMatch)
		}

		results = append(results, s1)
	}

	return results, nil
}

func verifyMatch(ctx context.Context, client *http.Client, secret string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.appointedd.com/v1/availability/slots", http.NoBody)
	if err != nil {
		return false, err
	}
	req.Header.Add("X-API-KEY", secret)
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	}()

	switch res.StatusCode {
	case http.StatusOK:
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return false, err
		}
		return strings.Contains(string(bodyBytes), "total"), nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return false, nil
	default:
		return false, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
}

func (s Scanner) Type() detectorspb.DetectorType {
	return detectorspb.DetectorType_Appointedd
}

func (s Scanner) Description() string {
	return "Appointedd provides online booking and scheduling services. The API key can be used to access and manage booking data."
}
