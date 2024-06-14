package utility

type Pair[P any, Q any] struct {
  Fst P 
  Snd Q
}

func NewPair[P any, Q any](fst P, snd Q) Pair[P, Q] {
  return Pair[P, Q]{fst, snd}
}
