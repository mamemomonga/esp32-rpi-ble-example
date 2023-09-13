package main

import (
	"context"
	"log"
	"sync"
	"time"

	"tinygo.org/x/bluetooth"
)

type BLE struct {
	adapter   *bluetooth.Adapter
	targets   BLETargets
	devices   *BLEDevices
	ledStatus bool
	ctx       context.Context
}

type BLETargets struct {
	LocalName      string
	Service        bluetooth.UUID
	Characteristic bluetooth.UUID
}

func NewBLE() (t *BLE) {
	t = new(BLE)
	t.adapter = bluetooth.DefaultAdapter
	t.adapter.SetConnectHandler(t.connectHandler)
	t.targets = BLETargets{}
	t.devices = NewBLEDevices()
	t.ledStatus = false
	return t
}

func (t *BLE) SetTargetDeviceLocalName(s string) {
	t.targets.LocalName = s
}

func (t *BLE) SetTargetUUIDService(s string) (err error) {
	uuid, err := bluetooth.ParseUUID(s)
	if err != nil {
		return
	}
	t.targets.Service = uuid
	return nil
}

func (t *BLE) SetTargetUUIDCharacteristic(s string) (err error) {
	uuid, err := bluetooth.ParseUUID(s)
	if err != nil {
		return
	}
	t.targets.Characteristic = uuid
	return nil
}

func (t *BLE) connectHandler(addr bluetooth.Address, connected bool) {
	if connected {
		log.Printf("[BLE] %s CONNECTED", addr.String())
		t.devices.Connected(addr)
	} else {
		log.Printf("[BLE] %s DISCONNECTED", addr.String())
		t.devices.Disconnected(addr)
		//		t.Stop()
	}
}

func (t *BLE) scanStart() {
	log.Print("[BLE] start scan")
	err := t.adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		select {
		case <-t.ctx.Done():
			adapter.StopScan()
			log.Print("[BLE] stop scan")
			return
		default:
		}
		if (device.LocalName() == t.targets.LocalName) && (device.RSSI != 0) {
			t.devices.Update(device.Address, device.RSSI)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (t *BLE) connectLoop(wg *sync.WaitGroup) {
	for {
		time.Sleep(time.Second * 1)
		select {
		case <-t.ctx.Done():
			t.devices.DisconnectAll()
			wg.Done()
			return
		default:
		}
		t.devices.Cleanup()
		t.connect()
		t.discover()
		t.devices.Show()
	}
}

func (t *BLE) connect() {
	d := t.devices.Connecting()
	if d == nil {
		return
	}

	var err error
	// connect
	d.device, err = t.adapter.Connect(
		d.address,
		bluetooth.ConnectionParams{
			ConnectionTimeout: bluetooth.NewDuration(time.Second * 5),
		},
	)
	if err != nil {
		log.Printf("[BLE] %s ERR [Connect] %s",
			d.address.String(),
			err,
		)
		t.devices.Disconnected(d.address)
		time.Sleep(time.Second * 1)
		return
	}
}

func (t *BLE) discover() {
	d := t.devices.Discovering()
	if d == nil {
		return
	}

	srvcs, err := d.device.DiscoverServices(nil)
	if err != nil {
		log.Printf("[BLE] %s ERR DiscoverServices %s",
			d.address.String(),
			err,
		)
		d.device.Disconnect()
		t.devices.Disconnected(d.address)
		time.Sleep(time.Second * 1)
		return
	}

	for _, srvc := range srvcs {
		log.Printf("[BLE] %s - service %s",
			d.address.String(),
			srvc.UUID().String(),
		)
		if srvc.UUID() == t.targets.Service {
			chars, err := srvc.DiscoverCharacteristics(nil)
			if err != nil {
				log.Printf("[BLE] %s ERR DiscoverCharacteristics %s",
					d.address.String(),
					err,
				)
				d.device.Disconnect()
				t.devices.Disconnected(d.address)
				time.Sleep(time.Second * 1)
				return
			}

			for _, char := range chars {
				log.Printf("[BLE] %s -- characteristic %s",
					d.address.String(),
					char.UUID().String(),
				)
				if char.UUID() == t.targets.Characteristic {
					d.characteristic = &char
					// notify
					d.characteristic.EnableNotifications(func(buf []byte) {
						t.notify(d.address, buf)
					})
					t.devices.Ready(d.address)
					return
				}
			}
		}
	}
}

/*
func (t *BLE) discover() (found bool, err error) {
	log.Println("[BLE] connecting")

	// connect
	t.device, err = t.adapter.Connect(
		*t.deviceAddress,
		bluetooth.ConnectionParams{},
	)
	if err != nil {
		log.Println(err)
		return false, err
	}

	select {
	case <-t.ctx.Done():
		return false, nil
	default:
	}

	// discover
	log.Println("[BLE] discovering services/characteristics")

	srvcs, err := t.device.DiscoverServices(nil)
	if err != nil {
		return false, err
	}

	for _, srvc := range srvcs {
		log.Println("[BLE] - service", srvc.UUID().String())
		if srvc.UUID() == t.targets.Service {
			chars, err := srvc.DiscoverCharacteristics(nil)
			if err != nil {
				return false, err
			}
			for _, char := range chars {
				log.Println("[BLE] -- characteristic", char.UUID().String())
				if char.UUID() == t.targets.Characteristic {
					t.characteristic = char
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func (t *BLE) Stop() {
	log.Println("[BLE] stop")
	if t.connected {
		t.connected = false
		t.characteristic.EnableNotifications(nil)
		t.device.Disconnect()
		log.Println("[BLE] disconnect")
	}
}
*/
