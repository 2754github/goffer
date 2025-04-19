package goffer

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type Item int

type fields struct {
	size            int
	items           []Item
	isNotSubscribed bool
	isClosed        bool
}

func TestNew(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
		want *buffer[Item]
	}{
		{
			name: "size == -1",
			args: args{
				size: -1,
			},
			want: &buffer[Item]{
				size:            1,
				items:           make([]Item, 0, 1),
				isNotSubscribed: true,
				isClosed:        false,
			},
		},
		{
			name: "size == 0",
			args: args{
				size: 0,
			},
			want: &buffer[Item]{
				size:            1,
				items:           make([]Item, 0, 1),
				isNotSubscribed: true,
				isClosed:        false,
			},
		},
		{
			name: "size == 1",
			args: args{
				size: 1,
			},
			want: &buffer[Item]{
				size:            1,
				items:           make([]Item, 0, 1),
				isNotSubscribed: true,
				isClosed:        false,
			},
		},
		{
			name: "size == 2",
			args: args{
				size: 2,
			},
			want: &buffer[Item]{
				size:            2,
				items:           make([]Item, 0, 2),
				isNotSubscribed: true,
				isClosed:        false,
			},
		},
		{
			name: fmt.Sprintf("size == %d", math.MaxInt32),
			args: args{
				size: math.MaxInt32,
			},
			want: &buffer[Item]{
				size:            math.MaxInt32,
				items:           make([]Item, 0, math.MaxInt32),
				isNotSubscribed: true,
				isClosed:        false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New[Item](tt.args.size).(*buffer[Item])

			check(t, got, tt.want)
		})
	}
}

func Test_buffer_Publish(t *testing.T) {
	type args struct {
		item Item
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       error
		wantBuffer *buffer[Item]
	}{
		{
			name: "buffer is closed",
			fields: fields{
				size:            3,
				items:           itemsWithCap([]Item{0}, 3),
				isNotSubscribed: false,
				isClosed:        true,
			},
			args: args{
				item: 0,
			},
			want: ErrClosed,
			wantBuffer: &buffer[Item]{
				size:            3,
				items:           itemsWithCap([]Item{0}, 3),
				isNotSubscribed: false,
				isClosed:        true,
			},
		},
		{
			name: "buffer is not subscribed",
			fields: fields{
				size:            3,
				items:           itemsWithCap([]Item{}, 3),
				isNotSubscribed: true,
				isClosed:        false,
			},
			args: args{
				item: 0,
			},
			want: ErrNotSubscribed,
			wantBuffer: &buffer[Item]{
				size:            3,
				items:           itemsWithCap([]Item{}, 3),
				isNotSubscribed: true,
				isClosed:        false,
			},
		},
		{
			name: "buffer is not full",
			fields: fields{
				size:            3,
				items:           itemsWithCap([]Item{0}, 3),
				isNotSubscribed: false,
				isClosed:        false,
			},
			args: args{
				item: 0,
			},
			want: nil,
			wantBuffer: &buffer[Item]{
				size:            3,
				items:           itemsWithCap([]Item{0, 0}, 3),
				isNotSubscribed: false,
				isClosed:        false,
			},
		},
		{
			name: "buffer fills up",
			fields: fields{
				size:            3,
				items:           itemsWithCap([]Item{0, 0}, 3),
				isNotSubscribed: false,
				isClosed:        false,
			},
			args: args{
				item: 0,
			},
			want: nil,
			wantBuffer: &buffer[Item]{
				size:            3,
				items:           itemsWithCap([]Item{}, 3),
				isNotSubscribed: false,
				isClosed:        false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &buffer[Item]{
				size:            tt.fields.size,
				items:           tt.fields.items,
				subscriber:      make(chan []Item, 1),
				isNotSubscribed: tt.fields.isNotSubscribed,
				isClosed:        tt.fields.isClosed,
			}
			if err := b.Publish(tt.args.item); !errors.Is(err, tt.want) {
				t.Errorf("buffer.Publish() error = %v, want %v", err, tt.want)
			}

			check(t, b, tt.wantBuffer)
		})
	}
}

func Test_buffer_Subscribe(t *testing.T) {
	tests := []struct {
		name       string
		fields     fields
		want       []Item
		wantBuffer *buffer[Item]
	}{
		{
			name: "buffer is not full",
			fields: fields{
				size:            3,
				items:           itemsWithCap([]Item{0}, 3),
				isNotSubscribed: true,
				isClosed:        false,
			},
			want: nil,
			wantBuffer: &buffer[Item]{
				size:            3,
				items:           itemsWithCap([]Item{0, 0}, 3),
				isNotSubscribed: false,
				isClosed:        false,
			},
		},
		{
			name: "buffer fills up",
			fields: fields{
				size:            3,
				items:           itemsWithCap([]Item{0, 0}, 3),
				isNotSubscribed: true,
				isClosed:        false,
			},
			want: []Item{0, 0, 0},
			wantBuffer: &buffer[Item]{
				size:            3,
				items:           itemsWithCap([]Item{}, 3),
				isNotSubscribed: false,
				isClosed:        false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &buffer[Item]{
				size:            tt.fields.size,
				items:           tt.fields.items,
				subscriber:      make(chan []Item, 1),
				isNotSubscribed: tt.fields.isNotSubscribed,
				isClosed:        tt.fields.isClosed,
			}

			go (func() {
				for got := range b.Subscribe() {
					if !reflect.DeepEqual(got, tt.want) {
						t.Errorf("buffer.Subscribe() = %v, want %v", got, tt.want)
					}
				}
			})()

			for b.Publish(0) != nil {
			}

			check(t, b, tt.wantBuffer)
		})
	}
}

