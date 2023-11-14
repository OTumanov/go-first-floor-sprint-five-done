package main

import (
	"fmt"
	"math"
	"time"
)

// Consts for calculates
const (
	MInKm      = 1000 // Meters in kilometers
	MinInHours = 60
	LenStep    = 0.65
)

// Training ...
type Training struct {
	TrainingType string
	Action       int
	LenStep      float64
	Duration     time.Duration
	Weight       float64
}

func (t Training) distance() float64 {
	return float64(t.Action) * t.LenStep / MInKm
}

func (t Training) meanSpeed() float64 {
	distance := t.distance()
	return distance / t.Duration.Hours()
}

// Calories ...
func (t Training) Calories() float64 {
	return 0
}

// TrainingInfo ...
func (t Training) TrainingInfo() InfoMessage {
	distance := t.distance()
	speed := t.meanSpeed()
	calories := t.Calories()
	return InfoMessage{
		TrainingType: t.TrainingType,
		Duration:     t.Duration,
		Distance:     distance,
		Speed:        speed,
		Calories:     calories,
	}
}

// InfoMessage ...
type InfoMessage struct {
	TrainingType string
	Duration     time.Duration
	Distance     float64
	Speed        float64
	Calories     float64
}

// Message returns info for training
func (i InfoMessage) Message() string {
	return fmt.Sprintf("Тип тренировки: %s; Длительность: %v мин; Дистанция: %.2f км.; Ср. скорость: %.2f; Потрачено ккал: %.2f", i.TrainingType, i.Duration.Minutes(), i.Distance, i.Speed, i.Calories)
}

// CaloriesCalculator ...
type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

// Consts ...
const (
	CaloriesMeanSpeedMultiplier = 18
	CaloriesMeanSpeedShift      = 1.79
)

// Running ...
type Running struct {
	Training
}

// Calories ...
func (r Running) Calories() float64 {
	return (CaloriesMeanSpeedMultiplier*r.meanSpeed() + CaloriesMeanSpeedShift) * r.Weight / MInKm * r.Duration.Hours() * MinInHours
}

// TrainingInfo ...
func (r Running) TrainingInfo() InfoMessage {
	return r.Training.TrainingInfo()
}

// Consts ...
const (
	CaloriesWeightMultiplier      = 0.035
	CaloriesSpeedHeightMultiplier = 0.029
	KmHInMsec                     = 0.278
	CmInM                         = 100
)

// Walking ...
type Walking struct {
	Training
	Height float64
}

// Calories ...
func (w Walking) Calories() float64 {
	return (CaloriesWeightMultiplier*w.Weight + (math.Pow(w.meanSpeed()*KmHInMsec, 2)/(w.Height/CmInM))*CaloriesSpeedHeightMultiplier*w.Weight) * w.Duration.Hours() * MinInHours
}

// TrainingInfo ...
func (w Walking) TrainingInfo() InfoMessage {
	return w.Training.TrainingInfo()
}

// Consts ...
const (
	SwimmingLenStep                  = 1.38
	SwimmingCaloriesMeanSpeedShift   = 1.1
	SwimmingCaloriesWeightMultiplier = 2
)

// Swimming ...
type Swimming struct {
	Training
	LengthPool int
	CountPool  int
}

func (s Swimming) meanSpeed() float64 {
	return float64(s.LengthPool) * float64(s.CountPool) / MInKm / s.Duration.Hours()
}

// Calories ...
func (s Swimming) Calories() float64 {
	return (s.meanSpeed() + SwimmingCaloriesMeanSpeedShift) * SwimmingCaloriesWeightMultiplier * s.Weight * s.Duration.Hours()
}

// TrainingInfo ...
func (s Swimming) TrainingInfo() InfoMessage {
	return s.Training.TrainingInfo()
}

// ReadData ...
func ReadData(training CaloriesCalculator) string {
	calories := training.Calories()

	info := training.TrainingInfo()
	info.Calories = calories

	return info.Message()
}

func main() {

	swimming := Swimming{
		Training: Training{
			TrainingType: "Плавание",
			Action:       2000,
			LenStep:      SwimmingLenStep,
			Duration:     time.Minute * 90,
			Weight:       85,
		},
		LengthPool: 50,
		CountPool:  5,
	}

	fmt.Println(ReadData(swimming))

	walking := Walking{
		Training: Training{
			TrainingType: "Ходьба",
			Action:       20000,
			LenStep:      LenStep,
			Duration:     time.Hour*3 + time.Minute*45,
			Weight:       85,
		},
		Height: 185,
	}

	fmt.Println(ReadData(walking))

	running := Running{
		Training: Training{
			TrainingType: "Бег",
			Action:       5000,
			LenStep:      LenStep,
			Duration:     time.Minute * 30,
			Weight:       85,
		},
	}

	fmt.Println(ReadData(running))

}
