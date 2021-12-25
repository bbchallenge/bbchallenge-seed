# BB Challenge

```
Usage of ./bbchallenge:
  -b int
    	simulation backend (0 for go, 1 for C)
  -divtask int
    	divides the size of the job by 1, 2, 4 or 8 (default 1)
  -mytask int
    	select which task bucket this run will do
  -n int
    	# of states (default 4)
  -nf
    	disable extra pruning of redundant machines from the enumeration
  -slim int
    	space limit after which machines are killed and marked as 'UNDECIDED_SPACE' (known values of Busy Beaver space are also used for early termination) (default 12289)
  -tlim int
    	time limit after which running machines are killed and marked as 'UNDECIDED_TIME' (known values of Busy Beaver are also used for early termination) (default 47176870)
  -v	displays infos about the current run on stdout
  -vf int
    	seconds between each stdout log in verbose mode (default 30)
```
