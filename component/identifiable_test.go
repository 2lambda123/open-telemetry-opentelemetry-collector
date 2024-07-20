// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package component

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalText(t *testing.T) {
	id := NewIDWithName(MustNewType("test"), MustNewName("name"))
	got, err := id.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, id.String(), string(got))
}

func TestUnmarshalText(t *testing.T) {
	validType := MustNewType("valid_type")
	validName := MustNewName("valid_name")
	var testCases = []struct {
		idStr       string
		expectedErr bool
		expectedID  ID
	}{
		{
			idStr:      "valid_type",
			expectedID: ID{typeVal: validType, nameVal: Name{}},
		},
		{
			idStr:      "valid_type/valid_name",
			expectedID: ID{typeVal: validType, nameVal: validName},
		},
		{
			idStr:      "   valid_type   /   valid_name  ",
			expectedID: ID{typeVal: validType, nameVal: validName},
		},
		{
			idStr:       "/valid_name",
			expectedErr: true,
		},
		{
			idStr:       "     /valid_name",
			expectedErr: true,
		},
		{
			idStr:       "valid_type/",
			expectedErr: true,
		},
		{
			idStr:       "valid_type/      ",
			expectedErr: true,
		},
		{
			idStr:       "      ",
			expectedErr: true,
		},
		{
			idStr:       "valid_type/not valid name",
			expectedErr: true,
		},
		{
			idStr:       "not valid type/valid_name",
			expectedErr: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.idStr, func(t *testing.T) {
			id := ID{}
			err := id.UnmarshalText([]byte(test.idStr))
			if test.expectedErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.expectedID, id)
			assert.Equal(t, test.expectedID.Type(), id.Type())
			assert.Equal(t, test.expectedID.Name(), id.Name())
			assert.Equal(t, test.expectedID.String(), id.String())
		})
	}
}
