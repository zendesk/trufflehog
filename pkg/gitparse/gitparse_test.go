	"bytes"
	"strings"
	"time"

	"github.com/trufflesecurity/trufflehog/v3/pkg/context"
	fails    [][]byte
			fails:    [][]byte{[]byte("notcorrect")},
			fails:    [][]byte{[]byte("notcorrect")},
			fails:    [][]byte{[]byte("notcorrect"), []byte("--- s"), []byte("short")},
			fails:    [][]byte{[]byte("notcorrect"), []byte("+++ s"), []byte("short")},
			fails:    [][]byte{[]byte("notcorrect")},
			fails:    [][]byte{[]byte("notcorrect")},
			fails:    [][]byte{[]byte("notcorrect")},
			fails:    [][]byte{[]byte("notcorrect")},
		for _, fail := range test.fails {
			if test.function(fail) {
				t.Errorf("%s: Parser did not recognize incorrect line.", name)
			}

func TestSingleCommitSingleDiff(t *testing.T) {
	r := bytes.NewReader([]byte(singleCommitSingleDiff))
	commitChan := make(chan Commit)
	date, _ := time.Parse(DateFormat, "Mon Mar 15 23:27:16 2021 -0700")
	content := bytes.NewBuffer([]byte(singleCommitSingleDiffDiff))
	builder := strings.Builder{}
	builder.Write([]byte(singleCommitSingleDiffMessage))
	expected := []Commit{
		{
			Hash:    "70001020fab32b1fcf2f1f0e5c66424eae649826",
			Author:  "Dustin Decker <humanatcomputer@gmail.com>",
			Date:    date,
			Message: builder,
			Diffs: []Diff{
				{
					PathB:     "aws",
					LineStart: 1,
					Content:   *content,
					IsBinary:  false,
				},
			},
		},
	}
	go func() {
		FromReader(context.TODO(), r, commitChan)
	}()
	i := 0
	for commit := range commitChan {
		if len(expected) < i {
			t.Errorf("Commit does not match. Wrong number of commits.")
		}

		if !commit.Equal(&expected[i]) {
			t.Errorf("Commit does not match. Got: %v, expected: %v", commit, expected)
		}
		i++
	}
}

func TestMultiCommitContextDiff(t *testing.T) {
	r := bytes.NewReader([]byte(singleCommitContextDiff))
	commitChan := make(chan Commit)
	dateOne, _ := time.Parse(DateFormat, "Mon Mar 15 23:27:16 2021 -0700")
	dateTwo, _ := time.Parse(DateFormat, "Wed Dec 12 18:19:21 2018 -0800")
	diffOneA := bytes.NewBuffer([]byte(singleCommitContextDiffDiffOneA))
	diffTwoA := bytes.NewBuffer([]byte(singleCommitContextDiffDiffTwoA))
	// diffTwoB := bytes.NewBuffer([]byte(singleCommitContextDiffDiffTwoB))
	messageOne := strings.Builder{}
	messageOne.Write([]byte(singleCommitContextDiffMessageOne))
	messageTwo := strings.Builder{}
	messageTwo.Write([]byte(singleCommitContextDiffMessageTwo))
	expected := []Commit{
		{
			Hash:    "70001020fab32b1fcf2f1f0e5c66424eae649826",
			Author:  "Dustin Decker <humanatcomputer@gmail.com>",
			Date:    dateOne,
			Message: messageOne,
			Diffs: []Diff{
				{
					PathB:     "aws",
					LineStart: 1,
					Content:   *diffOneA,
					IsBinary:  false,
				},
			},
		},
		{
			Hash:    "84e9c75e388ae3e866e121087ea2dd45a71068f2",
			Author:  "Dylan Ayrey <dxa4481@rit.edu>",
			Date:    dateTwo,
			Message: messageTwo,
			Diffs: []Diff{
				{
					PathB:     "aws",
					LineStart: 1,
					Content:   *diffTwoA,
					IsBinary:  false,
				},
			},
		},
	}
	go func() {
		FromReader(context.TODO(), r, commitChan)
	}()
	i := 0
	for commit := range commitChan {
		if len(expected) < i {
			t.Errorf("Commit does not match. Wrong number of commits.")
		}

		if !commit.Equal(&expected[i]) {
			t.Errorf("Commit does not match. Got: %v, expected: %v", commit, expected[i])
		}
		i++
	}
}

const singleCommitSingleDiff = `commit 70001020fab32b1fcf2f1f0e5c66424eae649826 (HEAD -> master, origin/master, origin/HEAD)
Author: Dustin Decker <humanatcomputer@gmail.com>
Date:   Mon Mar 15 23:27:16 2021 -0700

    Update aws

diff --git a/aws b/aws
index 2ee133b..12b4843 100644
--- a/aws
+++ b/aws
@@ -1,7 +1,5 @@
-blah blaj
-
-this is the secret: [Default]
-Access key Id: AKIAILE3JG6KMS3HZGCA
-Secret Access Key: 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
-
-okay thank you bye
+[default]
+aws_access_key_id = AKIAXYZDQCEN4B6JSJQI
+aws_secret_access_key = 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
+output = json
+region = us-east-2
`
const singleCommitSingleDiffMessage = `Update aws
`

const singleCommitSingleDiffDiff = `[default]
aws_access_key_id = AKIAXYZDQCEN4B6JSJQI
aws_secret_access_key = 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
output = json
region = us-east-2
`
const singleCommitContextDiff = `commit 70001020fab32b1fcf2f1f0e5c66424eae649826 (HEAD -> master, origin/master, origin/HEAD)
Author: Dustin Decker <humanatcomputer@gmail.com>
Date:   Mon Mar 15 23:27:16 2021 -0700

    Update aws

diff --git a/aws b/aws
index 2ee133b..12b4843 100644
--- a/aws
+++ b/aws
@@ -1,7 +1,5 @@
-blah blaj
-
-this is the secret: [Default]
-Access key Id: AKIAILE3JG6KMS3HZGCA
-Secret Access Key: 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
-
-okay thank you bye
+[default]
+aws_access_key_id = AKIAXYZDQCEN4B6JSJQI
+aws_secret_access_key = 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
+output = json
+region = us-east-2

commit 84e9c75e388ae3e866e121087ea2dd45a71068f2
Author: Dylan Ayrey <dxa4481@rit.edu>
Date:   Wed Dec 12 18:19:21 2018 -0800

    Update aws again

diff --git a/aws b/aws
index 239b415..2ee133b 100644
--- a/aws
+++ b/aws
@@ -1,5 +1,7 @@
 blah blaj
 
-this is the secret: AKIA2E0A8F3B244C9986
+this is the secret: [Default]
+Access key Id: AKIAILE3JG6KMS3HZGCA
+Secret Access Key: 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
 
-okay thank you bye
\ No newline at end of file
+okay thank you bye
`

const singleCommitContextDiffMessageOne = `Update aws
`

const singleCommitContextDiffMessageTwo = `Update aws again
`

const singleCommitContextDiffDiffOneA = `[default]
aws_access_key_id = AKIAXYZDQCEN4B6JSJQI
aws_secret_access_key = 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
output = json
region = us-east-2
`

const singleCommitContextDiffDiffTwoA = `

this is the secret: [Default]
Access key Id: AKIAILE3JG6KMS3HZGCA
Secret Access Key: 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

okay thank you bye
`