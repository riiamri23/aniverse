# Aniverse
An uneducated attempt at rewriting [Consumet](https://github.com/consumet/api.consumet.org?tab=readme-ov-file#installation) / [Anify](https://github.com/Eltik/Anify) in Go.

### TODO 
- Quite literally everything.

## Providers
- [AniList](https://anilist.co)
    

#### 1. Get Anime Information
Fetch information about an anime by its ID or title from AniList.
##### Query Parameters
- `provider` (optional): The provider to fetch the anime information from. Default is `anilist`.

##### Request Examples

##### By ID
```http
GET /info/166531?provider=anilist
```


##### By Title
```http
GET /info/oshi-no-ko-2nd-season??provider=anilist
```

##### Response
###### Success (200 OK)
```json
{
  "id": 166531,
  "title": {
    "romaji": "[Oshi no Ko] 2nd Season",
    "english": "Oshi no Ko Season 2",
    "native": "【推しの子】第2期",
    "synonyms": null
  },
  "coverImage": {
    "large": "https://s4.anilist.co/file/anilistcdn/media/anime/cover/medium/bx166531-dAL5MsqDHUkj.jpg",
    "extraLarge": "https://s4.anilist.co/file/anilistcdn/media/anime/cover/large/bx166531-dAL5MsqDHUkj.jpg",
    "color": "#e44343"
  },
  "description": "The second season of [Oshi no Ko]. Aqua’s desire for revenge takes center stage as he navigates the dark underbelly of the entertainment world alongside his twin sister, Ruby. While Ruby follows in their slain mother’s footsteps to become an idol, Aqua joins a famous theater troupe in hopes of uncovering clues to the identity of his father — the man who arranged their mother’s untimely death, and the man who once starred in the same troupe Aqua hopes to infiltrate. (Source: HIDIVE)",
  "status": "RELEASING",
  "episodes": 13,
  "duration": 24,
  "season": "SUMMER",
  "seasonYear": 2024,
  "genres": [
    "Drama",
    "Mystery",
    "Psychological",
    "Supernatural"
  ],
  "synonyms": [
    "我推的孩子"
  ],
  "averageScore": 81,
  "meanScore": 81,
  "popularity": 79741,
  "trailer": {
    "id": "QMuajQlx64c",
    "site": "youtube",
    "thumbnail": "https://i.ytimg.com/vi/QMuajQlx64c/hqdefault.jpg"
  },
  "bannerImage": "https://s4.anilist.co/file/anilistcdn/media/anime/banner/166531-vUu7iDwUkC67.jpg"
}   
```

###### Error (400 Bad Request)
```json
{
  "error": "Invalid ID"
}
```

###### Error (401 Unauthorized)
```json
{
  "error": "Access token not set. Please authenticate first."
}
```

###### Error (500 Internal Server Error)
```json
{
  "error": "Internal server error message"
}

```