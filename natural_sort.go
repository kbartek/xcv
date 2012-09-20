package main

import (
	"strconv"
)

func isDigit(a uint8) bool {
	if (a < 48) || (a > 57) {
		return false
	}
	return true
}

func getNumberAndLength(str string, pos int) (uint64, int) {
	var endpos int = pos
	for {
		if endpos == len(str) || isDigit(str[endpos]) == false {
			break
		}
		endpos += 1
	}
	num, err := strconv.ParseUint(str[pos:endpos], 10, 0)
	if err == nil {
		return num, endpos-pos
	}
	return 0, 0
}

func compareNatural(left, right string) bool {
	lpos, rpos := 0, 0
	var lend, rend bool
	var lnum, rnum uint64
	var lnumlen, rnumlen int
	for {
		if isDigit(left[lpos]) && isDigit(right[rpos]) {
			lnum, lnumlen = getNumberAndLength(left, lpos)
			rnum, rnumlen = getNumberAndLength(right, rpos)
			if lnum < rnum {
				return true
			} else if lnum > rnum {
				return false
			} else if lnumlen > 0 && rnumlen > 0 {
				lpos += lnumlen-1
				rpos += rnumlen-1
			}
		} else {
			if left[lpos:lpos+1] < right[rpos:rpos+1] {
				return true
			} else if left[lpos:lpos+1] > right[rpos:rpos+1] {
				return false
			}
		}
		lend, rend = lpos == len(left)-1, rpos == len(right)-1
		if lend == false && rend == false {
			lpos += 1
			rpos += 1
		} else if lend && rend == false {
			return true
		} else {
			break
		}
	}
	return false
}
