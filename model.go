// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import ()

// Model struct
type Model struct {
	// Module name, as DB name.
	Module string

	// table instance
	Table *Table
}

// Insert
func (m *Model) Insert() *Query {
	q := NewQuery(Servers[m.Module])
	q.Table = m.Table
	q.InsertInto(m.Table.Name)
	return q
}

// Update
func (m *Model) Update() *Query {
	q := NewQuery(Servers[m.Module])
	q.Table = m.Table
	q.Update(m.Table.Name)
	return q
}

// Delete
func (m *Model) Delete() *Query {
	q := NewQuery(Servers[m.Module])
	q.Table = m.Table
	q.DeleteFrom(m.Table.Name)
	return q
}

// Select
func (m *Model) Select(f ...string) *Query {
	if len(f) == 0 {
		f = m.Table.SelectFields
	}

	q := NewQuery(Servers[m.Module])
	q.Table = m.Table
	q.Select(f...)
	q.From(m.Table.Name)
	return q
}

// New Model
func NewModel(module string, table *Table) Model {
	return Model{Module: module, Table: table}
}
