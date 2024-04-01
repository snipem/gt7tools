package gtsport

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type TimeTrialResult struct {
	BoardID      string `json:"board_id"`
	DisplayRank  int    `json:"display_rank"`
	LineReplayID string `json:"line_replay_id"`
	RankingStats struct {
		CarCode int `json:"car_code"`
		Color   int `json:"color"`
		StyleID int `json:"style_id"`
	} `json:"ranking_stats"`
	ReplayID   int64     `json:"replay_id,omitempty"`
	Score      int       `json:"score"`
	UpdateTime time.Time `json:"update_time"`
	User       struct {
		AvatarPhotoID       int64  `json:"avatar_photo_id"`
		CountryCode         string `json:"country_code"`
		IsHiddenAvatar      bool   `json:"is_hidden_avatar"`
		IsHiddenNickname    bool   `json:"is_hidden_nickname"`
		IsValidUser         bool   `json:"is_valid_user"`
		NickName            string `json:"nick_name"`
		NpOnlineID          string `json:"np_online_id"`
		UserID              string `json:"user_id"`
		DriverRating        int    `json:"driver_rating"`
		IsStarPlayer        bool   `json:"is_star_player"`
		ManufacturerID      int    `json:"manufacturer_id"`
		SportsmanshipRating int    `json:"sportsmanship_rating"`
	} `json:"user"`
}

type OnlineResult struct {
	Result struct {
		List  []TimeTrialResult `json:"list"`
		Total int               `json:"total"`
	} `json:"result"`
}

func GetOnlineResult(eventId int, page int) (OnlineResult, error) {

	requestBody, err := json.Marshal(map[string]interface{}{
		"board_id": fmt.Sprintf("p_rt_100%d_001", eventId),
		"page":     page,
	})
	if err != nil {
		fmt.Println("Error marshaling request body:", err)
		return OnlineResult{}, err
	}

	// Send the request
	resp, err := http.Post("https://web-api.gt7.game.gran-turismo.com/ranking/get_list_by_page", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return OnlineResult{}, err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected response status:", resp.Status)
		return OnlineResult{}, err
	}

	// Decode the response body into OnlineResult struct
	var result OnlineResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding response body:", err)
		return OnlineResult{}, err
	}
	return result, nil
}

type Event struct {
	BeginDate           int    `json:"begin_date"`
	EndDate             int    `json:"end_date"`
	EventDescriptionKey string `json:"event_description_key"`
	EventTitleKey       string `json:"event_title_key"`
	FlyerImagePath      string `json:"flyer_image_path"`
	PlayerSelectType    string `json:"player_select_type"`
	RegistrationKey     string `json:"registration_key"`
}

type Online struct {
	BeginDate time.Time `json:"begin_date"`
	EndDate   time.Time `json:"end_date"`
	RankingID string    `json:"ranking_id"`
}

type EventResult struct {
	EventID    int `json:"event_id"`
	Parameters struct {
		Championship struct {
			GlobalRankingID         string `json:"global_ranking_id"`
			RoundID                 int    `json:"round_id"`
			SeasonID                int    `json:"season_id"`
			TotalRaceCount          int    `json:"total_race_count"`
			UseDivisorForValidRaces bool   `json:"use_divisor_for_valid_races"`
			ValidRaceCount          int    `json:"valid_race_count"`
			ValidRaceDivisor        int    `json:"valid_race_divisor"`
		} `json:"championship"`
		Event   Event  `json:"event"`
		Online  Online `json:"online"`
		Penalty struct {
			IsGuideCourseOff bool `json:"is_guide_course_off"`
			PenaltyLevel     int  `json:"penalty_level"`
		} `json:"penalty"`
		Pit struct {
			RefuelingSpeed int `json:"refueling_speed"`
		} `json:"pit"`
		Race struct {
			BehaviorDamage         string `json:"behavior_damage"`
			BehaviorSlipStreamType string `json:"behavior_slip_stream_type"`
			ConsumeFuel            int    `json:"consume_fuel"`
			ConsumeTire            int    `json:"consume_tire"`
			LimitFuelCapacity      int    `json:"limit_fuel_capacity"`
			PitConstraint          int    `json:"pit_constraint"`
			StartType              string `json:"start_type"`
		} `json:"race"`
		Regulation struct {
			Compounds         []string      `json:"compounds"`
			LimitCarSettings  bool          `json:"limit_car_settings"`
			LimitNos          string        `json:"limit_nos"`
			LimitWideFender   string        `json:"limit_wide_fender"`
			RequiredCompounds []interface{} `json:"required_compounds"`
			Tire              string        `json:"tire"`
			UseBalanceOfPower bool          `json:"use_balance_of_power"`
		} `json:"regulation"`
		Track struct {
			CourseCode   int    `json:"course_code"`
			CourseLabel  string `json:"course_label"`
			CourseLayout int    `json:"course_layout"`
			Hour         int    `json:"hour"`
			Minute       int    `json:"minute"`
			WeatherID    int    `json:"weather_id"`
		} `json:"track"`
	} `json:"parameters"`
}
type Folder struct {
	Result []EventResult `json:"result"`
}

func (result EventResult) isActive() bool {
	return time.Now().Before(result.Parameters.Online.EndDate)
}

func GetFolder(regionID, folderID int) (Folder, error) {
	// Prepare the request body
	requestBody, err := json.Marshal(map[string]interface{}{
		"region_id": regionID,
		"folder_id": folderID,
	})
	if err != nil {
		return Folder{}, fmt.Errorf("error marshaling request body: %v", err)
	}

	// Send the request
	resp, err := http.Post("https://web-api.gt7.game.gran-turismo.com/event/get_folder", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return Folder{}, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return Folder{}, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	// Decode the response body into Folder struct
	var folder Folder
	if err := json.NewDecoder(resp.Body).Decode(&folder); err != nil {
		return Folder{}, fmt.Errorf("error decoding response body: %v", err)
	}

	return folder, nil
}

func GetActiveTimeTrials() ([]EventResult, error) {
	// Region is Europe, folderId seems to be Timetrials
	folder, err := GetFolder(6, 249)

	if err != nil {
		return nil, err
	}

	var activeTimeTrials []EventResult
	for _, result := range folder.Result {
		if result.isActive() {
			activeTimeTrials = append(activeTimeTrials, result)
		}
	}
	return activeTimeTrials, nil
}

func GetImageUrl(url string) string {
	return fmt.Sprintf("https://asset.gt7.game.gran-turismo.com/delivery/raw/event/image/%s", url)
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
