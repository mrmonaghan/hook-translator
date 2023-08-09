package stitch

import (
	"reflect"
	"testing"

	"github.com/slack-go/slack"
)

func TestInitSlack(t *testing.T) {
	tests := []struct {
		name    string
		want    *slack.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InitSlack()
			if (err != nil) != tt.wantErr {
				t.Errorf("InitSlack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitSlack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlocks_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Blocks []slack.Block
	}
	type args struct {
		data []byte
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
			b := &Blocks{
				Blocks: tt.fields.Blocks,
			}
			if err := b.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Blocks.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
