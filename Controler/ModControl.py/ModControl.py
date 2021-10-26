# Constant Addresses
# Front Face
# Main | 1 | 2
#   3  | 4 | 5
FRONT_MOD_1 = 0x30
FRONT_MOD_2 = 0x31
FRONT_MOD_3 = 0x32
FRONT_MOD_4 = 0x33
FRONT_MOD_5 = 0x34
# Back Face
#   1 | 2 | Main
#   3 | 4 | 5
BACK_MOD_1 = 0x40
BACK_MOD_2 = 0x41
BACK_MOD_3 = 0x42
BACK_MOD_4 = 0x43
BACK_MOD_5 = 0x44

class ModuleControler:
    import smbus
    clsbusctrl = None
    def __init__(busctrl: smbus.SMBus):
        clsbusctrl = busctrl
    
    # Gets a modules type
    def getModuleType(address: int):
        return clsbusctrl.read_i2c_block_data(address, 0x00)
    
    # Get Module Solved Status
    def getModuleSolved(address: int):
        return clsbusctrl.read_byte_data(address, 0x01)

    # Sets a modules solved status
    def SetModuleSolved(address: int, value: int):
        clsbusctrl.write_byte_data(address, 0x11, value)
    
    # Sync Module Time
    def SyncModuleTime(address: int, value: int):
        clsbusctrl.write_byte_data(address, 0x12, value.to_bytes(4, byteorder='big'))

    # Set Strike Reduction Rate
    def SetStrikeReductionRate(address: int, value: int):
        import struct
        clsbusctrl.write_i2c_block_data(address, 0x13, struct.pack('f', value))
