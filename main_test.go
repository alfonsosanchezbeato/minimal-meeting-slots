// go test *.go
package main

import "testing"

func TestCompatibleMeetings(t *testing.T) {
	compatibleTests := []struct{ m1, m2 *meeting }{
		{
			&meeting{index: 0, title: "Meeting 1", assistants: []string{"a", "b"}},
			&meeting{index: 1, title: "Meeting 2", assistants: []string{"c", "d"}},
		},
		{
			&meeting{index: 0, title: "Meeting 1", assistants: []string{"a"}},
			&meeting{index: 1, title: "Meeting 2", assistants: []string{"b", "c", "d"}},
		},
	}
	for _, mt := range compatibleTests {
		if !compatibleMeetings(mt.m1, mt.m2) {
			t.Error(mt.m1, "incompatible with", mt.m2, "; expected compatible")
		}
	}

	incompatibleTests := []struct{ m1, m2 *meeting }{
		{
			&meeting{index: 0, title: "Meeting 1", assistants: []string{"a", "b"}},
			&meeting{index: 1, title: "Meeting 2", assistants: []string{"c", "a"}},
		},
		{
			&meeting{index: 0, title: "Meeting 1", assistants: []string{"a"}},
			&meeting{index: 1, title: "Meeting 2", assistants: []string{"b", "c", "a"}},
		},
	}
	for _, mt := range incompatibleTests {
		if compatibleMeetings(mt.m1, mt.m2) {
			t.Error(mt.m1, "compatible with", mt.m2, "; expected incompatible")
		}
	}
}
