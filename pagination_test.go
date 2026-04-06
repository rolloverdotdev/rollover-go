package rollover

import (
	"context"
	"testing"
)

func mockListFn(pages []Page[string]) func(context.Context, *ListOptions) (*Page[string], error) {
	return func(ctx context.Context, opts *ListOptions) (*Page[string], error) {
		idx := opts.Offset / opts.Limit
		if idx >= len(pages) {
			return &Page[string]{}, nil
		}
		return &pages[idx], nil
	}
}

func TestPagesIteratesAllPages(t *testing.T) {
	fn := mockListFn([]Page[string]{
		{Data: []string{"a", "b"}, Total: 5, Limit: 2},
		{Data: []string{"c", "d"}, Total: 5, Limit: 2},
		{Data: []string{"e"}, Total: 5, Limit: 2},
	})

	iter := Pages(fn, &ListOptions{Limit: 2})
	var all []string
	for iter.Next(context.Background()) {
		all = append(all, iter.Page().Data...)
	}
	if err := iter.Err(); err != nil {
		t.Fatal(err)
	}
	if len(all) != 5 {
		t.Errorf("expected 5 items, got %d", len(all))
	}
}

func TestPagesStopsOnEmptyPage(t *testing.T) {
	fn := mockListFn([]Page[string]{
		{Data: []string{}, Total: 0, Limit: 10},
	})

	iter := Pages(fn, nil)
	if iter.Next(context.Background()) {
		t.Error("expected Next to return false on empty page")
	}
}

func TestPagesStopsOnShortPage(t *testing.T) {
	calls := 0
	fn := func(ctx context.Context, opts *ListOptions) (*Page[string], error) {
		calls++
		return &Page[string]{
			Data:  []string{"a"},
			Total: 1,
			Limit: 100,
		}, nil
	}

	iter := Pages(fn, nil)
	iter.Next(context.Background())
	if iter.Next(context.Background()) {
		t.Error("expected iteration to stop after short page")
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestPagesDefaultLimit(t *testing.T) {
	iter := Pages(func(ctx context.Context, opts *ListOptions) (*Page[string], error) {
		if opts.Limit != 100 {
			t.Errorf("expected default limit 100, got %d", opts.Limit)
		}
		return &Page[string]{}, nil
	}, nil)
	iter.Next(context.Background())
}

func TestCollect(t *testing.T) {
	fn := mockListFn([]Page[string]{
		{Data: []string{"a", "b"}, Total: 3, Limit: 2},
		{Data: []string{"c"}, Total: 3, Limit: 2},
	})

	all, err := Collect(context.Background(), fn, &ListOptions{Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 3 {
		t.Errorf("expected 3 items, got %d", len(all))
	}
}

func TestCollectEmpty(t *testing.T) {
	fn := mockListFn([]Page[string]{
		{Data: []string{}, Total: 0, Limit: 10},
	})

	all, err := Collect(context.Background(), fn, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 0 {
		t.Errorf("expected 0 items, got %d", len(all))
	}
}
