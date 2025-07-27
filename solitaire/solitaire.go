package main

import (
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const windowMultiplier = 2.5

// 16:9 aspect ratio
const windowWidth = int(1280 * windowMultiplier)
const windowHeight = int(720 * windowMultiplier)

// standard card ratio is 5:7
const cardMultipler = windowMultiplier
const cardHeight = int(140 * cardMultipler)
const cardWidth = int(100 * cardMultipler)
const cardBorder = int(10 * cardMultipler)
const cardNumberHeight = int(30 * cardMultipler)

var green = color.RGBA{120, 200, 120, 0xff}
var leftMouseButton = ebiten.MouseButton0

type symbol int

const (
	Diamond symbol = iota
	Heart
	Spade
	Club
)

type value int

const (
	Ace value = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

var cardBackImage *ebiten.Image
var cardPlaceholder *ebiten.Image

// diamond images
var diamondAceImage *ebiten.Image
var diamondTwoImage *ebiten.Image
var diamondThreeImage *ebiten.Image
var diamondFourImage *ebiten.Image
var diamondFiveImage *ebiten.Image
var diamondSixImage *ebiten.Image
var diamondSevenImage *ebiten.Image
var diamondEightImage *ebiten.Image
var diamondNineImage *ebiten.Image
var diamondTenImage *ebiten.Image
var diamondJackImage *ebiten.Image
var diamondQueenImage *ebiten.Image
var diamondKingImage *ebiten.Image

// heart images
var heartAceImage *ebiten.Image
var heartTwoImage *ebiten.Image
var heartThreeImage *ebiten.Image
var heartFourImage *ebiten.Image
var heartFiveImage *ebiten.Image
var heartSixImage *ebiten.Image
var heartSevenImage *ebiten.Image
var heartEightImage *ebiten.Image
var heartNineImage *ebiten.Image
var heartTenImage *ebiten.Image
var heartJackImage *ebiten.Image
var heartQueenImage *ebiten.Image
var heartKingImage *ebiten.Image

// spade images
var spadeAceImage *ebiten.Image
var spadeTwoImage *ebiten.Image
var spadeThreeImage *ebiten.Image
var spadeFourImage *ebiten.Image
var spadeFiveImage *ebiten.Image
var spadeSixImage *ebiten.Image
var spadeSevenImage *ebiten.Image
var spadeEightImage *ebiten.Image
var spadeNineImage *ebiten.Image
var spadeTenImage *ebiten.Image
var spadeJackImage *ebiten.Image
var spadeQueenImage *ebiten.Image
var spadeKingImage *ebiten.Image

// club images
var clubAceImage *ebiten.Image
var clubTwoImage *ebiten.Image
var clubThreeImage *ebiten.Image
var clubFourImage *ebiten.Image
var clubFiveImage *ebiten.Image
var clubSixImage *ebiten.Image
var clubSevenImage *ebiten.Image
var clubEightImage *ebiten.Image
var clubNineImage *ebiten.Image
var clubTenImage *ebiten.Image
var clubJackImage *ebiten.Image
var clubQueenImage *ebiten.Image
var clubKingImage *ebiten.Image

var cardDeck []card

var foundations [4]*foundation
var tableaus [7]*tableau
var stockPile *stock
var wastePile *waste

func init() {
	var monitors []*ebiten.MonitorType
	monitors = ebiten.AppendMonitors(monitors)
	// TODO: Remove before release in order to avoid bugs for customers which own this monitor
	for _, m := range monitors {
		if m.Name() == "DELL U2720Q" {
			ebiten.SetMonitor(m)
		}
	}

	// load card images
	cardBackImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_card_back.png")
	cardPlaceholder, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_card_placeholder.png")

	// load diamond images
	diamondAceImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_diamond_one.png")
	diamondTwoImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_diamond_two.png")
	diamondThreeImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_diamond_three.png")
	diamondFourImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_diamond_four.png")
	diamondFiveImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_diamond_five.png")
	diamondSixImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_diamond_six.png")
	diamondSevenImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_diamond_seven.png")
	diamondEightImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_diamond_eight.png")
	diamondNineImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_diamond_nine.png")
	diamondTenImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_diamond_ten.png")
	diamondJackImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_diamond_jack.png")
	diamondQueenImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_diamond_queen.png")
	diamondKingImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_diamond_king.png")

	// load heart images
	heartAceImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_heart_one.png")
	heartTwoImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_heart_two.png")
	heartThreeImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_heart_three.png")
	heartFourImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_heart_four.png")
	heartFiveImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_heart_five.png")
	heartSixImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_heart_six.png")
	heartSevenImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_heart_seven.png")
	heartEightImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_heart_eight.png")
	heartNineImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_heart_nine.png")
	heartTenImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_heart_ten.png")
	heartJackImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_heart_jack.png")
	heartQueenImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_heart_queen.png")
	heartKingImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_heart_king.png")

	// load spade images
	spadeAceImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_spade_one.png")
	spadeTwoImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_spade_two.png")
	spadeThreeImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_spade_three.png")
	spadeFourImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_spade_four.png")
	spadeFiveImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_spade_five.png")
	spadeSixImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_spade_six.png")
	spadeSevenImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_spade_seven.png")
	spadeEightImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_spade_eight.png")
	spadeNineImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_spade_nine.png")
	spadeTenImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_spade_ten.png")
	spadeJackImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_spade_jack.png")
	spadeQueenImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_spade_queen.png")
	spadeKingImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_spade_king.png")

	// load club images
	clubAceImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_club_one.png")
	clubTwoImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_club_two.png")
	clubThreeImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_club_three.png")
	clubFourImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_club_four.png")
	clubFiveImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_club_five.png")
	clubSixImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_club_six.png")
	clubSevenImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_club_seven.png")
	clubEightImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_club_eight.png")
	clubNineImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_club_nine.png")
	clubTenImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_club_ten.png")
	clubJackImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_club_jack.png")
	clubQueenImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_club_queen.png")
	clubKingImage, _, _ = ebitenutil.NewImageFromFile("assets/solitaire_club_king.png")

	// create the 52 playing cards of four suits from ace-king. No jokers needed.

	// spades
	cardDeck = append(cardDeck, *newCard(Ace, Spade, *spadeAceImage))
	cardDeck = append(cardDeck, *newCard(Two, Spade, *spadeTwoImage))
	cardDeck = append(cardDeck, *newCard(Three, Spade, *spadeThreeImage))
	cardDeck = append(cardDeck, *newCard(Four, Spade, *spadeFourImage))
	cardDeck = append(cardDeck, *newCard(Five, Spade, *spadeFiveImage))
	cardDeck = append(cardDeck, *newCard(Six, Spade, *spadeSixImage))
	cardDeck = append(cardDeck, *newCard(Seven, Spade, *spadeSevenImage))
	cardDeck = append(cardDeck, *newCard(Eight, Spade, *spadeEightImage))
	cardDeck = append(cardDeck, *newCard(Nine, Spade, *spadeNineImage))
	cardDeck = append(cardDeck, *newCard(Ten, Spade, *spadeTenImage))
	cardDeck = append(cardDeck, *newCard(Jack, Spade, *spadeJackImage))
	cardDeck = append(cardDeck, *newCard(Queen, Spade, *spadeQueenImage))
	cardDeck = append(cardDeck, *newCard(King, Spade, *spadeKingImage))

	// clubs
	cardDeck = append(cardDeck, *newCard(Ace, Club, *clubAceImage))
	cardDeck = append(cardDeck, *newCard(Two, Club, *clubTwoImage))
	cardDeck = append(cardDeck, *newCard(Three, Club, *clubThreeImage))
	cardDeck = append(cardDeck, *newCard(Four, Club, *clubFourImage))
	cardDeck = append(cardDeck, *newCard(Five, Club, *clubFiveImage))
	cardDeck = append(cardDeck, *newCard(Six, Club, *clubSixImage))
	cardDeck = append(cardDeck, *newCard(Seven, Club, *clubSevenImage))
	cardDeck = append(cardDeck, *newCard(Eight, Club, *clubEightImage))
	cardDeck = append(cardDeck, *newCard(Nine, Club, *clubNineImage))
	cardDeck = append(cardDeck, *newCard(Ten, Club, *clubTenImage))
	cardDeck = append(cardDeck, *newCard(Jack, Club, *clubJackImage))
	cardDeck = append(cardDeck, *newCard(Queen, Club, *clubQueenImage))
	cardDeck = append(cardDeck, *newCard(King, Club, *clubKingImage))

	// hearts
	cardDeck = append(cardDeck, *newCard(Ace, Heart, *heartAceImage))
	cardDeck = append(cardDeck, *newCard(Two, Heart, *heartTwoImage))
	cardDeck = append(cardDeck, *newCard(Three, Heart, *heartThreeImage))
	cardDeck = append(cardDeck, *newCard(Four, Heart, *heartFourImage))
	cardDeck = append(cardDeck, *newCard(Five, Heart, *heartFiveImage))
	cardDeck = append(cardDeck, *newCard(Six, Heart, *heartSixImage))
	cardDeck = append(cardDeck, *newCard(Seven, Heart, *heartSevenImage))
	cardDeck = append(cardDeck, *newCard(Eight, Heart, *heartEightImage))
	cardDeck = append(cardDeck, *newCard(Nine, Heart, *heartNineImage))
	cardDeck = append(cardDeck, *newCard(Ten, Heart, *heartTenImage))
	cardDeck = append(cardDeck, *newCard(Jack, Heart, *heartJackImage))
	cardDeck = append(cardDeck, *newCard(Queen, Heart, *heartQueenImage))
	cardDeck = append(cardDeck, *newCard(King, Heart, *heartKingImage))

	// diamonds
	cardDeck = append(cardDeck, *newCard(Ace, Diamond, *diamondAceImage))
	cardDeck = append(cardDeck, *newCard(Two, Diamond, *diamondTwoImage))
	cardDeck = append(cardDeck, *newCard(Three, Diamond, *diamondThreeImage))
	cardDeck = append(cardDeck, *newCard(Four, Diamond, *diamondFourImage))
	cardDeck = append(cardDeck, *newCard(Five, Diamond, *diamondFiveImage))
	cardDeck = append(cardDeck, *newCard(Six, Diamond, *diamondSixImage))
	cardDeck = append(cardDeck, *newCard(Seven, Diamond, *diamondSevenImage))
	cardDeck = append(cardDeck, *newCard(Eight, Diamond, *diamondEightImage))
	cardDeck = append(cardDeck, *newCard(Nine, Diamond, *diamondNineImage))
	cardDeck = append(cardDeck, *newCard(Ten, Diamond, *diamondTenImage))
	cardDeck = append(cardDeck, *newCard(Jack, Diamond, *diamondJackImage))
	cardDeck = append(cardDeck, *newCard(Queen, Diamond, *diamondQueenImage))
	cardDeck = append(cardDeck, *newCard(King, Diamond, *diamondKingImage))

	// set start positions
	var firstRowHeight = cardBorder
	var secondRowHeight = (cardBorder * 2) + cardHeight

	var stockPileX = cardBorder * 2
	var stockPileY = firstRowHeight

	var wastePileX = (cardBorder * 3) + cardWidth
	var wastePileY = firstRowHeight

	var foundationPileOneX = (cardBorder * 5) + (cardWidth * 3)
	var foundationPileTwoX = (cardBorder * 6) + (cardWidth * 4)
	var foundationPileThreeX = (cardBorder * 7) + (cardWidth * 5)
	var foundationPileFourX = (cardBorder * 8) + (cardWidth * 6)

	var foundationPileOneY = firstRowHeight
	var foundationPileTwoY = firstRowHeight
	var foundationPileThreeY = firstRowHeight
	var foundationPileFourY = firstRowHeight

	var tableauPileOneX = cardBorder * 2
	var tableauPileTwoX = (cardBorder * 3) + cardWidth
	var tableauPileThreeX = (cardBorder * 4) + (cardWidth * 2)
	var tableauPileFourX = (cardBorder * 5) + (cardWidth * 3)
	var tableauPileFiveX = (cardBorder * 6) + (cardWidth * 4)
	var tableauPileSixX = (cardBorder * 7) + (cardWidth * 5)
	var tableauPileSevenX = (cardBorder * 8) + (cardWidth * 6)

	var tableauPileOneY = secondRowHeight
	var tableauPileTwoY = secondRowHeight
	var tableauPileThreeY = secondRowHeight
	var tableauPileFourY = secondRowHeight
	var tableauPileFiveY = secondRowHeight
	var tableauPileSixY = secondRowHeight
	var tableauPileSevenY = secondRowHeight

	// stock and waste setup
	stockPile = newStock(stockPileX, stockPileY)
	wastePile = newWaste(wastePileX, wastePileY)

	// foundation setup
	var foundationOne = newFoundation(foundationPileOneX, foundationPileOneY)
	var foundationTwo = newFoundation(foundationPileTwoX, foundationPileTwoY)
	var foundationThree = newFoundation(foundationPileThreeX, foundationPileThreeY)
	var foundationFour = newFoundation(foundationPileFourX, foundationPileFourY)
	foundations = [4]*foundation{foundationOne, foundationTwo, foundationThree, foundationFour}

	// tableau setup
	var tableauOne = newTableau(tableauPileOneX, tableauPileOneY)
	var tableauTwo = newTableau(tableauPileTwoX, tableauPileTwoY)
	var tableauThree = newTableau(tableauPileThreeX, tableauPileThreeY)
	var tableauFour = newTableau(tableauPileFourX, tableauPileFourY)
	var tableauFive = newTableau(tableauPileFiveX, tableauPileFiveY)
	var tableauSix = newTableau(tableauPileSixX, tableauPileSixY)
	var tableauSeven = newTableau(tableauPileSevenX, tableauPileSevenY)
	tableaus = [7]*tableau{tableauOne, tableauTwo, tableauThree, tableauFour, tableauFive, tableauSix, tableauSeven}

	// distribute cards NOTE: ported this code directly but it should be refactored to be more golang like
	for len(cardDeck) > 0 {
		var randomCardIndex = rand.Intn(len(cardDeck))
		var card = cardDeck[randomCardIndex]
		for ti, tableau := range tableaus {
			for len(tableau.cardPile.cards) <= ti {
				card.hidden = true
				card.x = tableau.cardPile.x
				card.y = tableau.cardPile.y
				tableau.cardPile.addCard(card)
				cardDeck = append(cardDeck[:randomCardIndex], cardDeck[randomCardIndex+1:]...)
				randomCardIndex = rand.Intn(len(cardDeck))
				card = cardDeck[randomCardIndex]
			}
		}
		card.x = stockPile.cardPile.x
		card.y = stockPile.cardPile.y
		stockPile.cardPile.addCard(card)
		cardDeck = append(cardDeck[:randomCardIndex], cardDeck[randomCardIndex+1:]...)
		if len(cardDeck) == 0 {
			break
		}
		randomCardIndex = rand.Intn(len(cardDeck))
		card = cardDeck[randomCardIndex]
	}

	for ti, tableau := range tableaus {
		for ci := range tableau.cardPile.cards {
			tableaus[ti].cardPile.cards[ci].active = false
		}
	}

	for _, tableau := range tableaus {
		tableau.cardPile.cards[len(tableau.cardPile.cards)-1].hidden = false
		tableau.cardPile.cards[len(tableau.cardPile.cards)-1].active = true
	}
}

type Game struct{}

func (g *Game) Draw(screen *ebiten.Image) {
	// background color
	screen.Fill(green)

	stockPile.cardPile.draw(screen)
	wastePile.cardPile.draw(screen)

	// draw foundations
	for _, f := range foundations {
		f.cardPile.draw(screen)
	}

	for _, t := range tableaus {
		t.cardPile.draw(screen)
	}

	// redraw selected card last so it always appears on top of other piles
	if selectedCardPile != nil {
		for _, card := range selectedCardPile.cards {
			card.draw(screen)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

var currentTime time.Time
var priorTime time.Time
var wait time.Duration = 1_000_000 * time.Second
var inputDelay time.Duration = 1000 * time.Millisecond
var WasteCardSelected = false
var priorMouseX int
var priorMouseY int
var selectedCardPile *cardPile = nil
var foundationSplitFrom *foundation
var foundationMergeTo *foundation
var tableauSplitFrom *tableau
var tableauMergeTo *tableau

func (g *Game) Update() error {
	currentTime = time.Now()
	if wait < inputDelay {
		wait = wait + currentTime.Sub(priorTime)
		return nil
	}

	var x, y = ebiten.CursorPosition()
	var mouseClicked = ebiten.IsMouseButtonPressed(leftMouseButton)

	if mouseClicked {
		// draw cards with a delay so a single click doesn't draw multiple
		if wait >= inputDelay && stockPile.cardPile.containsPoint(x, y) && selectedCardPile == nil {
			wait = 0
			var card, _ = stockPile.drawCard()
			if card != nil {
				card.hidden = false
				wastePile.cardPile.addCard(*card)
			} else {
				// redeal waste pile since stock is empty
				var cardToAdd, _ = wastePile.cardPile.removeCard()
				var emptyWastePile error
				for emptyWastePile == nil && cardToAdd != nil {
					cardToAdd.hidden = true
					stockPile.cardPile.addCard(*cardToAdd)
					cardToAdd, emptyWastePile = wastePile.cardPile.removeCard()
				}
			}

		}

		// select cards to move, no delay needed
		if selectedCardPile == nil {
			var wasteCard, _ = wastePile.cardPile.peek()
			if wasteCard != nil && wasteCard.containsPoint(x, y) {
				WasteCardSelected = true
				selectedCardPile = wastePile.cardPile.splitCardPile(*wasteCard)
			}

			for _, foundation := range foundations {
				for _, c := range foundation.cardPile.cards {
					if c.containsPoint(x, y) {
						selectedCardPile = foundation.cardPile.splitCardPile(c)
						foundationSplitFrom = foundation
						break
					}
				}
			}

			for _, tableau := range tableaus {
				for _, c := range tableau.cardPile.cards {
					if c.containsPoint(x, y) {
						selectedCardPile = tableau.cardPile.splitCardPile(c)
						tableauSplitFrom = tableau
						break
					}
				}
			}
		} else {
			// drag card already selected this update around the screen
			selectedCardPile = selectedCardPile.move(x-priorMouseX, y-priorMouseY)
		}
		// card dragging finished
	} else if !mouseClicked && selectedCardPile != nil {
		var selectedXMid = selectedCardPile.xMid()
		var selectedYMid = selectedCardPile.yMid()

		for _, foundation := range foundations {
			if foundation.cardPile.containsPoint(selectedXMid, selectedYMid) && foundationSplitFrom != foundation {
				foundationMergeTo = foundation
			}
		}

		for _, tableau := range tableaus {
			if tableau.cardPile.containsPoint(selectedXMid, selectedYMid) && tableauSplitFrom != tableau {
				tableauMergeTo = tableau
			}
		}

		for i := range selectedCardPile.cards {
			selectedCardPile.cards[i].x = selectedCardPile.cards[i].unselectedX
			selectedCardPile.cards[i].y = selectedCardPile.cards[i].unselectedY
		}

		var movedCardPile = false

		// card requested to be moved to valid position, so move it and update game state
		if foundationMergeTo != nil {
			movedCardPile = foundationMergeTo.addCardPileToFoundation(*selectedCardPile)
		}
		if tableauMergeTo != nil {
			movedCardPile = tableauMergeTo.addCardPileToTableau(*selectedCardPile)
		}
		// card requested to bne moved to invalid position, so move return it to original position and retain game state
		if WasteCardSelected && !movedCardPile {
			wastePile.cardPile.addCardPile(*selectedCardPile)
		}
		if foundationSplitFrom != nil && !movedCardPile {
			foundationSplitFrom.cardPile.addCardPile(*selectedCardPile)
		}
		if tableauSplitFrom != nil && !movedCardPile {
			tableauSplitFrom.cardPile.addCardPile(*selectedCardPile)
		}

		// ensure top foundation card is always active and other foundation cards are inactive
		for fi, foundation := range foundations {
			for ci := range foundation.cardPile.cards {
				foundations[fi].cardPile.cards[ci].active = false
			}
			var c, _ = foundation.cardPile.peek()
			if c != nil {
				c.active = true
			}
		}

		// activate newly uncovered tableau cards
		for _, tableau := range tableaus {
			var c, _ = tableau.cardPile.peek()
			if c != nil {
				c.hidden = false
				c.active = true
			}
		}

		selectedCardPile = nil
		WasteCardSelected = false
		foundationMergeTo = nil
		foundationSplitFrom = nil
		tableauMergeTo = nil
		tableauSplitFrom = nil
	}

	priorMouseX = x
	priorMouseY = y
	priorTime = currentTime

	return nil
}
