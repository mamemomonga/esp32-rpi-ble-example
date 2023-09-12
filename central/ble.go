package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"tinygo.org/x/bluetooth"
)

type BLE struct {
	adapter     *bluetooth.Adapter
	targets     BLETargets
	connections []*BLEConnections
	ledStatus   bool
	ctx         context.Context
}

type BLETargets struct {
	LocalName      string
	Service        bluetooth.UUID
	Characteristic bluetooth.UUID
}

type BLEConnections struct {
	deviceAddress  *bluetooth.Address
	device         *bluetooth.Device
	characteristic bluetooth.DeviceCharacteristic
	connected      bool
	updateTime     time.Time
}

func NewBLE() (t *BLE) {
	t = new(BLE)
	t.adapter = bluetooth.DefaultAdapter
	t.adapter.SetConnectHandler(t.connectHandler)
	t.targets = BLETargets{}
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

func (t *BLE) connectionShow() {
	for _, c := range t.connections {
		connected := "disconnected"
		if c.connected {
			connected = "connected"
		}
		fmt.Printf("[BLE]  Cstat %s %s %s\n",
			c.deviceAddress.String(),
			connected,
			time.Since(c.updateTime).String(),
		)
	}
}

func (t *BLE) connectionConnect(addr bluetooth.Address) {
	t.connections = append(t.connections, &BLEConnections{
		deviceAddress: &addr,
		updateTime:    time.Now(),
		connected:     true,
	})
	t.connectionShow()
}

func (t *BLE) connectionDisconnect(addr bluetooth.Address) {
	for _, c := range t.connections {
		if *c.deviceAddress == addr {
			c.updateTime = time.Now()
			c.connected = false
		}
		log.Printf("[BLE]  Cdisconnect %s %s",
			c.deviceAddress.String(),
			time.Since(c.updateTime).String(),
		)
	}
	t.connectionShow()
}

func (t *BLE) connectionCleanup() {
	for {
		select {
		case <-t.ctx.Done():
			return
		default:
		}

		fmt.Print("+")

		ncs := []*BLEConnections{}
		for _, c := range t.connections {
			if !c.connected {
				// 切断から10秒していたらエントリ削除
				if time.Since(c.updateTime).Seconds() > 10 {
					log.Printf("[BLE] Cremove: %s", c.deviceAddress.String())
					continue
				}
			}
			ncs = append(ncs, c)
		}
		t.connections = ncs

		time.Sleep(time.Second * 1)
	}
}

func (t *BLE) connectHandler(addr bluetooth.Address, connected bool) {
	if connected {
		log.Printf("[BLE] CONNECTED %s", addr.String())
		t.connectionConnect(addr)
	} else {
		log.Printf("[BLE] DISCONNECTED %s", addr.String())
		t.connectionDisconnect(addr)
		//		t.Stop()
	}
}

func (t *BLE) scan() {
	done := make(chan bool, 1)
	err := t.adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		select {
		case <-t.ctx.Done():
			adapter.StopScan()
			log.Print("SCAN STOP")
			done <- true
			return
		default:
		}
		fmt.Print(".")

		/*
			log.Printf("[BLE]   %s %ddB %s",
				device.Address.String(),
				device.RSSI,
				device.LocalName(),
			)
		*/

		if (device.LocalName() == t.targets.LocalName) && (device.RSSI != 0) {
			fmt.Print("!")

			exists := false
			for _, act := range t.connections {
				if *act.deviceAddress == device.Address {
					exists = true
				}
			}
			fmt.Print("1")
			if !exists {
				fmt.Print("#")
				fmt.Println()
				log.Printf("[BLE] %s %ddB %s",
					device.Address.String(),
					device.RSSI,
					device.LocalName(),
				)
				t.connect(&device.Address)
			}
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
func (t *BLE) connect(deviceAddress *bluetooth.Address) {
	log.Printf("[BLE] CONNECTING: %s", deviceAddress.String())

	// connect
	device, err := t.adapter.Connect(
		*deviceAddress,
		bluetooth.ConnectionParams{
			ConnectionTimeout: bluetooth.NewDuration(time.Second * 5),
		},
	)
	if err != nil {
		t.connectionDisconnect(*deviceAddress)
		log.Println(err)
		return
	}

	_ = device

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
