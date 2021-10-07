package instanceid

import (
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
)

func Test_parseMIID(t *testing.T) {
	type args struct {
		id string
	}

	log.SetLevel(log.TraceLevel)
	tests := []struct {
		name     string
		args     args
		wantMiid Miid
	}{
		{
			"simple",
			args{"msA/1.17/dev-123ab%3333s"},
			&StdMiid{
				sn: "msA",
				vn: "1.17",
				va: "dev-123ab",
				t:  3333,
			},
		},
		{
			"simple with minus",
			args{"msA/1.17/dev-123ab%-1s"},
			&StdMiid{
				sn: "msA",
				vn: "1.17",
				va: "dev-123ab",
				t:  -1,
			},
		},
		{
			"simple-short",
			args{"msA/1.17%3333s"},
			&StdMiid{
				sn: "msA",
				vn: "1.17",
				t:  3333,
			},
		},
		{
			"simple-short with minus",
			args{"msA/1.17%-1s"},
			&StdMiid{
				sn: "msA",
				vn: "1.17",
				t:  -1,
			},
		},

		{
			"simple-notSecond",
			args{"msA/1.17%3333"},
			&StdMiid{},
		},
		{
			"simple-notSecondNumber",
			args{"msA/1.17%333a"},
			&StdMiid{},
		},
		{
			"toomanydelimiters",
			args{"msA/1.17/addInfo/surplusInfo%333s"},
			&StdMiid{
				sn: "msA",
				vn: "1.17",
				va: "addInfo/surplusInfo",
				t:  333,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCiid := parseMIID(tt.args.id); !reflect.DeepEqual(gotCiid, tt.wantMiid) {
				t.Errorf("parseMIID() = %#v, want %#v", gotCiid, tt.wantMiid)
			}
		})
	}
}
func TestNewMiid(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantMiid Miid
	}{
		{
			"CiidNotExpected",
			args{"msA/1.1/feature-branch-22aabbcc%22s(msB+msC)"},
			&StdMiid{},
		},
		{
			"simple",
			args{"msA/1.1%22s"},
			&StdMiid{
				sn: "msA",
				vn: "1.1",
				va: "",
				t:  22,
			},
		},
		{
			"simple with minux",
			args{"msA/1.1%-1s"},
			&StdMiid{
				sn: "msA",
				vn: "1.1",
				va: "",
				t:  -1,
			},
		},
		{
			"complex",
			args{"msA/1.1/asdfasdf-asdfasdf%22s"},
			&StdMiid{
				sn: "msA",
				vn: "1.1",
				va: "asdfasdf-asdfasdf",
				t:  22,
			},
		},
		{
			"complex with minus",
			args{"msA/1.1/asdfasdf-asdfasdf%-1s"},
			&StdMiid{
				sn: "msA",
				vn: "1.1",
				va: "asdfasdf-asdfasdf",
				t:  -1,
			},
		},
		{
			"no clue",
			args{"This is some text"},
			&StdMiid{},
		},
		{
			"no clue",
			args{"(/)"},
			&StdMiid{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMiid := NewStdMiid(tt.args.id); !reflect.DeepEqual(gotMiid, tt.wantMiid) {
				t.Errorf("NewMiid() = %v, want %v", gotMiid, tt.wantMiid)
			}
		})
	}
}

func TestSanityCheck(t *testing.T) {
	type args struct {
		miid string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"correct",
			args{"a/a%22s"},
			true,
		},
		{
			"correct max",
			args{"a/b/c%22s"},
			true,
		},
		{
			"no sec",
			args{"a/b/xx%22"},
			false,
		},
		{
			"to many /",
			args{"a/b/22/s/s/"},
			false,
		},
		{
			"minimal valid",
			args{"abs/1.1%22s"},
			true,
		},
		{
			"minimal valid with minus",
			args{"abs/1.1%-1s"},
			true,
		},
		{
			"no /",
			args{"ab22s"},
			false,
		},
		{
			"has +",
			args{"ab/1.1%22s+ab/1.1%22s"},
			false,
		},
		{
			"has (",
			args{"ab/1.1%22s(A)s"},
			false,
		},
		{
			"has )",
			args{"ab/1.1%22s(A)+s"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanityCheck(tt.args.miid); got != tt.want {
				t.Errorf("SanityCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiid_Contains(t *testing.T) {

	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields string
		args   args
		want   bool
	}{
		{
			name:   "correct",
			fields: "msA/1.1/dev-1234%22s",
			args: args{
				s: "msA/1.1",
			},
			want: true,
		},
		{
			name:   "correct",
			fields: "msA/1.1/dev-1234%22s",
			args: args{
				s: "msA/1.1/dev-1234",
			},
			want: true,
		},
		{
			name:   "wrong service",
			fields: "msA/1.1/dev-1234%22s",
			args: args{
				s: "msB/1.1",
			},
			want: false,
		},
		{
			name:   "wrong version",
			fields: "msA/1.1/dev-1234%22s",
			args: args{
				s: "msA/1.2",
			},
			want: false,
		},
		{
			name:   "wrong service, correct dev",
			fields: "msA/1.1/dev-1234%22s",
			args: args{
				s: "dev-1234",
			},
			want: true,
		},
		{
			name:   "empty",
			fields: "msA/1.1/dev-1234%22s",
			args: args{
				s: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewStdMiid(tt.fields)
			if got := m.Contains(tt.args.s); got != tt.want {
				t.Errorf("Miid.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
