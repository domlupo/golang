package main

import (
	_ "image/png"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const windowMultipler = 1

// 16:9 aspect ratio
const windowWidth = int(1280 * windowMultipler)
const windowHeight = int(720 * windowMultipler)

// standard card ratio is 5:7
const cardMultipler = windowMultipler
const cardHeight = 140 * cardMultipler
const cardWidth = 100 * cardMultipler
const cardBorder = 10 * cardMultipler
const cardNumberHeight = 30 * cardMultipler

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

var img *ebiten.Image
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

var err error
var cp cardPile
var cardDeck []card

var foundations [4]*foundation
var tableaus [7]*tableau

func init() {
	// load card images
	cardBackImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_card_back.png")
	cardPlaceholder, _, err = ebitenutil.NewImageFromFile("assets/solitaire_card_placeholder.png")

	// load diamond images
	diamondAceImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_diamond_one.png")
	diamondTwoImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_diamond_two.png")
	diamondThreeImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_diamond_three.png")
	diamondFourImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_diamond_four.png")
	diamondFiveImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_diamond_five.png")
	diamondSixImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_diamond_six.png")
	diamondSevenImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_diamond_seven.png")
	diamondEightImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_diamond_eight.png")
	diamondNineImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_diamond_nine.png")
	diamondTenImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_diamond_ten.png")
	diamondJackImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_diamond_jack.png")
	diamondQueenImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_diamond_queen.png")
	diamondKingImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_diamond_king.png")

	// load heart images
	heartAceImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_heart_one.png")
	heartTwoImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_heart_two.png")
	heartThreeImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_heart_three.png")
	heartFourImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_heart_four.png")
	heartFiveImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_heart_five.png")
	heartSixImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_heart_six.png")
	heartSevenImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_heart_seven.png")
	heartEightImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_heart_eight.png")
	heartNineImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_heart_nine.png")
	heartTenImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_heart_ten.png")
	heartJackImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_heart_jack.png")
	heartQueenImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_heart_queen.png")
	heartKingImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_heart_king.png")

	// load spade images
	spadeAceImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_spade_one.png")
	spadeTwoImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_spade_two.png")
	spadeThreeImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_spade_three.png")
	spadeFourImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_spade_four.png")
	spadeFiveImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_spade_five.png")
	spadeSixImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_spade_six.png")
	spadeSevenImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_spade_seven.png")
	spadeEightImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_spade_eight.png")
	spadeNineImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_spade_nine.png")
	spadeTenImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_spade_ten.png")
	spadeJackImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_spade_jack.png")
	spadeQueenImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_spade_queen.png")
	spadeKingImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_spade_king.png")

	// load club images
	clubAceImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_club_one.png")
	clubTwoImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_club_two.png")
	clubThreeImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_club_three.png")
	clubFourImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_club_four.png")
	clubFiveImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_club_five.png")
	clubSixImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_club_six.png")
	clubSevenImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_club_seven.png")
	clubEightImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_club_eight.png")
	clubNineImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_club_nine.png")
	clubTenImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_club_ten.png")
	clubJackImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_club_jack.png")
	clubQueenImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_club_queen.png")
	clubKingImage, _, err = ebitenutil.NewImageFromFile("assets/solitaire_club_king.png")

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

	//var wastePileX = (cardBorder * 3) + cardWidth
	//var wastePileY = firstRowHeight

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
	var stockPile = newStock(stockPileX, stockPileY)
	//var wastePile = newStock(wastePileX, wastePileY)

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

	// distribute cards
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
	for _, tableau := range tableaus {
		tableau.cardPile.cards[len(tableau.cardPile.cards)-1].hidden = false
	}
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// background color

	// draw foundations
	for _, f := range foundations {
		f.cardPile.draw(screen)
	}

	for _, t := range tableaus {
		t.cardPile.draw(screen)
	}

	// redraw selected card last so it always appears on top of other cards

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
