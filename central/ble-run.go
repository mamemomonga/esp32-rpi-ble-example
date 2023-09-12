package main

import (
	"context"
	"log"
)

func (t *BLE) Run(ctx context.Context) (err error) {
	t.ctx = ctx
	err = t.adapter.Enable()
	if err != nil {
		return
	}
	go t.connectionCleanup()

	done := make(chan bool, 1)
	go func() {
		defer func() {
			// t.Stop()
			done <- true
		}()
		for {
			t.scan()
			if err != nil {
				log.Fatal(err)
			}
			select {
			case <-t.ctx.Done():
				return
			default:
			}
		}
	}()
	<-done
	close(done)
	return nil
}

/*
				found := false
				for !found {
					// scan
					err := t.scan()
					if err != nil {
						log.Fatal(err)
					}
					select {
					case <-t.ctx.Done():
						return
					default:
					}

					// discover
					found, err = t.discover()
					if err != nil {
						log.Fatal(err)
					}
					select {
					case <-t.ctx.Done():
						return
					default:
					}
				}
				if found {
					// notification
					err = t.characteristic.EnableNotifications(t.notify)
					if err != nil {
						log.Fatal(err)
					}

					// running
					for t.connected {
						select {
						case <-t.ctx.Done():
							return
						default:
						}
						t.blink()
						time.Sleep(time.Millisecond * 500)
					}
				}
			}
		}()

	<-done
	close(done)
	return nil
}

func (t *BLE) notify(buf []byte) {
	log.Printf("Notify: %s\n", string(buf))
}

func (t *BLE) blink() {
	var buf []byte
	if t.ledStatus {
		buf = []byte{0xFF}
	} else {
		buf = []byte{0x00}
	}
	_, err := t.characteristic.WriteWithoutResponse(buf)
	if err != nil {
		log.Printf("[BLE] Error Write: %v", err)
	}
	log.Printf("[BLE] Write: %0X", buf)
	t.ledStatus = !t.ledStatus
}
*/
