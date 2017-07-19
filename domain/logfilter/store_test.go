// Copyright 2017 The Serviced Authors.
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

// +build integration

package logfilter

import (
	"github.com/control-center/serviced/datastore"
	"github.com/control-center/serviced/datastore/elastic"
	. "gopkg.in/check.v1"
)

var _ = Suite(&S{
	ElasticTest: elastic.ElasticTest{
		Index:    "controlplane",
		Mappings: []elastic.Mapping{MAPPING},
	}})

type S struct {
	elastic.ElasticTest
	ctx   datastore.Context
	store Store
}

func (s *S) SetUpTest(c *C) {
	s.ElasticTest.SetUpTest(c)
	datastore.Register(s.Driver())
	s.ctx = datastore.Get()
	s.store = NewStore()
}

func (s *S) Test_LogFilterCRUD(c *C) {
	expected := &LogFilter{
		Name:    "filter name",
		Version: "1.0",
		Filter:  "filter it",
	}
	actual, err := s.store.Get(s.ctx, expected.Name, expected.Version)
	c.Assert(err, NotNil)
	c.Assert(actual, IsNil)
	c.Assert(datastore.IsErrNoSuchEntity(err), Equals, true)


	err = s.store.Put(s.ctx, expected)
	expected.DatabaseVersion++
	c.Assert(err, IsNil)

	actual, err = s.store.Get(s.ctx, expected.Name, expected.Version)
	c.Assert(err, IsNil)
	c.Assert(actual, DeepEquals, expected)

	err = s.store.Put(s.ctx, expected)
	expected.DatabaseVersion++
	c.Assert(err, IsNil)

	actual, err = s.store.Get(s.ctx, expected.Name, expected.Version)
	c.Assert(err, IsNil)
	c.Assert(actual, DeepEquals, expected)

	err = s.store.Delete(s.ctx, expected.Name, expected.Version)
	c.Assert(err, IsNil)

	actual, err = s.store.Get(s.ctx, expected.Name, expected.Version)
	c.Assert(err, NotNil)
	c.Assert(datastore.IsErrNoSuchEntity(err), Equals, true)
}

func (s *S) Test_GetRequiresNameAndVersion(c *C) {
	actual := &LogFilter{
		Name:    "filter name",
		Version: "1.0",
		Filter:  "filter it",
	}

	err := s.store.Put(s.ctx, actual)
	c.Assert(err, IsNil)

	result, err2 := s.store.Get(s.ctx, "", actual.Version)
	c.Assert(err2, NotNil)
	c.Assert(result, IsNil)
	c.Assert(datastore.IsErrNoSuchEntity(err2), Equals, true)

	result, err2 = s.store.Get(s.ctx, actual.Name, "")
	c.Assert(err2, NotNil)
	c.Assert(result, IsNil)
	c.Assert(datastore.IsErrNoSuchEntity(err2), Equals, true)
}
