package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

const (
	// Change these to match your router's LED directory names
	basePath        = "/sys/class/leds"
	defaultRedLED   = "red:status"
	defaultGreenLED = "green:status"
	defaultBlueLED  = "blue:status"

	// Tuning variables
	defaultStep  = "1"
	defaultDelay = "1000"
)

// Pre-computed LED file paths to avoid repeated filepath.Join calls.
var (
	redLED   = getWithDefault(os.LookupEnv, "RED_LED", defaultRedLED)
	greenLED = getWithDefault(os.LookupEnv, "GREEN_LED", defaultGreenLED)
	blueLED  = getWithDefault(os.LookupEnv, "BLUE_LED", defaultBlueLED)

	redTriggerPath   = filepath.Join(basePath, redLED, "trigger")
	greenTriggerPath = filepath.Join(basePath, greenLED, "trigger")
	blueTriggerPath  = filepath.Join(basePath, blueLED, "trigger")

	redBrightnessPath   = filepath.Join(basePath, redLED, "brightness")
	greenBrightnessPath = filepath.Join(basePath, greenLED, "brightness")
	blueBrightnessPath  = filepath.Join(basePath, blueLED, "brightness")

	step, _     = strconv.Atoi(getWithDefault(os.LookupEnv, "STEP", defaultStep))
	delayRaw, _ = strconv.Atoi(getWithDefault(os.LookupEnv, "DELAY", defaultDelay))
	delay       = time.Duration(delayRaw) * time.Millisecond
)

func getWithDefault(f func(string) (string, bool), key, fallback string) string {
	if value, exists := f(key); exists {
		return value
	}
	return fallback
}

// writeToFile writes a string to a file, ignoring errors just like `2>/dev/null`
func writeToFile(path, data string) {
	_ = os.WriteFile(path, []byte(data), 0644)
}

// clearTriggers removes system control from the LEDs so we can control them manually
func clearTriggers() {
	writeToFile(redTriggerPath, "none")
	writeToFile(greenTriggerPath, "none")
	writeToFile(blueTriggerPath, "none")
}

// setColors updates the brightness for all three channels
func setColors(r, g, b int) {
	writeToFile(redBrightnessPath, strconv.Itoa(r))
	writeToFile(greenBrightnessPath, strconv.Itoa(g))
	writeToFile(blueBrightnessPath, strconv.Itoa(b))
}

func main() {
	fmt.Println("Using LED paths:")
	fmt.Println("Red LED:", redLED)
	fmt.Println("Green LED:", greenLED)
	fmt.Println("Blue LED:", blueLED)

	// Set up channel to listen for Ctrl+C (SIGINT/SIGTERM) for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nStopping sweep. Resetting LEDs...")
		setColors(0, 0, 0)
		os.Exit(0)
	}()

	fmt.Println("Starting smooth RGB color sweep... Press [CTRL+C] to stop.")

	clearTriggers()

	// Initialize at pure Red
	r, g, b := 255, 0, 0
	setColors(r, g, b)

	for {
		// 1. Red to Yellow (Fade Green up)
		for g < 255 {
			g += step
			if g > 255 {
				g = 255
			} // Prevent overflow if step doesn't divide evenly
			setColors(r, g, b)
			time.Sleep(delay)
		}

		// 2. Yellow to Green (Fade Red down)
		for r > 0 {
			r -= step
			if r < 0 {
				r = 0
			}
			setColors(r, g, b)
			time.Sleep(delay)
		}

		// 3. Green to Cyan (Fade Blue up)
		for b < 255 {
			b += step
			if b > 255 {
				b = 255
			}
			setColors(r, g, b)
			time.Sleep(delay)
		}

		// 4. Cyan to Blue (Fade Green down)
		for g > 0 {
			g -= step
			if g < 0 {
				g = 0
			}
			setColors(r, g, b)
			time.Sleep(delay)
		}

		// 5. Blue to Magenta (Fade Red up)
		for r < 255 {
			r += step
			if r > 255 {
				r = 255
			}
			setColors(r, g, b)
			time.Sleep(delay)
		}

		// 6. Magenta to Red (Fade Blue down)
		for b > 0 {
			b -= step
			if b < 0 {
				b = 0
			}
			setColors(r, g, b)
			time.Sleep(delay)
		}
	}
}
