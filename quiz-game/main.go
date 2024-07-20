package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Questions struct {
	question string
	expectedAnswer string
	userAnswer string
	isCorrect bool
	wasAnswered bool
}


func buildQuestions(questionsArr [][]string, shuffle bool) []Questions {

	questionsObj := make([]Questions, len(questionsArr))

	if shuffle {
		rand.Shuffle(len(questionsArr), func(i, j int) { questionsArr[i], questionsArr[j] = questionsArr[j], questionsArr[i] })
	}

	for i, q := range questionsArr {
		questionsObj[i] = Questions{
			question: q[0],
			expectedAnswer: q[1],
		}
	}

	return questionsObj
}

func compileFlags() (string, bool, int) {
	fileName := flag.String("fileName", "problems.csv", "Problems file name, should contains the .csv. Default: problems.csv")
	shuffle := flag.Bool("shuffle", false, "Shuffle all questions. Default: true.")
	timer := flag.Int("timer", 30, "Set a timer to accomplish all questions in seconds. Default: 30.")
	flag.Parse()

	return *fileName, *shuffle, *timer
}

func compareQuizResponses(questionsObj []Questions, lastResponse int) int {

	correctAns := 0;

	for i := range questionsObj {
		expected := strings.TrimSpace(strings.ToLower(questionsObj[i].expectedAnswer))
		ans := strings.TrimSpace(strings.ToLower(questionsObj[i].userAnswer))

		if strings.Compare(expected, ans) == 0 {
			questionsObj[i].isCorrect = true
			correctAns = correctAns + 1
		}

		if i < lastResponse {
			questionsObj[i].wasAnswered = true
		}
	}

	return correctAns
}


func waitByUser(reader *bufio.Reader, totalQuestions int){
	fmt.Printf("This quiz has %d questions. Press ENTER when ready", totalQuestions)
	reader.ReadString('\n')
}

func runQuiz(questionsObj []Questions, timer int) int {
	totalOfQuestions := len(questionsObj)

	reader := bufio.NewReader(os.Stdin)
	waitByUser(reader, totalOfQuestions)


	timeout := time.After(time.Duration(timer)*time.Second)

	for i := range questionsObj {

		fmt.Printf("[%d/%d] %v: ", i+1, totalOfQuestions, questionsObj[i].question)
		ansChannel := make(chan string)


		go func() {
			userAnswer, _ := reader.ReadString('\n')
			userAnswer = strings.TrimSuffix(userAnswer, "\n")
			ansChannel <- userAnswer
		}()

		select {
			case userAnswer := <- ansChannel:
				questionsObj[i].userAnswer = userAnswer
			case <- timeout:
				fmt.Println("\nTime's up! Ending the quiz.")
				return i
		}
	}

	return totalOfQuestions
}

func printResults(questionsObj []Questions, totalOfCorrectAns int){

	fmt.Printf("Total of correct answer %d of %d\n", totalOfCorrectAns, len(questionsObj))
	var userAnswer string

	for _ , q := range questionsObj {
		if !q.isCorrect {

			if q.wasAnswered {
				userAnswer = q.userAnswer
			}else{
				userAnswer = "(Has no answer)"
			}
			fmt.Printf("Question: %v, Expected Answer: %v, Given Answer: %v\n", q.question, q.expectedAnswer, userAnswer)
		}
	}
}

func main(){


	fileName, shuffle, timer := compileFlags()
	
	
	file, err := os.Open(fileName)
	if err != nil{
		log.Fatalf("Error on opening file: %v", fileName)
	}
	defer file.Close()


	csvReader := csv.NewReader(file)
	questionsArr, err := csvReader.ReadAll()

	if err != nil {
		log.Fatalf("Error on reading CSV")
	}

	questionsObj := buildQuestions(questionsArr, shuffle)
	lastResponse := runQuiz(questionsObj, timer)
	totalOfCorrectAns := compareQuizResponses(questionsObj, lastResponse)
	printResults(questionsObj, totalOfCorrectAns)

}