func Test_buffer_Pull(t *testing.T) {
	tests := []struct {
		name       string
		fields     fields
		want       []Item
		wantBuffer *buffer[Item]
	}{
		{
			name: "buffer is closed",
			fields: fields{
				size:            3,
				items:           itemsWithCap([]Item{0}, 3),
				isNotSubscribed: false,
				isClosed:        true,
			},
			want: []Item{0},
			wantBuffer: &buffer[Item]{
				size:            3,
				items:           itemsWithCap([]Item{}, 3),
				isNotSubscribed: false,
				isClosed:        true,
			},
		},
		{
			name: "buffer is not subscribed",
			fields: fields{
				size:            3,
				items:           itemsWithCap([]Item{}, 3),
				isNotSubscribed: true,
				isClosed:        false,
			},
			want: []Item{},
			wantBuffer: &buffer[Item]{
				size:            3,
				items:           itemsWithCap([]Item{}, 3),
				isNotSubscribed: true,
				isClosed:        false,
			},
		},
		{
			name: "buffer is not full",
			fields: fields{
				size:            3,
				items:           itemsWithCap([]Item{0}, 3),
				isNotSubscribed: false,
				isClosed:        false,
			},
			want: []Item{0},
			wantBuffer: &buffer[Item]{
				size:            3,
				items:           itemsWithCap([]Item{}, 3),
				isNotSubscribed: false,
				isClosed:        false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &buffer[Item]{
				size:            tt.fields.size,
				items:           tt.fields.items,
				subscriber:      make(chan []Item, 1),
				isNotSubscribed: tt.fields.isNotSubscribed,
				isClosed:        tt.fields.isClosed,
			}
			if got := b.Pull(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buffer.Pull() = %v, want %v", got, tt.want)
			}

			check(t, b, tt.wantBuffer)
		})
	}
}

func Test_buffer_Close(t *testing.T) {
	tests := []struct {
		name       string
		fields     fields
		wantBuffer *buffer[Item]
	}{
		{
			name: "buffer is closed",
			fields: fields{
				size:            3,
				items:           itemsWithCap([]Item{0}, 3),
				isNotSubscribed: false,
				isClosed:        true,
			},
			wantBuffer: &buffer[Item]{
				size:            3,
				items:           itemsWithCap([]Item{0}, 3),
				isNotSubscribed: false,
				isClosed:        true,
			},
		},
		{
			name: "buffer is open",
			fields: fields{
				size:            3,
				items:           itemsWithCap([]Item{0}, 3),
				isNotSubscribed: false,
				isClosed:        false,
			},
			wantBuffer: &buffer[Item]{
				size:            3,
				items:           itemsWithCap([]Item{0}, 3),
				isNotSubscribed: false,
				isClosed:        true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &buffer[Item]{
				size:            tt.fields.size,
				items:           tt.fields.items,
				subscriber:      make(chan []Item, 1),
				isNotSubscribed: tt.fields.isNotSubscribed,
				isClosed:        tt.fields.isClosed,
			}
			b.Close()

			check(t, b, tt.wantBuffer)
		})
	}
}

func TestE2E(t *testing.T) {
	t.Run("should error without panic if Publish is called after Close", func(t *testing.T) {
		b := New[Item](3)
		defer b.Close()

		b.Close()
		err := b.Publish(0)

		if want := ErrClosed; !errors.Is(err, want) {
			t.Errorf("buffer.Publish() error = %v, want %v", err, want)
		}
	})

	t.Run("should error if Publish is called before Subscribe", func(t *testing.T) {
		b := New[Item](3)
		defer b.Close()

		err := b.Publish(0)

		if want := ErrNotSubscribed; !errors.Is(err, want) {
			t.Errorf("buffer.Publish() error = %v, want %v", err, want)
		}
	})

	t.Run("basic use case check", func(t *testing.T) {
		b := New[Item](3)
		defer b.Close()

		got := make([][]Item, 0, 2)
		var wg sync.WaitGroup
		wg.Add(1)
		go (func() {
			for items := range b.Subscribe() {
				got = append(got, items)
				wg.Done()
			}
		})()

		for i := range 5 {
			for b.Publish(Item(i)) != nil {
			}
		}

		wg.Wait()
		got = append(got, b.Pull())

		if want := [][]Item{{0, 1, 2}, {3, 4}}; !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	t.Run("should allow if Pull is called after Close", func(t *testing.T) {
		b := New[Item](3)
		defer b.Close()

		go (func() {
			for range b.Subscribe() {
			}
		})()

		for b.Publish(0) != nil {
		}

		b.Close()
		got := b.Pull()

		if want := []Item{0}; !reflect.DeepEqual(got, want) {
			t.Errorf("buffer.Pull() = %v, want %v", got, want)
		}
	})

	t.Run("should not panic if Close is called multiple times", func(t *testing.T) {
		b := New[Item](3)
		defer b.Close()

		b.Close()
		b.Close()
	})
}

func itemsWithCap(items []Item, cap int) []Item {
	r := make([]Item, len(items), cap)
	copy(r, items)

	return r
}

func check(t *testing.T, got *buffer[Item], want *buffer[Item]) {
	if diff := cmp.Diff(
		got,
		want,
		cmp.AllowUnexported(buffer[Item]{}),
		cmpopts.IgnoreFields(buffer[Item]{}, "mu", "subscriber"),
	); diff != "" {
		t.Error(diff)
	}

	if cap(got.items) != cap(want.items) {
		t.Errorf("cap(got.items) = %v, want %v", cap(got.items), cap(want.items))
	}

	if cap(got.subscriber) != 1 {
		t.Errorf("cap(got.subscriber) = %v, want %v", cap(got.subscriber), 1)
	}
}
