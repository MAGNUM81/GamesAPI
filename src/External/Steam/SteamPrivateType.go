package Steam

type basicUserSteamType struct {
	Response struct {
		Steamid string `json:"steamid"`
		Success int    `json:"success"`
	} `json:"response"`
}

type gameSteamType struct {
	Appid                    int `json:"appid"`
	Playtime_forever         int `json:"playtime_forever"`
	Playtime_windows_forever int `json:"playtime_windows_forever"`
	Playtime_mac_forever     int `json:"playtime_mac_forever"`
	Playtime_linux_forever   int `json:"playtime_linux_forever"`
}
type ownedGamesSteamType struct {
	Response struct {
		Game_count int `json:"game_count"`
		Games []gameSteamType `json:"games"`
	}
}
