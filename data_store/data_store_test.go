/*
 * Copyright (c) 2013-2014, Jeremy Bingham (<jbingham@gmail.com>)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package data_store

import (
	"testing"
	"io/ioutil"
	"fmt"
	"os"
)

type dsObj struct {
	Name string `json:"name"`
	JsonClass string `json:"json_class"`
	ChefType string `json:"chef_type"`
}

func makeDsObj() *dsObj {
	return &dsObj{ Name: "baz", JsonClass: "Chef::DsObj", ChefType: "ds_obj" }
}

func TestNew(t *testing.T){
	if d := New(); d == nil {
		t.Errorf("New() should have returned a data store object, but returned nil")
	}
}

func TestSet(t *testing.T){
	ds := New()
	baz := makeDsObj()
	ds.Set("foo", "bar", baz)
}

func TestGet(t *testing.T){
	ds := New()
	val, found := ds.Get("foo", "bar2")
	if found {
		t.Errorf("Get() returned a result improperly")
	}
	baz := makeDsObj()
	ds.Set("foo", "bar2", baz)
	val, found = ds.Get("foo", "bar2")
	if !found {
		t.Errorf("Get() did not return a result properly, got '%v' :: %v", val, found)
	}
}

func TestDelete(t *testing.T){
	ds := New()
	baz := makeDsObj()
	ds.Set("foo", "bar3", baz)
	val, found := ds.Get("foo", "bar3")
	if found == false {
		t.Errorf("Couldn't set bar3 baz")
	}
	ds.Delete("foo", "bar3")
	val, found = ds.Get("foo", "bar3")
	if found {
		t.Errorf("Delete() did not delete bar3, returned %v!", val)
	}
}

func TestGetList(t *testing.T){
	ds := New()
	complist := []string{ "baz", "moo" }
	baz := makeDsObj()
	moo := makeDsObj()
	moo.Name = "moo"
	ds.Set("foolist", "baz", baz)
	ds.Set("foolist", "moo", moo)
	dsl := ds.GetList("foolist")
	if dsl == nil || dsl[0] != complist[0] || dsl[1] != complist[1] {
		t.Errorf("GetList failed to return the expected list: returned %v, expected %v", dsl, complist)
	}
}

var dsTmpDir = dsTmpGen()

func dsTmpGen() string {
	tm, err := ioutil.TempDir("", "ds-test")
	if err != nil {
		panic("Couldn't create temporary directory!")
	}
	return tm
}

func TestSave(t *testing.T) {
	ds := New()
	tmpfile := fmt.Sprintf("%s/ds.bin", dsTmpDir)
	err := ds.Save(tmpfile)
	if err != nil {
		t.Errorf("Save() gave an error: %s", err)
	}
}

func TestLoad(t *testing.T) {
	ds := New()
	tmpfile := fmt.Sprintf("%s/ds.bin", dsTmpDir)
	err := ds.Load(tmpfile)
	if err != nil {
		t.Errorf("Load() save an error: %s", err)
	}
}

func TestSaveAndLoadData(t *testing.T) {
	ds := New()
	tmpfile := fmt.Sprintf("%s/ds2.bin", dsTmpDir)
	baz := makeDsObj()
	boo := makeDsObj()
	boo.Name = "boo"
	ds.Set("foo", "bar", baz)
	ds.Set("foo", "boo", boo)
	ds.Save(tmpfile)
	dsLoad := New()
	dsLoad.Load(tmpfile)
	bS, found := dsLoad.Get("foo", "bar")
	if !found {
		t.Errorf("Did not find bar!! dsLoad is: %v", dsLoad)
	}
	var bazSave *dsObj
	if bS != nil {
		bazSave = bS.(*dsObj)
	}
	if bazSave == nil {
		t.Errorf("Did not successfully retrieve baz from saved data store")
	} else if bazSave.Name != baz.Name {
		t.Errorf("Retrieved the wrong object! Expected %s, got %s", baz.Name, bazSave.Name)
	}
}

// clean up

func TestCleanup(t *testing.T) {
	os.RemoveAll(dsTmpDir)
}
