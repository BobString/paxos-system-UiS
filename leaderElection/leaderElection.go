package leaderElection

import ()

var (
	pSuspect = map[int]bool{}
	int      leader
)

func suspect(proc int) {
	pSuspect[proc] = true
}

func restore(proc int) {
	pSuspect[proc] = false
}

func getLeader() {
	for k, _ := range pSuspect {
		if !pSuspect[proc] {
			return proc
		}
	}
}
