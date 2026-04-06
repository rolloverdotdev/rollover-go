package rollover

import (
	"context"
	"net/http"
	"testing"
)

func TestListInvoices(t *testing.T) {
	c := testClient(t, orgHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("wallet") != "0xabc" {
			t.Errorf("expected wallet filter")
		}
		w.Write([]byte(`{"data":[{"id":"inv-1","wallet_address":"0xabc","status":"paid","total_amount":"9.99"}],"total":1,"limit":20,"offset":0}`))
	}))

	result, err := c.ListInvoices(context.Background(), &ListOptions{Wallet: "0xabc"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 1 || result.Data[0].Status != "paid" {
		t.Errorf("unexpected result: %+v", result)
	}
}
