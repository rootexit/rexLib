package rexUtils

import (
	"reflect"
	"sort"
	"testing"

	"github.com/rootexit/rexLib/rexCommon"
)

func TestDiffStringSlice(t *testing.T) {
	type args struct {
		oldList []string
		newList []string
	}
	tests := []struct {
		name         string
		args         args
		wantToDelete []string
		wantToAdd    []string
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				oldList: []string{"1", "2", "3", "4", "5"},
				newList: []string{"4", "5", "6", "7"},
			},
			wantToDelete: []string{"1", "2", "3"},
			wantToAdd:    []string{"6", "7"},
		},
		{
			name: "test2",
			args: args{
				oldList: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15"},
				newList: []string{"2", "3", "5", "11", "12", "16", "17"},
			},
			wantToDelete: []string{"1", "4", "6", "7", "8", "9", "10", "13", "14", "15"},
			wantToAdd:    []string{"16", "17"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToDelete, gotToAdd := DiffStringSlice(tt.args.oldList, tt.args.newList)
			sort.Slice(gotToDelete, func(i, j int) bool {
				return rexCommon.Str2Int(gotToDelete[i]) < rexCommon.Str2Int(gotToDelete[j])
			})
			sort.Slice(gotToAdd, func(i, j int) bool {
				return rexCommon.Str2Int(gotToAdd[i]) < rexCommon.Str2Int(gotToAdd[j])
			})
			if !reflect.DeepEqual(gotToDelete, tt.wantToDelete) {
				t.Errorf("DiffStringSlice() gotToDelete = %v, want %v", gotToDelete, tt.wantToDelete)
			}
			if !reflect.DeepEqual(gotToAdd, tt.wantToAdd) {
				t.Errorf("DiffStringSlice() gotToAdd = %v, want %v", gotToAdd, tt.wantToAdd)
			}
		})
	}
}
