package main

import "context"

var ble *BLE

func main() {
	ctx := sigIntContext(context.Background())
	ble = NewBLE()
	ble.SetTargetDeviceLocalName("ESP32-BLE-Remote")
	ble.SetTargetUUIDService("D7A74D4C-A077-445B-87E7-0B5DD5CE859D")
	ble.SetTargetUUIDCharacteristic("B4E7E29F-000E-469C-B3A4-A2FDFC578BBA")
	ble.Run(ctx)
}
