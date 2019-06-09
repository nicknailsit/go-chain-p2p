package repo

import (
	"encoding/json"
	"fmt"
	"github.com/BluntSporks/readability"
	"github.com/mb-14/gomarkov"
	"sync"

	"io/ioutil"
	"strings"
)



//not in use yet SMOG text analysis ?
func calculateDifficultyOfTarget(target string) (score float64) {

	score = read.Smog(target)
	return
}



type QuoteTrainingSet []Quote


type Quote struct {
	quoteText string
	quoteAuthor string
}

func GetTrainingSet() map[int]string {
	data, err := ioutil.ReadFile("./repo/quotes.json")
	if err != nil {
		panic(err)
	}

	mapper := make(map[int]string, len(data))

	splitted := strings.Split(string(data), "\n")

	for i := 0; i < len(splitted); i++ {

		mapper[i] = string(splitted[i])

	}
	return mapper
}


//training the markov chain to be able to generate sentences for pow
func TrainForPow() string {


	chain, _ := buildModel()
	saveModel(chain)

	chain, _ = loadModel()
	return generatePow(chain)


}

func buildModel() (*gomarkov.Chain, error) {

	quotes := GetTrainingSet()

	chain := gomarkov.NewChain(1)
	var wg sync.WaitGroup
	wg.Add(len(quotes))

	for _, quote := range quotes {
		go func() {

			defer wg.Done()

			chain.Add(strings.Split(quote, " "))

		}()
	}
	wg.Wait()

	return chain, nil

}

func saveModel(chain *gomarkov.Chain) {
	jsonObj, _ := json.Marshal(chain)
	err := ioutil.WriteFile("model.json", jsonObj, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func loadModel() (*gomarkov.Chain, error) {
	var chain gomarkov.Chain
	data, err := ioutil.ReadFile("model.json")
	if err != nil {
		return &chain, err
	}
	err = json.Unmarshal(data, &chain)
	if err != nil {
		return &chain, err
	}
	return &chain, nil
}

func generatePow(chain *gomarkov.Chain) string {
	tokens := []string{gomarkov.StartToken}
	for tokens[len(tokens)-1] != gomarkov.EndToken {
		next, _ := chain.Generate(tokens[(len(tokens) - 1):])
		tokens = append(tokens, next)
	}
	return fmt.Sprintf(strings.Join(tokens[1:len(tokens)-1], " "))
}