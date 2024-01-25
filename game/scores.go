package game

const HIGHEST_SCORE = 9999999999

type PlayerScore struct {
	PlayerName string `json:"player"`
	Score      uint32 `json:"score"`
}

type ScoreBoard struct {
	FirstPosted bool           `json:"firstPosted"`
	Scores      [5]PlayerScore `json:"scores"`
}

func HighestScore(highScores [5]PlayerScore, newScore PlayerScore) int {
	highScorePosition := 0

	for index, highScore := range highScores {
		if highScore.Score < newScore.Score || highScore.Score == newScore.Score {
			highScorePosition = index
			break
		}
	}

	return highScorePosition
}
