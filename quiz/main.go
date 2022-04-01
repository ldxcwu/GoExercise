package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

/* quiz game: 在一个CSV文件中有着question,answer的数据
读取CSV文件抛出问题并接受输入，检测answer是否正确
给定答题时间，最后得出正确解答的问题数量 */
func main() {
	//1. 读取命令行参数（CSV文件、时间限制limit）并解析
	csvFilename := flag.String("csv", "problems.csv", "A simple CSV file in the format 'qustion,answer'")
	limit := flag.Int("limit", 30, "time limit")

	//2. 打开CSV文件
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open CSV file: %s\n", err))
	}

	//3. 使用csv包提供的工具读取csv文件：创建一个Reader，并进行读取
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to read CSV file: %s\n", err))
	}

	//3. 解析读取的数据，
	problems := parseLines(lines)

	//4. 开启一个接受用户输入答案的channel，以及一个时间限制的channel
	timer := time.NewTimer(time.Duration(*limit) * time.Second)
	answerCh := make(chan string)

	count := 0

	for i, p := range problems {
		fmt.Printf("Question #%d: %s = ", i+1, p.q)
		go func() {
			var ans string
			fmt.Scanf("%s", &ans)
			answerCh <- ans
		}()
		select {
		case <-timer.C:
			fmt.Printf("\n You score %d out of %d\n", count, len(problems))
			close(answerCh)
			return
		case ans := <-answerCh:
			if strings.TrimSpace(ans) == p.a {
				count++
			}
		}
	}
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: line[1],
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
