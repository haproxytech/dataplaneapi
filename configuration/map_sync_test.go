// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package configuration

import (
	"strconv"
	"testing"

	"github.com/haproxytech/client-native/v5/models"
)

func data(differentAtIndex ...int) (fileEntries models.MapEntries, runtimeEntries models.MapEntries) {
	for i := 0; i < 50; i++ {
		fe := &models.MapEntry{Key: "k" + strconv.Itoa(i), Value: "v" + strconv.Itoa(i)}
		re := &models.MapEntry{Key: "k" + strconv.Itoa(i), Value: "v" + strconv.Itoa(i)}
		fileEntries = append(fileEntries, fe)
		runtimeEntries = append(runtimeEntries, re)
	}
	if len(differentAtIndex) > 0 {
		fileEntries[differentAtIndex[0]].Key = "abc"
		// fileEntries[differentAtIndex[0]].Value = "abc"
	}
	return fileEntries, runtimeEntries
}

func Test_equalSomeEntries(t *testing.T) {
	index := 10
	feDifferent, re := data(index)
	feSame, reSame := data()

	type args struct {
		fEntries models.MapEntries
		rEntries models.MapEntries
		index    []int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Same Entries",
			args: args{
				fEntries: models.MapEntries{&models.MapEntry{Key: "1", Value: "1"}},
				rEntries: models.MapEntries{&models.MapEntry{Key: "1", Value: "1"}},
				index:    []int{0},
			},
			want: true,
		},
		{
			name: "Same Entries 2",
			args: args{
				fEntries: feSame,
				rEntries: reSame,
			},
			want: true,
		},
		{
			name: "Both Empty",
			args: args{
				fEntries: models.MapEntries{},
				rEntries: models.MapEntries{},
			},
			want: true,
		},
		{
			name: "Different Entries",
			args: args{
				fEntries: models.MapEntries{&models.MapEntry{Key: "1", Value: "1"}},
				rEntries: models.MapEntries{&models.MapEntry{Key: "0", Value: "0"}},
				index:    []int{0},
			},
			want: false,
		},
		{
			name: "Different Length",
			args: args{
				fEntries: models.MapEntries{&models.MapEntry{Key: "1", Value: "1"}},
				rEntries: models.MapEntries{},
			},
			want: false,
		},
		{
			name: "Different At Index",
			args: args{
				fEntries: feDifferent,
				rEntries: re,
				index:    []int{index},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := equalSomeEntries(tt.args.fEntries, tt.args.rEntries, tt.args.index...); got != tt.want {
				t.Errorf("equalSomeEntries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_equal(t *testing.T) {
	index := 25
	feDifferent, re := data(index)
	feSame, reSame := data()

	type args struct {
		fEntries models.MapEntries
		rEntries models.MapEntries
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Same Entries",
			args: args{
				fEntries: models.MapEntries{&models.MapEntry{Key: "1", Value: "1"}, &models.MapEntry{Key: "2", Value: "2"}},
				rEntries: models.MapEntries{&models.MapEntry{Key: "1", Value: "1"}, &models.MapEntry{Key: "2", Value: "2"}},
			},
			want: true,
		},
		{
			name: "Same Entries 2",
			args: args{
				fEntries: feSame,
				rEntries: reSame,
			},
			want: true,
		},
		{
			name: "Both Empty",
			args: args{
				fEntries: models.MapEntries{},
				rEntries: models.MapEntries{},
			},
			want: true,
		},
		{
			name: "Different Entries",
			args: args{
				fEntries: models.MapEntries{&models.MapEntry{Key: "1", Value: "1"}},
				rEntries: models.MapEntries{&models.MapEntry{Key: "0", Value: "0"}},
			},
			want: false,
		},
		{
			name: "Different Length",
			args: args{
				fEntries: models.MapEntries{&models.MapEntry{Key: "1", Value: "1"}},
				rEntries: models.MapEntries{},
			},
			want: false,
		},
		{
			name: "Different At Index",
			args: args{
				fEntries: feDifferent,
				rEntries: re,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := equal(tt.args.fEntries, tt.args.rEntries); got != tt.want {
				t.Errorf("equalHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
