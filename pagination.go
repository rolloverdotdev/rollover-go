package rollover

import "context"

// Iter pages through a list endpoint one page at a time, fetching the next
// page only when Next is called.
type Iter[T any] struct {
	fn   func(context.Context, *ListOptions) (*Page[T], error)
	opts ListOptions
	page *Page[T]
	err  error
	done bool
}

// Pages returns an iterator that lazily fetches one page at a time from any
// list endpoint, defaulting to 100 items per page when opts is nil or Limit
// is zero.
func Pages[T any](
	fn func(context.Context, *ListOptions) (*Page[T], error),
	opts *ListOptions,
) *Iter[T] {
	it := &Iter[T]{fn: fn}
	if opts != nil {
		it.opts = *opts
	}
	if it.opts.Limit <= 0 {
		it.opts.Limit = 100
	}
	return it
}

// Next fetches the next page and returns true if there are results to read.
// Returns false when all pages have been consumed, the page is empty, or an
// error occurs.
func (it *Iter[T]) Next(ctx context.Context) bool {
	if it.err != nil || it.done {
		return false
	}

	it.page, it.err = it.fn(ctx, &it.opts)
	if it.err != nil || it.page == nil {
		return false
	}

	count := len(it.page.Data)
	it.opts.Offset += count

	if count < it.opts.Limit || it.opts.Offset >= it.page.Total {
		it.done = true
	}

	return count > 0
}

// Page returns the most recently fetched page.
func (it *Iter[T]) Page() *Page[T] { return it.page }

// Err returns the first error encountered during iteration.
func (it *Iter[T]) Err() error { return it.err }

// Collect fetches all pages from a list endpoint and returns every item in
// a single slice, handling pagination automatically.
func Collect[T any](
	ctx context.Context,
	fn func(context.Context, *ListOptions) (*Page[T], error),
	opts *ListOptions,
) ([]T, error) {
	iter := Pages(fn, opts)
	var all []T

	for iter.Next(ctx) {
		if all == nil {
			all = make([]T, 0, iter.Page().Total)
		}
		all = append(all, iter.Page().Data...)
	}

	return all, iter.Err()
}
