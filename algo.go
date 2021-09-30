package main

import (
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
)

func MoveWater(waterMap, waterMapNew, groundMap map[Vec3i32]bool) {
	for key := range waterMap {
		if waterMap[key] {
			if !waterMap[Vec3i32{key[0], key[1] - 1, key[2]}] && key[1]-1 >= 0 &&
				!groundMap[Vec3i32{key[0], key[1] - 1, key[2]}] &&
				!waterMapNew[Vec3i32{key[0], key[1] - 1, key[2]}] {

				waterMapNew[Vec3i32{key[0], key[1] - 1, key[2]}] = true
				waterMap[key] = false
			} else {
				var nextPosArray [4]Vec3i32 = [4]Vec3i32{
					{key[0] + 1, key[1], key[2]},
					{key[0] - 1, key[1], key[2]},
					{key[0], key[1], key[2] + 1},
					{key[0], key[1], key[2] - 1},
				}
				rand.Shuffle(len(nextPosArray), func(i, j int) { nextPosArray[i], nextPosArray[j] = nextPosArray[j], nextPosArray[i] })
				isFoundMed := false
				foundedMed := Vec3i32{0, 0, 0}

				didFound, isBest, posFound := stupidSearch(waterMap, waterMapNew, groundMap, nextPosArray[0])
				if didFound {
					isFoundMed = true
					foundedMed = posFound
				}
				if isBest {
					waterMapNew[posFound] = true
					waterMap[key] = false
				} else {
					didFound, isBest, posFound = stupidSearch(waterMap, waterMapNew, groundMap, nextPosArray[1])
					if didFound {
						isFoundMed = true
						foundedMed = posFound
					}
					if isBest {
						waterMapNew[posFound] = true
						waterMap[key] = false
					} else {
						didFound, isBest, posFound = stupidSearch(waterMap, waterMapNew, groundMap, nextPosArray[2])
						if didFound {
							isFoundMed = true
							foundedMed = posFound
						}
						if isBest {
							waterMapNew[posFound] = true
							waterMap[key] = false
						} else {
							didFound, isBest, posFound = stupidSearch(waterMap, waterMapNew, groundMap, nextPosArray[3])
							if didFound {
								isFoundMed = true
								foundedMed = posFound
							}
							if isBest {
								waterMapNew[posFound] = true
								waterMap[key] = false
							}
						}
					}
				}
				if !isFoundMed {
					waterMapNew[key] = true
				} else {
					waterMapNew[foundedMed] = true
					waterMap[key] = false
				}

			}
		}
	}
}

func MoveWaterVS(waterMap, waterMapNew, groundMap map[Vec3i32]bool) {
	for key := range waterMap {
		if waterMap[key] {
			if !waterMap[Vec3i32{key[0], key[1] - 1, key[2]}] && key[1]-1 >= 0 &&
				!groundMap[Vec3i32{key[0], key[1] - 1, key[2]}] &&
				!waterMapNew[Vec3i32{key[0], key[1] - 1, key[2]}] {

				waterMapNew[Vec3i32{key[0], key[1] - 1, key[2]}] = true
				waterMap[key] = false
			} else {
				var nextPosArray [4]Vec3i32 = [4]Vec3i32{
					{key[0] + 1, key[1], key[2]},
					{key[0] - 1, key[1], key[2]},
					{key[0], key[1], key[2] + 1},
					{key[0], key[1], key[2] - 1},
				}
				rand.Shuffle(len(nextPosArray), func(i, j int) { nextPosArray[i], nextPosArray[j] = nextPosArray[j], nextPosArray[i] })
				didFound, posFound := veryStupidSearch(waterMap, waterMapNew, groundMap, nextPosArray[0])
				if didFound {
					waterMapNew[posFound] = true
					waterMap[key] = false
				} else {
					waterMapNew[key] = true
				}
			}
		}
	}
}

