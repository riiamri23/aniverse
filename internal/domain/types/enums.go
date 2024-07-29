package types

// Define enums and constants
type (
	MediaFormat     int
	SortCriterion   int
	SortDirection   int
	MediaType       int
	Season          int
	MediaStatus     int
	AudioType       int
	StreamingServer int
	Genre           int
	Flag            int
)

const (
	MediaTypeAnime MediaType = iota
	MediaTypeManga
)

const (
	SortCriterionScore SortCriterion = iota
	SortCriterionPopularity
	SortCriterionTitle
	SortCriterionYear
	SortCriterionTotalEpisodes
	SortCriterionTotalChapters
	SortCriterionTotalVolumes
)

const (
	SortDirectionAsc SortDirection = iota
	SortDirectionDesc
)

const (
	MediaFormatTV MediaFormat = iota
	MediaFormatTVShort
	MediaFormatMovie
	MediaFormatSpecial
	MediaFormatOVA
	MediaFormatONA
	MediaFormatMusic
	MediaFormatManga
	MediaFormatNovel
	MediaFormatOneShot
	MediaFormatUnknown
)

const (
	SeasonWinter Season = iota
	SeasonSpring
	SeasonSummer
	SeasonFall
	SeasonUnknown
)

const (
	MediaStatusFinished MediaStatus = iota
	MediaStatusReleasing
	MediaStatusNotYetReleased
	MediaStatusCancelled
	MediaStatusHiatus
)

const (
	AudioTypeDub AudioType = iota
	AudioTypeSub
)

const (
	StreamingServerAsianLoad StreamingServer = iota
	StreamingServerGogoCDN
	StreamingServerStreamSB
	StreamingServerMixDrop
	StreamingServerUpCloud
	StreamingServerVidCloud
	StreamingServerStreamTape
	StreamingServerVizCloud
	StreamingServerMyCloud
	StreamingServerFilemoon
	StreamingServerVidStreaming
	StreamingServerAllAnime
	StreamingServerFPlayer
	StreamingServerKwik
	StreamingServerDuckStream
	StreamingServerDuckStreamV2
	StreamingServerBirdStream
	StreamingServerAnimeFlix
)

const (
	GenreAction Genre = iota
	GenreAdventure
	GenreAnimeInfluenced
	GenreAvantGarde
	GenreAwardWinning
	GenreBoysLove
	GenreCards
	GenreComedy
	GenreDementia
	GenreDemons
	GenreDoujinshi
	GenreDrama
	GenreEcchi
	GenreErotica
	GenreFamily
	GenreFantasy
	GenreFood
	GenreFriendship
	GenreGame
	GenreGenderBender
	GenreGirlsLove
	GenreGore
	GenreGourmet
	GenreHarem
	GenreHentai
	GenreHistorical
	GenreHorror
	GenreIsekai
	GenreKids
	GenreMagic
	GenreMahouShoujo
	GenreMartialArts
	GenreMecha
	GenreMedical
	GenreMilitary
	GenreMusic
	GenreMystery
	GenreParody
	GenrePolice
	GenrePolitical
	GenrePsychological
	GenreRacing
	GenreRomance
	GenreSamurai
	GenreSchool
	GenreSciFi
	GenreShoujoAi
	GenreShounenAi
	GenreSliceOfLife
	GenreSpace
	GenreSports
	GenreSuperPower
	GenreSupernatural
	GenreSuspense
	GenreThriller
	GenreVampire
	GenreWorkplace
	GenreYaoi
	GenreYuri
	GenreZombies
)

const (
	FlagCORSAllowed Flag = iota
)
