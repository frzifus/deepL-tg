package DeepL

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sort"
)

type lang struct {
	UserPreferredLangs     []string `json:"user_preferred_langs"`
	SourceLangUserSelected string   `json:"source_lang_user_selected"`
	TargetLang             string   `json:"target_lang"`
}

// Job -
type Job struct {
	Kind          string `json:"kind"`
	RawEnSentence string `json:"raw_en_sentence"`
}

type params struct {
	Jobs     []Job `json:"jobs"`
	lang     `json:"lang"`
	Priority int `json:"priority"`
}

// DeepL -
type DeepL struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	ID      int    `json:"id"`

	params `json:"params"`
}

// Beams -
type Beams struct {
	NumSymbols            int     `json:"num_symbols"`
	PostprocessedSentence string  `json:"postprocessed_sentence"`
	Score                 float64 `json:"score"`
	TotalLogProb          float64 `json:"totalLogProb"`
}

// DeepLResponse -
type DeepLResponse struct {
	ID      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		SourceLang            string `json:"source_lang"`
		SourceLangIsConfident int    `json:"source_lang_is_confident"`
		TargetLang            string `json:"target_lang"`
		Translations          []struct {
			Beams                    []Beams `json:"beams"`
			TimeAfterPreprocessing   int     `json:"timeAfterPreprocessing"`
			TimeReceivedFromEndpoint int     `json:"timeReceivedFromEndpoint"`
			TimeSentToEndpoint       int     `json:"timeSentToEndpoint"`
			TotalTimeEndpoint        int     `json:"total_time_endpoint"`
		} `json:"translations"`
	} `json:"result"`
}

// BeamsRes -
type BeamsRes []Beams

func (b BeamsRes) Len() int {
	return len(b)
}

func (b BeamsRes) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b BeamsRes) Less(i, j int) bool {
	return b[i].Score < b[j].Score
}

// NewDeepL -
func NewDeepL() *DeepL {
	return &DeepL{
		Jsonrpc: "2.0",
		Method:  "LMT_handle_jobs",
		ID:      1,
		params: params{
			lang: lang{
				UserPreferredLangs:     []string{"EN", "DE"},
				SourceLangUserSelected: "auto",
				TargetLang:             "DE", // target_lang
			},
			Priority: 1,
		},
	}
}

// SupportedLang returns a string array with all supported languages.
func (d *DeepL) SupportedLang() []string {
	return []string{"EN", "DE"}
}

// SetTargetLang - Sets target language by country abbreviation.
func (d *DeepL) SetTargetLang(lang string) {
	d.params.lang.TargetLang = lang
}

// AddJob - creates a job to be done using a record.
func (d *DeepL) AddJob(rawSentence string) {
	j := Job{
		Kind:          "default",
		RawEnSentence: rawSentence,
	}
	d.params.Jobs = append(d.params.Jobs, j)
}

// ResetJobs - resets the job queue.
func (d *DeepL) ResetJobs() {
	d.params.Jobs = d.params.Jobs[:0]
}

// Request - Send the request and wait for a result.
func (d *DeepL) Request() (*DeepLResponse, error) {
	// TODO: add timeout
	url := "https://deepl.com/jsonrpc"
	jsonStr, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}

	req.Header.Set("X-Custom-Header", "")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return decodeResponse(string(body))
}

func decodeResponse(responseJSON string) (*DeepLResponse, error) {
	res := &DeepLResponse{}
	return res, json.Unmarshal([]byte(responseJSON), res)
}

func (d *DeepLResponse) sortByScore() {
	sort.Sort(BeamsRes(d.Result.Translations[0].Beams))
}

// Translation - Returns the best translation result.
func (d *DeepLResponse) Translation() string {
	d.sortByScore()
	return d.Result.Translations[0].Beams[0].PostprocessedSentence
}
