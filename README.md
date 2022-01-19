# BB Challenge Seed

This program enumerates all 5-state 2-symbol Turing machines up to isomorphism, according to the technique developed in [http://turbotm.de/~heiner/BB/mabu90.html](http://turbotm.de/~heiner/BB/mabu90.html), in order to seed the database of undecided 5-state 2-symbol Turing machines. Undecided Turing machines are such that it is not known whether they halt or not.

The program was run in December 2021 and it detected `88,664,064` undecided 5-state 2-symbol machines out of `125,479,953` enumerated machines. A machine was declared to be undecided when it ran for more than `47,176,870` steps (current estimate for BB(5)) or if it visited more than `12,289` memory cells (current estimate for BB_SPACE(5)).

All these undecided machines are available here: [https://dna.hamilton.ie/tsterin/all_5_states_undecided_machines_with_global_header.zip](https://dna.hamilton.ie/tsterin/all_5_states_undecided_machines_with_global_header.zip). 

Once un-zipped you are left with a 2,28 Go binary file with the following structure:

- The first 30 bytes are a header which is currently mainly empty apart from beginning with the three following 4-byte int:
  1. `74,342,035`: The number of machines that are undecided because they exceeded `12,289` memory cells
  2. `14,322,029`: The number of machines that are undecided because they exceeded `47,176,870` steps
  3. `88,664,064`: The total number of machines, which is the sum of the two above numbers

- Then, each one of the `88,664,064` undecided machines is successively encoded in the file using 30 bytes each. Machines that exceeded the space limit of `12,289` cells come first and then come the machines that exceeded the time limit of `47,176,870` steps.
- The 30-byte encoding for a 5-state 2-symbol Turing machine can be understood looking at the following example which is the current BB(5) winner:

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

With `R = 0` and `L = 1`. Note that states are indexed starting at `A=1` as the state value `0` is used to encoded undefined transitions.

### Usage

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
