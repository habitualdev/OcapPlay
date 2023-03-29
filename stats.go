package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"
)

func GenerateStats(filename string) {
	ocapBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	ocapData := OcapData{}
	totalShots := 0
	groups := map[string][]string{}
	weaponToShotsFired := make(map[string]int)
	players := make(map[int]PlayerStats)
	err = json.Unmarshal(ocapBytes, &ocapData)
	if err != nil {
		log.Fatalln(err)
	}
	for _, entity := range ocapData.Entities {
		if entity.Side == "WEST" && entity.IsPlayer == 1 {
			numKills := 0
			weapons := make(map[string]int)
			if groups[entity.Group] == nil {
				groups[entity.Group] = []string{}
			}
			groups[entity.Group] = append(groups[entity.Group], entity.Name)
			for _, event := range ocapData.Events {
				if len(event) == 5 {
					if reflect.TypeOf(event[3].([]interface{})[0]) == reflect.TypeOf(float64(0)) {
						if event[3].([]interface{})[0].(float64) == float64(entity.ID) {
							if event[1].(string) == "killed" {
								weapons[event[3].([]interface{})[1].(string)]++
								numKills++
							}

						}
					}

				}
			}
			players[entity.ID] = PlayerStats{
				Name:            entity.Name,
				ID:              entity.ID,
				Group:           entity.Group,
				TotalShotsFired: len(entity.FramesFired),
				TotalKills:      numKills,
				Weapons:         weapons,
			}

			for weapon, _ := range weapons {
				weaponToShotsFired[weapon] += len(entity.FramesFired)
			}

			totalShots += len(entity.FramesFired)
		}
	}

	for n, player := range players {
		mostShotWeapon := ""
		mostShots := 0

		for weapon, shots := range player.Weapons {
			if shots > mostShots {
				mostShotWeapon = weapon
				mostShots = shots
			}
		}

		players[n].Weapons[mostShotWeapon] = players[n].TotalShotsFired
		playerTemp := PlayerStats{
			Name:            players[n].Name,
			ID:              players[n].ID,
			Group:           players[n].Group,
			PrimaryWeapon:   mostShotWeapon,
			TotalShotsFired: players[n].TotalShotsFired,
			TotalKills:      players[n].TotalKills,
			Weapons:         players[n].Weapons,
		}
		players[n] = playerTemp
	}
	for _, marker := range ocapData.Markers {
		_, ok := players[int(marker[4].(float64))]
		if ok {
			tempWeapons := players[int(marker[4].(float64))].Weapons
			if strings.Contains(marker[0].(string), "magIcons") {
				addShot := 0
				tempWeapons[marker[1].(string)]++
				if strings.Contains(marker[1].(string), "Grenade") {
					tempWeapons[players[int(marker[4].(float64))].PrimaryWeapon] = tempWeapons[players[int(marker[4].(float64))].PrimaryWeapon] + 1
					addShot = 1
				}
				playerTemp := PlayerStats{
					Name:            players[int(marker[4].(float64))].Name,
					ID:              players[int(marker[4].(float64))].ID,
					TotalShotsFired: players[int(marker[4].(float64))].TotalShotsFired + addShot,
					TotalKills:      players[int(marker[4].(float64))].TotalKills,
					PrimaryWeapon:   players[int(marker[4].(float64))].PrimaryWeapon,
					Weapons:         tempWeapons,
				}
				nonPrimaryShots := 0
				for weapon, shots := range tempWeapons {
					if weapon != players[int(marker[4].(float64))].PrimaryWeapon {
						nonPrimaryShots += shots
					}
				}

				players[int(marker[4].(float64))] = playerTemp
			}
		}
	}
	sortedGroups := []string{}
	for groupnName, _ := range groups {
		sortedGroups = append(sortedGroups, groupnName)
	}
	sort.Strings(sortedGroups)
	for _, groupName := range sortedGroups {
		fmt.Println(" ======== ", groupName, " ======== ")
		for _, player := range players {
			if stringInGroup(player.Name, groups[groupName]) {
				fmt.Println(player.Name, " - Kills: ", player.TotalKills, " - Total Shots Fired: ", player.TotalShotsFired)
				if player.PrimaryWeapon == "" {
					if len(player.Weapons) > 1 {
						for weapon, shots := range player.Weapons {
							if weapon != "" && shots == player.TotalShotsFired {
								player.PrimaryWeapon = weapon
								break
							}
						}
						if player.PrimaryWeapon == "" {
							player.PrimaryWeapon = "Generic Rifle"
						}
						fmt.Println("Primary Weapon: ", player.PrimaryWeapon)
						fmt.Println("Primary Weapon Shots: ", player.Weapons[player.PrimaryWeapon])
					} else {
						player.PrimaryWeapon = "Generic Rifle"
						fmt.Println("Primary Weapon: ", player.PrimaryWeapon)
						fmt.Println("Primary Weapon Shots: ", player.Weapons[""])
					}

				} else {
					if player.PrimaryWeapon == "" {
						player.PrimaryWeapon = "Generic Rifle"
					}
					fmt.Println("Primary Weapon: ", player.PrimaryWeapon)
					fmt.Println("Primary Weapon Shots: ", player.Weapons[player.PrimaryWeapon])
				}
				if len(player.Weapons) > 1 {
					fmt.Println("Alternate Weapon Shots: ")
					for weapon, shots := range player.Weapons {
						if weapon != player.PrimaryWeapon && weapon != "" {
							fmt.Println(" - ", weapon, ":", shots)
						}
					}
				}
				if player.TotalKills > 0 {
					fmt.Println("Kills per Weapon: ")
					killMap := map[string]int{}
					for weapon, _ := range player.Weapons {
						for _, event := range ocapData.Events {

							if len(event) == 5 {
								if reflect.TypeOf(event[3].([]interface{})[0]) == reflect.TypeOf(float64(0)) {
									if event[3].([]interface{})[0].(float64) == float64(player.ID) {
										if event[1].(string) == "killed" {
											if event[3].([]interface{})[1].(string) == weapon {
												killMap[event[3].([]interface{})[1].(string)]++
											}
										}

									}
								}

							}
						}
					}
					for weapon, kills := range killMap {
						fmt.Println(" - ", weapon, ":", kills)
					}
				}
				fmt.Println("---------------------------")
			}
		}
		fmt.Println()
	}
	fmt.Println("Total Shots: ", totalShots)
}
