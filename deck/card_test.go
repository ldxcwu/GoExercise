package deck

import (
	"fmt"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Rank: Ace, Suit: Heart})
	fmt.Println(Card{Rank: Two, Suit: Spade})
	fmt.Println(Card{Rank: Jack, Suit: Club})
	fmt.Println(Card{Rank: Queen, Suit: Diamond})
	fmt.Println(Card{Suit: Joker})

	//Output:
	//Ace of Hearts
	//Two of Spades
	//Jack of Clubs
	//Queen of Diamonds
	//Joker
}

func TestNew(t *testing.T) {
	cards := New()
	//13 Rank * 4 Suit
	if len(cards) != 13*4 {
		t.Errorf("Wrong number of cards in a new deck.")
	}
}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)
	first := Card{Rank: Ace, Suit: Spade}
	last := Card{Rank: King, Suit: Heart}
	if cards[0] != first || cards[len(cards)-1] != last {
		t.Errorf("Expected Ace of Spades as first card. Received:%v\n", cards[0])
		t.Errorf("Expected King of Heart as last card. Received:%v\n", cards[len(cards)-1])
	}
}

func TestSort(t *testing.T) {
	cards := New(Sort(Less))
	first := Card{Rank: Ace, Suit: Spade}
	last := Card{Rank: King, Suit: Heart}
	if cards[0] != first || cards[len(cards)-1] != last {
		t.Errorf("Expected Ace of Spades as first card. Received:%v\n", cards[0])
		t.Errorf("Expected King of Heart as last card. Received:%v\n", cards[len(cards)-1])
	}
}

func TestJokers(t *testing.T) {
	cards := New(Jokers(3))
	cnt := 0
	for _, c := range cards {
		if c.Suit == Joker {
			cnt++
		}
	}
	if cnt != 3 {
		t.Error("Expected 3 Jokers. Received:", cnt)
	}
}
