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
// and the second with the simultaneous meetings happening in that slot.
// Steps:
// 1. Read meetings from input
// 2. Calculate adjacency matrix
// 3. Create one slot per meeting
// 4. Apply removal strategy iteratively until the number of slots
//    does not decrease anymore
func main() {
	if len(os.Args) < 3 || len(os.Args) > 4 {
		fmt.Printf("Usage: %s <csv_file_in> <csv_file_out> [dot_out]\n", os.Args[0])
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
	meetings, err := createMeetingsFromCSVData(csvData)
	if err != nil {
		fmt.Printf("Error reading csv meeting data from %s: %v\n", os.Args[1], err)
		os.Exit(1)
	}

	fmt.Println("Running the algorithm...")
	slots, adj, err := findMinimalSlots(meetings)

	if err := writeSolution(os.Args[2], slots); err != nil {
		fmt.Printf("Cannot write solution: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) > 3 {
		if err := writeDot(os.Args[3], adj, slots); err != nil {
			fmt.Printf("Cannot write dot file: %v\n", err)
			os.Exit(1)
		}
	}
}

func writeDot(outPath string, adj adjacencyMat, slots []slot) error {
	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("cannot open %s for writing: %v", outPath, err)
	}
	defer outFile.Close()

	var sb strings.Builder
	sb.WriteString("strict graph {\n")
	sb.WriteString("  graph [overlap=false outputorder=edgesfirst]\n")
	sb.WriteString("  node [fontcolor=black shape=circle fixedsize=true style=filled]\n")
	numMeets := len(adj[0])
	// Same color for meetings in same slot
	for sl_i, sl := range slots {
		for _, mt := range sl.meetings {
			sb.WriteString(fmt.Sprintf("  %d [color=%s]\n",
				mt.index, slotColor[sl_i%len(slotColor)]))
		}
	}
	// Connect adjacent meetings
	for m_i := 0; m_i < numMeets; m_i++ {
		for m_j := m_i + 1; m_j < numMeets; m_j++ {
			if adj[m_i][m_j] != 1 {
				continue
			}
			sb.WriteString(fmt.Sprintf("  %d -- %d\n", m_i, m_j))
		}
	}
	// Finish graph
	sb.WriteString("}\n")

	if _, err := outFile.WriteString(sb.String()); err != nil {
		return err
	}

	return nil
}

var slotColor = [...]string{
	"aqua",
	"aquamarine",
	"beige",
	"brown",
	"cadetblue",
	"coral",
	"cornflowerblue",
	"crimson",
	"darkgreen",
	"darkorchid",
	"darkseagreen4",
	"cyan4",
	"darkorange",
	"darkolivegreen4",
	"deeppink",
	"firebrick1",
	"dimgray",
	"darkviolet",
	"darkseagreen",
	"floralwhite",
	"darkseagreen2",
	"darkred",
}

func writeSolution(outPath string, slots []slot) error {
	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("cannot open %s for writing: %v", outPath, err)
	}
	defer outFile.Close()

	csvWriter := csv.NewWriter(outFile)
	for sl_i, sl := range slots {
		for _, mt := range sl.meetings {
			assistants := ""
			for i := 0; i < len(mt.assistants)-1; i++ {
				assistants += mt.assistants[i]
				assistants += ", "
			}
			assistants += mt.assistants[len(mt.assistants)-1]
			if err := csvWriter.Write([]string{fmt.Sprintf("slot %d", sl_i),
				mt.title,
				assistants}); err != nil {
				return fmt.Errorf("cannot write to %s: %v", outPath, err)
			}
		}
	}
	csvWriter.Flush()
	fmt.Println("Finished")

	return nil
}

func createMeetingsFromCSVData(csvData [][]string) ([]meeting, error) {
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
		meetings[i] = mt
	}

	return meetings, nil
}

func findMinimalSlots(meetings []meeting) ([]slot, adjacencyMat, error) {
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

	// Remove emptied slots
	var consolidatedSlots []slot
	for _, sl := range slots {
		if len(sl.meetings) != 0 {
			consolidatedSlots = append(consolidatedSlots, sl)
		}
	}

	return consolidatedSlots, adj, nil
}

// Implements the removal strategy, one iteration across current slots.
// Returns true if at least one slot has been consolidated.
func removeSlotsIteration(adj adjacencyMat, slots *[]slot) (removed bool) {
	for s_i, sl := range *slots {
		// If no meetings, the slot was already removed
		if len(sl.meetings) == 0 {
			continue
		}
		// For each meeting in the slot, check if it can be moved to another slot
		compatibleSlots := make([]int, len(sl.meetings))
		allCompatible := true
		for m_i, m := range sl.meetings {
			compatibleSlots[m_i] = -1
			for s_j, otherSl := range *slots {
				// Ignore if same slot or if removed slot
				if s_i == s_j || len(otherSl.meetings) == 0 {
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
				compatSl := compatibleSlots[m_i]
				(*slots)[compatSl].meetings =
					append((*slots)[compatSl].meetings, m)
			}
			(*slots)[s_i] = slot{}
			removed = true
		}
	}

	return removed
}

// Will be compatible if all meetings already in the slot are adjacent to m
func slotCompatibleWithMeeting(sl slot, m meeting, adj adjacencyMat) bool {
	for _, slotMeeting := range sl.meetings {
		if adj[slotMeeting.index][m.index] != 1 {
			return false
		}
	}
	return true
}

// calcAdjacencyMatrix finds out the adjacency matrix for a set of
// meetings, that is, matrix values are 1 if row/column meetings are
// compatible, 0 otherwise.
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
	for _, p1 := range m1.assistants {
		for _, p2 := range m2.assistants {
			if p1 == p2 {
				return false
			}
		}
	}
	return true
}
