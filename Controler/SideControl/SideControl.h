#include <unistd.h>				//Needed for I2C port
#include <fcntl.h>				//Needed for I2C port
#include <sys/ioctl.h>			//Needed for I2C port
#include <linux/i2c-dev.h>		//Needed for I2C port

class SideControl {
    public:
        // Side Panel addresses
        const uint8_t TOP_PANEL = 0x50;
        const uint8_t RIGHT_PANEL = 0x51;
        const uint8_t BOTTOM_PANEL = 0x52;
        const uint8_t LEFT_PANEL = 0x53;
        
        SideControl(int *busname);
        
        // Set Serial Number
        int8_t setSerialNumber(char serialnumber[8]);
        // Set Lit Indicator
        // Max 2 per side
        // Once one it set caling this function again will set the other indicator for that panlel
        // If both are set the last one will be replaced with the new value
        int8_t setIndicator(uint8_t panel, bool lit, char indlabel[3]);
        // Set Side Art
        // Will start setting art on that panel then once the panel is full it will overwrite the last panel
        // first bit equals art type Battery or Port
        // 0 = Battery
        // 1 = Port
        // for Battery the last 2 bits sets the battery type
        // 0 = Battery
        // 0 = not used
        // 0 = not used
        // 0 = not used
        // 0 = not used
        // 0 = not used
        // 1 = 2xAA 
        // 1 = D
        // only one of the two last bits can be set
        // for Port the last six bits computes what ports are shown
        // 1 = Port
        // 0 = not used
        // 1 = DVI
        // 1 = Parallel
        // 1 = PS/2
        // 1 = RJ45
        // 1 = Serial
        // 1 = SteroRCA
        // 0x8b would be a parallel port with PS/2 and RJ45
        int8_t setSideArt(uint8_t panel, uint8_t artcode);

        int8_t clearSerialNumber();
        int8_t clearAllIndicator(uint8_t panel);
        int8_t clearAllSideArt(uint8_t panel);
    protected:
        int file_i2c;
}