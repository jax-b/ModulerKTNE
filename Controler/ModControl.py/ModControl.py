# Constant Addresses
# Front Face
# Main | 1 | 2
#   3  | 4 | 5
from typing import ValuesView


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

class ModuleController:
    import smbus
    clsbusctrl = None
    def __init__ (self, busctrl: smbus.SMBus):
        self.clsbusctrl = busctrl
    
    # Stop
    ## Tells the module that the game is over and to stop active gameplay functions
    def stopGame (self, address: int):
        self.clsbusctrl.write_byte_data(address, 0x40)
    
    # Start
    ## Tells the module to start the game and start its internal timer
    ## if no game seed is set a random seed will be generated
    def startGame (self, address: int):
        self.clsbusctrl.write_byte_data(address, 0x30)
    
    # Clear
    ## Clears out the gameplay serial numbner 
    def clearGameSerialNumber (self, address: int):
        self.clsbusctrl.write_byte_data(address, 0x24)
    ## Clears out all gameplay lit indicators
    def clearGameLitIndicators (self, address: int):
        self.clsbusctrl.write_byte_data(address, 0x25)
    ## Clears out the number of batteries
    def clearGameNumBateries (self, address: int):
        self.clsbusctrl.write_byte_data(address, 0x26)
    ## Clears out all ports from the module
    def clearGamePortIDS (self, address: int):
        self.clsbusctrl.write_byte_data(address, 0x27)
    ## Clears out the game seed from the module
    def clearGameSeed (self, address: int):
        self.clsbusctrl.write_byte_data(address, 0x28)

    # Set
    ## Sets the solved status of the module
    ## 0 = Unsolved
    ## 1 = Solved
    ## anything negitive is the number of strikes upto -128
    def setSolvedStatus (self, address: int, value: int):
        if value > 1 or value > -128:
            raise ValueError("Value must be between -128 and 1")
        self.clsbusctrl.write_byte_data(address, 0x11, value)
    ## Syncs the game time of the module with the provided value
    def SyncModuleTime (self, address: int, value: int):
        self.clsbusctrl.write_byte_data(address, 0x12, value.to_bytes(4, byteorder='big'))
    ## Sets how fast the module should accelerate per strike
    ## should be 0 <= x < 1
    ## defaults to 0.25
    def SetStrikeReductionRate (self, address: int, value: float):
        if value <= 0 or value > 1:
            raise ValueError("Value must be between 0 and 1")
        import ctypes
        value = ctypes.c_float(value).value
        self.clsbusctrl.write_i2c_block_data(address, 0x13, value)
    ## Sets the game serial number
    ## max length ogf serial number is 8, the rest will get cut off if provided 
    def setGameSerialNumber (self, address: int, value: str):
        if len(value) > 8:
            value = value[:8]
        self.clsbusctrl.write_byte_data(address, 0x14, value)
    ## Sets a lit indicator, only provide indicators that are lit
    ## Each time a lit indicator is sent, the module will append it to the list of indicators
    ## There is a max of 6 indicators that can be lit at a time
    ## once all are sent the last indicator will be overwritten if more data is sent
    ## use clearGameLitIndicators to clear the list
    ## the indicator label is exactly 3 characters long
    def setGameLitIndicator (self, address: int, label: str):
        if len(label) != 3:
            raise ValueError("Label must be exactly 3 characters long")
        self.clsbusctrl.write_byte_data(address, 0x15, label)
    ## Sets the number of batteries
    ## 0 - 255
    def setGameNumBatteries (self, address: int, value: int):
        if value > 255 or value < 0:
            raise ValueError("Value must be between 0 and 255")
        self.clsbusctrl.write_byte_data(address, 0x16, value)
    ## Sets the port ID
    ## 0x1: DVI-D, 0x2: Parallel, 0x3: PS2, 0x4: RJ-45, 0x5: Serial, 0x6: Stereo RCA
    ## There is a max of 6 ports that can be set at a time
    ## Once all ports are set the last port will be overwritten if more ports are set
    ## use clearGamePortIDS to clear the list
    ## Only send that specific port ID once, thats all that matters for the game logic
    def setGamePortID (self, address: int, portid: int):
        if portid < 0x1 or portid > 0x6:
            raise ValueError("Value must be between 0x1 and 0x6")
        self.clsbusctrl.write_byte_data(address, 0x17, portid)
    ## Sets the game seed
    ## The seed is a 2 byte number, 1-65535
    def setGameSeed (self, address: int, value: int):
        if value > 65535 or value < 1:
            raise ValueError("Value must be between 0 and 65535")
        import ctypes
        value = ctypes.c_uint16(value).value
        self.clsbusctrl.write_byte_data(address, 0x18, value)

    # Gets
    def getModuleType (self, address: int):
        return self.clsbusctrl.read_i2c_block_data(address, 0x00)
    def getModuleSolved (self, address: int):
        return self.clsbusctrl.read_byte_data(address, 0x01)


    # Users Functions
    ## Clears all game specific data from the module
    def clearGameFromMod (self, address: int):
        self.clearGameSerialNumber(address)
        self.clearGameLitIndicators(address)
        self.clearGameNumBateries(address)
        self.clearGamePortIDS(address)
        self.clearGameSeed(address)
        self.setSolvedStatus(address, 0)
        self.stopGame(address)
    
    ## Send all game data to the module after clearing it
    ## Defult values are specified but it would be a boring game if they were the same
    def setupGame (self, address: int, serial: str, numBatteries = 0, portID = [], seed = 0, litIndicators = []):
        self.clearGameFromMod(address)
        self.setGameSerialNumber(address, serial)
        self.setGameNumBatteries(address, numBatteries)
        for i in portID:
            self.setGamePortID(address, i)
        for i in litIndicators:
            self.setGameLitIndicator(address, i)
        self.setGameSeed(address, seed)
        self.setSolvedStatus(address, 0)