package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".envファイルが読み込めませんでした。")
	}

	workDuration := getEnvDuration("WORK_DURATION", 25*time.Minute)
	shortBreakDuration := getEnvDuration("SHORT_BREAK_DURATION", 5*time.Minute)
	longBreakDuration := getEnvDuration("LONG_BREAK_DURATION", 30*time.Minute)
	pomodoroCount := getEnvInt("POMODORO_COUNT", 4)

	workSoundFile := os.Getenv("WORK_SOUND_FILE")
	if workSoundFile == "" {
		workSoundFile = "work.mp3"
	}

	shortBreakSoundFile := os.Getenv("SHORT_BREAK_SOUND_FILE")
	if shortBreakSoundFile == "" {
		shortBreakSoundFile = "shortbreak.mp3"
	}

	longBreakSoundFile := os.Getenv("LONG_BREAK_SOUND_FILE")
	if longBreakSoundFile == "" {
		longBreakSoundFile = "longbreak.mp3"
	}

	for i := 0; i < pomodoroCount; i++ {
		fmt.Println("よーい...")
		time.Sleep(3 * time.Second)

		fmt.Println("作業セッション開始です！")
		runTimer("作業時間", workDuration)
		playSound(workSoundFile)

		if i < pomodoroCount-1 {
			runTimer("小休憩", shortBreakDuration)
			playSound(shortBreakSoundFile)
		} else {
			runTimer("大休憩", longBreakDuration)
			playSound(longBreakSoundFile)
		}
	}
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(valueStr)
	if err != nil {
		fmt.Printf("以下のEnvの値を適切に読み取れませんでした。 %s, デフォルト設定を使います: %v\n", key, defaultValue)
		return defaultValue
	}
	return duration
}

func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Printf("以下のEnvの値を適切に読み取れませんでした。 %s, デフォルト設定を使います: %v\n", key, defaultValue)
		return defaultValue
	}
	return value
}

func runTimer(label string, duration time.Duration) {
	fmt.Printf(" %s を開始します。時間は %s 分間です。\n", label, duration)
	time.Sleep(duration)
	fmt.Printf("%s 終了です。お疲れ様でした。\n", label)
}

func playSound(soundFile string) {
	f, err := os.Open(soundFile)
	if err != nil {
		fmt.Printf("サウンドファイルの読み取りに失敗しました。: %v\n", err)
		return
	}
	defer f.Close()

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		fmt.Printf("サウンドファイルのデコードに失敗しました。: %v\n", err)
		return
	}
	defer streamer.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		fmt.Printf("スピーカーをいい感じに設定できませんでした。: %v\n", err)
		return
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}
