package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"

	"github.com/andviro/goldie"
	"github.com/go-mixins/metadata"
	mdHTTP "github.com/go-mixins/metadata/http"
)

var testCases = `[
	{
		"X-B3-Sampled": [
			"1"
		],
		"X-B3-Spanid": [
			"f53aed610dfe10ec"
		],
		"X-B3-Traceid": [
			"ffa74529a27a9e7fcdfa858be1245cfd"
		],
		"X-Meta-From-Service": [
			"some-service"
		]
	}
]`

func TestFromHeader(t *testing.T) {
	var tcs []http.Header
	if err := json.Unmarshal([]byte(testCases), &tcs); err != nil {
		t.Fatal(err)
	}
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "\t")
	for _, tc := range tcs {
		ctx := mdHTTP.FromHeader(context.TODO(), tc)
		enc.Encode(metadata.From(ctx))
		for _, key := range []string{
			"from-service",
			"From-Service",
			"From_Service",
		} {
			fmt.Fprintf(buf, "\n%q: %q", key, metadata.Get(ctx, key))

		}
	}
	goldie.Assert(t, "from-header", buf.Bytes())
}

func TestToHeader(t *testing.T) {
	ctx := context.TODO()
	for key, value := range map[string]string{
		"some-key":       "some-value 1",
		"Some-Key":       "some-value 2",
		"Some_Key":       "some-value 3",
		"Some_Other_Key": "some-other-value",
	} {
		ctx = metadata.Set(ctx, key, value)
	}
	hdr := make(http.Header)
	mdHTTP.ToHeader(ctx, hdr)
	jd, _ := json.MarshalIndent(hdr, "", "\t")
	goldie.Assert(t, "to-header", jd)
}

func TestRoundtripper_Handler(t *testing.T) {
	c := &http.Client{
		Transport: &mdHTTP.Transport{},
	}
	srv := httptest.NewServer(
		&mdHTTP.Handler{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				td, _ := httputil.DumpRequest(r, false)
				t.Logf("%s", td)
				fromService := metadata.Get(r.Context(), "from-service")
				if fromService != "test" {
					t.Errorf("invalid metadata received: %q", fromService)
				}
			}),
		},
	)
	defer srv.Close()
	req, err := http.NewRequest("POST", srv.URL, strings.NewReader("test"))
	if err != nil {
		t.Fatal(err)
	}
	ctx := metadata.Set(req.Context(), "from-service", "test")
	req = req.WithContext(ctx)
	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("invalid code: %d", resp.StatusCode)
	}
}
