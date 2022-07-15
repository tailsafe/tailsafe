package __0

import "time"

type StepMonitoring struct {
	stageNumber  int
	startingTime time.Time
	endingTime   time.Time
}

func NewStepMonitoring(stageNumber int) *StepMonitoring {
	return &StepMonitoring{
		stageNumber:  stageNumber,
		startingTime: time.Now(),
	}
}

func (s *StepMonitoring) Reset() {
	s.startingTime = time.Now()
}

func (s *StepMonitoring) GetStage() int {
	return s.stageNumber
}

func (s *StepMonitoring) GetStageDuration() time.Duration {
	return s.endingTime.Sub(s.startingTime)
}

func (s *StepMonitoring) End() {
	s.endingTime = time.Now()
}
