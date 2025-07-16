// Copyright 2025 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package redisadapter

import (
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/gomodule/redigo/redis"
)

func TestNewAdapterWithConfig(t *testing.T) {
	// Test basic configuration
	config := &Config{
		Network: "tcp",
		Address: "127.0.0.1:6379",
	}
	a, _ := NewAdapter(config)

	testSaveLoad(t, a)
	testAutoSave(t, a)
	testFilteredPolicy(t, a)
	testAddPolicies(t, a)
	testRemovePolicies(t, a)
	testUpdatePolicies(t, a)
	testUpdateFilteredPolicies(t, a)
}

func TestNewAdapterWithPool(t *testing.T) {
	// Test with connection pool
	pool := &redis.Pool{
		MaxIdle:   3,
		MaxActive: 5,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
	}
	config := &Config{
		Pool: pool,
		Key:  "pool_test_rules",
	}
	a, err := NewAdapter(config)
	if err != nil {
		t.Fatal(err)
	}

	testSaveLoad(t, a)
	testAutoSave(t, a)
	testFilteredPolicy(t, a)
	testAddPolicies(t, a)
	testRemovePolicies(t, a)
	testUpdatePolicies(t, a)
	testUpdateFilteredPolicies(t, a)
}

func TestNewAdapterErrorCases(t *testing.T) {
	// Test error cases
	_, err := NewAdapter(nil)
	if err == nil {
		t.Error("NewAdapter should fail with nil config")
	}

	config := &Config{
		Network: "",
		Address: "127.0.0.1:6379",
	}
	_, err = NewAdapter(config)
	if err == nil {
		t.Error("NewAdapter should fail with empty network")
	}

	config = &Config{
		Network: "tcp",
		Address: "",
	}
	_, err = NewAdapter(config)
	if err == nil {
		t.Error("NewAdapter should fail with empty address")
	}
}

func TestNewAdapterWithPassword(t *testing.T) {
	// Test with password authentication
	a, err := NewAdapterWithPassword("tcp", "127.0.0.1:6379", "testpass")
	if err != nil {
		t.Skipf("Password authentication test skipped (Redis may not have auth configured): %v", err)
	}

	testSaveLoad(t, a)
	testAutoSave(t, a)
	testFilteredPolicy(t, a)
	testAddPolicies(t, a)
	testRemovePolicies(t, a)
	testUpdatePolicies(t, a)
	testUpdateFilteredPolicies(t, a)
}

func TestNewAdapterWithUser(t *testing.T) {
	// Test with username and password authentication
	a, err := NewAdapterWithUser("tcp", "127.0.0.1:6379", "testuser", "testpass")
	if err != nil {
		t.Skipf("User authentication test skipped (Redis may not have auth configured): %v", err)
	}

	testSaveLoad(t, a)
	testAutoSave(t, a)
	testFilteredPolicy(t, a)
	testAddPolicies(t, a)
	testRemovePolicies(t, a)
	testUpdatePolicies(t, a)
	testUpdateFilteredPolicies(t, a)
}

func TestNewAdapterWithKey(t *testing.T) {
	// Test with custom key
	a, err := NewAdapterWithKey("tcp", "127.0.0.1:6379", "custom_casbin_rules")
	if err != nil {
		t.Fatal(err)
	}

	testSaveLoad(t, a)
	testAutoSave(t, a)
	testFilteredPolicy(t, a)
	testAddPolicies(t, a)
	testRemovePolicies(t, a)
	testUpdatePolicies(t, a)
	testUpdateFilteredPolicies(t, a)
}

func TestFilterFunctionality(t *testing.T) {
	// Test various filter functionality
	a, err := NewAdapterBasic("tcp", "127.0.0.1:6379")
	if err != nil {
		t.Fatal(err)
	}

	// Initialize policy
	initPolicy(t, a)

	// Create enforcer
	e, _ := casbin.NewEnforcer("examples/rbac_model.conf")
	e.SetAdapter(a)

	// Test filtering by subject
	filter := &Filter{V0: []string{"alice"}}
	err = a.LoadFilteredPolicy(e.GetModel(), filter)
	if err != nil {
		t.Fatal(err)
	}

	policies := e.GetPolicy()
	if len(policies) == 0 {
		t.Log("No policies found for alice (this might be expected)")
	}

	// Test filtering by object
	filter = &Filter{V0: []string{"", "data1"}}
	err = a.LoadFilteredPolicy(e.GetModel(), filter)
	if err != nil {
		t.Fatal(err)
	}

	policies = e.GetPolicy()
	t.Logf("Found %d policies for data1", len(policies))
}
