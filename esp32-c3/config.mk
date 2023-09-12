# シリアルポート
# AR_PORT=/dev/ttyACM0
AR_PORT=/dev/cu.usbmodem844401

# モニタボーレート
AR_MONITOR_BAUD=115200

# スケッチ
# AR_SKETCH=$(notdir $(basename $(shell find . -name '*.ino')))
AR_SKETCH=app

# -----------------
#   M5Camera
# AR_FQBN=esp32:esp32:esp32wrover:PartitionScheme=huge_app
# -----------------

# -----------------
#   mamemo 0002-AVR (ATmega328PB,Bootloader)
# AR_FQBN=MiniCore:avr:328:bootloader=uart0,variant=modelPB,clock=16MHz_external
# -----------------

# -----------------
#   mamemo 0005-WROOM
# AR_FQBN=esp8266:esp8266:generic:baud=3000000,ResetMethod=nodemcu,FlashMode=qio,eesz=4M2M,led=4
# -----------------

# -----------------
#   mamemo 0011-ESP32(ESP32-WROOM-32E)
# AR_FQBN=esp32:esp32:esp32:UploadSpeed=921600,PSRAM=enabled,FlashSize=16M,PartitionScheme=rainmaker
# -----------------

# -----------------
#  mamemo 0033-WROOM_C3_MODULE(ESP23-WROOM-C3)
AR_FQBN=esp32:esp32:esp32c3:DebugLevel=info,CDCOnBoot=cdc
# -----------------

