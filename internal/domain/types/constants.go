package types

type Format string
type Sort string
type SortDirection string
type Type string
type Season string
type MediaStatus string
type SubType string
type StreamingServer string
type Genre string

const (
	// Media Types
	TypeAnime Type = "ANIME"
	TypeManga Type = "MANGA"

	// Sorting Criteria
	SortScore         Sort = "averageRating"
	SortPopularity    Sort = "averagePopularity"
	SortTitle         Sort = "title"
	SortYear          Sort = "year"
	SortTotalEpisodes Sort = "episodes"
	SortTotalChapters Sort = "chapters"
	SortTotalVolumes  Sort = "volumes"

	// Sort Directions
	SortAsc  SortDirection = "ASC"
	SortDesc SortDirection = "DESC"

	// Provider Types
	ProviderTypeAnime       Type = "ANIME"
	ProviderTypeManga       Type = "MANGA"
	ProviderTypeMeta        Type = "META"
	ProviderTypeInformation Type = "INFORMATION"
	ProviderTypeBase        Type = "BASE"

	// Media Formats
	FormatTV      Format = "TV"
	FormatTVShort Format = "TV_SHORT"
	FormatMovie   Format = "MOVIE"
	FormatSpecial Format = "SPECIAL"
	FormatOVA     Format = "OVA"
	FormatONA     Format = "ONA"
	FormatMusic   Format = "MUSIC"
	FormatManga   Format = "MANGA"
	FormatNovel   Format = "NOVEL"
	FormatOneShot Format = "ONE_SHOT"
	FormatUnknown Format = "UNKNOWN"

	// Seasons
	SeasonWinter  Season = "WINTER"
	SeasonSpring  Season = "SPRING"
	SeasonSummer  Season = "SUMMER"
	SeasonFall    Season = "FALL"
	SeasonUnknown Season = "UNKNOWN"

	// Media Statuses
	MediaStatusFinished       = "FINISHED"
	MediaStatusReleasing      = "RELEASING"
	MediaStatusNotYetReleased = "NOT_YET_RELEASED"
	MediaStatusCancelled      = "CANCELLED"
	MediaStatusHiatus         = "HIATUS"

	// Sub Types
	SubTypeDub = "dub"
	SubTypeSub = "sub"

	// Streaming Servers
	StreamingServerAsianLoad    = "asianload"
	StreamingServerGogoCDN      = "gogocdn"
	StreamingServerStreamSB     = "streamsb"
	StreamingServerMixDrop      = "mixdrop"
	StreamingServerUpCloud      = "upcloud"
	StreamingServerVidCloud     = "vidcloud"
	StreamingServerStreamTape   = "streamtape"
	StreamingServerVizCloud     = "vidplay"
	StreamingServerMyCloud      = "mycloud"
	StreamingServerFilemoon     = "filemoon"
	StreamingServerVidStreaming = "vidstreaming"
	StreamingServerAllAnime     = "allanime"
	StreamingServerFPlayer      = "fplayer"
	StreamingServerKwik         = "kwik"
	StreamingServerDuckStream   = "duckstream"
	StreamingServerDuckStreamV2 = "duckstreamv2"
	StreamingServerBirdStream   = "birdstream"
	StreamingServerAnimeFlix    = "animeflix"

	// Genres
	GenreAction          Genre = "Action"
	GenreAdventure       Genre = "Adventure"
	GenreAnimeInfluenced Genre = "Anime Influenced"
	GenreAvantGarde      Genre = "Avant Garde"
	GenreAwardWinning    Genre = "Award Winning"
	GenreBoysLove        Genre = "Boys Love"
	GenreCards           Genre = "Cards"
	GenreComedy          Genre = "Comedy"
	GenreDementia        Genre = "Dementia"
	GenreDemons          Genre = "Demons"
	GenreDoujinshi       Genre = "Doujinshi"
	GenreDrama           Genre = "Drama"
	GenreEcchi           Genre = "Ecchi"
	GenreErotica         Genre = "Erotica"
	GenreFamily          Genre = "Family"
	GenreFantasy         Genre = "Fantasy"
	GenreFood            Genre = "Food"
	GenreFriendship      Genre = "Friendship"
	GenreGame            Genre = "Game"
	GenreGenderBender    Genre = "Gender Bender"
	GenreGirlsLove       Genre = "Girls Love"
	GenreGore            Genre = "Gore"
	GenreGourmet         Genre = "Gourmet"
	GenreHarem           Genre = "Harem"
	GenreHentai          Genre = "Hentai"
	GenreHistorical      Genre = "Historical"
	GenreHorror          Genre = "Horror"
	GenreIsekai          Genre = "Isekai"
	GenreKids            Genre = "Kids"
	GenreMagic           Genre = "Magic"
	GenreMahouShoujo     Genre = "Mahou Shoujo"
	GenreMartialArts     Genre = "Martial Arts"
	GenreMecha           Genre = "Mecha"
	GenreMedical         Genre = "Medical"
	GenreMilitary        Genre = "Military"
	GenreMusic           Genre = "Music"
	GenreMystery         Genre = "Mystery"
	GenreParody          Genre = "Parody"
	GenrePolice          Genre = "Police"
	GenrePolitical       Genre = "Political"
	GenrePsychological   Genre = "Psychological"
	GenreRacing          Genre = "Racing"
	GenreRomance         Genre = "Romance"
	GenreSamurai         Genre = "Samurai"
	GenreSchool          Genre = "School"
	GenreSciFi           Genre = "Sci-Fi"
	GenreShoujoAi        Genre = "Shoujo Ai"
	GenreShounenAi       Genre = "Shounen Ai"
	GenreSliceOfLife     Genre = "Slice of Life"
	GenreSpace           Genre = "Space"
	GenreSports          Genre = "Sports"
	GenreSuperPower      Genre = "Super Power"
	GenreSupernatural    Genre = "Supernatural"
	GenreSuspense        Genre = "Suspense"
	GenreThriller        Genre = "Thriller"
	GenreVampire         Genre = "Vampire"
	GenreWorkplace       Genre = "Workplace"
	GenreYaoi            Genre = "Yaoi"
	GenreYuri            Genre = "Yuri"
	GenreZombies         Genre = "Zombies"
)

// Valid Sorts
var Sorts = []Sort{
	SortScore,
	SortPopularity,
	SortTitle,
	SortYear,
	SortTotalEpisodes,
	SortTotalChapters,
	SortTotalVolumes,
}

// Valid Formats
var Formats = []Format{
	FormatTV,
	FormatTVShort,
	FormatMovie,
	FormatSpecial,
	FormatOVA,
	FormatONA,
	FormatMusic,
	FormatManga,
	FormatNovel,
	FormatOneShot,
	FormatUnknown,
}
