package models

import (
	"reflect"
)

type SleepResults []SleepResult

type SleepResult struct {
	Night               PolarTimeForSleep   `json:"night,omitempty"`
	SleepScoreResult    SleepScoreResult    `json:"sleepScoreResult,omitempty"`
	SleepScoreBaselines SleepScoreBaselines `json:"sleepScoreBaselines,omitempty"`
}

type SleepScoreBaselines struct {
	ContinuityScoreAverage              float64 `json:"continuityScoreAverage,omitempty"`
	ContinuityScoreBaseline             float64 `json:"continuityScoreBaseline,omitempty"`
	EfficiencyPercentAverage            float64 `json:"efficiencyPercentAverage,omitempty"`
	EfficiencyScoreBaseline             float64 `json:"efficiencyScoreBaseline,omitempty"`
	GroupRefreshScoreBaseline           float64 `json:"groupRefreshScoreBaseline,omitempty"`
	GroupSolidityScoreBaseline          float64 `json:"groupSolidityScoreBaseline,omitempty"`
	LongInterruptionsAverageTimeMinutes int     `json:"longInterruptionsAverageTimeMinutes,omitempty"`
	LongInterruptionsScoreBaseline      float64 `json:"longInterruptionsScoreBaseline,omitempty"`
	N3PercentAverage                    float64 `json:"n3PercentAverage,omitempty"`
	N3ScoreBaseline                     float64 `json:"n3ScoreBaseline,omitempty"`
	RemPercentAverage                   float64 `json:"remPercentAverage,omitempty"`
	RemScoreBaseline                    float64 `json:"remScoreBaseline,omitempty"`
	SleepScoreBaseline                  float64 `json:"sleepScoreBaseline,omitempty"`
	SleepTimeAverageMinutes             int     `json:"sleepTimeAverageMinutes,omitempty"`
	SleepTimeScoreBaseline              float64 `json:"sleepTimeScoreBaseline,omitempty"`
}

type SleepScoreResult struct {
	ContinuityScore              float64 `json:"continuityScore,omitempty"`
	EfficiencyScore              float64 `json:"efficiencyScore,omitempty"`
	GroupDurationScore           float64 `json:"groupDurationScore,omitempty"`
	GroupRefreshScore            float64 `json:"groupRefreshScore,omitempty"`
	GroupSolidityScore           float64 `json:"groupSolidityScore,omitempty"`
	LongInterruptionsScore       float64 `json:"longInterruptionsScore,omitempty"`
	N3Score                      float64 `json:"n3Score,omitempty"`
	RemScore                     float64 `json:"remScore,omitempty"`
	ScoreRate                    int     `json:"scoreRate,omitempty"`
	SleepScore                   float64 `json:"sleepScore,omitempty"`
	SleepTimeOwnTargetScore      float64 `json:"sleepTimeOwnTargetScore,omitempty"`
	SleepTimeRecommendationScore float64 `json:"sleepTimeRecommendationScore,omitempty"`
}

func (s *SleepResult) ToInflux() map[string]interface{} {
	fields := map[string]interface{}{}
	baseLine := s.SleepScoreBaselines
	score := s.SleepScoreResult
	for suffix, SR := range map[string]interface{}{"baseline": baseLine, "score": score} {
		v := reflect.ValueOf(SR)
		t := reflect.TypeOf(SR)

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)
			fields[field.Name+"_"+suffix] = value.Interface()
		}
	}
	return fields
}
