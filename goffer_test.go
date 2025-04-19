package goffer

import (
	"reflect"
	"testing"
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New[Item](tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buffer_Publish(t *testing.T) {
	type args struct {
		item Item
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if err := b.Publish(tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("buffer.Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_buffer_Subscribe(t *testing.T) {
	tests := []struct {
		name   string
		fields fields
		want   <-chan []Item
	}{
		// TODO: Add test cases.
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
			if got := b.Subscribe(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buffer.Subscribe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buffer_Pull(t *testing.T) {
	tests := []struct {
		name   string
		fields fields
		want   []Item
	}{
		// TODO: Add test cases.
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
		})
	}
}

func Test_buffer_Close(t *testing.T) {
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
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
		})
	}
}
