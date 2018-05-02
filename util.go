package main

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	IP      string `json:"ip"`
	Port    string `json:"port"`
	Public  string `json:"public"`
	Private string `json:"private"`
	Token   string `json:"token"`
}

func (c *config) getPrivateKey() (string, error) {
	privateKey, err := ioutil.ReadFile(c.Private)
	if err != nil {
		return "", err
	}
	return string(privateKey), nil
}

func (c *config) getPublicKey() (string, error) {
	publicKey, err := ioutil.ReadFile(c.Public)
	if err != nil {
		return "", err
	}
	return string(publicKey), nil
}

func (c *config) getWebHookURL() string {
	return c.IP + c.Port + "/" + c.Token
}

func loadConfig(file string) (*config, error) {
	c := &config{}
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(raw, &c); err != nil {
		return nil, err
	}
	return c, nil
}
