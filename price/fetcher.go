package price

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type timeRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

func Fetch(fromTime time.Time, toTime time.Time, out chan *HistoricalPrice, logger *zap.Logger) error {
	var chunks []*timeRange
	timeDiff := toTime.Sub(fromTime)
	chunksCount := timeDiff.Nanoseconds() / int64(24*time.Hour)
	logger.Info("chunks count:", zap.Int64("count", chunksCount))

	if chunksCount == 0 {
		chunks = append(chunks, &timeRange{
			From: fromTime,
			To:   toTime,
		})
	} else {
		startTime := fromTime
		for range chunksCount {
			endTime := startTime.Add(24 * time.Hour)
			chunks = append(chunks, &timeRange{
				From: startTime,
				To:   endTime,
			})
			startTime = endTime
		}
	}

	for _, chunk := range chunks {
		prices, err := fetch(chunk.From, chunk.To, "5m", logger)
		if err != nil {
			return fmt.Errorf("fetching chunk: %w", err)
		}
		firstPrice := prices[0]
		fixed := false
		for _, p := range prices {
			if !fixed && p.UnixTime != chunk.From.Unix() {
				diff := p.UnixTime - chunk.From.Unix()
				fixCount := diff / int64(5*time.Minute)
				for i := int64(0); i <= fixCount; i++ {
					out <- &HistoricalPrice{
						UnixTime: chunk.From.Add(5 * time.Minute * time.Duration(i)).Unix(),
						Value:    firstPrice.Value,
					}
				}
				fixed = true
			}
			out <- p
		}
		time.Sleep(100 * time.Millisecond) //prevent rate limit 13 request sec
	}

	close(out)

	fmt.Println("price channel closed")

	return nil
}

func fetch(fromTime time.Time, toTime time.Time, period string, logger *zap.Logger) ([]*HistoricalPrice, error) {
	fmt.Println("fetching from:", fromTime, "to:", toTime)
	client := &http.Client{}
	url := fmt.Sprintf(
		"https://public-api.birdeye.so/defi/history_price?address=4vMsoUT2BWatFweudnQM1xedRLfJgJ7hswhcpz4xgBTy&address_type=token&type=%s&time_from=%d&time_to=%d",
		period,
		fromTime.Unix(),
		toTime.Unix(),
	)
	logger.Info("calling:", zap.String("url", url))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	//todo: before making the we need to set this to a env var
	req.Header.Add("X-Api-Key", "c28c10a82a614930af83a48a3b47189d")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}
	response := &Response{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("unsuccess request: %s", string(body))
	}
	return response.Data.Items, err

}
