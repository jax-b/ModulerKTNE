#include "SideControl.h"

SideControl::SideControl(int *busname)
{
    SideControl::file_i2c = &busname;
}
SideControl::setSerialNumber(char serialnumber[8]) {
    uint8_t buffer[9];
    buffer[0] = 0x10;
    bytestosend = 1;
    for (int i = 0;i < 8; i++){
        if (serialNumber[i] != 0){
            buffer[i+1] = serialNumber[i];
            bytestosend++;
        }
    }   
    if (ioctl(SideControl::file_i2c, I2C_SLAVE, SideControl::RIGHT_PANEL)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(SideControl::file_i2c, buffer, bytestosend) != bytestosend) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}
SideControl::clearSerialNumber() {
    if (ioctl(SideControl::file_i2c, I2C_SLAVE, SideControl::RIGHT_PANEL)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(SideControl::file_i2c, {0x20}, 1) != 1) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}
SideControl::setIndicator(uint8_t panel, bool lit, char indicator[3]) {\
    uint8_t buffer[5];
    buffer[0] = 0x11;
    buffer[1] = lit;
    for (int i = 0;i < 3; i++){
        buffer[i+2] = indicator[i];
    }
    if (ioctl(SideControl::file_i2c, I2C_SLAVE, panel)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(SideControl::file_i2c, buffer, 5) != 5) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}
SideControl::clearAllIndicator(uint8_t panel) {
    if (ioctl(SideControl::file_i2c, I2C_SLAVE, panel)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(SideControl::file_i2c, {0x21}, 1) != 1) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}
SideControl::setSideArt(uint8_t panel, uint8_t art]) {
    if (ioctl(SideControl::file_i2c, I2C_SLAVE, panel)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(SideControl::file_i2c, {0x12, art}, 2) != 2) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}
SideControl::clearAllSideArt(uint8_t panel) {
    if (ioctl(SideControl::file_i2c, I2C_SLAVE, panel)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(SideControl::file_i2c, {0x22}, 1) != 1) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}