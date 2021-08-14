package db

import (
	"reflect"
	"testing"
	"time"

	"github.com/go-pg/pg/urlvalues"
)

type testStruct struct {
	tableName struct{} `sql:"tests"`

	TestId     int64  `json:"id" sql:"id,pk"`
	SomeNumber int64  `json:"someNumber" sql:"some_number,unique"`
	TestSelf   string `json:"inviteSelf" sql:"-"`

	CodeName string `json:"codeName" sql:"code_name"`
	SendTest bool   `json:"sendTest" sql:"is_send_test,notnull"`
	DbField  string `json:"-" sql:"db_field"`

	//  standard DB fields
	CreatedAt  time.Time  `json:"createdAt" sql:"created_at,default:now()"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty" sql:"updated_at"`
	ArchivedAt *time.Time `json:"archivedAt,omitempty" sql:"archived_at"`
}

func TestPrepareFiltersByModel(t *testing.T) {
	tests := []struct {
		name         string
		filters      urlvalues.Values
		model        interface{}
		wantPrepared urlvalues.Values
		wantErr      bool
	}{
		{
			name: "good case",
			filters: urlvalues.Values{
				"codeName":       []string{"test"},
				"someNumber__gt": []string{"100"},
			},
			model: testStruct{},
			wantPrepared: urlvalues.Values{
				"code_name":       []string{"test"},
				"some_number__gt": []string{"100"},
			},
			wantErr: false,
		},
		{
			name:    "nil model case",
			model:   nil,
			wantErr: true,
		},
		{
			name: "unknown field case",
			filters: urlvalues.Values{
				"unknown": []string{"test"},
			},
			model:        testStruct{},
			wantPrepared: urlvalues.Values{},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrepared, err := PrepareFiltersByModel(tt.filters, tt.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrepareFiltersByModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPrepared, tt.wantPrepared) {
				t.Errorf("PrepareFiltersByModel() = %v, want %v", gotPrepared, tt.wantPrepared)
			}
		})
	}
}