func MoveWaterFloodfill(waterMap, waterMapNew, groundMap map[Vec3i32]bool) {
	for key := range waterMap {
		if waterMap[key] {
			if !waterMap[Vec3i32{key[0], key[1] - 1, key[2]}] && key[1]-1 >= 0 &&
				!groundMap[Vec3i32{key[0], key[1] - 1, key[2]}] &&
				!waterMapNew[Vec3i32{key[0], key[1] - 1, key[2]}] {

				waterMapNew[Vec3i32{key[0], key[1] - 1, key[2]}] = true
				waterMap[key] = false
			} else if groundMap[Vec3i32{key[0], key[1] - 1, key[2]}] &&
				groundMap[Vec3i32{key[0], key[1] - 1, key[2] + 1}] &&
				groundMap[Vec3i32{key[0], key[1] - 1, key[2] - 1}] &&
				groundMap[Vec3i32{key[0] + 1, key[1] - 1, key[2]}] &&
				groundMap[Vec3i32{key[0] - 1, key[1] - 1, key[2]}] &&

				groundMap[Vec3i32{key[0] - 1, key[1] - 1, key[2] + 1}] &&
				groundMap[Vec3i32{key[0] - 1, key[1] - 1, key[2] - 1}] &&
				groundMap[Vec3i32{key[0] + 1, key[1] - 1, key[2] + 1}] &&
				groundMap[Vec3i32{key[0] + 1, key[1] - 1, key[2] - 1}] {
				waterMapNew[key] = true
			} else {
				visitedMap := make(map[Vec3i32]bool)
				var nextPosArray [4]Vec3i32 = [4]Vec3i32{
					{key[0] + 1, key[1], key[2]},
					{key[0] - 1, key[1], key[2]},
					{key[0], key[1], key[2] + 1},
					{key[0], key[1], key[2] - 1},
				}
				rand.Shuffle(len(nextPosArray), func(i, j int) { nextPosArray[i], nextPosArray[j] = nextPosArray[j], nextPosArray[i] })
				didFound, posFound := search(waterMap, waterMapNew, groundMap, visitedMap, key, key, nextPosArray[0], 0, 100)
				if didFound {
					if !waterMapNew[posFound] && !waterMap[posFound] && !groundMap[posFound] {
						waterMapNew[posFound] = true
						waterMap[key] = false
					} else {
						waterMapNew[key] = true
					}
				} else {
					didFound, posFound := search(waterMap, waterMapNew, groundMap, visitedMap, key, key, nextPosArray[1], 0, 100)
					if didFound {
						if !waterMapNew[posFound] && !waterMap[posFound] && !groundMap[posFound] {
							waterMapNew[posFound] = true
							waterMap[key] = false
						} else {
							waterMapNew[key] = true
						}
					} else {
						didFound, posFound := search(waterMap, waterMapNew, groundMap, visitedMap, key, key, nextPosArray[2], 0, 100)
						if didFound {
							if !waterMapNew[posFound] && !waterMap[posFound] && !groundMap[posFound] {
								waterMapNew[posFound] = true
								waterMap[key] = false
							} else {
								waterMapNew[key] = true
							}
						} else {
							didFound, posFound := search(waterMap, waterMapNew, groundMap, visitedMap, key, key, nextPosArray[3], 0, 100)
							if didFound {
								if !waterMapNew[posFound] && !waterMap[posFound] && !groundMap[posFound] {
									waterMapNew[posFound] = true
									waterMap[key] = false
								} else {
									waterMapNew[key] = true
								}
							} else {
								waterMapNew[key] = true
							}
						}
					}
				}
			}
		}
	}
}

func MoveWaterSnow(waterMap, waterMapNew, groundMap map[Vec3i32]bool) {
	cWater = mgl32.Vec3{1, 1, 1}
	for key := range waterMap {
		if waterMap[key] {
			if !waterMap[Vec3i32{key[0], key[1] - 1, key[2]}] && key[1]-1 >= 0 &&
				!groundMap[Vec3i32{key[0], key[1] - 1, key[2]}] &&
				!waterMapNew[Vec3i32{key[0], key[1] - 1, key[2]}] {

				waterMapNew[Vec3i32{key[0], key[1] - 1, key[2]}] = true
				waterMap[key] = false
			} else if groundMap[Vec3i32{key[0], key[1] - 1, key[2]}] {
				waterMapNew[key] = true
			} else {
				visitedMap := make(map[Vec3i32]bool)
				var nextPosArray [4]Vec3i32 = [4]Vec3i32{
					{key[0] + 1, key[1], key[2]},
					{key[0] - 1, key[1], key[2]},
					{key[0], key[1], key[2] + 1},
					{key[0], key[1], key[2] - 1},
				}
				rand.Shuffle(len(nextPosArray), func(i, j int) { nextPosArray[i], nextPosArray[j] = nextPosArray[j], nextPosArray[i] })
				didFound, posFound := search(waterMap, waterMapNew, groundMap, visitedMap, key, key, nextPosArray[0], 0, 5)
				if didFound {
					if !waterMapNew[posFound] && !waterMap[posFound] && !groundMap[posFound] {
						waterMapNew[posFound] = true
						waterMap[key] = false
					} else {
						waterMapNew[key] = true
					}
				} else {
					didFound, posFound := search(waterMap, waterMapNew, groundMap, visitedMap, key, key, nextPosArray[1], 0, 5)
					if didFound {
						if !waterMapNew[posFound] && !waterMap[posFound] && !groundMap[posFound] {
							waterMapNew[posFound] = true
							waterMap[key] = false
						} else {
							waterMapNew[key] = true
						}
					} else {
						didFound, posFound := search(waterMap, waterMapNew, groundMap, visitedMap, key, key, nextPosArray[2], 0, 5)
						if didFound {
							if !waterMapNew[posFound] && !waterMap[posFound] && !groundMap[posFound] {
								waterMapNew[posFound] = true
								waterMap[key] = false
							} else {
								waterMapNew[key] = true
							}
						} else {
							didFound, posFound := search(waterMap, waterMapNew, groundMap, visitedMap, key, key, nextPosArray[3], 0, 5)
							if didFound {
								if !waterMapNew[posFound] && !waterMap[posFound] && !groundMap[posFound] {
									waterMapNew[posFound] = true
									waterMap[key] = false
								} else {
									waterMapNew[key] = true
								}
							} else {
								waterMapNew[key] = true
							}
						}
					}
				}
			}
		}
	}
}

