# Typed Notes
## ModAdders
### I2C Command List
 - 0x00 Get Module Type
 - 0x01 Get Module Config (Gets current scenario id and progress)
 - 0x02 Get Module Status (Should contain if solved)
 - 0x11 Set Module Solved Status (Sets if solved)
 - 0x12 Time Sync

### Time Keeping
~~- the current countdown time should be sent out to all modules using the 0xFF address~~
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
|| 52 ||
## SideCase
 - Each Side Should have its own PCB
 - Each Side Should have some indicators of each type (excluding serial and second factor)
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
 - warning indicator should just be two led's with diffuses and then a mask applied to the front of the diffuse to be stars
 - Inside of the device should be a extra sd card conatining stuff like WiFi Information (Parameters that cannot be set through the web interface)
## ModSlot
### Pin Order
| | | | | | |
|:-:|:-:|:-:|:-:|:-:|:-:|
|VCC|Ground|Addr|Interrupt|Clock|Data|
 - Interupt is active low
 - Addr is a voltage divider
 - Addr's voltage must be present at time of module boot in order to get address
 - Addr Should tie to the first analog input on the chip of the Module
## New Notes
- Any information that does not change throughout the course of a defusal session should be displayed on a eInk Display
- Information that does change should use either a LCD or OLED

## Gameplay port ID's
There are 6 types of ports in the game here is the id's we use to keep track of them
```
  1 = Port
  0 = not used
  1 = DVI
  1 = Parallel
  1 = PS/2
  1 = RJ45
  1 = Serial
  1 = SteroRCA
```