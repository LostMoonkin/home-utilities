#include <Arduino.h>
#include "setup_wifi.h"

void setup() {
    Serial.begin(115200);
    Serial.printf("Deafult free size: %d\n", heap_caps_get_free_size(MALLOC_CAP_DEFAULT));
    Serial.printf("PSRAM free size: %d\n", heap_caps_get_free_size(MALLOC_CAP_SPIRAM));
    Serial.printf("Flash size: %d\n", ESP.getFlashChipSize());
// write your initialization code here
}

void loop() {
// write your code here
}