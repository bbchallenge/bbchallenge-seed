package bbchallenge

import (
	"os"
	"sync"
	"testing"
)

// Here we want to test if file appends of 30 bytes are atomic or not
func TestAppendAtomic(t *testing.T) {
	testFile := InitAppendFile("test-atomic", "")

	var tms [2]TM = [2]TM{
		TM{'A', 'A', 'A', 'A', 'A', 'A',
			'A', 'A', 'A', 'A', 'A', 'A',
			'A', 'A', 'A', 'A', 'A', 'A',
			'A', 'A', 'A', 'A', 'A', 'A',
			'A', 'A', 'A', 'A', 'A', 'A'},
		TM{'B', 'B', 'B', 'B', 'B', 'B',
			'B', 'B', 'B', 'B', 'B', 'B',
			'B', 'B', 'B', 'B', 'B', 'B',
			'B', 'B', 'B', 'B', 'B', 'B',
			'B', 'B', 'B', 'B', 'B', 'B'}}

	var wg sync.WaitGroup
	for i := 0; i < 1000000; i += 1 {
		wg.Add(1)
		go func(i int) {
			testFile.Write([]byte(tms[i%2][:]))
			wg.Done()
		}(i)
	}

	wg.Wait()

	testFile.Close()

	testFile, err := os.OpenFile("test-atomic", os.O_RDONLY, 0644)
	if err != nil {
		t.Fail()
	}

	var buffer [30]byte
	err = nil
	for err == nil {
		_, err = testFile.Read(buffer[:])
		for i := 0; i < 30; i += 1 {
			if buffer[i] != buffer[0] {
				t.Fail()
			}
		}
	}

	os.Remove("test-atomic")
}
