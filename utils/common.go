package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var LOG_LEVELS = map[string]int{
	"ERROR": 0,
	"INFO":  1,
	"DEBUG": 2,
}

func LogMe(level string, info string, message interface{}) {
	log_file, cond1 := os.LookupEnv("UCMDB_PROVIDER_LOG")
	log_level, cond2 := os.LookupEnv("UCMDB_PROVIDER_LOG_LEVEL")

	if cond1 && cond2 {
		log_file = strings.TrimSpace(log_file)
		ll_index, ll_found := LOG_LEVELS[log_level]

		l_index, l_found := LOG_LEVELS[level]

		if ll_found && l_found && l_index <= ll_index {
			currentTime := time.Now()
			file, err := os.OpenFile(log_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				log.Fatal(err)
			}

			data := fmt.Sprintf("--- %v %s: %s: ---\n%v\n", currentTime.Format("2006-01-02 15:04:05"), level, info, message)

			_, err1 := file.WriteString(data)
			if err1 != nil {
				log.Fatal(err)
			}
			file.Close()
		}
	}
}

var reStringID = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func validStringID(v *string) bool {
	return v != nil && reStringID.MatchString(*v)
}

func UnmarshalJson(data string, filename string, v interface{}) (interface{}, error) {
	var b []byte
	if strings.TrimSpace(data) != "" {
		b = []byte(data)
	} else if strings.TrimSpace(data) == "" && strings.TrimSpace(filename) != "" {
		var err error
		b, err = os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("UnmarshalJson: %v", err)
		}
	} else {
		return nil, fmt.Errorf("json string is nil")
	}

	err := json.Unmarshal(b, v)
	if err != nil {
		return nil, fmt.Errorf("UnmarshalJson: %v", err)
	}
	return v, nil
}
