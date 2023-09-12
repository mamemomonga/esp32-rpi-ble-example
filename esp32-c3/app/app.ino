#include "config.h"
#include "BLEApp.h"

int buttonCount = 0;
unsigned long blinker = 0;

void setup() {
 	pinMode(LED_BUILTIN, OUTPUT);
 	pinMode(LED1, OUTPUT);
 	pinMode(BTN1, INPUT_PULLUP);

	digitalWrite(LED_BUILTIN, HIGH);
	Serial.begin(115200);
	Serial.setDebugOutput(true);

	BLEApp.init();	

	delay(100);
	digitalWrite(LED_BUILTIN, LOW);
}

void loop() {
	BLEApp.handler();

	if(!BLEApp.deviceConnected) {
		unsigned long ms = millis();
		if((ms - blinker) > 1000) {
			digitalWrite(LED_BUILTIN,HIGH);
			delay(5);
			digitalWrite(LED_BUILTIN,LOW);
			blinker=ms;
		}
		return;
	}

	// LED
	digitalWrite(LED1,BLEApp.ledOn ? HIGH:LOW);

	// BUTTON
	if(digitalRead(BTN1) == LOW){
		digitalWrite(LED_BUILTIN,HIGH);
		buttonCount++;
		String str = " BTN: "+String(buttonCount);
		Serial.println(str);
		BLEApp.notify(str);
		while(digitalRead(BTN1) == LOW) delay(10);
		delay(50);
		digitalWrite(LED_BUILTIN,LOW);
	}
}
