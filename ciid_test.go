package instanceid

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
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

func TestReversibleCiid(t *testing.T) {
	fileName := "test/iidtestsetValid.txt"
	file, err := os.Open(fileName)

	if err != nil {
		t.Errorf("failed to open: %v", fileName)

	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	i := 0
	for scanner.Scan() {
		i++
		ttwant := scanner.Text()
		ttname := fmt.Sprintf("[%v:%v]", fileName, i)
		// Filtering # comments and empty lines
		if strings.HasPrefix(ttwant, "#") || ttwant == "" {
			continue
		}
		t.Run(ttname, func(t *testing.T) {
			m := NewStdCiid(ttwant)
			if got := m.String(); got != ttwant {
				t.Errorf("Iid not reversible = %v, want %v", got, ttwant)
			}
		})
	}

}

func TestInvalidCiid(t *testing.T) {
	fileName := "test/iidtestsetInvalid.txt"
	file, err := os.Open(fileName)

	if err != nil {
		t.Errorf("failed to open: %v", fileName)

	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	i := 0
	for scanner.Scan() {
		i++
		invalidIid := scanner.Text()
		ttname := fmt.Sprintf("[%v:%v]", fileName, i)
		// Filtering # comments and empty lines
		if strings.HasPrefix(invalidIid, "#") || invalidIid == "" {
			continue
		}
		t.Run(ttname, func(t *testing.T) {
			m := NewStdCiid(invalidIid)
			if got := m.String(); got != "" {
				t.Errorf("Iid %v should not be parseble. Parsed to %v", invalidIid, got)
			}
		})
	}

}
