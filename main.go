package main

import (
	"fmt"
	"math"
)

// Afstandsmultiplier berekenen op basis van de afstand
func getDistanceMultiplier(distance float64) float64 {
	if distance < 0.1 {
		return 0.8
	} else if distance >= 0.1 && distance < 1 {
		return 0.8 + (distance-0.1)*(0.9-0.8)/0.9
	} else if distance >= 1 && distance <= 5 {
		// Lineaire schaal tussen 1 km en 5 km
		return 0.9 + (distance-1)*(1.0-0.9)/4
	} else if distance <= 10 {
		return 1.0 + (distance-5)*(1.1-1.0)/5
	} else if distance <= 15 {
		return 1.1 + (distance-10)*(1.15-1.1)/5
	} else if distance <= 21.1 {
		return 1.15 + (distance-15)*(1.2-1.15)/6.1
	} else if distance <= 42.3 {
		return 1.2 + (distance-21.1)*(1.3-1.2)/21.2
	} else if distance < 100 {
		return 1.3 + (distance-42.3)*(1.4-1.3)/57.7
	} else {
		return 1.4 // Max multiplier voor afstanden groter dan 100 km
	}
}

// Deze functie converteert snelheid naar rating
func getRatingForSpeed(speed, distance float64) float64 {
	multiplier := getDistanceMultiplier(distance)
	baseRating := calculateBaseRating(speed) // Hier moet je zorgen dat de logica correct is
	return baseRating * multiplier
}

// Deze functie converteert rating naar snelheid
func getSpeedForRating(rating, distance float64) float64 {
	multiplier := getDistanceMultiplier(distance)
	adjustedRating := rating / multiplier
	return calculateSpeedFromRating(adjustedRating) // Zorg dat deze functie bestaat en correct is
}

// Helperfunctie om basisrating te berekenen (zonder multiplier)
func calculateBaseRating(speed float64) float64 {
	if speed > 15 {
		return 2000 + (speed-15)*100
	}
	// Lineaire interpolatie tussen bekende punten
	points := []struct {
		speed  float64
		rating float64
	}{
		{0, 0},
		{5, 600},
		{11, 1500},
		{15, 2000},
	}

	// Zoek het juiste segment voor interpolatie
	for i := 0; i < len(points)-1; i++ {
		if speed >= points[i].speed && speed <= points[i+1].speed {
			low := points[i]
			high := points[i+1]
			baseRating := low.rating + (speed-low.speed)*(high.rating-low.rating)/(high.speed-low.speed)
			return baseRating
		}
	}
	return 0
}

// Helperfunctie om snelheid te berekenen van rating
func calculateSpeedFromRating(adjustedRating float64) float64 {
	if adjustedRating > 2000 {
		return 15 + (adjustedRating-2000)/100
	}
	points := []struct {
		speed  float64
		rating float64
	}{
		{0, 0},
		{5, 600},
		{11, 1500},
		{15, 2000},
	}

	for i := 0; i < len(points)-1; i++ {
		if adjustedRating >= points[i].rating && adjustedRating <= points[i+1].rating {
			low := points[i]
			high := points[i+1]
			interpolatedSpeed := low.speed + (adjustedRating-low.rating)*(high.speed-low.speed)/(high.rating-low.rating)
			return interpolatedSpeed
		}
	}
	return 0
}

// Tijd berekenen voor een gegeven afstand en snelheid
func calculateTime(distance, speed float64) (int, int, int) {
	if speed == 0 {
		return 0, 0, 0
	}
	totalSeconds := int(math.Round(distance / speed * 3600))
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	return hours, minutes, seconds
}

// Hoofdfunctie met menu voor gebruikerskeuze
func main() {
	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1. Get Rating")
		fmt.Println("2. Get Speed/Time")
		fmt.Println("3. Rating table")
		var choice int
		fmt.Print("Enter choice: ")
		fmt.Scan(&choice)

		var distance float64
		var hours, minutes int
		var seconds, rating, startr, endr, stepr float64

		switch choice {
		case 1:
			fmt.Print("Enter distance (in km): ")
			fmt.Scan(&distance)
			fmt.Print("Enter time (hours minutes seconds): ")
			fmt.Scan(&hours, &minutes, &seconds)
			timeInHours := float64(hours) + float64(minutes)/60 + seconds/3600
			speed := distance / timeInHours
			rating := getRatingForSpeed(speed, distance)
			fmt.Printf("Your speed is: %.2f km/h\n", speed)
			fmt.Printf("Your running rating is: %.0f ELO\n", rating)
		case 2:
			fmt.Print("Enter distance (in km): ")
			fmt.Scan(&distance)
			fmt.Print("Enter desired rating: ")
			fmt.Scan(&rating)
			speed := getSpeedForRating(rating, distance)
			h, m, s := calculateTime(distance, speed)
			fmt.Printf("Rating of %.0f ELO for a distance of %.2f km, Average speed of %.2f km/h.\n", rating, distance, speed)
			fmt.Printf("Time of %d hours %d minutes %d seconds.\n", h, m, s)
		case 3:
			fmt.Print("Enter distance (in km): ")
			fmt.Scan(&distance)
			fmt.Print("Start rating: ")
			fmt.Scan(&startr)
			fmt.Print("End rating: ")
			fmt.Scan(&endr)
			fmt.Print("Each step is x rating: ")
			fmt.Scan(&stepr)
			for i := startr; i <= endr; i += stepr {
				speed := getSpeedForRating(i, distance)
				h, m, s := calculateTime(distance, speed)
				fmt.Printf("\nRating %.0f ELO, Distance %.2f km, Speed %.2f km/h\n", i, distance, speed)
				fmt.Printf("Time %d hours %d minutes %d seconds.\n", h, m, s)
			}
		}
	}
}
