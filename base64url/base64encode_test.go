// package base64url provides base64url encoding/decoding support
package base64url

import (
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "simlpe string",
			args: args{
				"U2VuZCByZWluZm9yY2VtZW50cw",
			},
			want:    []byte("Send reinforcements"),
			wantErr: false,
		},
		{
			name: "longer string",
			args: args{"Tm93IGlzIHRoZSB0aW1lIGZvciBhbGwgZ29vZCBjb2RlcnN0byBsZWFybiBSdWJ5"},
			want: []byte("Now is the time for all good coders" +
				"to learn Ruby"),
			wantErr: false,
		},
		{
			name:    "zero",
			args:    args{"AA"},
			want:    []byte{00},
			wantErr: false,
		}, {
			name:    "zero zero",
			args:    args{"AAA"},
			want:    []byte{00, 00},
			wantErr: false,
		},
		{
			name:    "zero zero zero",
			args:    args{"AAAA"},
			want:    []byte{00, 00, 00},
			wantErr: false,
		},
		{
			name:    "LF",
			args:    args{"Cg"},
			want:    []byte("\n"),
			wantErr: false,
		},
		{
			name:    "CR LF",
			args:    args{"DQo"},
			want:    []byte{13, 10},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple string",
			args: args{[]byte("Send reinforcements")},
			want: "U2VuZCByZWluZm9yY2VtZW50cw",
		},
		{
			name: "longer string",
			args: args{[]byte("Now is the time for all good coders" +
				"to learn Ruby")},
			want: "Tm93IGlzIHRoZSB0aW1lIGZvciBhbGwgZ29vZCBjb2RlcnN0byBsZWFybiBSdWJ5",
		},
		{
			name: "multiline string",
			args: args{[]byte("This is line one" +
				"This is line two" +
				"This is line three" +
				"And so on..." +
				"")},
			want: "VGhpcyBpcyBsaW5lIG9uZVRoaXMgaXMgbGluZSB0d29UaGlzIGlzIGxpbmUgdGhyZWVBbmQgc28gb24uLi4",
		}, {
			name: "zero",
			args: args{
				data: []byte{00},
			},
			want: "AA",
		},
		{
			name: "zero zero",
			args: args{
				data: []byte{00, 00},
			},
			want: "AAA",
		},
		{
			name: "zero zero zero",
			args: args{
				data: []byte{00, 00, 00},
			},
			want: "AAAA",
		},
		{
			name: "LF",
			args: args{
				data: []byte{10},
			},
			want: "Cg",
		},
		{
			name: "CR",
			args: args{
				data: []byte{13},
			},
			want: "DQ",
		},
		{
			name: "CR LF",
			args: args{
				data: []byte{13, 10},
			},
			want: "DQo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Encode(tt.args.data); got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestDecompress(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name           string
		args           args
		wantDecompData []byte
		wantErr        bool
	}{
		{
			name: "empty, shortest possible",
			args: args{
				data: []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 0, 0, 0, 255, 255, 1, 0, 0, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			wantDecompData: []byte(""),
			wantErr:        false,
		},
		{
			name: "empty struct",
			args: args{
				data: []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 170, 174, 5, 0, 0, 0, 255, 255, 1, 0, 0, 255, 255, 67, 191, 166, 163, 2, 0, 0, 0},
			},
			wantDecompData: []byte("{}"),
			wantErr:        false,
		},
		{
			name: "simple JSON",
			args: args{
				data: []byte{31, 139, 8, 0, 0, 0, 0, 0, 0,
					255, 170, 86, 10, 72, 76, 79, 85,
					178, 50, 212, 81, 114, 43, 42, 205,
					44, 41, 86, 178, 138, 86, 74, 44,
					40, 200, 73, 85, 210, 81, 42, 72,
					77, 76, 206, 128, 208, 69, 74, 177,
					181, 0, 0, 0, 0, 255, 255, 1, 0, 0,
					255, 255, 247, 140, 57, 118, 44, 0, 0, 0},
			},
			wantDecompData: []byte("{\"Page\":1,\"Fruits\":[\"apple\",\"peach\",\"pear\"]}"),
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDecompData, err := Decompress(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decompress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDecompData, tt.wantDecompData) {
				t.Errorf("Decompress() = %v, want %v", gotDecompData, tt.wantDecompData)
			}
		})
	}
}

func TestCompress(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name               string
		args               args
		wantCompressedData []byte
		wantErr            bool
	}{
		{
			name: "empty, shortest possible",
			args: args{
				data: []byte(""),
			},
			wantCompressedData: []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 0, 0, 0, 255, 255, 1, 0, 0, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:            false,
		},
		{
			name: "empty struct",
			args: args{
				data: []byte("{}"),
			},
			wantCompressedData: []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 170, 174, 5, 0, 0, 0, 255, 255, 1, 0, 0, 255, 255, 67, 191, 166, 163, 2, 0, 0, 0},
			wantErr:            false,
		},
		{
			name: "simple JSON",
			args: args{
				data: []byte("{\"Page\":1,\"Fruits\":[\"apple\",\"peach\",\"pear\"]}"),
			},
			wantCompressedData: []byte{31, 139, 8, 0, 0, 0, 0, 0, 0,
				255, 170, 86, 10, 72, 76, 79, 85,
				178, 50, 212, 81, 114, 43, 42, 205,
				44, 41, 86, 178, 138, 86, 74, 44,
				40, 200, 73, 85, 210, 81, 42, 72,
				77, 76, 206, 128, 208, 69, 74, 177,
				181, 0, 0, 0, 0, 255, 255, 1, 0, 0,
				255, 255, 247, 140, 57, 118, 44, 0, 0, 0},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCompressedData, err := Compress(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCompressedData, tt.wantCompressedData) {
				t.Errorf("Compress() = %v, want %v", gotCompressedData, tt.wantCompressedData)
			}
		})
	}
}

func TestEncodeDecode(t *testing.T) {
	tests := []struct {
		name string
		args string
	}{
		{
			name: "empty",
			args: "",
		},
		{
			name: "simple",
			args: "this is a simple test",
		},
		{
			name: "json",
			args: "{\"Page\":1,\"Fruits\":[\"apple\",\"peach\",\"pear\"]}",
		},
		{
			name: "CR LF",
			args: "\r\n",
		},
		{
			name: "spaces and comma",
			args: "  ,,,,",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := []byte(tt.args)
			c1, _ := Compress(want)
			e1 := Encode(c1)
			d1, _ := Decode(e1)
			got, _ := Decompress(d1)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Encode/Decode.got = %v, want %v", got, want)
			}
		})
	}
}
