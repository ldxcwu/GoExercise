# **Deck of Cards**

##  **Description**
Create a package that can used to build decks of cards, which may contains the following features:  
1. Sorting

---

## **Gains**  
---
###   **Use** ```golang.org/x/tools/cmd/stringer``` **package**
With stringer package you can generate String() func quickly.  

**Usage:**  
1. Install this package.(use ```go install``` intead of ```go get``` after 1.16version)
   ```go install golang.org/x/tools/cmd/stringer```
2. Define your own type or struct in go file starts with specific heading like:
   ```//go:generate stringer -type=Suit,Rank```
3. Run command like:  
   ```go generate```
   or
   ```stringer -type=Suit,Rank```
4. Then you will see a new file named ```xxx_string.go```
---
### **Use Go Example Func**
**Usage:**   
1. Create a new file named ```xxx_test.go```
2. Define a func which name has the ```Example``` prefix
3. Write your code in it.
4. Write comments which have the following format:  
   ``` // Output: ```   
   ``` // Some result ```
5. Run command ``` go test ```
6. The Go Program will test the ```ExampleXXX ``` func's result wheather equals the comments or not.
---
### **Use** ```type Options func(xxx) xxx ``` **to process our data when they are initializing.**
```go
//The cards may have a few options when they are initializing.
func New(opts ...func([]Card) []Card) []Card {
	var cards []Card
	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}

	for _, opt := range opts {
		opt(cards)
	}
	return cards
}

func DefaultSort(cards []Card) []Card {
	//sort.Slice(slice, func) takes in an slice and sort func
	//and sort the slice with the func.
	sort.Slice(cards, Less(cards))
	return cards
}

func Less(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return absRank(cards[i]) <= absRank(cards[j])
	}
}

func absRank(card Card) int {
	return int(card.Suit)*int(maxRank) + int(card.Rank)
}
```
```go
//Then we can use it like:
New(DefaultSort)
```
---
### **Combine Option and Customize**
```go 
CustomizeSort(less func) Option
```
```go
//提供一个方法用以让用户自定义排序规则，用户只需要提供比较的方法即可进行对应规则初始化
//Sort is an func that make user custom rank rule possibly,
//user only need to realize compare func then the data will be related initializing.
func Sort(less func(cards []Card) func(i, j int) bool) func(cards []Card) []Card {
	return func(cards []Card) []Card {
		sort.Slice(cards, less(cards))
		return cards
	}
}
```
```go
//Then we can use it like:
cards := New(Sort(Less))
//or
cards := New(Sort(func(...){...}))
```
---
### **Shuffle**
1. Create a rand seed
   ```go
   r := rand.New(rand.NewSource(time.Now().Unix()))
   ```
2. Use ```rand.Perm(n int)``` to build a shuffled slice
   ```go
   s := r.Perm(n)
   //s : A [0~n) slice that has been shuffled.
   ```
3. Swap our data with the slice above.  
```go
//Example:
func Shuffle(cards []Card) []Card {
	ret := make([]Card, len(cards))
	r := rand.New(rand.NewSource(time.Now().Unix()))
	//r.Perm takes in an integer n and return a slice [0~n) which has been shuffled.
	perm := r.Perm(len(cards))
	for i, j := range perm {
		ret[i] = cards[j]
	}
	return ret
}
```
---
