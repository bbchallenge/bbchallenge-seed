# BB Challenge Seed

**Update 06/03/22**. To improve the reproducibility of the results we have decided to lexicographically sort the database computed by this program (`all_5_states_undecided_machines_with_global_header`). The first `14,322,029` undecided machines (47M time limit exceeded) were lexicographically sorted independently of the next `74,342,035` undecided machines (12k space limit exceeded). 

This program enumerates all 5-state 2-symbol Turing machines up to isomorphism, according to the technique developed in [http://turbotm.de/~heiner/BB/mabu90.html](http://turbotm.de/~heiner/BB/mabu90.html), in order to seed the database of undecided 5-state 2-symbol Turing machines. Undecided Turing machines are such that it is not known whether they halt or not.

The program was run in December 2021 and it detected `88,664,064` undecided 5-state 2-symbol machines out of `125,479,953` enumerated machines. A machine was declared to be undecided when it ran for more than `47,176,870` steps (current estimate for BB(5)) or if it visited more than `12,289` memory cells (current estimate for BB_SPACE(5)).

## Progam usage

Although the seed database has already been computed, you are invited to reproduce the results. See [https://bbchallenge.org/contribute](https://bbchallenge.org/contribute).

Note that it took 92 computing hours, splitted on 4 computers, to enumerate the *useful* space of all 5-state 2-symbol Turing machines. 

You can find more statistics about this seminal run here: [https://bbchallenge.org/method#metrics](https://bbchallenge.org/method#metrics)

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

## Database

All these undecided machines are available at these mirrors: 

- [https://dna.hamilton.ie/tsterin/all_5_states_undecided_machines_with_global_header.zip](https://dna.hamilton.ie/tsterin/all_5_states_undecided_machines_with_global_header.zip)
- [ipfs://QmcgucgLRjAQAjU41w6HR7GJbcte3F14gv9oXcf8uZ8aFM](ipfs://QmcgucgLRjAQAjU41w6HR7GJbcte3F14gv9oXcf8uZ8aFM)
- [https://ipfs.prgm.dev/ipfs/QmcgucgLRjAQAjU41w6HR7GJbcte3F14gv9oXcf8uZ8aFM](https://ipfs.prgm.dev/ipfs/QmcgucgLRjAQAjU41w6HR7GJbcte3F14gv9oXcf8uZ8aFM)

Database shasum: 
  1. all_5_states_undecided_machines_with_global_header.zip: `2576b647185063db2aa3dc2f5622908e99f3cd40`
  2. all_5_states_undecided_machines_with_global_header: `e57063afefd900fa629cfefb40731fd083d90b5e`

### Database Format

Once un-zipped you are left with a 2,28 Go binary file with the following structure:

- The first 30 bytes are a header which is currently mainly empty apart from beginning with the three following 4-byte int followed by a 1-byte bool:
  1. `14,322,029`: The number of machines that are undecided because they exceeded `47,176,870` steps
  2. `74,342,035`: The number of machines that are undecided because they exceeded `12,289` memory cells
  4. `88,664,064`: The total number of machines, which is the sum of the two above numbers
  5. `1`: the database has been lexicographically sorted. The first `14,322,029` undecided machines (47M time limit exceeded) were lexicographically sorted independently of the next `74,342,035` undecided machines (12k space limit exceeded). 

Then, each one of the `88,664,064` undecided machines is successively encoded in the file using 30 bytes each. Machines that exceeded the time limit of `47,176,870` steps come first and then come the machines that exceeded the space limit of `12,289` cells.

The 30-byte encoding for a 5-state 2-symbol Turing machine can be understood looking at the following example which is the current BB(5) winner:

```
+---+-----+-----+
| - |  0  |  1  |
+---+-----+-----+
| A | 1RB | 1LC |
| B | 1RC | 1RB |
| C | 1RD | 0LE |
| D | 1LA | 1LD |
| E | 1RH | 0LA |
+---+-----+-----+
```

Is encoded by the following successive 30 bytes:

```
1, R, 2, 1, L, 3,
1, R, 3, 1, R, 2,
1, R, 4, 0, L, 5,
1, L, 1, 1, L, 4,
1, R, 6, 0, L, 1
```

With `R = 0` and `L = 1`. Note that states are indexed starting at `A=1` as the state value `0` is used to encode undefined transitions.

### Use the database

Here are some routines that you will let you interact with the database in different languages:

#### Python

```python

def get_header(machine_db_path):
    with open(machine_db_path, "rb") as f:
        return f.read(30)
        
def get_machine_i(machine_db_path, i, db_has_header=True):
    with open(machine_db_path, "rb") as f:
        c = 1 if db_has_header else 0
        f.seek(30*(i+c))
        return f.read(30)
```

More at [https://github.com/bbchallenge/bbchallenge-py/](https://github.com/bbchallenge/bbchallenge-py/).

#### Go

```go
func GetMachineI(db []byte, i int, hasHeader bool) (tm TM, err error) {

	if i < 0 || i > len(db)/30 {
		err := errors.New("invalid db index")
		return tm, err
	}

	offset := 0
	if hasHeader {
		offset = 1
	}

	copy(tm[:], db[30*(i+offset):30*(i+offset+1)])
	return tm, nil
}
```

More at [https://github.com/bbchallenge/bbchallenge-go/](https://github.com/bbchallenge/bbchallenge-go/).
