package helpers

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

type (
	testSourceStruct struct {
		A *string `json:"field_a"`
		B *bool   `json:"field_b"`
		C *int    `json:"field_c"`
	}
	testTargetStruct struct {
		A string `json:"field_a"`
		B bool   `json:"field_b"`
		C int    `json:"field_c"`
	}

	privateFieldSourceStruct struct {
		a *string
	}
	privateFieldTargetStruct struct {
		a string
	}

	benchSourceStruct struct {
		A0 *string
		B0 *bool
		C0 *int
		A1 *string
		B1 *bool
		C1 *int
		A2 *string
		B2 *bool
		C2 *int
		A3 *string
		B3 *bool
		C3 *int
		A4 *string
		B4 *bool
		C4 *int
		A5 *string
		B5 *bool
		C5 *int
	}
	benchTargetStruct struct {
		A0 string `json:"field_a0"`
		B0 bool   `json:"field_b0"`
		C0 int    `json:"field_c0"`
		A1 string `json:"field_a1"`
		B1 bool   `json:"field_b1"`
		C1 int    `json:"field_c1"`
		A2 string `json:"field_a2"`
		B2 bool   `json:"field_b2"`
		C2 int    `json:"field_c2"`
		A3 string `json:"field_a3"`
		B3 bool   `json:"field_b3"`
		C3 int    `json:"field_c3"`
		A4 string `json:"field_a4"`
		B4 bool   `json:"field_b4"`
		C4 int    `json:"field_c4"`
		A5 string `json:"field_a5"`
		B5 bool   `json:"field_b5"`
		C5 int    `json:"field_c5"`
	}
)

type StringAliasType string

const (
	StringAliasType1 StringAliasType = "1"
	StringAliasType2 StringAliasType = "2"
)

var (
	a, b, c, u = "world", true, 10, uint64(10)
	testSource = testSourceStruct{
		A: &a,
		B: &b,
		C: nil,
	}

	testTarget = testTargetStruct{
		A: "hello",
		B: false,
		C: c,
	}

	pTSource = privateFieldSourceStruct{
		a: &a,
	}

	pTTarget = privateFieldTargetStruct{
		a: "hello",
	}

	benchSource = benchSourceStruct{
		A0: &a,
		B0: &b,
		C0: nil,
		A1: &a,
		B1: &b,
		C1: nil,
		A2: &a,
		B2: &b,
		C2: nil,
		A3: &a,
		B3: &b,
		C3: nil,
		A4: &a,
		B4: &b,
		C4: nil,
		A5: &a,
		B5: &b,
		C5: nil,
	}

	benchTarget = benchTargetStruct{
		A0: "hello",
		B0: false,
		C0: c,
		A1: "hello",
		B1: false,
		C1: c,
		A2: "hello",
		B2: false,
		C2: c,
		A3: "hello",
		B3: false,
		C3: c,
		A4: "hello",
		B4: false,
		C4: c,
		A5: "hello",
		B5: false,
		C5: c,
	}
)

func TestNullableFieldsToStruct(t *testing.T) {
	type args struct {
		source interface{}
		target interface{}
	}
	tests := []struct {
		name        string
		args        args
		wantIsFound bool
		wantErr     bool
	}{
		{
			name: "good case",
			args: args{
				source: testSource,
				target: &testTarget,
			},
			wantIsFound: true,
			wantErr:     false,
		},
		{
			name: "target not a pointer case",
			args: args{
				source: testSource,
				target: testTarget,
			},
			wantIsFound: false,
			wantErr:     true,
		},
		{
			name: "source is not a struct case",
			args: args{
				source: c,
				target: &testTarget,
			},
			wantIsFound: false,
			wantErr:     true,
		},
		{
			name: "target is not a struct case",
			args: args{
				source: testSource,
				target: &c,
			},
			wantIsFound: false,
			wantErr:     true,
		},
		{
			name: "different fields case",
			args: args{
				source: struct {
					A *string
					B *bool
					C *int
					D *int
				}{
					A: &a,
					B: &b,
					C: &c,
					D: &c,
				},
				target: &testTarget,
			},
			wantIsFound: true,
			wantErr:     false,
		},
		{
			name: "cannot set case",
			args: args{
				source: pTSource,
				target: &pTTarget,
			},
			wantIsFound: true,
			wantErr:     true,
		},
		{
			name: "different types case",
			args: args{
				source: struct {
					A *string
					B *bool
					C *string
				}{
					A: &a,
					B: &b,
					C: &a,
				},
				target: &testTarget,
			},
			wantIsFound: true,
			wantErr:     true,
		},
		{
			name: "map field case",
			args: args{
				source: struct {
					A *map[string]string
				}{
					A: &map[string]string{"one": "1", "two": "two"},
				},
				target: &struct {
					A map[string]string
				}{},
			},
			wantIsFound: true,
			wantErr:     false,
		},
		{
			name: "nested struct field case",
			args: args{
				source: struct {
					A *struct {
						N string
					}
				}{
					A: &struct {
						N string
					}{
						N: "nested",
					},
				},
				target: &struct {
					A struct {
						N string
					}
				}{},
			},
			wantIsFound: true,
			wantErr:     false,
		},
		{
			name: "array field case",
			args: args{
				source: struct {
					A *[]string
				}{
					A: &[]string{"1", "2"},
				},
				target: &struct {
					A []string
				}{},
			},
			wantIsFound: true,
			wantErr:     false,
		},
		{
			name: "array field with type alias case",
			args: args{
				source: struct {
					A *[]StringAliasType
				}{
					A: &[]StringAliasType{StringAliasType1, StringAliasType2},
				},
				target: &struct {
					A []StringAliasType
				}{},
			},
			wantIsFound: true,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsFound, err := NullableFieldsToStruct(tt.args.source, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("FillFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsFound != tt.wantIsFound {
				t.Errorf("FillFields() = %v, want %v", gotIsFound, tt.wantIsFound)
			}
			//TODO check the fields somehow
			fmt.Printf("%+v\n", tt.args.target)
		})
	}
}

