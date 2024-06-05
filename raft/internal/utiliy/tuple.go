package utiliy

type Pair[P any, Q any] struct{
    Fst P
    Snd Q
}

type Triple[P any, Q any, R any] struct{
    Pair[P,Q]
    Trd R
}
