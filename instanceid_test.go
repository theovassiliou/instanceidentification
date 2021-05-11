package instanceid

import (
	"io/ioutil"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
)

func Test_splitOnPlus(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"one plain arg",
			args{" a"},
			[]string{" a"},
		},
		{
			"one plain arg with spaces",
			args{" a"},
			[]string{" a"},
		},
		{
			"one plain arg longer name",
			args{"someId"},
			[]string{"someId"},
		},
		{
			"two plain args",
			args{"a+b"},
			[]string{"a", "b"},
		},
		{
			"two longer plain args",
			args{"one+two"},
			[]string{"one", "two"},
		},

		{
			"three plain args",
			args{"a+b+c"},
			[]string{"a", "b", "c"},
		},
		{
			"three longer args",
			args{"one+two+three"},
			[]string{"one", "two", "three"},
		},
		{
			"three plain args with ()-1",
			args{"a+b+c()"},
			[]string{"a", "b", "c()"},
		},
		{
			"4 plain args with",
			args{"a+b+c+d"},
			[]string{"a", "b", "c", "d"},
		},
		{
			"three plain args with ()2",
			args{"a+(b()+c)"},
			[]string{"a", "(b()+c)"},
		},
		{
			"three structured args with ()",
			args{"a+(b()+c+ff(xx+zz))"},
			[]string{"a", "(b()+c+ff(xx+zz))"},
		},
		{
			"nested sum",
			args{"(a+b)"},
			[]string{"(a+b)"},
		},
		{
			"2 nested sum",
			args{"(x+y)+(a+b)"},
			[]string{"(x+y)", "(a+b)"},
		},
		{
			"invalid +",
			args{"+a"},
			[]string{"a"},
		},
		{
			"invalid postfix+",
			args{"a+"},
			[]string{"a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitOnPlus(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitOnPlus() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func Test_parseMIIDFromFile(t *testing.T) {
	type args struct {
		fileName string
	}

	log.SetLevel(log.TraceLevel)
	tests := []struct {
		name     string
		args     args
		wantMiid Miid
	}{
		{
			"complex1",
			args{"test/instanceId1.txt"},
			&StdMiid{},
		},
		{
			"complex2",
			args{"test/miid1.txt"},
			&StdMiid{
				sn: "MsA",
				vn: "1.1",
				va: "xxx",
				t:  22,
			},
		},
		{
			"Nonsens",
			args{"README.md"},
			&StdMiid{
				sn: "",
				vn: "",
				va: "",
				t:  0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := ioutil.ReadFile(tt.args.fileName)
			if err != nil {
				t.Fatal(err)
			}
			if gotMiid := parseMIID(string(file)); !reflect.DeepEqual(gotMiid, tt.wantMiid) {
				t.Errorf("parseMIID() = %v, want %v", gotMiid, tt.wantMiid)
			}
		})
	}
}

func Test_parseCiid(t *testing.T) {
	log.SetLevel(log.TraceLevel)

	A := &StdMiid{
		sn: "A",
		vn: "1.1",
		va: "",
		t:  22,
	}

	B := &StdMiid{
		sn: "B",
		vn: "1.1",
		va: "",
		t:  22,
	}
	C := &StdMiid{
		sn: "C",
		vn: "1.1",
		va: "",
		t:  22,
	}
	D := &StdMiid{
		sn: "D",
		vn: "1.1",
		va: "",
		t:  22,
	}
	E := &StdMiid{
		sn: "E",
		vn: "1.1",
		va: "",
		t:  22,
	}

	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantCiid Ciid
	}{
		{
			"simple",
			args{
				"msA/1.17/dev-123ab%3333s",
			},
			&StdCiid{
				miid: &StdMiid{
					sn: "msA",
					vn: "1.17",
					va: "dev-123ab",
					t:  3333,
				},
			},
		},
		{
			"one Call",
			args{
				"msA/1.17/dev-123ab%3333s(A/1.1%22s)",
			},
			&StdCiid{
				miid: &StdMiid{sn: "msA", vn: "1.17", va: "dev-123ab", t: 3333},
				ciids: Stack{
					&StdCiid{
						miid: &StdMiid{
							sn: "A",
							vn: "1.1",
							va: "",
							t:  22,
						},
					},
				},
			},
		},
		{
			"one Call Plus another one",
			args{
				"msA/1.17/dev-123ab%3333s(A/1.1%22s+B/1.1%22s)",
			},
			&StdCiid{
				miid: &StdMiid{sn: "msA", vn: "1.17", va: "dev-123ab", t: 3333},
				ciids: Stack{
					&StdCiid{
						miid: &StdMiid{
							sn: "A",
							vn: "1.1",
							va: "",
							t:  22,
						},
					},
					&StdCiid{
						miid: &StdMiid{
							sn: "B",
							vn: "1.1",
							va: "",
							t:  22,
						},
					},
				},
			},
		},
		{
			"one Call Plus another one and one call",
			args{
				"msA/1.17/dev-123ab%3333s(A/1.1%22s+B/1.1%22s(C/1.1%22s))",
			},
			&StdCiid{
				miid: &StdMiid{sn: "msA", vn: "1.17", va: "dev-123ab", t: 3333},
				ciids: Stack{
					&StdCiid{
						miid: &StdMiid{
							sn: "A",
							vn: "1.1",
							va: "",
							t:  22,
						},
					},
					&StdCiid{
						miid: &StdMiid{
							sn: "B",
							vn: "1.1",
							va: "",
							t:  22,
						},
						ciids: Stack{
							&StdCiid{
								miid: &StdMiid{
									sn: "C",
									vn: "1.1",
									va: "",
									t:  22,
								},
							},
						},
					},
				},
			},
		},
		{
			"simple",
			args{
				"A/1.1%22s(B/1.1%22s+C/1.1%22s(D/1.1%22s+E/1.1%22s))",
			},
			&StdCiid{miid: A, ciids: Stack{
				&StdCiid{miid: B},
				&StdCiid{miid: C, ciids: Stack{
					&StdCiid{miid: D},
					&StdCiid{miid: E},
				}}}},
		},
		{
			"simple",
			args{
				"A/1.1%22s(B/1.1%22s)",
			},
			&StdCiid{miid: A,
				ciids: Stack{
					&StdCiid{miid: B},
				},
			},
		},
		{
			"simple",
			args{
				"A/1.1%22s(B/1.1%22s+C/1.1%22s)",
			},
			&StdCiid{miid: A,
				ciids: Stack{
					&StdCiid{miid: B},
					&StdCiid{miid: C},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCiid := parseCiid(tt.args.id); !reflect.DeepEqual(gotCiid, tt.wantCiid) {
				t.Errorf("parseCiid() = %v, want %v", gotCiid, tt.wantCiid)
			}
		})
	}
}
func Test_seperateFNameFromArg(t *testing.T) {
	type args struct {
		signature string
	}
	tests := []struct {
		name     string
		args     args
		wantName string
		wantArg  string
	}{
		{
			"simple",
			args{"A(B)"},
			"A",
			"B",
		},
		{
			"no Arg",
			args{"A"},
			"A", "",
		},
		{
			"empty Parenthesis",
			args{"A()"},
			"A",
			"",
		},
		{
			"no name",
			args{"(B)"},
			"",
			"B",
		},
		{
			"more complex",
			args{"A(B+C)"},
			"A",
			"B+C",
		},
		{
			"more complex, neste functions",
			args{"A(B(D)+C(E(F)))"},
			"A",
			"B(D)+C(E(F))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotArg := seperateFNameFromArg(tt.args.signature)
			if gotName != tt.wantName {
				t.Errorf("seperateFNameFromArg() gotName = %v, want %v", gotName, tt.wantName)
			}
			if gotArg != tt.wantArg {
				t.Errorf("seperateFNameFromArg() gotArg = %v, want %v", gotArg, tt.wantArg)
			}
		})
	}
}

