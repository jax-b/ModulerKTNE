# Constant Addresses
TOP_PANEL = 0x50
RIGHT_PANEL = 0x51
BOTTOM_PANEL = 0x52
LEFT_PANEL = 0x53

class SideController:
    import smbus
    clsbusctrl = None
    def __init__(busctrl: smbus.SMBus):
        clsbusctrl = busctrl
    
    # Resets a indicators value to nothing
    def clearIndicator(address: int, indicatornum: int):
        if address not in range(256) or indicatornum not in range(2):
            return -1
        # write data to the specified side panel
        # first byte is thje indicator number, second byte is if it is lit, following 3 bytes are char
        # 0x11 is set indictator
        # the first 4 bits of byte 3 are the indicator adddress and bits 567 are not used and bit 8 is whether or not the indicator is lit
        clsbusctrl.write_i2c_block_data(address, 0x11, indicatornum << 4, "")
        
    def setIndicator(address: int, indicatornum: int, indicatortext: str, lit: bool):
        if address not in range(256) or indicatornum not in range(2) or len(str) > 3:
            return -1
        # write data to the specified side panel
        # first byte is thje indicator number, second byte is if it is lit, following 3 bytes are char
        # 0x11 is set indictator
        # the first 4 bits of byte 3 are the indicator adddress and bits 567 are not used and bit 8 is whether or not the indicator is lit
        indicatorNumlit = (indicatornum << 4) + lit
        clsbusctrl.write_i2c_block_data(address, 0x11, indicatorNumlit, indicatortext)

    def clearSerialNumber():
        clsbusctrl.write_i2c_block_data(BOTTOM_PANEL, 0x12, "")

    def setSerialNumber(SerialNumber: str):
        if SerialNumber > 6:
            return -1
        # 0x12 is set Serial Number
        clsbusctrl.write_i2c_block_data(BOTTOM_PANEL, 0x12, SerialNumber)