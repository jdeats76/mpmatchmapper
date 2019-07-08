// Author(s): (1) Immersha Entertainment, (2)Jeremy Deats
// Copyright 2019
//
// License:
// By using this source code and any and all derivied products that may come from this
// soruce code, you agree to the following terms-
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.



package main

import(
	"net/http"
	"fmt"
	"time"
	"math/rand"
	"strconv"
	"strings"
)


// globals
var m (map[string]string)
var m_timestamp (map[string]time.Time)
var m_timestamp_threshold float64
var m_count int
var m_count_threshold int
var key string


// main
func main() {
	// build hashes, set global values
	m = make(map[string]string)
	m_timestamp = make(map[string]time.Time)
	m_timestamp_threshold = 240  // in minutes
	m_count = 1
	m_count_threshold = 500
	key = "testkey"

	// set up handlers
	http.HandleFunc("/writemap", writeMap)
	http.HandleFunc("/readmap", readMap)
	http.HandleFunc("/removemap", removeMap)

	// start server
	http.ListenAndServe(":8087", nil)
}


// generates a random four character label
func getRandomLabel() string {
	charList := "ABCDEFGHIJKLMNOPQRSTUVWZYZ0123456789"
	charListSlices := strings.Split(charList, "")

	// char1
	rnd_offset := rand.Intn(34)
	char1 := string(charListSlices[rnd_offset])

	// char2
	rnd_offset = rand.Intn(34)
	char2 := string(charListSlices[rnd_offset])

	// char3
	rnd_offset = rand.Intn(34)
	char3 := string(charListSlices[rnd_offset])

	result := char1 + char2 + char3
	return result
}


// server provides auto-cleanup, but if our clients want to manually remove maps
// they can
func removeMap(w http.ResponseWriter, r *http.Request) {
	match_key := r.URL.Query().Get("key");
	map_id :=  r.URL.Query().Get("mapid");
	var json_send string

	if (match_key == key) {
		// authorized
		_, isInList := m[map_id]
		if (isInList == false) {
			// if key not found return error
			json_send = "{\"message\":\"mapid not found\"}"
			fmt.Fprint(w, json_send)
			return
		} else {
			// remove from hash
			 delete(m, map_id)
			// return message
			json_send = "{\"message\":\"success\"}"
			fmt.Fprint(w, json_send)
			return
		}
	} else {
		// unauthorized
		json_send = "{\"message\":\"unauthorized\"}"
		fmt.Fprint(w, json_send)
	}
}



//readMap
func readMap(w http.ResponseWriter, r *http.Request) {
	match_key := r.URL.Query().Get("key");
	map_id :=  r.URL.Query().Get("mapid");
	var json_send string

	doCleanup()

	if (match_key == key) {
		// authorized
		_, isInList := m[map_id]
		if (isInList == false) {
			// key already exist in hash, return error
			json_send = "{\"message\":\"mapid not found\"}"
			fmt.Fprint(w, json_send)
			return
		} else {
			// get next count in iteratioe
			result_label := m[map_id]
			// return to client new label
			json_send = "{\"matchid\":\"" + result_label + "\"}"
			fmt.Fprint(w, json_send)
			return
		}
	} else {
		// unauthorized
		json_send = "{\"message\":\"unauthorized\"}"
		fmt.Fprint(w, json_send)
	}
}


// writeMap
func writeMap(w http.ResponseWriter, r *http.Request) {
	match_key := r.URL.Query().Get("key");
	match_id :=  string(r.URL.Query().Get("matchid"));
	var json_send string

	doCleanup()

	if (match_key == key) {
		// authorized
		if (isMatchIdInHash(match_id) == true) {
			// key already exist in hash, return error
			json_send = "{\"message\":\"duplicate matchid\"}"
			fmt.Fprint(w, json_send)
			return
		} else {
			// get next count in iteratioe
			m_count++
			// copy value so we aren't referencing global counter
			ref_count := m_count

			// cleanup. If we're past our threshold restart. This is to flush old session maps
			// adjust treshold as needed.
			if (m_count > m_count_threshold) {
				m_count = 1
			}
			// convert to string
			new_label := strconv.Itoa(ref_count)
			// tag on our random label
			new_label = new_label + getRandomLabel()
			// store new label in hash
			m[new_label] = match_id
			m_timestamp[new_label] = time.Now().UTC()
			// return to client new label
			json_send = "{\"mapid\":\"" + new_label + "\"}"
			fmt.Fprint(w, json_send)
			return
		}
	} else {
		// unauthorized
		json_send = "{\"message\":\"unauthorized\"}"
		fmt.Fprint(w, json_send)
	}
}

// private methods
// detect is value is found in hash
func isMatchIdInHash(matchid string) bool {
	result := false
	for _, v := range m {
		if (v == matchid) {
			result = true
			break
		}
	}
	return result
}


// cleanup map entries that have gone beyond threshold
func doCleanup() {

	currTime := time.Now().UTC();
	removeStamps := make([]string, 1)

	// remove from our primary hash object (m)
	for k, v := range m_timestamp {
		dur := currTime.Sub(v)
		// is the entry expired?
		if (dur.Minutes() > m_timestamp_threshold) {
			// remove expired entry
			delete(m, k)
			// add to slice hash keys for what we're removing from m
			// so we can also remove it from m_timestamp
			removeStamps = append(removeStamps, k)
		}
	}

	// remove m_timestamp entries
	for i, _ := range removeStamps  {
		delete(m_timestamp, removeStamps[i])
	}
}


