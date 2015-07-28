package main

import "testing"

func TestEmptyScrapeStore(t *testing.T) {
	store := newScrapeStore()

	if store.current() != nil {
		t.Error("store.current() expected nil, got ", store.current())
	}
	if store.prev() != nil {
		t.Error("store.prev() expected nil, got ", store.prev())
	}
}

func TestScrapeStore(t *testing.T) {
	store := newScrapeStore()
	store.add([]byte("first"))

	if string(store.current()) != "first" {
		t.Fail()
	}

	store.add([]byte("second"))
	store.add([]byte("third"))

	if string(store.current()) != "third" {
		t.Fail()
	}

	if string(store.prev()) != "second" {
		t.Fail()
	}
}

func TestScrapeStoreDiff(t *testing.T) {
	store := newScrapeStore()

	store.add([]byte("a\nb\nc\n"))
	store.add([]byte("a\nc\nb\n"))
	diff, _ := store.htmlDiffPrev()

	expected := `<table><tr><td class="line-num">1</td><td><pre>a</pre></td><td><pre>a</pre></td><td class="line-num">1</td></tr>
<tr><td class="line-num">2</td><td class="deleted"><pre>b</pre></td><td></td><td></td></tr>
<tr><td class="line-num">3</td><td><pre>c</pre></td><td><pre>c</pre></td><td class="line-num">2</td></tr>
<tr><td class="line-num"></td><td></td><td class="added"><pre>b</pre></td><td class="line-num">3</td></tr>
<tr><td class="line-num">4</td><td><pre></pre></td><td><pre></pre></td><td class="line-num">4</td></tr>
</html>`

	if diff != expected {
		t.Error("HTML diff expected output: ", expected, "got output: ", diff)
	}
}
