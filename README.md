# esp32-rpi-ble-example

* ESP32-C3(ベリフェラル)とRaspberry Pi4(セントラル)を使ったBLE送受信のサンプル。
* ベリフェラルのLocalNameとServiceUUID, CharacteristicUUIDでマッチする。
* ベリフェラル側のボタンを押すとカウントアップされ、セントラル側にNotifyで数値が通知される。
* セントラル側はWriteでベリフェラル側に数値が送られる。0x00でなければLEDが点灯する。
* 同一設定の複数のベリフェラルに対応。(二台でテスト)

![photo](./resource/photo.jpg)

# 環境

## ベリフェラル(子機)

[ソース](./esp32-c3)

* [ESP32-C3(RISC-V)](https://akizukidenshi.com/catalog/g/gM-17493/)を使った自作モジュール
* Arduino(C++),Arduino-CLI
* [Seeed Studio XIAO ESP32C3](https://akizukidenshi.com/catalog/g/gM-17454/)はピンアサインを変更すれば利用可能。

回路図

![schematics](./resource/schematics/board.png)

[PDF](./resource/schematics/board.pdf)

## セントラル(親機)

[ソース](./central)

* Raspberry Pi 4(aarch64)
* Raspberry Pi OS Lite(64Bit)
* Golang, tinygo bluetooth
* Linux専用

tinygo-bluetoothはクラスプラットフォームだが、macOSの場合はWriteWithoutResponseが正常に動作しなかった。

# TIPS

* macOSでUUIDを生成する場合は `uuidgen` を使うとよい

## Raspberry Pi OS上にarduino-cliを導入する

    $ cd
    $ curl -fsSL https://raw.githubusercontent.com/arduino/arduino-cli/master/install.sh | sh
    $ echo 'export PATH=$HOME/bin:$PATH' >> $HOME/.bashrc
    $ exec $SHELL
    $ arduino-cli config init

    $ sh -xe << 'EOS'
    arduino-cli config add board_manager.additional_urls 'https://arduino.esp8266.com/stable/package_esp8266com_index.json'
    arduino-cli config add board_manager.additional_urls 'https://espressif.github.io/arduino-esp32/package_esp32_index.json'
    arduino-cli core install esp8266:esp8266
    arduino-cli core install esp32:esp32
    arduino-cli board list
    arduino-cli board listall
    EOS

    $ sudo apt install python3-serial

## 参考URL

* [arduino-esp32 ble examples](https://github.com/espressif/arduino-esp32/tree/master/libraries/BLE/examples)
* [ESP32でBLEデバイスを作る](https://zenn.dev/luup/articles/iot-yamaguchi-20221204)
* [go bluetooth](https://github.com/tinygo-org/bluetooth)

