package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/fatih/color"
	emoji "github.com/jayco/go-emoji-flag"
	gt7 "github.com/snipem/go-gt7-telemetry/lib"
	"github.com/snipem/gt7-utils/lib/gtsport"
	"image"
	_ "image/png"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetDurationFromGT7Time(gt7time int32) time.Duration {
	seconds := gt7time / 1000
	milliseconds := gt7time % 1000

	return time.Duration(seconds)*time.Second + time.Duration(milliseconds*int32(time.Millisecond))

}

func cleanUpGtSportOnlineTime(oldtime string) (time.Duration, error) {

	// 1:30.766 -> 1m30s766ms
	oldtime = strings.Replace(oldtime, ":", "m", -1)
	oldtime = strings.Replace(oldtime, ".", "s", -1)
	oldtime += "ms"
	parsed, err := time.ParseDuration(oldtime)
	return parsed, err

}

func GetSportFormat(duration time.Duration) string {
	minusString := " "
	if duration < 0 {
		minusString = "-"

	}

	hours := int(math.Abs(duration.Hours()))
	minutes := int(math.Abs(duration.Minutes())) % 60
	seconds := int(math.Abs(duration.Seconds())) % 60
	milliseconds := duration.Milliseconds() % 1000
	milliseconds = int64(math.Abs(float64(milliseconds)))

	// If hours are present, accumulate them into minutes
	minutes += hours * 60
	return fmt.Sprintf("%s%02d:%02d.%03d", minusString, minutes, seconds, milliseconds)

}

func main() {

	var worldRecord time.Duration
	var err error

	topPercentage := 3
	if len(os.Args) >= 2 {
		topPercentage, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	}

	activeTimetrials, err := gtsport.GetActiveTimeTrials()
	if err != nil {
		log.Fatal(err)
	}

	for i, timetrial := range activeTimetrials {
		trackname, err := getTrackNameByNumber(timetrial.Parameters.Track.CourseCode)
		if err != nil {
			log.Fatal(err)
		}
		imageUrl := gtsport.GetImageUrl(timetrial.Parameters.Event.FlyerImagePath)
		daysToGo := time.Until(timetrial.Parameters.Online.EndDate).Hours() / 24
		fmt.Printf("%d: %s - %.2f days to go\n%s\n", i+1, trackname, daysToGo, imageUrl)
		//img, err := downloadImage(imageUrl)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//err = kittyimg.Fprint(os.Stdout, img)
		//if err != nil {
		//	println(err)
		//}

	}
	fmt.Print("Select event: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	fmt.Println(input.Text())

	k, err := strconv.Atoi(input.Text())
	onlineResult, err := gtsport.GetOnlineResult(activeTimetrials[k-1].EventID, 0)
	if err != nil {
		log.Fatal(err)
	}
	worldRecord = GetDurationFromGT7Time(int32(onlineResult.Result.List[0].Score))
	fmt.Printf("Got World Record: %s by %s %s (%s)\n", GetSportFormat(worldRecord),
		onlineResult.Result.List[0].User.CountryCode,
		onlineResult.Result.List[0].User.NickName,
		onlineResult.Result.List[0].User.NpOnlineID)

	if err != nil {
		log.Fatal(err)
	}

	//FIX-Me varies
	goldTime := worldRecord * (100 + time.Duration(topPercentage)) / 100
	silverTime := worldRecord * 105 / 100
	bronzeTime := worldRecord * 110 / 100

	fmt.Printf("World Record: %s\n", GetSportFormat(worldRecord))
	fmt.Printf("Gold Time %d%%: %s\n", topPercentage, GetSportFormat(goldTime))
	fmt.Printf("Silver Time : %s\n", GetSportFormat(silverTime))
	fmt.Printf("Bronze Time : %s\n", GetSportFormat(bronzeTime))

	gt7c := gt7.NewGT7Communication("255.255.255.255")
	go gt7c.Run()

	lastLap := int16(0)

	for {
		if gt7c.LastData.CurrentLap != lastLap {
			lastLap = gt7c.LastData.CurrentLap

			if gt7c.LastData.LastLap > 0 {

				lastLapTime := GetDurationFromGT7Time(gt7c.LastData.LastLap)

				fmt.Printf("\nLast Lap #%2d : %s\n", gt7c.LastData.CurrentLap, GetSportFormat(lastLapTime))
				fmt.Printf("🌍WR Diff    : ")
				printTimeInColor(worldRecord - lastLapTime)
				fmt.Printf("🥇Gold Diff  : ")
				printTimeInColor(goldTime - lastLapTime)
				fmt.Printf("🥈Silver Diff: ")
				printTimeInColor(silverTime - lastLapTime)
				fmt.Printf("🥉Bronze Diff: ")
				printTimeInColor(bronzeTime - lastLapTime)
				fmt.Println("")
			}
		}

		time.Sleep(100 * time.Millisecond)
	}

}

func countryCodeToFlag(code string) any {
	return emoji.GetFlag(code)
}

func printTimeInColor(td time.Duration) {

	if td > 0 {
		color.Green("%s", GetSportFormat(td))
	} else if td < 0 {
		color.Red("%s", GetSportFormat(td))
	}
}
func downloadImage(url string) (image.Image, error) {
	// Fetch the image from the URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching image: %v", err)
	}
	defer resp.Body.Close()

	// Create a temporary file to store the image
	tempFile, err := os.CreateTemp("", "image_*.png") // Adjust the file extension as needed
	if err != nil {
		return nil, fmt.Errorf("error creating temporary file: %v", err)
	}
	defer tempFile.Close()

	// Write the image data to the temporary file
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error writing image data to file: %v", err)
	}

	// Seek back to the beginning of the file
	if _, err := tempFile.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("error seeking to the beginning of the file: %v", err)
	}

	// Decode the image from the temporary file
	img, _, err := image.Decode(tempFile)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %v", err)
	}

	return img, nil
}

func getTrackNameByNumber(number int) (string, error) {
	// Download the CSV file
	resp, err := http.Get("https://github.com/ddm999/gt7info/raw/web-new/_data/db/course.csv")
	if err != nil {
		return "", fmt.Errorf("error downloading CSV file: %v", err)
	}
	defer resp.Body.Close()

	// Parse the CSV data
	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return "", fmt.Errorf("error parsing CSV data: %v", err)
	}

	// Find the track name for the provided number
	for _, record := range records {
		if len(record) >= 2 {
			trackNumber := strings.TrimSpace(record[0])
			trackName := strings.TrimSpace(record[1])
			if trackNumber == fmt.Sprintf("%d", number) {
				return trackName, nil
			}
		}
	}

	return "", fmt.Errorf("track with number %d not found", number)
}
