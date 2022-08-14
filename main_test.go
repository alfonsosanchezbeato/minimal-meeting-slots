// go test *.go
package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestCompatibleMeetings(t *testing.T) {
	compatibleTests := []struct{ m1, m2 *meeting }{
		{
			&meeting{index: 0, title: "Meeting 0", assistants: []string{"a", "b"}},
			&meeting{index: 1, title: "Meeting 1", assistants: []string{"c", "d"}},
		},
		{
			&meeting{index: 0, title: "Meeting 0", assistants: []string{"a"}},
			&meeting{index: 1, title: "Meeting 1", assistants: []string{"b", "c", "d"}},
		},
	}
	for _, mt := range compatibleTests {
		if !compatibleMeetings(mt.m1, mt.m2) {
			t.Error(mt.m1, "incompatible with", mt.m2, "; expected compatible")
		}
	}

	incompatibleTests := []struct{ m1, m2 *meeting }{
		{
			&meeting{index: 0, title: "Meeting 0", assistants: []string{"a", "b"}},
			&meeting{index: 1, title: "Meeting 1", assistants: []string{"c", "a"}},
		},
		{
			&meeting{index: 0, title: "Meeting 0", assistants: []string{"a"}},
			&meeting{index: 1, title: "Meeting 1", assistants: []string{"b", "c", "a"}},
		},
	}
	for _, mt := range incompatibleTests {
		if compatibleMeetings(mt.m1, mt.m2) {
			t.Error(mt.m1, "compatible with", mt.m2, "; expected incompatible")
		}
	}
}

func TestCalcAdjacencyMatrix(t *testing.T) {
	mts := []meeting{
		{index: 0, title: "Meeting 0", assistants: []string{"a", "b"}},
		{index: 1, title: "Meeting 1", assistants: []string{"c", "d"}},
		{index: 2, title: "Meeting 2", assistants: []string{"e", "a"}},
		{index: 3, title: "Meeting 3", assistants: []string{"e", "d"}},
	}
	expMat := [][]uint8{
		{0, 1, 0, 1},
		{1, 0, 1, 0},
		{0, 1, 0, 0},
		{1, 0, 0, 0},
	}
	adjMat := calcAdjacencyMatrix(mts)
	if len(expMat) != len(adjMat) {
		t.Error(expMat, "different size from", adjMat, "; expected equal")
	}
	for i := range expMat[0] {
		if len(expMat[i]) != len(adjMat[i]) {
			t.Error(expMat, "different size from", adjMat, "; expected equal")
		}
		for j := range expMat[0] {
			if expMat[i][j] != adjMat[i][j] {
				t.Error(expMat, "different from", adjMat, "; expected equal")
			}
		}
	}
}

func TestSlotCompatibleWithMeeting(t *testing.T) {
	m0 := meeting{index: 0, title: "Meeting 0", assistants: []string{"a", "b"}}
	m1 := meeting{index: 1, title: "Meeting 1", assistants: []string{"c", "d"}}
	m2 := meeting{index: 2, title: "Meeting 2", assistants: []string{"e", "g"}}
	m3 := meeting{index: 3, title: "Meeting 3", assistants: []string{"e", "h"}}
	mts := []meeting{m0, m1, m2, m3}
	adjMat := calcAdjacencyMatrix(mts)
	var sl slot
	sl = slot{meetings: []meeting{m0, m1}}
	if !slotCompatibleWithMeeting(sl, m2, adjMat) {
		t.Error(sl, "incompatible with", m2, "; expected compatible")
	}
	if !slotCompatibleWithMeeting(sl, m3, adjMat) {
		t.Error(sl, "incompatible with", m3, "; expected compatible")
	}

	sl = slot{meetings: []meeting{m0, m1, m2}}
	if slotCompatibleWithMeeting(sl, m3, adjMat) {
		t.Error(sl, "compatible with", m3, "; expected incompatible")
	}
}

func TestRemoveSlotsIteration1(t *testing.T) {
	m0 := meeting{index: 0, title: "Meeting 0", assistants: []string{"a", "b"}}
	m1 := meeting{index: 1, title: "Meeting 1", assistants: []string{"c", "d"}}
	m2 := meeting{index: 2, title: "Meeting 2", assistants: []string{"e", "g"}}
	m3 := meeting{index: 3, title: "Meeting 3", assistants: []string{"e", "h"}}
	mts := []meeting{m0, m1, m2, m3}
	adjMat := calcAdjacencyMatrix(mts)
	slots := make([]slot, len(mts))
	for m_i, m := range mts {
		slots[m_i] = slot{meetings: []meeting{m}}
	}

	i := 0
	for slotsRemoved := true; slotsRemoved; i++ {
		slotsRemoved = removeSlotsIteration(adjMat, &slots)
	}
	if expIt := 2; i != expIt {
		t.Error(i, "iterations: expected", expIt)
	}
	expSlots := []slot{
		{},
		{},
		{meetings: []meeting{m2, m1, m0}},
		{meetings: []meeting{m3}},
	}
	if !reflect.DeepEqual(slots, expSlots) {
		t.Error(slots, "final slots: expected", expSlots)
	}
}

