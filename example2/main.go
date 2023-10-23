package main

import (
	"image/color"
	"machine"
	"math"
	"math/rand"
	"time"

	"tinygo.org/x/drivers/uc8151"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

// Preferences
const (
	// [DEFAULT, MEDIUM, FAST, TURBO]
	// DEFAULT is the best quality, but takes the longest to draw
	// TURBO is very fast, but *by far* the lowest quality
	EPDSPEED = uc8151.FAST

	// How many shapes to draw & what's the max/min size range
	SHAPECOUNT   = 20
	SHAPESIZEMAX = 50
	SHAPESIZEMIN = 20

	// How many pixels to move each frame & how long to sleep between frames
	MAXMOVESIZE   = 10
	SLEEPYSECONDS = 10

	// If the shape goes out of bounds, should it wrap around?
	WRAPAROUND = true

	// Text at the bottom of the screen
	// leave both/either line blank if unwanted
	LINEA1 = "H3's badge"
	LINEA2 = "Tinkering again 2023-10-23"

	LINEB1 = "I love you"
	LINEB2 = "Still tinkering"

	LINEC1 = "Humphrey Drummond"
	LINEC2 = "Atom Senior Backend Engineer"
)

// Shapes to include (comment out a shape to disable it)
var SHAPES = []ShapeType{
	CIRCLE,
	RECTANGLE,
	TRIANGLE,
}

//

var led machine.Pin
var display uc8151.Device

// Color constants
var (
	WHITE = color.RGBA{0, 0, 0, 0}
	BLACK = color.RGBA{1, 1, 1, 255}
)

func main() {
	// rand.Seed(time.Now().UnixNano())

	led = machine.LED
	led.Configure(
		machine.PinConfig{Mode: machine.PinOutput},
	)

	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 12000000,
		SCK:       machine.SPI0_SCK_PIN,
		SDO:       machine.SPI0_SDO_PIN,
	})

	display = uc8151.New(
		machine.SPI0,
		machine.EPD_CS_PIN,
		machine.EPD_DC_PIN,
		machine.EPD_RESET_PIN,
		machine.EPD_BUSY_PIN,
	)
	display.Configure(uc8151.Config{
		Blocking: true,
		Rotation: uc8151.ROTATION_270,
		Speed:    EPDSPEED,
	})
	displayWidth, displayHeight := display.Size()

	shapes := make([]Shape, SHAPECOUNT)
	for i := 0; i < cap(shapes); i++ {
		size := int16(math.Max(SHAPESIZEMAX, 1))
		if SHAPESIZEMAX > SHAPESIZEMIN {
			size = int16(rand.Intn(int(SHAPESIZEMAX-SHAPESIZEMIN))) + SHAPESIZEMIN
		}
		x := int16(rand.Intn(int(displayWidth)))
		y := int16(rand.Intn(int(displayHeight)))
		shape := SHAPES[rand.Intn(len(SHAPES))]
		switch shape {
		case CIRCLE:
			shapes[i] = NewCircle(x, y, size, BLACK, WHITE)
			break
		case RECTANGLE:
			shapes[i] = NewRectangle(x, y, size, size, BLACK, WHITE)
			break
		case TRIANGLE:
			shapes[i] = NewTriangle(x, y, size, BLACK, WHITE)
			break
		}
	}

	for {
		// led.Set(!led.Get()) // toggle LED but it doesn't work
		// as the LED is on the wireless chip
		// https://github.com/tinygo-org/tinygo/issues/2947#user-content-fn-A-9d58b855fdc06c9f1071082446c3d201

		display.ClearBuffer()
		display.ClearDisplay()
		for i, shape := range shapes {
			// Move shape.X and shape.Y
			shape.X += int16(rand.Intn(MAXMOVESIZE*2)) - MAXMOVESIZE
			shape.Y += int16(rand.Intn(MAXMOVESIZE*2)) - MAXMOVESIZE
			hitbox := shape.Hitbox()

			// Reposition X if out of bounds
			if hitbox.MinX < 0 {
				if WRAPAROUND {
					shape.X = displayWidth - (hitbox.SizeX - hitbox.RegistrationX)
				} else {
					shape.X = 0
				}
			} else if hitbox.MaxX > displayWidth {
				if WRAPAROUND {
					shape.X = hitbox.RegistrationX
				} else {
					shape.X = shape.X - (hitbox.MaxX - displayWidth)
				}
			}

			// Reposition Y if out of bounds
			if hitbox.MinY < 0 {
				if WRAPAROUND {
					shape.Y = displayHeight - (hitbox.SizeY - hitbox.RegistrationY)
				} else {
					shape.Y = 0
				}
			} else if hitbox.MaxY > displayHeight {
				if WRAPAROUND {
					shape.Y = hitbox.RegistrationY
				} else {
					shape.Y = shape.Y - (hitbox.MaxY - displayHeight)
				}
			}

			// Save shape changes
			shapes[i] = shape

			// Draw shape
			shape.Plot(&display)
		}

		if LINEA1 != "" {
			y := displayHeight - 34
			if LINEA2 == "" {
				y += 20
			}
			tinyfont.WriteLineRotated(&display, &freemono.Bold24pt7b, 8, y, LINEA1, BLACK, tinyfont.NO_ROTATION)
		}
		if LINEA2 != "" {
			tinyfont.WriteLineRotated(&display, &freemono.Bold9pt7b, 10, displayHeight-10, LINEA2, BLACK, tinyfont.NO_ROTATION)
		}

		display.Display()

		time.Sleep(SLEEPYSECONDS * time.Second)
	}
}
