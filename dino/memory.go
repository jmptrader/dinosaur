package dino

import (
	"errors"
	"fmt"
)

type Memory []*Process

type MemoryLayout []*MemoryBlock
type MemoryBlock struct {
	Start int
	Size  int
	Name  string
}

const FREE_BLOCK = string('▓')

func (m Memory) HasSpace(size int) bool {
	_, _, err := m.WorstFit(size)
	return err == nil
}

func (m Memory) WorstFit(sizeToFit int) (start, offset int, err error) {
	bestStart := -1
	bestSize := 0

	currentStart := -1
	currentSize := 0
	previousWasEmpty := false

	for i := range m {
		if m[i] == nil {
			if previousWasEmpty == false { // an empty string started
				currentStart = i
				currentSize = 1
				previousWasEmpty = true
			} else { // an empty string is continued
				currentSize++
			}
		} else if previousWasEmpty == true { // an empty string just ended
			previousWasEmpty = false
			if currentSize > bestSize {
				bestStart = currentStart
				bestSize = currentSize
			}
		}
	}
	if currentSize > bestSize {
		bestStart = currentStart
		bestSize = currentSize
	}

	if sizeToFit > bestSize {
		err = errors.New("There's not enough contiguous free space")
	}

	return bestStart, bestSize, err
}

func (m Memory) isEmpty(start, offset int) bool {
	if err := m.checkBounds(start, offset); err != nil {
		return false
	}

	for i := start; i < start+offset; i++ {
		if m[i] != nil {
			return false
		}
	}
	return true

}

func (m Memory) checkBounds(start, offset int) error {
	if start < 0 {
		return errors.New("Cannot allocate -- start index should be non-negative")
	} else if start+offset > len(m) {
		return errors.New("Cannot allocate -- out of memory bound")
	}
	return nil
}

func (m Memory) Allocate(p *Process, start int) (err error) {
	if p == nil {
		return errors.New("Cannot allocate -- nil process")
	} else if p.IsAllocated {
		return errors.New("Cannot allocate -- process already in memory")
	} else if err = m.checkBounds(start, p.SizeInKB); err != nil {
		return err
	} else if !m.isEmpty(start, p.SizeInKB) {
		return errors.New("Cannot allocate -- space already occupied")
	} else if p.ID == "" {
		return errors.New("Cannot allocate -- please assign a (unique) ID to all your processes to unsafe memory operations")
	}

	for i := start; i < start+p.SizeInKB; i++ {
		m[i] = p
	}
	p.IsAllocated = true
	p.MemoryAddress = start
	return nil
}

func (m Memory) AllocateWorstFit(p *Process) (err error) {
	if p == nil {
		return errors.New("Cannot allocate -- nil process")
	}
	start, _, err := m.WorstFit(p.SizeInKB)
	if err != nil {
		return err
	}

	err = m.Allocate(p, start)
	return err
}

func (m Memory) hardRelease(start, offset int) (err error) {
	if err = m.checkBounds(start, offset); err != nil {
		return err
	}

	for i := start; i < start+offset; i++ {
		m[i] = nil
	}
	return nil
}

func (m Memory) ReleaseProcess(p *Process) (bool, error) {
	start := p.MemoryAddress
	offset := p.SizeInKB

	if err := m.checkBounds(start, offset); err != nil {
		return false, err
	}

	beenReleased := false

	if p.ID == "" {
		return false, errors.New("Please assign a (unique) ID to all your processes to unsafe deletions.")
	}

	for i := start; i < start+offset; i++ {
		if m[i].ID == p.ID {
			m[i] = nil
			beenReleased = true
		} else {
			errorStr := fmt.Sprintf("Unsafe delete -- Memory occupied by another process with ID '%s') /nProcess information: ID (%s), MemoryAdress (%d), SizeInKB (%d) \n", m[i].ID, p.ID, p.MemoryAddress, p.SizeInKB)
			if beenReleased == false {
				return false, errors.New(errorStr)
			} else {
				panic(errorStr)
			}
		}
	}

	p.IsAllocated = false
	p.MemoryAddress = -1
	return beenReleased, nil
}

func (m Memory) Layout() MemoryLayout {
	layout := make(MemoryLayout, 0)

	var currentBlock *MemoryBlock

	previousWasEmpty := false

	for i, _ := range m {
		if m[i] == nil && previousWasEmpty == false { // starting empty block
			currentBlock = &MemoryBlock{Start: i, Size: 0, Name: FREE_BLOCK}
			layout = append(layout, currentBlock)
			previousWasEmpty = true
		} else if m[i] != nil && i == 0 || m[i-1] != m[i] { // starting nonempty block
			currentBlock = &MemoryBlock{Start: i, Size: 0, Name: m[i].Name}
			layout = append(layout, currentBlock)
			previousWasEmpty = false
		}
		currentBlock.Size++
	}
	return layout
}

func (m Memory) TotalFree() int {
	total := 0
	for i := range m {
		if m[i] == nil {
			total++
		}
	}
	return total
}

func (ml MemoryLayout) String() string {
	str := "\n\t\t------------ MemoryLayout ------------\n"
	str += fmt.Sprintf("\t\t\t[init, size,  end]\t-\towner\n")
	for i, _ := range ml {
		str += fmt.Sprintf("\t\t\t[%4d, %4d, %4d]\t-\t%2s\n", ml[i].Start, ml[i].Size, ml[i].Start+ml[i].Size-1, ml[i].Name)
	}
	return str
}
