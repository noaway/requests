package requests

import (
	"testing"
)

func TestHTTPClient_Get(t *testing.T) {
	c := HTTPClient{}
	d, e := c.Get("http://ip.cn").String()
	if e != nil {
		t.Error(e)
		return
	}
	t.Log(d)
}