func search(waterMap, waterMapNew, groundMap, visitedMap map[Vec3i32]bool, initialPos, oldPos, currPos Vec3i32, nbIteration, maxIteration int) (bool, Vec3i32) {
	if visitedMap[currPos] || groundMap[currPos] || waterMap[currPos] || waterMapNew[currPos] ||
		currPos[0] < 0 || currPos[1] < 0 || currPos[2] < 0 ||
		currPos[0] >= simulationSize || currPos[2] >= simulationSize ||
		nbIteration > maxIteration {
		return false, Vec3i32{0, 0, 0}
	} else {
		visitedMap[currPos] = true
	}

	if !waterMap[Vec3i32{currPos[0], currPos[1] - 1, currPos[2]}] &&
		!groundMap[Vec3i32{currPos[0], currPos[1] - 1, currPos[2]}] &&
		!waterMapNew[Vec3i32{currPos[0], currPos[1] - 1, currPos[2]}] &&
		currPos[1]-1 >= 0 {

		if !waterMapNew[Vec3i32{currPos[0], currPos[1] - 1, currPos[2]}] {
			return true, Vec3i32{currPos[0], currPos[1] - 1, currPos[2]}
		} else {
			return false, Vec3i32{0, 0, 0}
		}
	} else {
		var nextPosArray [4]Vec3i32 = [4]Vec3i32{
			{currPos[0] + 1, currPos[1], currPos[2]},
			{currPos[0] - 1, currPos[1], currPos[2]},
			{currPos[0], currPos[1], currPos[2] + 1},
			{currPos[0], currPos[1], currPos[2] - 1},
		}
		rand.Shuffle(len(nextPosArray), func(i, j int) { nextPosArray[i], nextPosArray[j] = nextPosArray[j], nextPosArray[i] })

		newPos := nextPosArray[0]
		didFound, posFound := search(waterMap, waterMapNew, groundMap, visitedMap, initialPos, currPos, newPos, nbIteration+1, maxIteration)
		if didFound {
			return true, Vec3i32{posFound[0], posFound[1], posFound[2]}
		} else {
			newPos := nextPosArray[1]
			didFound, posFound := search(waterMap, waterMapNew, groundMap, visitedMap, initialPos, currPos, newPos, nbIteration+1, maxIteration)
			if didFound {
				return true, Vec3i32{posFound[0], posFound[1], posFound[2]}
			} else {
				newPos := nextPosArray[2]
				didFound, posFound := search(waterMap, waterMapNew, groundMap, visitedMap, initialPos, currPos, newPos, nbIteration+1, maxIteration)
				if didFound {
					return true, Vec3i32{posFound[0], posFound[1], posFound[2]}
				} else {
					newPos := nextPosArray[3]
					didFound, posFound := search(waterMap, waterMapNew, groundMap, visitedMap, initialPos, currPos, newPos, nbIteration+1, maxIteration)
					if didFound {
						return true, Vec3i32{posFound[0], posFound[1], posFound[2]}
					} else {
						return false, Vec3i32{0, 0, 0}
					}
				}
			}
		}
	}
}

