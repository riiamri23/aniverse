package types

type SkipTime struct {
	ID       string        `json:"id"`
	Episodes []EpisodeTime `json:"episodes"`
}

type EpisodeTime struct {
	Intro  TimeRange `json:"intro"`
	Outro  TimeRange `json:"outro"`
	Number int       `json:"number"`
}

type TimeRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}
