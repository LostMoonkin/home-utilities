//
// Created by moonkin on 2023/9/12.
//

#ifndef ESP32_FAN_CONTROLLER_SETUP_WIFI_H
#define ESP32_FAN_CONTROLLER_SETUP_WIFI_H

#endif //ESP32_FAN_CONTROLLER_SETUP_WIFI_H

void startAPMode(const IPAddress* apIP, const IPAddress* subnet, const char* ssid);
void startSTAMode(const char* ssid, const char* password, int timeoutMills);
