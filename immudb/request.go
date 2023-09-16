package immudb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

func dumpRequest(req *http.Request, prefix string) {
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\nREQUEST:\n%s\n", prefix, dump)
}

func dumpResponse(resp *http.Response, suffix string) {
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("RESPONSE:\n%s\n%s\n", dump, suffix)
}

func (l *Ledger) request(method string, pathSuffix string, in, out any) error {
	var body io.Reader
	if in != nil {
		data, err := json.Marshal(in)
		if err != nil {
			return fmt.Errorf("json.Marshal: %w", err)
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, l.baseUrl+pathSuffix, body)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Set("X-API-Key", l.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if l.dumpHttp {
		dumpRequest(req, "----------------------------------->")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http.Client.Do: %w", err)
	}
	if l.dumpHttp {
		dumpResponse(resp, "<-----------------------------------")
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return newHttpError(method, pathSuffix, resp)
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(&out)
}
