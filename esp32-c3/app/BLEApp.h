#ifndef _BLEAPP_H_
#define _BLEAPP_H_

#include <Arduino.h>
#include <BLEDevice.h>
#include <BLEServer.h>
#include <BLEUtils.h>
#include <BLE2902.h>

class BLEAppClass {
    public:
        BLEAppClass();
        void init();
        void handler();
        void notify(String);
        bool deviceConnected = false;
        bool ledOn = false;
    private:
        BLEServer *pServer = NULL;
        BLECharacteristic *pCharacteristic = NULL;
        bool deviceConnectedPrev = false;
};

extern BLEAppClass BLEApp;
#endif