package rss

import "testing"

func TestRSSStruct(t *testing.T) {
	rts, err := RSSStruct("https://habr.com/ru/rss/hub/go/all/?fl=ru")
	if err != nil {
		t.Fatal(err)
	}
	if len(rts) == 0 {
		t.Fatal("Данные отсуствуют")
	}
	t.Logf("получено %d новостей\n%+v", len(rts), rts)
	rts, err = RSSStruct("https://habr.com/ru/rss/best/daily/?fl=ru")
	if err != nil {
		t.Fatal(err)
	}
	if len(rts) == 0 {
		t.Fatal("данные не раскодированы")
	}
	t.Logf("получено %d новостей\n%+v", len(rts), rts)
}
