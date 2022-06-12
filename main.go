package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type adjacencyMat [][]uint8

type meeting struct {
	index      int
	title      string
	assistants []string
	// Slot currently assigned
	slot int
}

// slot contains all meetings happening at the same time
type slot struct {
	meetings []meeting
}

// Input must be in the form of a csv file, with the first column being
// the meeting title and the second the list of assistants.
// Output will be a csv file with two columns, one with the slot number
// and the second with the simutaneous meetings happening in that slot.
// Steps:
// 1. Read meetings from input
// 2. Calculate adjacency matrix
// 3. Create one slot per meeting
// 4. Apply removal strategy iteratively until the number of slots
//    does not decrease anymore
func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <csv_file>\n", os.Args[0])
		os.Exit(1)
	}

	csvFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Cannot open %s: %v\n", os.Args[1], err)
		os.Exit(1)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvData, err := csvReader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading csv data from %s: %v\n", os.Args[1], err)
		os.Exit(1)
	}
	meetings, err := createMeetingsFromCVSData(csvData)
	if err != nil {
		fmt.Printf("Error reading csv meeting data from %s: %v\n", os.Args[1], err)
		os.Exit(1)
	}
	slots, err := findMinimumSlots(meetings)
	for sl_i, sl := range slots {
		for _, mt := range sl.meetings {
			fmt.Printf("Slot %d, %q, %q\n", sl_i, mt.title, mt.assistants)
		}
	}
}

func createMeetingsFromCVSData(csvData [][]string) ([]meeting, error) {
	meetings := make([]meeting, len(csvData))
	for i, line := range csvData {
		if len(line) != 2 {
			return nil, fmt.Errorf("wrong number of fields in line %d", i)
		}
		mt := meeting{index: i, title: line[0]}
		assistants := strings.Split(line[1], ",")
		for _, person := range assistants {
			mt.assistants = append(mt.assistants, strings.TrimSpace(person))
		}
		meetings = append(meetings, mt)
	}

	return meetings, nil
}

func findMinimumSlots(meetings []meeting) ([]slot, error) {
	adj := calcAdjacencyMatrix(meetings)

	// Initially, one slot per meeting
	slots := make([]slot, len(meetings))
	for i := range slots {
		slots[i].meetings = append(slots[i].meetings, meetings[i])
		meetings[i].index = i
		meetings[i].slot = i
	}

	// Iterated until there is no slot consolidation
	for removed := true; removed == true; removed = removeSlotsIteration(adj, &slots) {
	}

	return slots, nil
}

// Implements the removal strategy, one iteration across current slots.
// Returns true if at least one slot has been consolidated.
func removeSlotsIteration(adj adjacencyMat, slots *[]slot) (removed bool) {
	for s_i, sl := range *slots {
		if len(sl.meetings) == 0 {
			continue
		}
		// For each meeting in the slot, check if it can be moved to another slot
		compatibleSlots := make(map[int]int, len(sl.meetings))
		allCompatible := true
		for m_i, m := range sl.meetings {
			compatibleSlots[m_i] = -1
			for s_j, otherSl := range *slots {
				if s_i == s_j {
					continue
				}
				if slotCompatibleWithMeeting(otherSl, m, adj) {
					compatibleSlots[m_i] = s_j
					break
				}
			}
			if compatibleSlots[m_i] == -1 {
				allCompatible = false
				break
			}
		}
		// If we have found compatible slots for all meetings, transfer
		// the meetings and remove the slot.
		if allCompatible {
			for m_i, m := range sl.meetings {
				newSlot := compatibleSlots[m_i]
				(*slots)[newSlot].meetings =
					append((*slots)[newSlot].meetings, m)
			}
			(*slots)[s_i] = slot{}
			removed = true
		}
	}

	return removed
}

// Will be compatible if all meetings already in the slot have distance one to m
func slotCompatibleWithMeeting(sl slot, m meeting, adj adjacencyMat) bool {
	for _, slotMeeting := range sl.meetings {
		if adj[slotMeeting.index][m.index] != 1 {
			return false
		}
	}
	return true
}

func calcAdjacencyMatrix(meetings []meeting) adjacencyMat {
	numMeets := len(meetings)
	adj := make([][]uint8, numMeets)
	for i := range adj {
		adj[i] = make([]uint8, numMeets)
	}

	for m_i := range meetings {
		for m_j := m_i + 1; m_j < numMeets; m_j++ {
			if compatibleMeetings(&meetings[m_i], &meetings[m_j]) {
				adj[m_i][m_j] = 1
				adj[m_j][m_i] = 1
			}
		}
	}
	return adj
}

// compatibleMeetings returns true if two meetings are compatible,
// that is, will return true if there are no coincident assistants
// between the two meetings, otherwise it will return false.
func compatibleMeetings(m1, m2 *meeting) bool {
	for p1 := range m1.assistants {
		for p2 := range m2.assistants {
			if p1 == p2 {
				return false
			}
		}
	}
	return true
}
