# Raspberry Pi OS上にarduino-cliを導入する

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