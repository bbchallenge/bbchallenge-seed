// Here we define filters that prune redundant TMs
package bbchallenge

func pruneTM(nbStates byte, tm TM, state byte, read byte) bool {
	// Returns true if the machine should be ditched
	return pruneEquivalentStates(nbStates, tm, state) ||
		pruneRedundantTransition(nbStates, tm, state, read)
}

func areStatesEquivalent(tm TM, state1 byte, state2 byte) bool {
	// Returns if states 1 and 2 are equivalent, assuming they are fully defined
	write1_0 := tm[6*state1+0]
	move1_0 := tm[6*state1+1]
	goto1_0 := tm[6*state1+2]

	write1_1 := tm[6*state1+3]
	move1_1 := tm[6*state1+4]
	goto1_1 := tm[6*state1+5]

	write2_0 := tm[6*state2+0]
	move2_0 := tm[6*state2+1]
	goto2_0 := tm[6*state2+2]

	write2_1 := tm[6*state2+3]
	move2_1 := tm[6*state2+4]
	goto2_1 := tm[6*state2+5]

	minState := byte(MinI(int(state1+1), int(state2+1)))

	if goto1_0 == state1+1 || goto1_0 == state2+1 {
		goto1_0 = minState
	}

	if goto1_1 == state1+1 || goto1_1 == state2+1 {
		goto1_1 = minState
	}

	if goto2_0 == state1+1 || goto2_0 == state2+1 {
		goto2_0 = minState
	}

	if goto2_1 == state1+1 || goto2_1 == state2+1 {
		goto2_1 = minState
	}

	return write1_0 == write2_0 && write1_1 == write2_1 &&
		move1_0 == move2_0 && move1_1 == move2_1 &&
		goto1_0 == goto2_0 && goto1_1 == goto2_1
}

func pruneEquivalentStates(nbStates byte, tm TM, state byte) bool {
	// Quote from http://turbotm.de/~heiner/BB/mabu90.html#Enumeration:
	// "If there are two states which are (syntactically) equivalent,
	//  these two can be identified (Sigma(N+1) > Sigma(N)).
	//  Example: (x,0)->(x,0,R), (x,1)->(z,0,L), (y,0)->(y,0,R)
	//           and (y,1)->(z,0,L) imply that states x and y are equivalent."
	// Returns true if the machine should be ditched (i.e. no eq states)

	i := state - 1
	if tm[6*i+2] == 0 || tm[6*i+5] == 0 { // state not fully defined
		return false
	}

	for j := byte(0); j < nbStates; j += 1 {
		if i == j {
			continue
		}

		if tm[6*j+2] == 0 || tm[6*j+5] == 0 { // state not fully defined
			continue
		}

		if areStatesEquivalent(tm, i, j) {
			return true
		}
	}

	return false
}

func pruneRedundantTransition(nbStates byte, tm TM, state byte, read byte) bool {
	// Quote from http://turbotm.de/~heiner/BB/mabu90.html#Enumeration:
	// "If a sequence of three transitions is guaranteed to have the same effect
	// as a single transition, only one of both constructions need be inspected.
	// Example: Let x,y,z and s be states with x!=y, a, b, and c arbitrary symbols,
	//          and D from {L,R}, then (x,0)->(y,0,L), (x,1)->(y,1,L), and (y,b)->(z,c,D)
	//          implies that (s,a)->(x,b,R) and (s,a)->(z,c,D) have the same effect."

	move := tm[6*(state-1)+3*read+1]
	goto_ := tm[6*(state-1)+3*read+2]

	if tm[6*(goto_-1)+3*0+2] == 0 || tm[6*(goto_-1)+3*1+2] == 0 { // goto state is not fully defined
		return false
	}

	if tm[6*(goto_-1)+3*0] != 0 || tm[6*(goto_-1)+3*1] != 1 { // goto state is not copying
		return false
	}

	if tm[6*(goto_-1)+3*0+1] != (1-move) || tm[6*(goto_-1)+3*1+1] != (1-move) { // goto state is not coming back
		return false
	}

	return true

}
