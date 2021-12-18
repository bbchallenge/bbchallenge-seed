#ifndef DEF_SIMULATE_H
#define DEF_SIMULATE_H

unsigned char simulate(unsigned char* tm,
                       int limit_time,
                       int limit_space,
                       unsigned char* ret_state,  // output
                       unsigned char* ret_read,   // output
                       int* ret_steps_count,      // output
                       int* ret_space_count);     // output

#endif