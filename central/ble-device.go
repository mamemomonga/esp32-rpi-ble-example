package main

import (
	"log"
	"time"

	"tinygo.org/x/bluetooth"
)

const (
	BLEDeviceStateNone = iota
	BLEDeviceStateConnecting
	BLEDeviceStateConnected
	BLEDeviceStateDiscovering
	BLEDeviceStateDisconnected
	BLEDeviceStateReady
)

type BLEDevices struct {
	devices []*BLEDevice
}

type BLEDevice struct {
	address        bluetooth.Address
	device         *bluetooth.Device
	characteristic *bluetooth.DeviceCharacteristic
	rssi           int16
	state          int
	connectOk      bool
	updateTime     time.Time
}

func NewBLEDevices() (t *BLEDevices) {
	t = new(BLEDevices)
	t.devices = []*BLEDevice{}
	return t
}

func (t *BLEDevices) Update(addr bluetooth.Address, rssi int16) {

	for _, d := range t.devices {
		// すでにあるならRSSIのみ更新
		if d.address == addr {
			d.rssi = rssi
			return
		}
	}
	t.devices = append(t.devices, &BLEDevice{
		address:    addr,
		rssi:       rssi,
		updateTime: time.Now(),
		state:      BLEDeviceStateNone,
	})
	log.Printf("[BLEDevices] %s add", addr.String())
}

func (t *BLEDevices) Connecting() *BLEDevice {
	for _, d := range t.devices {
		if d.state == BLEDeviceStateNone {
			d.state = BLEDeviceStateConnecting
			d.device = nil
			d.characteristic = nil
			d.connectOk = false
			d.updateTime = time.Now()
			log.Printf("[BLEDevices] %s connecting", d.address.String())
			return d
		}
	}
	return nil
}

func (t *BLEDevices) Discovering() *BLEDevice {
	for _, d := range t.devices {
		if d.state == BLEDeviceStateConnected {
			d.state = BLEDeviceStateDiscovering
			d.characteristic = nil
			d.updateTime = time.Now()
			log.Printf("[BLEDevices] %s discovering", d.address.String())
			return d
		}
	}
	return nil
}

func (t *BLEDevices) Connected(addr bluetooth.Address) {
	d := t.Device(addr)
	if d == nil {
		return
	}
	d.state = BLEDeviceStateConnected
	d.updateTime = time.Now()
	log.Printf("[BLEDevices] %s connected", addr.String())
}

func (t *BLEDevices) Disconnected(addr bluetooth.Address) {
	d := t.Device(addr)
	if d == nil {
		return
	}
	if d.state == BLEDeviceStateReady {
		d.characteristic.EnableNotifications(nil)
	}
	if d.connectOk {
		d.device.Disconnect()
	}
	d.state = BLEDeviceStateDisconnected
	d.updateTime = time.Now()
	d.connectOk = false
	log.Printf("[BLEDevices] %s disconnected", addr.String())
}

func (t *BLEDevices) Ready(addr bluetooth.Address) {
	d := t.Device(addr)
	if d == nil {
		return
	}
	d.state = BLEDeviceStateReady
	d.updateTime = time.Now()
	log.Printf("[BLEDevices] %s ready", addr.String())
}

func (t *BLEDevices) Device(addr bluetooth.Address) *BLEDevice {
	for _, d := range t.devices {
		if d.address == addr {
			return d
		}
	}
	return nil
}

func (t *BLEDevices) DevicesReady() []*BLEDevice {
	dr := []*BLEDevice{}
	for _, d := range t.devices {
		if d.state == BLEDeviceStateReady {
			dr = append(dr, d)
		}
	}
	return dr
}

func (t *BLEDevices) DisconnectAll() {
	for _, d := range t.devices {
		t.Disconnected(d.address)
	}
}

func (t *BLEDevices) Cleanup() {
	nd := []*BLEDevice{}
	for _, d := range t.devices {
		// 切断されて5秒たったデバイス
		if d.state == BLEDeviceStateDisconnected {
			if time.Since(d.updateTime).Seconds() > 5 {
				log.Printf("[BLEDevices] %s remove", d.address.String())
				continue
			}
		}
		nd = append(nd, d)
	}
	t.devices = nd
}

func (t *BLEDevices) Show() {
	for n, c := range t.devices {
		state := "NONE"
		switch c.state {
		case BLEDeviceStateConnecting:
			state = "CONNECTING"
		case BLEDeviceStateDiscovering:
			state = "DISCOVERING"
		case BLEDeviceStateConnected:
			state = "CONNECTED"
		case BLEDeviceStateDisconnected:
			state = "DISCONECTED"
		case BLEDeviceStateReady:
			state = "READY"
		}

		log.Printf("[BLEDevices]  ** STATE %s ** %02d %ddB %s %s\n",
			c.address.String(),
			n+1,
			c.rssi,
			time.Since(c.updateTime).String(),
			state,
		)
	}
}
