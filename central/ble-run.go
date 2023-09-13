package main

import (
	"context"
	"log"
	"sync"
	"time"

	"tinygo.org/x/bluetooth"
)

func (t *BLE) Run(ctx context.Context) (err error) {
	t.ctx = ctx
	err = t.adapter.Enable()
	if err != nil {
		return
	}

	var wg sync.WaitGroup

	go t.scanStart()

	wg.Add(1)
	go t.connectLoop(&wg)

	go func() {
		for {
			t.blink()
			select {
			case <-t.ctx.Done():
				return
			default:
				time.Sleep(time.Millisecond * 500)
			}
		}
	}()

	wg.Wait()

	return nil
}

func (t *BLE) notify(addr bluetooth.Address, buf []byte) {
	log.Printf("[BLE] %s Notify: %s\n",
		addr.String(),
		string(buf),
	)
}

func (t *BLE) blink() {
	var buf []byte
	if t.ledStatus {
		buf = []byte{0xFF}
	} else {
		buf = []byte{0x00}
	}

	for _, d := range t.devices.DevicesReady() {
		_, err := d.characteristic.WriteWithoutResponse(buf)
		if err != nil {
			log.Printf("[BLE] %s Error Write: %v",
				d.address.String(),
				err,
			)
		}
		log.Printf("[BLE] %s Write: %0X",
			d.address.String(),
			buf,
		)
	}

	t.ledStatus = !t.ledStatus
}