func BenchmarkSetFields(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = NullableFieldsToStruct(testSource, &testTarget)
	}
}

func BenchmarkSetFieldsManually(b *testing.B) {
	for n := 0; n < b.N; n++ {
		if testSource.A != nil {
			testTarget.A = *testSource.A
		}
		if testSource.B != nil {
			testTarget.B = *testSource.B
		}
		if testSource.C != nil {
			testTarget.C = *testSource.C
		}
	}
}

func BenchmarkSetFieldsALot(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = NullableFieldsToStruct(benchSource, &benchTarget)
	}
}

func BenchmarkSetFieldsManuallyALot(b *testing.B) {
	for n := 0; n < b.N; n++ {
		if benchSource.A0 != nil {
			benchTarget.A0 = *benchSource.A0
		}
		if benchSource.B0 != nil {
			benchTarget.B0 = *benchSource.B0
		}
		if benchSource.C0 != nil {
			benchTarget.C0 = *benchSource.C0
		}

		if benchSource.A1 != nil {
			benchTarget.A1 = *benchSource.A1
		}
		if benchSource.B1 != nil {
			benchTarget.B1 = *benchSource.B1
		}
		if benchSource.C1 != nil {
			benchTarget.C1 = *benchSource.C1
		}

		if benchSource.A2 != nil {
			benchTarget.A2 = *benchSource.A2
		}
		if benchSource.B2 != nil {
			benchTarget.B2 = *benchSource.B2
		}
		if benchSource.C2 != nil {
			benchTarget.C2 = *benchSource.C2
		}

		if benchSource.A3 != nil {
			benchTarget.A3 = *benchSource.A3
		}
		if benchSource.B3 != nil {
			benchTarget.B3 = *benchSource.B3
		}
		if benchSource.C3 != nil {
			benchTarget.C3 = *benchSource.C3
		}

		if benchSource.A4 != nil {
			benchTarget.A4 = *benchSource.A4
		}
		if benchSource.B4 != nil {
			benchTarget.B4 = *benchSource.B4
		}
		if benchSource.C4 != nil {
			benchTarget.C4 = *benchSource.C4
		}

		if benchSource.A5 != nil {
			benchTarget.A5 = *benchSource.A5
		}
		if benchSource.B5 != nil {
			benchTarget.B5 = *benchSource.B5
		}
		if benchSource.C5 != nil {
			benchTarget.C5 = *benchSource.C5
		}
	}
}

type TestModelStatus string
type TestModel struct {
	tableName struct{}        `sql:"?SHARD.jobs"`
	Id        int64           `json:"id" pg:"id,pk"`
	Status    TestModelStatus `json:"status" pg:"status" sql:"type: ?SHARD.jobs_status"`
}

func TestSetTagsSqlTypeShard(t *testing.T) {
	type args struct {
		shardId int64
		target  interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "base usage", args: struct {
			shardId int64
			target  interface{}
		}{shardId: 1, target: &TestModel{
			Id:     1,
			Status: "new",
		}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}
			var err error
			if result, err = SetTagsSqlTypeShard(tt.args.shardId, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("SetTagsSqlTypeShard() error = %v, wantErr %v", err, tt.wantErr)
			}
			targetFields := reflect.TypeOf(result)
			targetFieldsNum := targetFields.NumField()
			for i := 0; i < targetFieldsNum; i++ {
				field := targetFields.Field(i)
				tag := field.Tag
				sqlTag := tag.Get("sql")
				if len(sqlTag) > 0 {
					var re = regexp.MustCompile(`type:\s*(\?SHARD)`)
					if re.MatchString(sqlTag) {
						t.Errorf("SetTagsSqlTypeShard() is NOT replaced field %s tag %s with value %s", field.Name, tag, sqlTag)
					}
				}
			}
		})
	}
}
