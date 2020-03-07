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

package main

import (
	"encoding/json"
	"fmt"
	"github.com/savaki/automerge/encoding"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/savaki/automerge"
	"github.com/urfave/cli"
)

var opts struct {
	file string
}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "file",
			Usage:       "data file",
			Value:       "testdata/sample.json",
			Destination: &opts.file,
		},
	}
	app.Action = apply
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

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

const tick = 25e3

func benchmark(callback func() (*automerge.Object, error)) error {
	var obj *automerge.Object
	var err error
	defer func(begin time.Time) {
		fmt.Println()
		fmt.Println("edits:   ", obj.RowCount())
		fmt.Println("bytes:   ", obj.Size())
		fmt.Println("elapsed: ", time.Now().Sub(begin).Round(time.Millisecond))
		fmt.Println()

	}(time.Now())

	obj, err = callback()
	return err
}

func apply(_ *cli.Context) error {
	data, err := ioutil.ReadFile(opts.file)
	if err != nil {
		return fmt.Errorf("unable to benchmark file, %v: %w", opts.file, err)
	}

	fmt.Println()
	fmt.Printf("applying %v (%v bytes)\n", opts.file, len(data))

	switch ext := filepath.Ext(opts.file); ext {
	case ".json":
		return applyJSON(data)

	default:
		return applyText(data)
	}
}

func applyJSON(data []byte) error {
	var edits []Edit
	if err := json.Unmarshal(data, &edits); err != nil {
		return fmt.Errorf("unable to unmarshal edits: %w", err)
	}

	return benchmark(func() (*automerge.Object, error) {
		var (
			obj     = automerge.NewObject(encoding.RawTypeVarInt)
			actor   = []byte("abc")
			begin   = time.Now()
			counter = int64(1)
		)

		for i, edit := range edits {
			switch edit.OpType {
			case 0: // insert
				for _, r := range edit.Value {
					ref := actor
					if counter == 1 {
						ref = nil
					}
					op := automerge.Op{
						ID:    automerge.NewID(counter, actor),
						Ref:   automerge.NewID(counter-1, ref),
						Type:  edit.OpType,
						Value: encoding.RuneValue(r),
					}
					if err := obj.Insert(op); err != nil {
						return nil, err
					}
					counter++
				}

			case 1: // delete
				op := automerge.Op{
					ID:    automerge.NewID(counter, actor),
					Ref:   automerge.NewID(counter-1, actor),
					Type:  edit.OpType,
					Value: encoding.RuneValue('_'),
				}
				if err := obj.Insert(op); err != nil {
					return nil, err
				}

			default:
				return nil, fmt.Errorf("got unknown op type, %v", edit.OpType)
			}

			if row := i + 1; row%tick == 0 {
				if row == tick {
					fmt.Println()
				}
				now := time.Now()
				elapsed := float64(now.Sub(begin) / time.Microsecond)
				begin = now
				fmt.Printf("%6d: %6d bytes, %3.1f µs/op\n", row, obj.Size(), elapsed/tick)
			}
		}

		return obj, nil
	})
}

func applyText(data []byte) error {
	s := string(data)

	return benchmark(func() (*automerge.Object, error) {
		var (
			obj     = automerge.NewObject(encoding.RawTypeVarInt)
			actor   = []byte("abc")
			begin   = time.Now()
			counter = int64(1)
		)

		for i, r := range s {
			ref := actor
			if counter == 1 {
				ref = nil
			}
			op := automerge.Op{
				ID:    automerge.NewID(counter, actor),
				Ref:   automerge.NewID(counter-1, ref),
				Type:  0,
				Value: encoding.RuneValue(r),
			}
			if err := obj.Insert(op); err != nil {
				return nil, err
			}
			counter++

			if row := i + 1; row%tick == 0 {
				now := time.Now()
				elapsed := float64(now.Sub(begin) / time.Microsecond)
				begin = now
				fmt.Printf("%6d: %6d bytes, %3.1f µs/op\n", row, obj.Size(), elapsed/tick)
			}
		}

		return obj, nil
	})
}
