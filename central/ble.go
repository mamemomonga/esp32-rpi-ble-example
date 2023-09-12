package main

import (
	"context"
	"log"
	"time"

	"tinygo.org/x/bluetooth"
)

type BLE struct {
	adapter        *bluetooth.Adapter
	deviceAddress  *bluetooth.Address
	device         *bluetooth.Device
	characteristic bluetooth.DeviceCharacteristic
	targets        BLETargets
	ledStatus      bool
	connected      bool
	ctx            context.Context
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
	t.ledStatus = false
	t.connected = false
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

func (t *BLE) connectHandler(device bluetooth.Address, connected bool) {
	if connected {
		log.Printf("[BLE] CONNECTED %s", device.String())
		t.connected = true
	} else {
		log.Printf("[BLE] DISCONNECTED %s", device.String())
		t.Stop()
	}
}

func (t *BLE) scan() (err error) {
	ch := make(chan bool, 1)

	time.Sleep(time.Second * 1)
	log.Println("[BLE] scanning")

	err = t.adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		select {
		case <-t.ctx.Done():
			adapter.StopScan()
			ch <- true
			return
		default:
		}
		//		log.Println("[BLE] ",
		//			device.Address.String(),
		//			device.RSSI,
		//			device.LocalName(),
		//		)
		if (device.LocalName() == t.targets.LocalName) && (device.RSSI != 0) {
			adapter.StopScan()
			t.deviceAddress = &device.Address
			ch <- true
		}
	})
	if err != nil {
		return err
	}
	<-ch
	return nil

}

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
