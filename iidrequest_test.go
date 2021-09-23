package instanceid

import (
	"reflect"
	"testing"
)

func TestNewIRequestFromValueValid(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want *IRequest
	}{
		{
			name: "no key",
			args: args{
				v: "",
			},
			want: &IRequest{
				key:     "empty",
				options: map[string]Option{},
			},
		},

		{
			name: "empty key",
			args: args{
				v: "empty",
			},
			want: &IRequest{
				key:     "empty",
				options: map[string]Option{},
			},
		},
		{
			name: "empty key;",
			args: args{
				v: "empty;",
			},
			want: &IRequest{
				key:     "empty",
				options: map[string]Option{},
			},
		},
		{
			name: "simple key",
			args: args{
				v: "key=asdf",
			},
			want: &IRequest{
				key:     "asdf",
				options: map[string]Option{},
			},
		},
		{
			name: "simple key;",
			args: args{
				v: "key=asdf;",
			},
			want: &IRequest{
				key:     "asdf;",
				options: map[string]Option{},
			},
		},
		{
			name: "simple key ;",
			args: args{
				v: "key=asdf",
			},
			want: &IRequest{
				key:     "asdf",
				options: map[string]Option{},
			},
		},
		{
			name: "simple key tab",
			args: args{
				v: "key=asdf \t",
			},
			want: &IRequest{
				key:     "asdf",
				options: map[string]Option{},
			},
		},
		{
			name: "longer key",
			args: args{
				v: "key=1234-4444-asdf-234234-23423423",
			},
			want: &IRequest{
				key:     "1234-4444-asdf-234234-23423423",
				options: map[string]Option{},
			},
		},
		{
			name: "longer key;",
			args: args{
				v: "key=1234-4444-asdf-234234-23423423;",
			},
			want: &IRequest{
				key:     "1234-4444-asdf-234234-23423423;",
				options: map[string]Option{},
			},
		},
		{
			name: "empty and option",
			args: args{
				v: "empty options=v",
			},
			want: &IRequest{
				key: "empty",
				options: map[string]Option{
					"v": IOption{
						commandName: "v",
					},
				},
			},
		},
		{
			name: "empty and options",
			args: args{
				v: "empty options=vc",
			},
			want: &IRequest{
				key: "empty",
				options: map[string]Option{
					"v": IOption{
						commandName: "v",
					},
					"c": IOption{
						commandName: "c",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIRequestFromValue(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIRequestFromValue() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNewIRequestFromValueInvalid(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want *IRequest
	}{
		{
			name: "two empty key",
			args: args{
				v: "empty empty",
			},
			want: &IRequest{
				key:     "empty",
				options: map[string]Option{},
			},
		},
		{
			name: "two empty key",
			args: args{
				v: "empty; empty",
			},
			want: &IRequest{
				key:     "empty",
				options: map[string]Option{},
			},
		},
		{
			name: "two empty key reversed",
			args: args{
				v: "empty empty;",
			},
			want: &IRequest{
				key:     "empty",
				options: map[string]Option{},
			},
		},
		{
			name: "empty;",
			args: args{
				v: "empty;",
			},
			want: &IRequest{
				key:     "empty",
				options: map[string]Option{},
			},
		},
		{
			name: "two key",
			args: args{
				v: "key=ab; key=cd",
			},
			want: &IRequest{
				key:     "cd",
				options: map[string]Option{},
			},
		},
		{
			name: "empty and key",
			args: args{
				v: "empty; key=ab",
			},
			want: &IRequest{
				key:     "ab",
				options: map[string]Option{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIRequestFromValue(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIRequestFromValue() = %#v, want %+v", got, tt.want)
			}
		})
	}
}
func TestIRequest_SetIidAuth(t *testing.T) {
	type fields struct {
		key string
	}
	type args struct {
		in0 string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   IidRequest
	}{
		{
			name: "delete key",
			fields: fields{
				key: "aaa",
			},
			args: args{
				in0: "",
			},
			want: &IRequest{
				key:     "empty",
				options: map[string]Option{},
			},
		},
		{
			name: "set empty key",
			fields: fields{
				key: "aaa",
			},
			args: args{
				in0: "empty",
			},
			want: &IRequest{
				key:     "empty",
				options: map[string]Option{},
			},
		},
		{
			name: "set key short",
			fields: fields{
				key: "empty",
			},
			args: args{
				in0: "1234",
			},
			want: &IRequest{
				key:     "1234",
				options: map[string]Option{},
			},
		},
		{
			name: "set key long",
			fields: fields{
				key: "empty",
			},
			args: args{
				in0: "1234-1234-1234-1234-1234-1234",
			},
			want: &IRequest{
				key:     "1234-1234-1234-1234-1234-1234",
				options: map[string]Option{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &IRequest{
				key:     tt.fields.key,
				options: map[string]Option{},
			}
			if got := r.SetIidAuth(tt.args.in0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IRequest.SetIidAuth() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestIRequest_GetIidAuth(t *testing.T) {
	type fields struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "empty",
			fields: fields{
				key: "empty",
			},
			want: "empty",
		},
		{
			name: "empty",
			fields: fields{
				key: "",
			},
			want: "",
		},
		{
			name: "short",
			fields: fields{
				key: "1",
			},
			want: "1",
		},
		{
			name: "longer",
			fields: fields{
				key: "1234-1234-1234-1234-1234-1234",
			},
			want: "1234-1234-1234-1234-1234-1234",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := IRequest{
				key: tt.fields.key,
			}
			if got := r.GetIidAuth(); got != tt.want {
				t.Errorf("IRequest.GetIidAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIRequest_Options(t *testing.T) {
	type fields struct {
		key     string
		options Options
	}
	tests := []struct {
		name   string
		fields fields
		want   Options
	}{
		{
			name: "one option",
			fields: fields{
				key: "1234",
				options: Options{
					"v": IOption{commandName: "v"},
				},
			},
			want: Options{
				"v": IOption{commandName: "v"},
			},
		},
		{
			name: "two option",
			fields: fields{
				key: "1234",
				options: Options{
					"v": IOption{commandName: "v"},
					"c": IOption{commandName: "c"},
				},
			},
			want: Options{
				"v": IOption{commandName: "v"},
				"c": IOption{commandName: "c"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := IRequest{
				key:     tt.fields.key,
				options: tt.fields.options,
			}
			if got := r.Options(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IRequest.Options() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIRequest_SetOption(t *testing.T) {

	v := IOption{commandName: "v"}
	c := IOption{commandName: "c"}

	type fields struct {
		key     string
		options Options
	}
	type args struct {
		in0 Option
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   IidRequest
	}{
		{
			name: "add empty option",
			fields: fields{
				key: "",
				options: Options{
					"v": v,
				},
			},
			args: args{
				in0: nil,
			},
			want: &IRequest{
				key: "",
				options: Options{
					"v": v,
				},
			},
		},
		{
			name: "add one option",
			fields: fields{
				key: "",
				options: Options{
					"v": v,
				},
			},
			args: args{
				in0: c,
			},
			want: &IRequest{
				key: "",
				options: Options{
					"v": v,
					"c": c,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &IRequest{
				key: tt.fields.key,
			}
			if (tt.fields.options) != nil {
				r.options = tt.fields.options
			}
			if got := r.SetOption(tt.args.in0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IRequest.SetOption() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestIRequest_String(t *testing.T) {
	v := IOption{commandName: "v"}
	c := IOption{commandName: "c"}
	s := IOption{commandName: "s"}

	type fields struct {
		key     string
		options Options
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple request",
			fields: fields{
				key: "empty",
			},
			want: "empty",
		},
		{
			name: "simple key request",
			fields: fields{
				key: "1234-1234",
			},
			want: "key=1234-1234",
		},
		{
			name: "simple key request with one option",
			fields: fields{
				key: "1234-1234",
				options: Options{
					"v": v,
				},
			},
			want: "key=1234-1234 options=v",
		},
		{
			name: "simple key request with two options",
			fields: fields{
				key: "1234-1234",
				options: Options{
					"v": v,
					"c": c,
				},
			},
			want: "key=1234-1234 options=cv",
		},
		{
			name: "empty with three options",
			fields: fields{
				options: Options{
					"v": v,
					"c": c,
					"s": s,
				},
			},
			want: "empty options=csv",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := IRequest{
				key: tt.fields.key,
			}
			if (tt.fields.options) != nil {
				r.options = tt.fields.options
			}
			if got := r.String(); got != tt.want {
				t.Errorf("IRequest.String() = %v, want %v", got, tt.want)
			}

			if got := r.GetHeader(); got != XINSTANCEID+": "+tt.want {
				t.Errorf("IRequest.GetHeader() = %v, want %v", got, XINSTANCEID+": "+tt.want)
			}

		})
	}
}

func TestIRequest_parseIidRequest(t *testing.T) {
	type fields struct {
		key string
	}
	type args struct {
		in0 string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   IidRequest
	}{
		{
			name: "empty",
			fields: fields{
				key: "",
			},
			args: args{
				in0: "empty",
			},
			want: &IRequest{
				key:     "empty",
				options: map[string]Option{},
			},
		},
		{
			name: "empty=;",
			fields: fields{
				key: "",
			},
			args: args{
				in0: "empty=;",
			},
			want: &IRequest{
				key:     "empty",
				options: map[string]Option{},
			},
		},
		{
			name: "empty=",
			fields: fields{
				key: "",
			},
			args: args{
				in0: "empty=",
			},
			want: &IRequest{
				key:     "empty",
				options: map[string]Option{},
			},
		},
		{
			name: "no key",
			fields: fields{
				key: "",
			},
			args: args{
				in0: "empty",
			},
			want: &IRequest{
				key:     "empty",
				options: map[string]Option{},
			},
		},
		{
			name: "no key;",
			fields: fields{
				key: "",
			},
			args: args{
				in0: "empty;",
			},
			want: &IRequest{
				key:     "empty",
				options: map[string]Option{},
			},
		}, {
			name: "simple key",
			fields: fields{
				key: "",
			},
			args: args{
				in0: `key="asdf"`,
			},
			want: &IRequest{
				key:     `"asdf"`,
				options: map[string]Option{},
			},
		},
		{
			name: "simple key",
			fields: fields{
				key: "empty",
			},
			args: args{
				in0: `key="asdf"`,
			},
			want: &IRequest{
				key:     `"asdf"`,
				options: map[string]Option{},
			},
		},
		{
			name: "two keys",
			fields: fields{
				key: "",
			},
			args: args{
				in0: `key=asdf key=jkl`,
			},
			want: &IRequest{
				key:     `jkl`,
				options: map[string]Option{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := IRequest{
				key: tt.fields.key,
			}
			if got := r.parseIidRequest(tt.args.in0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IRequest.parseIidRequest() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestIOption_Command(t *testing.T) {
	type fields struct {
		commandName string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := IOption{
				commandName: tt.fields.commandName,
			}
			if got := o.Command(); got != tt.want {
				t.Errorf("IOption.Command() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIRequestFromValue(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want *IRequest
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIRequestFromValue(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIRequestFromValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseOption(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want Options
	}{
		{
			name: "v",
			args: args{
				id: "v",
			},
			want: Options{
				"v": IOption{commandName: "v"},
			},
		},
		{
			name: "vv",
			args: args{
				id: "vv",
			},
			want: Options{
				"v": IOption{commandName: "v"},
			},
		},
		{
			name: "vc",
			args: args{
				id: "vc",
			},
			want: Options{
				"v": IOption{commandName: "v"},
				"c": IOption{commandName: "c"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseOption(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseOption() = %v, want %v", got, tt.want)
			}
		})
	}
}
