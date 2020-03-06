// Copyright 2020 Matt Ho
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package automerge_benchmark

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/savaki/automerge"
	"github.com/savaki/automerge/encoding"
	"github.com/tj/assert"
)

type Edit struct {
	Pos    int64 // Pos to insert or update at
	OpType int64 // 0 insert; 1 delete
	Value  string
}

func (e *Edit) UnmarshalJSON(data []byte) error {
	var record []interface{}
	if err := json.Unmarshal(data, &record); err != nil {
		return err
	}

	length := len(record)
	if length != 2 && length != 3 {
		return fmt.Errorf("expected 2 or 3 columns; got %v", length)
	}

	counter, ok := record[0].(float64)
	if !ok {
		return fmt.Errorf("expected counter; got %v (%T)", record[0], record[0])
	}

	actor, ok := record[1].(float64)
	if !ok {
		return fmt.Errorf("expected counter; got %v (%T)", record[1], record[1])
	}

	var value string
	if length >= 3 {
		v, ok := record[2].(string)
		if !ok {
			return fmt.Errorf("expected counter; got %v (%T)", record[2], record[2])
		}
		value = v
	}

	*e = Edit{
		Pos:    int64(counter),
		OpType: int64(actor),
		Value:  value,
	}

	return nil
}

func TestPerformance(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/sample.json")
	assert.Nil(t, err)

	var edits []Edit
	err = json.Unmarshal(data, &edits)
	assert.Nil(t, err)
	assert.Len(t, edits, 259778)

	fmt.Println("applying edits ...")

	node := automerge.NewObject(encoding.RawTypeVarInt)

	actor := []byte("abc")
	begin := time.Now()
	const tick = 25e3
	counter := int64(1)
	for i, edit := range edits {
		switch edit.OpType {
		case 0: // insert
			for _, r := range edit.Value {
				ref := actor
				if counter == 1 {
					ref = nil
				}
				op := automerge.Op{
					OpCounter:  counter,
					OpActor:    actor,
					RefCounter: counter - 1,
					RefActor:   ref,
					Type:       edit.OpType,
					Value:      encoding.RuneValue(r),
				}
				err := node.Insert(op)
				if err != nil {
					t.Fatalf("got %v; want nil", err)
				}
				counter++
			}

		case 1: // delete
			op := automerge.Op{
				OpCounter:  counter,
				OpActor:    actor,
				RefCounter: counter - 1,
				RefActor:   actor,
				Type:       edit.OpType,
				Value:      encoding.RuneValue('_'),
			}
			err := node.Insert(op)
			assert.Nil(t, err)

		default:
			t.Fatalf("got unknown op type, %v", edit.OpType)
		}

		if row := i + 1; row%tick == 0 {
			now := time.Now()
			elapsed := float64(now.Sub(begin) / time.Microsecond)
			fmt.Printf("%6d: %6d bytes, %3.1f µs/op\n", row, node.Size(), elapsed/tick)
			begin = now
		}
	}
	fmt.Println()
	fmt.Println("edits ->", len(edits))
	fmt.Println("bytes ->", node.Size())
}
