//go:build detectors
// +build detectors

package hasura

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/trufflesecurity/trufflehog/v3/pkg/common"
	"github.com/trufflesecurity/trufflehog/v3/pkg/detectors"
	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/detectorspb"
)

func TestHasura_FromData(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	testSecrets, err := common.GetSecret(ctx, "trufflehog-testing", "detectors5")
	if err != nil {
		t.Fatalf("could not get test secrets from GCP: %s", err)
	}

	secret := testSecrets.MustGetField("HASURA")
	inactiveSecret := testSecrets.MustGetField("HASURA_INACTIVE")
	domain := testSecrets.MustGetField("HASURA_DOMAIN")

	type args struct {
		ctx    context.Context
		data   []byte
		verify bool
	}
	tests := []struct {
		name                string
		s                   Scanner
		args                args
		want                []detectors.Result
		wantErr             bool
		wantVerificationErr bool
	}{
		{
			name: "found, verified",
			s:    Scanner{},
			args: args{
				ctx:    context.Background(),
				data:   []byte(fmt.Sprintf("You can find a hasura secret %s within hasura domain %s", secret, domain)),
				verify: true,
			},
			want: []detectors.Result{
				{
					DetectorType: detectorspb.DetectorType_Hasura,
					Verified:     true,
					Raw:          []byte(secret),
					RawV2:        []byte(fmt.Sprintf("%s:%s", domain, secret)),
					ExtraData:    map[string]string{"domain": domain},
				},
			},
			wantErr:             false,
			wantVerificationErr: false,
		},
		{
			name: "found, unverified",
			s:    Scanner{},
			args: args{
				ctx:    context.Background(),
				data:   []byte(fmt.Sprintf("You can find a hasura secret %s within hasura domain %s but not valid", inactiveSecret, domain)),
				verify: true,
			},
			want: []detectors.Result{
				{
					DetectorType: detectorspb.DetectorType_Hasura,
					Verified:     false,
					Raw:          []byte(inactiveSecret),
					RawV2:        []byte(fmt.Sprintf("%s:%s", domain, inactiveSecret)),
					ExtraData:    map[string]string{"domain": domain},
				},
			},
			wantErr:             false,
			wantVerificationErr: false,
		},
		{
			name: "not found",
			s:    Scanner{},
			args: args{
				ctx:    context.Background(),
				data:   []byte("You cannot find the secret within"),
				verify: true,
			},
			want:                nil,
			wantErr:             false,
			wantVerificationErr: false,
		},
		{
			name: "found, would be verified if not for timeout",
			s:    Scanner{client: common.SaneHttpClientTimeOut(1 * time.Microsecond)},
			args: args{
				ctx:    context.Background(),
				data:   []byte(fmt.Sprintf("You can find a hasura secret %s within %s", secret, domain)),
				verify: true,
			},
			want: []detectors.Result{
				{
					DetectorType: detectorspb.DetectorType_Hasura,
					Verified:     false,
					Raw:          []byte(secret),
					RawV2:        []byte(fmt.Sprintf("%s:%s", domain, secret)),
					ExtraData:    nil,
				},
			},
			wantErr:             false,
			wantVerificationErr: true,
		},
		{
			name: "found, unexpected api response",
			s:    Scanner{client: common.ConstantResponseHttpClient(500, "")},
			args: args{
				ctx:    context.Background(),
				data:   []byte(fmt.Sprintf("You can find a hasura secret %s within %s", secret, domain)),
				verify: true,
			},
			want: []detectors.Result{
				{
					DetectorType: detectorspb.DetectorType_Hasura,
					Verified:     false,
					Raw:          []byte(secret),
					RawV2:        []byte(fmt.Sprintf("%s:%s", domain, secret)),
					ExtraData:    map[string]string{"domain": domain},
				},
			},
			wantErr:             false,
			wantVerificationErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.FromData(tt.args.ctx, tt.args.verify, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := range got {
				if len(got[i].Raw) == 0 {
					t.Fatalf("no raw secret present: \n %+v", got[i])
				}

				if (got[i].VerificationError() != nil) != tt.wantVerificationErr {
					t.Fatalf("wantVerificationErr = %v, got verification error = %v", tt.wantVerificationErr, got[i].VerificationError())
				}
			}
			ignoreOpts := cmpopts.IgnoreFields(detectors.Result{}, "verificationError", "primarySecret")
			if diff := cmp.Diff(tt.want, got, ignoreOpts); diff != "" {
				t.Errorf("Hasura.FromData() %s - diff: (-got +want)\n%s", tt.name, diff)
			}
		})
	}
}

func BenchmarkFromData(benchmark *testing.B) {
	ctx := context.Background()
	s := Scanner{}
	for name, data := range detectors.MustGetBenchmarkData() {
		benchmark.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_, err := s.FromData(ctx, false, data)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
