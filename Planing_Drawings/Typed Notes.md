# Typed Notes
## ModAdders
### I2C Command List
 - 0x00 Get Module Type
 - 0x01 Get Module Config (Gets current scenario id and progress)
 - 0x02 Get Module Status (Should contain if solved)
 - 0x11 Set Module Config (Should contain the scenario id for that module, ie blue abort button with 1 bat and blue strip)
### TimeKeeping
    # the current countdown time should be sent out to all modules using the 0xFF address
### I2C Address Locations
#### Modules
##### Front
| | | |
|:-:|:-:|:-:|
| Main | 30 | 31 |
| 32 | 33 | 34 |
#### Back
| | | |
|:-:|:-:|:-:|
| 40 | 41 | Main |
| 42 | 43 | 44 |
#### Sides
| | | |
|:-:|:-:|:-:|
||50||
|53| Front Face | 51|
|| 54 ||
## SideCase
 - Each Side Should have its own PCB
 - Each Side Should have some indicators of each type (exluding serial and second factor)
 - Serial and Second Factor should be on the bottom
## FrontOfDevice
### Dimensions
 - The Device
  - 25" x 18" x 5"
 - The Modules
  - 8" x 8.365" (Depth has not been determined as of writing)
  - Locking Screw should be 1/4 20
  - PoGo pins should be centered at the top of the module
  - 6 PoGo pins are required
### Master Module
 - 7 segment display for time
 - warning indicator should just be two led's with difusers and then a mask applied to the front of the defuser to be stars
 - Inside of the device should be a extra sd card conatining stuff like WiFi Information (Parameters that cannot be set through the web interface)
## ModSlot
### Pin Order
| | | | | | |
|:-:|:-:|:-:|:-:|:-:|:-:|
|VCC|Ground|Addr|Interupt|Clock|Data|
 - Interupt is active low
 - Addr is a voltage devider
 - Addr's voltage must be present at time of module boot in order to get address
 - Addr Should tie to the first analog input on the chip of the Module
## New Notes
- Any information that does not change throughout the course of a defusal session should be displayed on a eInk Display
- Information that does change should use either a LCD or OLED
