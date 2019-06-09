package chain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/minio/blake2b-simd"
	"math/rand"
	"swaggp2p/repo"
	"sync"
	"time"
)

type DNA struct {
	sync.Mutex
	ChainDNA
}

var dna *DNA
var testTarget []byte
var params *DNAParams
var O Organism
var A *Ancestor

func init() {

	dna = &DNA{}
	testTarget = []byte("test target")
	params = O.defaultParams(testTarget)
	A = &Ancestor{}


}

func setTarget() string {
	target := repo.TrainForPow()
	buf := bytes.NewBuffer(make([]byte, 255))
	i := 0
	if len(target) > 256 {
		tBytes := []byte(target)
		for i < 255 {
			buf.WriteByte(tBytes[i])

			i++
		}

		target = string(buf.Bytes())

	}
	return target
}


func (d Organism) defaultParams(target []byte) (params *DNAParams) {

	params = &DNAParams{
		mutationRate: 0.0005,
		populationSize: 5000,
		target: target,
	}

	return

}


func createOrganism(target []byte) (organism Organism) {



	ba := make([]byte, len(target))
	for i := 0; i < len(target); i++ {
		ba[i] = byte(rand.Intn(95) + 32)
	}

	organism = Organism{
		DNA: ba,
		Fitness: 0,
	}

	organism.calculateFitness(target)
	return

}


func createPopulation(target []byte) (population []Organism){


	population = make([]Organism, params.populationSize)
	for i := 0; i < int(params.populationSize); i++ {
		population[i] = createOrganism(target)
	}
	return
}


func(d *Organism) calculateFitness(target []byte) {

	score :=0
	for i := 0; i < len(d.DNA); i++ {
		if d.DNA[i] == target[i] {
			score++
		}
	}
	d.Fitness = float64(score) / float64(len(d.DNA))
	return

}

func createGenePool(population []Organism, target []byte, maxFitness float64) (pool []Organism) {


	pool = make([]Organism, 0)
	for i := 0; i < len(population); i++ {
		population[i].calculateFitness(target)
		num := int((population[i].Fitness / maxFitness) * 100)
		for n := 0; n < num; n++ {
			pool = append(pool, population[i])
		}
	}

	A.genePool = pool

	return
}

func naturalSelection(pool []Organism, population []Organism, target []byte) []Organism {



	next := make([]Organism, len(population))
	for i := 0; i < len(population); i++ {
		r1, r2 := rand.Intn(len(pool)), rand.Intn(len(pool))
		a := pool[r1]
		b := pool[r2]

		pa, _ := json.Marshal(a)
		pb, _ := json.Marshal(b)
		next[i].ParentA = pa
		next[i].ParentB = pb

		child := crossover(a, b)
		child.mutate()
		child.calculateFitness(target)
		next[i] = child


	}
	return next

}

func crossover(d1 Organism, d2 Organism) Organism {
	child := Organism{
		DNA: make([]byte, len(d1.DNA)),
		Fitness: 0,
	}

	mid := rand.Intn(len(d1.DNA))
	for i := 0; i < len(d1.DNA); i++ {
		if i > mid {
			child.DNA[i] = d1.DNA[i]
		} else {
			child.DNA[i] = d2.DNA[i]
		}
	}
	return child
}

func (d *Organism) mutate() {
	for i := 0; i < len(d.DNA); i++ {
		if rand.Float64() < params.mutationRate {
			d.DNA[i] = byte(rand.Intn(95) + 32)
		}
	}
}

func getBest(population []Organism) Organism {
	best := 0.0
	index := 0
	for i := 0; i < len(population); i++ {
		if population[i].Fitness > best {
			index = i
			best = population[i].Fitness
		}
	}
	return population[index]
}



func GetDNA() (solution *DNASolution) {

	target := []byte(setTarget())
	start := time.Now()
	rand.Seed(time.Now().UTC().UnixNano())
	population := createPopulation(target)

	found := false
	generation := 0

	for !found {
		generation++
		bestOrganism := getBest(population)
		fmt.Printf("\r generation: %d | %s | fitness: %2f", generation, string(bestOrganism.DNA), bestOrganism.Fitness)

		if bytes.Compare(bestOrganism.DNA, target) == 0 {
			found = true

			oBytes, _ := json.Marshal(bestOrganism)
			pop, _ := json.Marshal(population)
			h := blake2b.New256()
			h.Reset()
			h.Sum(pop)
			h.Sum(oBytes)

			Dh := blake2b.New512()
			Dh.Write(bestOrganism.DNA)

			solution = &DNASolution{
				a:h.Sum(nil),
				o:oBytes,
				D:Dh.Sum(nil),
				t: target,
			}

		} else {
			maxFitness := bestOrganism.Fitness
			pool := createGenePool(population, target, maxFitness)
			population = naturalSelection(pool, population, target)
		}


	}

	elapsed := time.Since(start)
	fmt.Printf("\nTime taken: %s\n", elapsed)


	return
}