func TestRemoveSlotsIteration2(t *testing.T) {
	m0 := meeting{index: 0, title: "Meeting 0", assistants: []string{"a", "e"}}
	m1 := meeting{index: 1, title: "Meeting 1", assistants: []string{"b", "f"}}
	m2 := meeting{index: 2, title: "Meeting 2", assistants: []string{"c", "g"}}
	m3 := meeting{index: 3, title: "Meeting 3", assistants: []string{"d", "h"}}
	m4 := meeting{index: 4, title: "Meeting 4", assistants: []string{"b", "c", "d"}}
	m5 := meeting{index: 5, title: "Meeting 5", assistants: []string{"a", "c", "d"}}
	m6 := meeting{index: 6, title: "Meeting 6", assistants: []string{"a", "b", "d"}}
	m7 := meeting{index: 7, title: "Meeting 7", assistants: []string{"a", "b", "c"}}
	mts := []meeting{m0, m1, m2, m3, m4, m5, m6, m7}
	adjMat := calcAdjacencyMatrix(mts)
	slots := make([]slot, len(mts))
	for m_i, m := range mts {
		slots[m_i] = slot{meetings: []meeting{m}}
	}

	i := 0
	for slotsRemoved := true; slotsRemoved; i++ {
		slotsRemoved = removeSlotsIteration(adjMat, &slots)
	}
	if expIt := 2; i != expIt {
		t.Error(i, "iterations: expected", expIt)
	}
	expSlots := []slot{
		{},
		{},
		{},
		{},
		{meetings: []meeting{m4, m0}},
		{meetings: []meeting{m5, m1}},
		{meetings: []meeting{m6, m2}},
		{meetings: []meeting{m7, m3}},
	}
	if !reflect.DeepEqual(slots, expSlots) {
		t.Error(slots, "final slots: expected", expSlots)
	}
}

func TestFindMinimalSlots(t *testing.T) {
	m0 := meeting{index: 0, title: "Meeting 0", assistants: []string{"a", "b"}}
	m1 := meeting{index: 1, title: "Meeting 1", assistants: []string{"c", "d"}}
	m2 := meeting{index: 2, title: "Meeting 2", assistants: []string{"e", "g"}}
	m3 := meeting{index: 3, title: "Meeting 3", assistants: []string{"e", "h"}}
	mts := []meeting{m0, m1, m2, m3}
	slots, _, err := findMinimalSlots(mts)
	if err != nil {
		t.Error("unexpected error while calling findMinimalSlots:", err)
	}
	expSlots := []slot{
		{meetings: []meeting{m2, m1, m0}},
		{meetings: []meeting{m3}},
	}
	if !reflect.DeepEqual(slots, expSlots) {
		t.Error(slots, "final slots: expected", expSlots)
	}
}

func TestCreateMeetingsFromCSVData(t *testing.T) {
	csvData := [][]string{
		{"Meeting 0", "Maria Barker, John Garcia"},
		{"Meeting 1", "Antonio Sanchez, Ana de Austria"},
	}
	meetings, err := createMeetingsFromCSVData(csvData)
	if err != nil {
		t.Error("unexpected error while calling createMeetingsFromCSVData:", err)
	}
	expMeets := []meeting{
		{0, "Meeting 0", []string{"Maria Barker", "John Garcia"}, 0},
		{1, "Meeting 1", []string{"Antonio Sanchez", "Ana de Austria"}, 0},
	}
	if !reflect.DeepEqual(meetings, expMeets) {
		t.Error(meetings, "csv data read: expected", expMeets)
	}
}

func TestMain(t *testing.T) {
	out, err := ioutil.TempFile("", "slots*.csv")
	if err != nil {
		t.Error("cannot create temp file:", err)
	}
	defer os.Remove(out.Name())
	out.Close()

	os.Args = []string{"./main", "test-data/test1.csv", out.Name()}
	main()
	outData, err := os.ReadFile(out.Name())
	if err != nil {
		t.Error("cannot read result from processing of test1:", err)
	}
	expected, err := os.ReadFile("test-data/res1.csv")
	if err != nil {
		t.Error("cannot read res1.csv:", err)
	}
	if bytes.Compare(outData, expected) != 0 {
		t.Errorf("expected:\n%q\n, obtained:\n%q\n", string(expected), string(outData))
	}
}
