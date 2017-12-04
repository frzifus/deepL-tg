package deepL

import (
    "net/http"
    "sort"
    "io/ioutil"
    "encoding/json"
    "bytes"
    "errors"
)

type lang struct {
    UserPreferredLangs []string `json:"user_preferred_langs"`
    SourceLangUserSelected string `json:"source_lang_user_selected"`
    TargetLang string `json:"target_lang"`
}

// Job -
type Job struct {
    Kind string `json:"kind"`
    RawEnSentence string `json:"raw_en_sentence"`
} 

type params struct {
    Jobs []Job `json:"jobs"`
    lang `json:"lang"`
    Priority int `json:"priority"`
}


type deepL struct {
    Jsonrpc string `json:"jsonrpc"`
    Method string `json:"method"`
    ID int `json:"id"`

    params `json:"params"`
}

// Beams -
type Beams struct {
    NumSymbols            int     `json:"num_symbols"`
    PostprocessedSentence string  `json:"postprocessed_sentence"`
    Score                 float64 `json:"score"`
    TotalLogProb          float64 `json:"totalLogProb"`
} 

type deepLResponse struct {
	ID      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		SourceLang            string `json:"source_lang"`
		SourceLangIsConfident int    `json:"source_lang_is_confident"`
		TargetLang            string `json:"target_lang"`
		Translations          []struct {
            Beams []Beams `json:"beams"`
			TimeAfterPreprocessing   int `json:"timeAfterPreprocessing"`
			TimeReceivedFromEndpoint int `json:"timeReceivedFromEndpoint"`
			TimeSentToEndpoint       int `json:"timeSentToEndpoint"`
			TotalTimeEndpoint        int `json:"total_time_endpoint"`
		} `json:"translations"`
	} `json:"result"`
}

type dconfig struct {
    Path    string
	IP      string `json:"ip"`
	Port    string `json:"port"`
	Public  string `json:"public"`
	Private string `json:"private"`
	Token   string `json:"token"`
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

// LoadConfig -
func LoadConfig(path string) (*dconfig, error) {
    d := &dconfig{}
    raw, err := ioutil.ReadFile(path + "./bot.json")
    if err != nil {
        return d, err
    }
    err = json.Unmarshal(raw, &d)
    if err != nil {
        return d, err
    }
    d.Public = path + "./" + d.Public
    d.Private = path + "./" + d.Private
    return d, nil
}

// NewDeepL -
func NewDeepL() *deepL {
    return &deepL{
        Jsonrpc : "2.0",
        Method : "LMT_handle_jobs",
        ID : 1,
        params : params{
            lang : lang{
            UserPreferredLangs : []string{"EN", "DE"},
            SourceLangUserSelected : "auto",
            TargetLang : "DE",  // target_lang
            },
          Priority : 1,  
        },
    }
}

func (d *deepL) SupportedLang() []string {
    return []string{"EN", "DE"}
}

func (d *deepL) SetTargetLang(lang string) {
    d.params.lang.TargetLang = lang
}

func (d *deepL) AddJob(rawSentence string) {
    j := Job{
        Kind : "default",
        RawEnSentence: rawSentence,
    }
    d.params.Jobs = append(d.params.Jobs, j)
}

func (d *deepL) ResetJobs() {
    d.params.Jobs = d.params.Jobs[:0]
}

// TODO: add timeout
func (d *deepL) Request() (*deepLResponse, error) {
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
    body, _ := ioutil.ReadAll(resp.Body)

    return decodeResponse(string(body))
}

func decodeResponse(responseJSON string) (*deepLResponse, error) {
    res := &deepLResponse{}
    return res, json.Unmarshal([]byte(responseJSON), res)    
}

func (d *deepLResponse) sortByScore() {
    sort.Sort(BeamsRes(d.Result.Translations[0].Beams))
}

func (d *deepLResponse) Translation() string {
    d.sortByScore()
    return d.Result.Translations[0].Beams[0].PostprocessedSentence
}
