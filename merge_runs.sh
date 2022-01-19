#!/usr/bin/env bash
make write_header
OUTPUT=all_5_states_undecided_machines_with_global_header
./write_header >> $OUTPUT
for var in "$@"
do
	cat output/${var}_undecided_time >> $OUTPUT
done
for var in "$@"
do
	cat output/${var}_undecided_space >> $OUTPUT
done
