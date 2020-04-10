package kreedz

import (
	"math/big"
	"regexp"
	"strings"
)

//type SteamID string   // STEAM_0:0:33403241
//type SteamID64 uint64 // 76561198132612090

func SteamID64ToSteamID(steam64 string) string {
	steamID, _ := new(big.Int).SetString(steam64, 10)
	magic, _ := new(big.Int).SetString("76561197960265728", 10)
	steamID = steamID.Sub(steamID, magic)
	isServer := new(big.Int).And(steamID, big.NewInt(1))
	steamID = steamID.Sub(steamID, isServer)
	steamID = steamID.Div(steamID, big.NewInt(2))
	return "STEAM_0:" + isServer.String() + ":" + steamID.String()
}

func SteamIDToSteamID64(steamID string) string {
	idParts := strings.Split(string(steamID), ":")
	magic, _ := new(big.Int).SetString("76561197960265728", 10)
	steam64, _ := new(big.Int).SetString(idParts[2], 10)
	steam64 = steam64.Mul(steam64, big.NewInt(2))
	steam64 = steam64.Add(steam64, magic)
	auth, _ := new(big.Int).SetString(idParts[1], 10)
	return steam64.Add(steam64, auth).String()
}

func IsSteamID(steamID string) bool {
	if len(steamID) > 19 {
		return false
	}
	match, _ := regexp.MatchString("STEAM_0:[0-1]:[0-9]{6}", steamID)
	return match
}