package bitcoinwalletwatcher

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/blockcypher/gobcy"
)

// Watcher struct
type Watcher struct {
	info      *InfoFile
	loopInSec time.Duration
	addresses []string
	callbacks []WatcherCallback
	api       gobcy.API
}

// WatcherCallback method
type WatcherCallback func(string, int)

var (
	// ErrMissingBlockCyperToken error
	ErrMissingBlockCyperToken = errors.New("blockcyper.io API token is not provied")
)

// NewWatcher returns new bitcoin wallet watcher
func NewWatcher(cfg *Config) (*Watcher, error) {
	if cfg.BlockCyperToken == "" {
		return nil, ErrMissingBlockCyperToken
	}

	info, err := NewInfoStorage(cfg.InfoFile)
	if err != nil {
		return nil, err
	}

	if cfg.Coin == "" {
		cfg.Coin = "btc"
	}

	if cfg.Chain == "" {
		cfg.Coin = "main"
	}

	api := gobcy.API{
		Coin:  cfg.Coin,
		Chain: cfg.Chain,
		Token: cfg.BlockCyperToken,
	}

	return &Watcher{
		info:      info,
		loopInSec: time.Duration(cfg.DefaultLoopSec),
		addresses: cfg.Adresses,
		api:       api,
	}, nil
}

// SetCallback sets a new callback
func (w *Watcher) SetCallback(callback WatcherCallback) {
	w.callbacks = append(w.callbacks, callback)
}

// AddNewAddress adds new address to the list in order to listen it
func (w *Watcher) AddNewAddress(address string) {
	w.addresses = append(w.addresses, address)
}

// IsAddressExists checks the address list if the given address exists
func (w *Watcher) IsAddressExists(address string) bool {
	for _, addr := range w.addresses {
		if addr == address {
			return true
		}
	}
	return false
}

// RemoveAddress removes the given address from the list
func (w *Watcher) RemoveAddress(address string) {
	var temp []string
	index := 0
	for _, addr := range w.addresses {
		if addr != address {
			temp[index] = addr
			index++
		}
	}
	w.addresses = temp
}

// Run runs the watcher
func (w *Watcher) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			w.Stop()
			return
		default:
			w.readBlock()
			time.Sleep(time.Second * w.loopInSec)
		}
	}

}

func (w *Watcher) readBlock() {
	block, err := w.api.GetBlock(w.info.CurrentBlock, "", nil)
	if err != nil {
		log.Printf("Error getting block: %v\n", err)
		return
	}

	for _, tx := range block.TXids {
		w.checkTx(tx)
	}

	w.info.Update(block.Height + 1)
}

func (w *Watcher) checkTx(txHash string) {
	tx, err := w.api.GetTX(txHash, nil)
	if err != nil {
		log.Printf("Error getting transaction %v\n", err)
		return
	}

	for _, out := range tx.Outputs {
		for _, outAddress := range out.Addresses {
			if contains(w.addresses, outAddress) {
				for _, c := range w.callbacks {
					c(outAddress, out.Value)
				}
			}
		}
	}
}

// Stop runs when watcher stopping
func (w *Watcher) Stop() {
	if err := w.info.Save(); err != nil {
		log.Fatalf("error saving file %v", err)
	}
}
