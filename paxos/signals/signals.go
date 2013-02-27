package signals

/// FOR LOCAL USE ONLY ///
type PrepareType struct {
	RoundNum int
	From int
}

type PromiseType struct {
	Current int
	LastRound int
	LastValue string
}

type AcceptType struct {
	RoundNum int
	Value string
}

type LearnType struct {
	RoundNum int
	Value string
}