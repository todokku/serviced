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

package logfilter

import (
	"strings"

	"github.com/control-center/serviced/datastore"
	"fmt"
)


// Store is the database for the LogFilters
type Store interface {
	// Get a LogFilter by name and version. Return ErrNoSuchEntity if not found
	Get(ctx datastore.Context, name, version string) (*LogFilter, error)

	// Put adds or updates a LogFilter
	Put(ctx datastore.Context, lf *LogFilter) error

	// Delete removes the a LogFilter if it exists
	Delete(ctx datastore.Context, name, version string) error

}

type storeImpl struct {
	ds datastore.DataStore
}

// NewStore creates a Store for LogFilters
func NewStore() Store {
	return &storeImpl{}
}

// Get a LogFilter by id.  Return ErrNoSuchEntity if not found
func (s *storeImpl) Get(ctx datastore.Context, name, version string) (*LogFilter, error) {
	defer ctx.Metrics().Stop(ctx.Metrics().Start("LogFilterStore.Get"))
	val := &LogFilter{}
	if err := s.ds.Get(ctx, Key(name, version), val); err != nil {
		return nil, err
	}
	return val, nil
}

// Put adds/updates a LogFilter
func (s *storeImpl) Put(ctx datastore.Context, lf *LogFilter) error {
	defer ctx.Metrics().Stop(ctx.Metrics().Start("LogFilterStore.Put"))
	return s.ds.Put(ctx, Key(lf.Name, lf.Version), lf)
}

// Delete removes a LogFilter
func (s *storeImpl) Delete(ctx datastore.Context, name, version string) error {
	defer ctx.Metrics().Stop(ctx.Metrics().Start("LogFilterStore.Delete"))
	return s.ds.Delete(ctx, Key(name, version))
}

//Key creates a Key suitable for getting, putting and deleting LogFilters
func Key(name, version string) datastore.Key {
	name = strings.TrimSpace(name)
	version = strings.TrimSpace(version)
	return datastore.NewKey(kind, buildID(name, version))
}

func buildID(name, version string) string {
	return fmt.Sprintf("%s-%s", name, version)
}
