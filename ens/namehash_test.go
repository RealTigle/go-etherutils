// Copyright 2017 Orinoco Payments
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

package ens

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNameHashEmpty(t *testing.T) {
	expected := "0000000000000000000000000000000000000000000000000000000000000000"
	actual, err := NameHash("")
	assert.Nil(t, err, "Failed to hash")
	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
}

func TestNameHashTLD(t *testing.T) {
	expected := "93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae"
	actual, err := NameHash("eth")
	assert.Nil(t, err, "Failed to hash")
	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
}

func TestNameHashDottedTLD(t *testing.T) {
	expected := "8cc9f31a5e7af6381efc751d98d289e3f3589f1b6f19b9b989ace1788b939cf7"
	actual, err := NameHash(".eth")
	assert.Nil(t, err, "Failed to hash")
	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
}

func TestNameHashResolverEth(t *testing.T) {
	expected := "fdd5d5de6dd63db72bbc2d487944ba13bf775b50a80805fe6fcaba9b0fba88f5"
	actual, err := NameHash("resolver.eth")
	assert.Nil(t, err, "Failed to hash")
	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
}

func TestNameHashFooEth(t *testing.T) {
	expected := "de9b09fd7c5f901e23a3f19fecc54828e9c848539801e86591bd9801b019f84f"
	actual, err := NameHash("foo.eth")
	assert.Nil(t, err, "Failed to hash")
	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
}

func TestNameHashFooDotEth(t *testing.T) {
	expected := "4143a5b2f547838d3b49982e3f2ec6a26415274e5b9c3ffeb21971bbfdfaa052"
	actual, err := NameHash("foo..eth")
	assert.Nil(t, err, "Failed to hash")
	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
}

func TestNameHashNickJohnsonEth(t *testing.T) {
	expected := "25cfe90ad9477590acf268bb3ad00ab18465ecce12760be3d8eac81c9f329995"
	actual, err := NameHash("nickjohnson.eth")
	assert.Nil(t, err, "Failed to hash")
	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
}

func TestLabelHashEth(t *testing.T) {
	expected := "4f5b812789fc606be1b3b16908db13fc7a9adf7ca72641f84d75b47069d3d7f0"
	actual, err := LabelHash("icotrust")
	assert.Nil(t, err, "Failed to hash")
	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
}

func TestLabelHashFoo(t *testing.T) {
	expected := "41b1a0649752af1b28b3dc29a1556eee781e4a4c3a1f7f53f90fa834de098c4d"
	actual, err := LabelHash("foo")
	assert.Nil(t, err, "Failed to hash")
	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
}

func TestNameHashReverse(t *testing.T) {
	expected := "91d1777781884d03a6757a803996e38de2a42967fb37eeaca72729271025a9e2"
	actual, err := NameHash("addr.reverse")
	assert.Nil(t, err, "Failed to hash")
	assert.Equal(t, expected, hex.EncodeToString(actual[:]), "Did not receive expected result")
}

func TestNameHashNormalize1(t *testing.T) {
	first, err := NameHash("foo.eth")
	assert.Nil(t, err, "Failed to hash")
	second, err := NameHash("FOO.eth")
	assert.Nil(t, err, "Failed to hash")
	assert.Equal(t, hex.EncodeToString(first[:]), hex.EncodeToString(second[:]), "Did not receive expected result")
}

func TestNormalizeFoo(t *testing.T) {
	expected := "foo"
	actual, err := normalize("FOO")
	assert.Nil(t, err, "Failed to normalize")
	assert.Equal(t, expected, actual, "Did not receive expected result")
}

func TestNormalizeCase(t *testing.T) {
	expected := "foo"
	actual, err := normalize("FOO")
	assert.Nil(t, err, "Failed to normalize")
	assert.Equal(t, expected, actual, "Did not receive expected result")
}

func TestNormalize2(t *testing.T) {
	expected := ".eth"
	actual, err := normalize(".eth")
	assert.Nil(t, err, "Failed to normalize")
	assert.Equal(t, expected, actual, "Did not receive expected result")
}
