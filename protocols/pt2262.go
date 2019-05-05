package protocols

type PT2262 struct {
	PulseLength    int
	SyncFactor     HighLow
	Zero           HighLow
	One            HighLow
	Gap            int // max time between signals (microseconds)
	InvertedSignal bool
}

func GetPT2262Protocol() PT2262 {
	return PT2262{
		350,
		HighLow{1, 31},
		HighLow{1, 3},
		HighLow{3, 1},
		4300,
		false,
	}
}
