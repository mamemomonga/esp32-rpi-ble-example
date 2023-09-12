#include "BLEApp.h"
#include "config.h"

// ------------------------------------
// Callbacks
// ------------------------------------
class BLEAppServerCallbacks: public BLEServerCallbacks {
	void onConnect(BLEServer *pServer) {
		BLEApp.deviceConnected = true;
		Serial.println("[BLEAPP] connect");

		digitalWrite(LED_BUILTIN,HIGH);
		delay(500);
		digitalWrite(LED_BUILTIN,LOW);

	};
	void onDisconnect(BLEServer *pServer) {
		BLEApp.deviceConnected = false;
		Serial.println("[BLEAPP] disconnect");
 
		digitalWrite(LED_BUILTIN,HIGH);
		delay(500);
		digitalWrite(LED_BUILTIN,LOW);

	};

};

class BLEAppCharacteristicCallbacks: public BLECharacteristicCallbacks {
	void onWrite(BLECharacteristic *pCharacteristic) {
		digitalWrite(LED_BUILTIN,HIGH);
		delay(5);
		digitalWrite(LED_BUILTIN,LOW);

		Serial.println("[BLEAPP] write");
		std::string rxValue = pCharacteristic->getValue();
		if( rxValue.length() > 0 ){
			BLEApp.ledOn = rxValue[0]!=0;
			Serial.print("[BLEAPP] received: ");
			for(int i=0; i<rxValue.length(); i++ ){
				Serial.print(rxValue[i],HEX);
			}
			Serial.println();
		}
	};
};

// ------------------------------------
// Class
// ------------------------------------
BLEAppClass::BLEAppClass() {}

void BLEAppClass::init() {
    Serial.println("[BLEApp] init");
    BLEDevice::init(BLE_DEVICE_NAME);

	// Server
	pServer = BLEDevice::createServer();
    pServer->setCallbacks(new BLEAppServerCallbacks());

  	// Service
	BLEService *pService = pServer->createService(SERVICE_UUID);

	// Characteristic
	pCharacteristic = pService->createCharacteristic(
		CHARACTERISTIC_UUID,
		BLECharacteristic::PROPERTY_WRITE  |
		BLECharacteristic::PROPERTY_NOTIFY
 	);
 	pCharacteristic->setCallbacks(new BLEAppCharacteristicCallbacks());
 	pCharacteristic->addDescriptor(new BLE2902());

	// Start
	pService->start();

	// Advertising
	BLEAdvertising *pAdvertising = BLEDevice::getAdvertising();
	pAdvertising->addServiceUUID(SERVICE_UUID);
	pAdvertising->setScanResponse(false);
	pAdvertising->setMinPreferred(0x0);
	BLEDevice::startAdvertising();
	Serial.println("[BLEApp] start advertising");

}

void BLEAppClass::handler() {
	// disconnecting
	if(!deviceConnected && deviceConnectedPrev) {
		delay(500);
		pServer->startAdvertising();
		Serial.println("[BLEAPP] restart advertising");
		deviceConnectedPrev = deviceConnected;
	}

	// connecting
	if(deviceConnected && !deviceConnectedPrev){
		deviceConnectedPrev = deviceConnected;
	}
}

void BLEAppClass::notify(String str) {
	if( !deviceConnected ) return;
	std::string txValue = str.c_str();
	pCharacteristic->setValue(txValue);
	pCharacteristic->notify();
}

BLEAppClass BLEApp;
