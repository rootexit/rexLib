package rexUtils

import (
	"reflect"
	"sort"
	"testing"
)

func TestDiffUintSlice(t *testing.T) {
	type args struct {
		oldList []uint
		newList []uint
	}
	tests := []struct {
		name         string
		args         args
		wantToDelete []uint
		wantToAdd    []uint
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				oldList: []uint{1, 2, 3, 4, 5},
				newList: []uint{4, 5, 6, 7},
			},
			wantToDelete: []uint{1, 2, 3},
			wantToAdd:    []uint{6, 7},
		},
		{
			name: "test2",
			args: args{
				oldList: []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
				newList: []uint{2, 3, 5, 11, 12, 16, 17},
			},
			wantToDelete: []uint{1, 4, 6, 7, 8, 9, 10, 13, 14, 15},
			wantToAdd:    []uint{16, 17},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToDelete, gotToAdd := DiffUintSlice(tt.args.oldList, tt.args.newList)
			sort.Slice(gotToDelete, func(i, j int) bool {
				return gotToDelete[i] < gotToDelete[j]
			})
			sort.Slice(gotToAdd, func(i, j int) bool {
				return gotToAdd[i] < gotToAdd[j]
			})
			if !reflect.DeepEqual(gotToDelete, tt.wantToDelete) {
				t.Errorf("DiffUintSlice() gotToDelete = %v, want %v", gotToDelete, tt.wantToDelete)
			}
			if !reflect.DeepEqual(gotToAdd, tt.wantToAdd) {
				t.Errorf("DiffUintSlice() gotToAdd = %v, want %v", gotToAdd, tt.wantToAdd)
			}
		})
	}
}
