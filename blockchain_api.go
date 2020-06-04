package bitcoinwalletwatcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	rawBlockEndpoint = "https://blockchain.info/rawblock/%d"
)

type blockInformation struct {
	Hash   string              `json:"hash"`
	Height uint64              `json:"height"`
	Tx     []transactionDetail `json:"tx"`
}

type transactionDetail struct {
	Hash   string
	Inputs []transactionInputDetail
	Out    []transactionOutDetail
}

type transactionInputDetail struct {
	PrevOut struct {
		Amount  uint64 `json:"value"`
		Address string `json:"addr"`
	} `json:"prev_out,omitempty"`
}

type transactionOutDetail struct {
	Amount uint64 `json:"value"`
	Addr   string `json:"addr"`
}

func getBlockInformation(height uint64) (*blockInformation, error) {
	resp, err := http.Get(fmt.Sprintf(rawBlockEndpoint, height))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info blockInformation
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}

	return &info, nil
}