/* func Test_printCiid(t *testing.T) {
	type args struct {
		ciid Ciid
	}
	log.SetLevel(log.TraceLevel)
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"simple",
			args{parseCiid("A/1.1%22s(B/1.1%22s+C/1.1%22s)")},
			79,
		},
		{
			"simple2",
			args{parseCiid("A/1.1%22s(B/1.1%22s+C/1.1%22s(D/1.1%22s))")},
			110,
		},
		{
			"iid1",
			args{parseCiid("msA/1.1/abs%22s(msB/2.0/xxxx%333s+C/1.1%22s(D/1.1%22s))")},
			115,
		},
		{
			"iid2",
			args{parseCiid("msA/1.1/abs%22s(msB/2.0/xxxx%333s+msC/1.1%22s(D/1.1%22s))")},
			117,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(tt.args.ciid.PrintCiid()); got != tt.want {
				t.Errorf("printCiid() = %v, want %v", got, tt.want)
				t.Errorf("theTree() = \n%v", tt.args.ciid.PrintCiid())
			}
		})
	}
} */

func TestCiid_String(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"simpleMiid",
			"msA/1.1%22s",
		},
		{
			"simpleMiid zero seconds",
			"SS/1.2/YY%0s",
		},
		{
			"simpleMiid minus seconds",
			"SS/1.2/YY%-1s",
		},
		{
			"fullMiid",
			"msA/1.1/feature-branch-22aabbcc%22s",
		},
		{
			"emptyMiid",
			"",
		},
		{
			"justSimpleMidd",
			"A/1%22s",
		},
		{
			"fullMiidOneCiid",
			"msA/1.1/feature-branch-22aabbcc%22s(msB/2.2%33s)",
		},
		{
			"fullMiidTwoCiid",
			"msA/1.1/feature-branch-22aabbcc%22s(msB/xx%333s+msC/222%444s)",
		},
		{
			"complexFunc",
			"A/1.1%22s(B/1.1%22s(C/1.1%22s+D/1.1%22s)+D/1.1%22s(E/1.1%22s))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewStdCiid(tt.want)
			if got := mock.String(); got != tt.want {
				t.Errorf("Ciid.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCiid(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantCiid Ciid
	}{
		{
			"simpleMiid",
			args{"msA/1.1%22s"},
			&StdCiid{
				miid: &StdMiid{
					sn: "msA",
					vn: "1.1",
					va: "",
					t:  22,
				},
				ciids: nil,
			},
		},
		{
			"simpleMiid2",
			args{"SS/1.2/YY%0s"},
			&StdCiid{
				miid: &StdMiid{
					sn: "SS",
					vn: "1.2",
					va: "YY",
					t:  0,
				},
				ciids: nil,
			},
		},
		{
			"fullMiid",
			args{"msA/1.1/feature-branch-22aabbcc%22s"},
			&StdCiid{
				miid: &StdMiid{
					sn: "msA",
					vn: "1.1",
					va: "feature-branch-22aabbcc",
					t:  22,
				},
				ciids: nil,
			},
		},
		{
			"emptyMiid",
			args{""},
			&StdCiid{
				miid: &StdMiid{
					sn: "",
					vn: "",
					va: "",
					t:  0,
				},
				ciids: nil,
			},
		},
		{
			"fullMiidOneCiid",
			args{"msA/1.1/feature-branch-22aabbcc%22s(msB/1.1%22s)"},
			&StdCiid{
				miid: &StdMiid{
					sn: "msA",
					vn: "1.1",
					va: "feature-branch-22aabbcc",
					t:  22,
				},
				ciids: Stack{
					&StdCiid{
						miid: &StdMiid{
							sn: "msB",
							vn: "1.1",
							va: "",
							t:  22,
						},
						ciids: nil,
					},
				},
			},
		},
		{
			"fullMiidTwoCiid",
			args{"msA/1.1/feature-branch-22aabbcc%22s(msB/1.1%22s+msC/1.1%22s)"},
			&StdCiid{
				miid: &StdMiid{
					sn: "msA",
					vn: "1.1",
					va: "feature-branch-22aabbcc",
					t:  22,
				},
				ciids: Stack{
					&StdCiid{
						miid: &StdMiid{
							sn: "msB",
							vn: "1.1",
							va: "",
							t:  22,
						},
						ciids: nil,
					},
					&StdCiid{
						miid: &StdMiid{
							sn: "msC",
							vn: "1.1",
							va: "",
							t:  22,
						},
						ciids: nil,
					},
				},
			},
		},
		{
			"complexFunc",
			args{"A/1.1%22s(B/1.1%22s(C/1.1%22s+D/1.1%22s)+D/1.1%22s(E/1.1%22s)"},
			&StdCiid{
				miid: &StdMiid{
					sn: "A",
					vn: "1.1",
					va: "",
					t:  22,
				},
				ciids: Stack{
					&StdCiid{
						miid: &StdMiid{
							sn: "B",
							vn: "1.1",
							va: "",
							t:  22,
						},
						ciids: Stack{
							&StdCiid{
								miid: &StdMiid{
									sn: "C",
									vn: "1.1",
									va: "",
									t:  22,
								},
								ciids: nil,
							},
							&StdCiid{
								miid: &StdMiid{
									sn: "D",
									vn: "1.1",
									va: "",
									t:  22,
								},
								ciids: nil,
							},
						},
					},
					&StdCiid{
						miid: &StdMiid{
							sn: "D",
							vn: "1.1",
							va: "",
							t:  22,
						},
						ciids: Stack{
							&StdCiid{
								miid: &StdMiid{
									sn: "E",
									vn: "1.1",
									va: "",
									t:  22,
								},
								ciids: nil,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCiid := NewStdCiid(tt.args.id); !reflect.DeepEqual(gotCiid, tt.wantCiid) {
				t.Errorf("NewCiid() = %#v, want %#v", gotCiid, tt.wantCiid)
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

func TestCiid_Contains(t *testing.T) {

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
			fields: "MsA/1.1/xxx%22s(msC/1.4%5555s+msD/2.2%23234s)",
			args: args{
				s: "MsA/1.1",
			},
			want: true,
		},
		{
			name:   "correct",
			fields: "MsA/1.1/xxx%22s(msC/1.4%5555s+msD/2.2%23234s)",
			args: args{
				s: "msC",
			},
			want: true,
		},
		{
			name:   "correct",
			fields: "MsA/1.1/xxx%22s(msC/1.4%5555s+msD/2.2%23234s)",
			args: args{
				s: "msC/1.4",
			},
			want: true,
		},
		{
			name:   "correct",
			fields: "MsA/1.1/xxx%22s(msC/1.4%5555s+msD/2.2%23234s)",
			args: args{
				s: "msC/1.3",
			},
			want: false,
		},
		{
			name:   "correct",
			fields: "MsA/1.1/xxx%22s(msC/1.4%5555s+msD/2.2%23234s)",
			args: args{
				s: "msD/1.4",
			},
			want: false,
		},
		{
			name:   "correct",
			fields: "MsA/1.1/xxx%22s(msC/1.4%5555s+msD/2.2%23234s)",
			args: args{
				s: "msD/2.2",
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
			m := NewStdCiid(tt.fields)
			if got := m.Contains(tt.args.s); got != tt.want {
				t.Errorf("Ciid.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
