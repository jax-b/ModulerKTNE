#include "ModControl.h"

ModControl::ModControl(int *busname){
    ModControl::file_i2c = &busname;
}
ModControl::stopGame(uint8_t address){
    if (ioctl(ModControl::file_i2c, I2C_SLAVE, address)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(ModControl::file_i2c, {0x40}, 1) != 1) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}
ModControl::startGame(uint8_t address){
    if (ioctl(ModControl::file_i2c, I2C_SLAVE, address)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(ModControl::file_i2c, {0x30}, 1) != 1) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}
ModControl::clearGameSerialNumber(uint8_t address){
    if (ioctl(ModControl::file_i2c, I2C_SLAVE, address)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(ModControl::file_i2c, {0x24}, 1) != 1) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}
ModControl::clearGameLitIndicator(uint8_t address){
    if (ioctl(ModControl::file_i2c, I2C_SLAVE, address)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(ModControl::file_i2c, {0x25}, 1) != 1) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}
ModControl::clearGameNumBatteries(uint8_t address){
    if (ioctl(ModControl::file_i2c, I2C_SLAVE, address)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(ModControl::file_i2c, {0x26}, 1) != 1) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}
ModControl::clearGamePortIDS(uint8_t address){
    if (ioctl(ModControl::file_i2c, I2C_SLAVE, address)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(ModControl::file_i2c, {0x27}, 1) != 1) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}
ModControl::clearGameSeed(uint8_t address){
    if (ioctl(ModControl::file_i2c, I2C_SLAVE, address)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(ModControl::file_i2c, {0x28}, 1) != 1) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}
ModControl::setSolvedStatus(uint8_t address, uint8_t status){
    if (ioctl(ModControl::file_i2c, I2C_SLAVE, address)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(ModControl::file_i2c, {0x11, status}, 2) != 2) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}
ModControl::syncGameTime(uint8_t address, unsigned long time){
    uint8_t buffer[5];
    buffer[0] = 0x12;
    buffer[1] = time >> 24;
    buffer[2] = time & 0xFF0000>> 16;
    buffer[3] = time & 0xFF00>> 8;
    buffer[4] = time & 0xFF;

    if (ioctl(ModControl::file_i2c, I2C_SLAVE, address)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(ModControl::file_i2c, buffer, 5) != 5) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}
ModControl::setStrikeReductionRate(uint8_t address, uint8_t rate){
    uint8_t buffer[5];
    buffer[0] = 0x13;
    buffer[1] = rate >> 24;
    buffer[2] = rate & 0xFF0000>> 16;
    buffer[3] = rate & 0xFF00>> 8;
    buffer[4] = rate & 0xFF;

    if (ioctl(ModControl::file_i2c, I2C_SLAVE, address)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(ModControl::file_i2c, buffer, 5) != 5) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}
ModControl::setGameSerialNumber(uint8_t address, uint8_t serialNumber[8]){
    uint8_t buffer[9];
    buffer[0] = 0x14;
    bytestosend = 1;
    for (int i = 0;i < 8; i++){
        if (serialNumber[i] != 0){
            buffer[i+1] = serialNumber[i];
            bytestosend++;
        }
    }   
    if (ioctl(ModControl::file_i2c, I2C_SLAVE, address)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(ModControl::file_i2c, buffer, bytestosend) != bytestosend) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}
ModControl::setGameLitIndicator(uint8_t address, char indlabel[3]){
    uint8_t buffer[4];
    buffer[0] = 0x15;
    buffer[1] = indlabel[0];
    buffer[2] = indlabel[1];
    buffer[3] = indlabel[2];
    if (ioctl(ModControl::file_i2c, I2C_SLAVE, address)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(ModControl::file_i2c, buffer, 4) != 4) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}
ModControl::setGameNumBatteries(uint8_t address, uint8_t num){
    if (ioctl(ModControl::file_i2c, I2C_SLAVE, address)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(ModControl::file_i2c, {0x16, num}, 2) != 2) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}
ModControl::setGamePortID(uint8_t address, uint8_t portID){
    if (ioctl(ModControl::file_i2c, I2C_SLAVE, address)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(ModControl::file_i2c, {0x17, portID}, 2) != 2) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}
ModControl::setGameSeed(uint8_t address, uint16_t seed){
    uint8_t buffer[3];
    buffer[0] = 0x18;
    buffer[1] = seed >> 8;
    buffer[2] = seed & 0xFF;
    if (ioctl(ModControl::file_i2c, I2C_SLAVE, address)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(ModControl::file_i2c, buffer, 3) != 3) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    return 1;
}
ModControl::getModType(uint8_t address){
    char buffer[4];

    if (ioctl(ModControl::file_i2c, I2C_SLAVE, address)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(ModControl::file_i2c, {0x19}, 1) != 1) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    if (read(ModControl::file_i2c, buffer, 4) != 4) {
        printf("Failed to read from the i2c bus.\n");
        return -1;
    }
    static char out[] = buffer;
    return *out;
}
ModControl::getModuleSolved(uint8_t address){
    if (ioctl(ModControl::file_i2c, I2C_SLAVE, address)) {
        printf("Failed to acquire bus access and/or talk to slave.\n");
        return -1;
    }
    if (write(ModControl::file_i2c, {0x1A}, 1) != 1) {
        printf("Failed to write to the i2c bus.\n");
        return -1;
    }
    uint8_t buffer[1];
    if (read(ModControl::file_i2c, buffer, 1) != 1) {
        printf("Failed to read from the i2c bus.\n");
        return -1;
    }
    return int8_t(buffer[1]);
}
ModControl::clearAllGameData(uint8_t address){
    if (ModControl::clearGameSerialNumber(address) == -1) return -1;
    if (ModControl::clearGameLitIndicator(address) == -1) return -1;
    if (ModControl::clearGameNumBatteries(address) == -1) return -1;
    if (ModControl::clearGamePortIDS(address) == -1) return -1;
    if (ModControl::clearGameSeed(address) == -1) return -1;
    if (ModControl::setSolvedStatus(address, 0) == -1) return -1;
    if (ModControl::stopGame(address) == -1) return -1;
    return 1;   
}
ModControl::setupAllGameData(uint8_t address, char serialNumber[], char litIndicators[][3], uint8_t numBatteries, uint8_t portIDs[], uint16_t seed = NULL){
    if (ModControl::clearGameFromMod(address) == -1) return -1;
    if (ModControl::setGameSerialNumber(address, serialNumber) == -1) return -1;
    for (int i = 0; i < length(indlabel); i++){
        if (ModControl::setGameLitIndicator(address, indlabel[i]) == -1) return -1;
    }
    if (ModControl::setGameNumBatteries(address, numBatteries) == -1) return -1;
    for  = 0; i < length(portIDs); i++){
        if (ModControl::setGamePortID(address, portIDs[i]) == -1) return -1;
    }
    if (seed != NULL) {
        if (ModControl::setGameSeed(address, seed) == -1) return -1;
    }
    ModControl::setSolvedStatus(address, 0);
    return 1;
}