func stupidSearch(waterMap, waterMapNew, groundMap map[Vec3i32]bool, currPos Vec3i32) (bool, bool, Vec3i32) {
	if groundMap[currPos] || waterMap[currPos] || waterMapNew[currPos] ||
		currPos[0] < 0 || currPos[1] < 0 || currPos[2] < 0 ||
		currPos[0] >= simulationSize || currPos[2] >= simulationSize {
		return false, false, Vec3i32{0, 0, 0}
	}
	var newPos Vec3i32 = Vec3i32{currPos[0], currPos[1] - 1, currPos[2]}
	if groundMap[newPos] || waterMap[newPos] || waterMapNew[newPos] ||
		newPos[0] < 0 || newPos[1] < 0 || newPos[2] < 0 ||
		newPos[0] >= simulationSize || newPos[1] >= simulationSize || newPos[2] >= simulationSize {
		return true, true, currPos
	}

	newPos = Vec3i32{currPos[0] - 1, currPos[1] - 1, currPos[2]}
	if !groundMap[newPos] && !waterMap[newPos] && !waterMapNew[newPos] &&
		newPos[0] >= 0 && newPos[1] >= 0 && newPos[2] >= 0 &&
		newPos[0] < simulationSize && newPos[1] < simulationSize && newPos[2] < simulationSize {
		newPos = Vec3i32{currPos[0] - 1, currPos[1], currPos[2]}
		if !groundMap[newPos] && !waterMap[newPos] && !waterMapNew[newPos] &&
			newPos[0] >= 0 && newPos[1] >= 0 && newPos[2] >= 0 &&
			newPos[0] < simulationSize && newPos[1] < simulationSize && newPos[2] < simulationSize {
			return true, true, newPos
		}
	}

	newPos = Vec3i32{currPos[0] + 1, currPos[1] - 1, currPos[2]}
	if !groundMap[newPos] && !waterMap[newPos] && !waterMapNew[newPos] &&
		newPos[0] >= 0 && newPos[1] >= 0 && newPos[2] >= 0 &&
		newPos[0] < simulationSize && newPos[1] < simulationSize && newPos[2] < simulationSize {
		newPos = Vec3i32{currPos[0] + 1, currPos[1], currPos[2]}
		if !groundMap[newPos] && !waterMap[newPos] && !waterMapNew[newPos] &&
			newPos[0] >= 0 && newPos[1] >= 0 && newPos[2] >= 0 &&
			newPos[0] < simulationSize && newPos[1] < simulationSize && newPos[2] < simulationSize {
			return true, true, newPos
		}
	}

	newPos = Vec3i32{currPos[0], currPos[1] - 1, currPos[2] + 1}
	if !groundMap[newPos] && !waterMap[newPos] && !waterMapNew[newPos] &&
		newPos[0] >= 0 && newPos[1] >= 0 && newPos[2] >= 0 &&
		newPos[0] < simulationSize && newPos[1] < simulationSize && newPos[2] < simulationSize {
		newPos = Vec3i32{currPos[0], currPos[1], currPos[2] + 1}
		if !groundMap[newPos] && !waterMap[newPos] && !waterMapNew[newPos] &&
			newPos[0] >= 0 && newPos[1] >= 0 && newPos[2] >= 0 &&
			newPos[0] < simulationSize && newPos[1] < simulationSize && newPos[2] < simulationSize {
			return true, true, newPos
		}
	}

	newPos = Vec3i32{currPos[0], currPos[1] - 1, currPos[2] - 1}
	if !groundMap[newPos] && !waterMap[newPos] && !waterMapNew[newPos] &&
		newPos[0] >= 0 && newPos[1] >= 0 && newPos[2] >= 0 &&
		newPos[0] < simulationSize && newPos[1] < simulationSize && newPos[2] < simulationSize {
		newPos = Vec3i32{currPos[0], currPos[1], currPos[2] - 1}
		if !groundMap[newPos] && !waterMap[newPos] && !waterMapNew[newPos] &&
			newPos[0] >= 0 && newPos[1] >= 0 && newPos[2] >= 0 &&
			newPos[0] < simulationSize && newPos[1] < simulationSize && newPos[2] < simulationSize {
			return true, true, newPos
		}
	}

	return true, false, currPos
}

func veryStupidSearch(waterMap, waterMapNew, groundMap map[Vec3i32]bool, currPos Vec3i32) (bool, Vec3i32) {
	if groundMap[currPos] || waterMap[currPos] || waterMapNew[currPos] ||
		currPos[0] < 0 || currPos[1] < 0 || currPos[2] < 0 ||
		currPos[0] >= simulationSize || currPos[2] >= simulationSize {
		return false, Vec3i32{0, 0, 0}
	}

	return true, currPos
}
