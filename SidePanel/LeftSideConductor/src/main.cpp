#include <SPI.h>
#include <GxEPD2_3C.h>
#include <Fonts/FreeMonoBold12pt7b.h>
#include <ArduinoJson.h>
#include <Adafruit_GFX.h>
#include "BatteryPictures.h"
#include "PortPictures.h"

/// Buffers and data tracking for I2C communication
byte incomeingI2CData[10];
byte outgoingI2CData[10];
uint8_t bytesToSend = 0;
uint8_t bytesReceived = 0;

#define PCB_VERSION 1

#if PCB_VERSION == 1
#define EPD_RESET_PIN 5
#define EPD_DC_PIN 4
#define EPD_MOSI_PIN 3
#define EPD_SCK_PIN 2
#define EPD_BUS_PIN 18
#define DSP0_CS_PIN 6
#define DSP0_LED_PIN 7
#define DSP1_CS_PIN 14
#define DSP1_LED_PIN 9
#define DSP2_CS_PIN 10
#define DSP3_CS_PIN 11
#define DSP4_CS_PIN 12
#define DSP4_LED_PIN 13
#define DSP5_CS_PIN 8
#define DSP5_LED_PIN 15
#define S2S_SERIAL_RX 17
#define S2S_SERIAL_TX 16
#define S2S_SERIAL_SPEED 115200
#elif PCB_VERSION == 2
#define EPD_RESET_PIN 5
#define EPD_DC_PIN 4
#define EPD_MOSI_PIN 3
#define EPD_SCK_PIN 2
#define EPD_BUS_PIN 18
#define DSP0_CS_PIN 6
#define DSP0_LED_PIN 7
#define DSP1_CS_PIN 8
#define DSP1_LED_PIN 9
#define DSP2_CS_PIN 10
#define DSP3_CS_PIN 11
#define DSP4_CS_PIN 12
#define DSP4_LED_PIN 13
#define DSP5_CS_PIN 14
#define DSP5_LED_PIN 15
#define S2S_SERIAL_RX 17
#define S2S_SERIAL_TX 16
#define S2S_SERIAL_SPEED 115200
#endif

GxEPD2_290_C90c IndicatorDSP0(/*CS=*/DSP0_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/EPD_BUS_PIN); // Indicator
GxEPD2_290_C90c IndicatorDSP1(/*CS=*/DSP1_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/EPD_BUS_PIN); // Ind
GxEPD2_290_C90c ArtDSP0(/*CS=*/DSP2_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/EPD_BUS_PIN);
GxEPD2_290_C90c ArtDSP1(/*CS=*/DSP3_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/EPD_BUS_PIN);
GxEPD2_290_C90c IndicatorDSP2(/*CS=*/DSP4_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/EPD_BUS_PIN);
GxEPD2_290_C90c IndicatorDSP3(/*CS=*/DSP5_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/EPD_BUS_PIN);

GFXcanvas1 BlackCanvas(IndicatorDSP0.WIDTH, IndicatorDSP0.HEIGHT);
GFXcanvas1 RedCanvas(IndicatorDSP0.WIDTH, IndicatorDSP0.HEIGHT);

// Tracking Variables
uint8_t indiNumber = 0;
const uint8_t ARTDSPRANGE[] = {1, 5};
const uint8_t MAXDSPNUM = ARTDSPRANGE[1];
uint8_t batDSPDrawCycle = 0;
uint8_t portDSPDrawCycle = 0;
uint8_t portBatDSPMap[2][5] = {
    {1, 2}, // Port DSPS
    {3, 4}  // Bat DSPS
};
uint8_t indicatorLEDMap[6] = { 0 };

// Initializes and clears the connected displays
void setupDisplays()
{
  // Tell remote to startup
  Serial1.println("{\"command\":\"startup\"}");

  // Setup DSP Pins
  pinMode(EPD_RESET_PIN, OUTPUT);
  pinMode(EPD_DC_PIN, OUTPUT);
  pinMode(DSP0_CS_PIN, OUTPUT);
  pinMode(DSP1_CS_PIN, OUTPUT);
  pinMode(DSP2_CS_PIN, OUTPUT);
  pinMode(DSP3_CS_PIN, OUTPUT);
  pinMode(DSP4_CS_PIN, OUTPUT);
  pinMode(DSP5_CS_PIN, OUTPUT);
  pinMode(EPD_BUS_PIN, INPUT);

  // Release All Displays
  digitalWrite(DSP0_CS_PIN, true);
  digitalWrite(DSP1_CS_PIN, true);
  digitalWrite(DSP2_CS_PIN, true);
  digitalWrite(DSP3_CS_PIN, true);
  digitalWrite(DSP4_CS_PIN, true);
  digitalWrite(DSP5_CS_PIN, true);

  // Reset All Displays
  Serial.println("DSP RESET");
  digitalWrite(EPD_RESET_PIN, true);
  delay(100);
  digitalWrite(EPD_RESET_PIN, false);
  delay(100);
  digitalWrite(EPD_RESET_PIN, true);

  // Init and attempt to fill the displays buffers
  Serial.println("DSP INIT");
  IndicatorDSP0.init();
  IndicatorDSP0.writeScreenBuffer(GxEPD_WHITE);
  IndicatorDSP0.init();
  IndicatorDSP0.writeScreenBuffer(GxEPD_WHITE);
  ArtDSP0.init();
  ArtDSP0.writeScreenBuffer(GxEPD_WHITE);
  ArtDSP1.init();
  ArtDSP1.writeScreenBuffer(GxEPD_WHITE);
  IndicatorDSP1.init();
  IndicatorDSP1.writeScreenBuffer(GxEPD_WHITE);
  IndicatorDSP2.init();
  IndicatorDSP2.writeScreenBuffer(GxEPD_WHITE);

  // Grab all dislays and clear them then release them
  digitalWrite(DSP0_CS_PIN, false);
  digitalWrite(DSP1_CS_PIN, false);
  digitalWrite(DSP2_CS_PIN, false);
  digitalWrite(DSP3_CS_PIN, false);
  digitalWrite(DSP4_CS_PIN, false);
  digitalWrite(DSP5_CS_PIN, false);
  IndicatorDSP0.clearScreen(GxEPD_WHITE);
  digitalWrite(DSP0_CS_PIN, true);
  digitalWrite(DSP1_CS_PIN, true);
  digitalWrite(DSP2_CS_PIN, true);
  digitalWrite(DSP3_CS_PIN, true);
  digitalWrite(DSP4_CS_PIN, true);
  digitalWrite(DSP5_CS_PIN, true);

  delay(10);
}

// Draws a Battery to the specified display
// True is AA false is D
void writeBattery(uint8_t displayNum, bool AAorD)
{
  // Code needs to be updated to support multiple displays
  // displays[displayNum].fillScreen(GxEPD_WHITE);
  
  switch (displayNum)
  {
  case 0:
    if (AAorD)
    {
      ArtDSP0.writeScreenBuffer(GxEPD_WHITE);
      ArtDSP0.writeImage(epd_bitmap_AA_Black, epd_bitmap_AA_Black, 0, 6, 250, 122, true, false, true);
    }
    else
    {
      ArtDSP0.writeScreenBuffer(GxEPD_WHITE);
      ArtDSP0.writeImage(epd_bitmap_D_Black, epd_bitmap_D_Red, 0, 6, 250, 122, true, false, true);
    }
    break;
  case 1:
    if (AAorD)
    {
      ArtDSP1.writeScreenBuffer(GxEPD_WHITE);
      ArtDSP1.writeImage(epd_bitmap_AA_Black, epd_bitmap_AA_Black, 0, 6, 250, 122, true, false, true);
    }
    else
    {
      ArtDSP1.writeScreenBuffer(GxEPD_WHITE);
      ArtDSP1.writeImage(epd_bitmap_D_Black, epd_bitmap_D_Red, 0, 6, 250, 122, true, false, true);
    }
    break;
  }
}

// Draws the SerialNumber to dsp 0
void writeSerialNum(String text)
{
  DynamicJsonDocument doc(50);
  doc["write"]["SerialNumber"] = text;
  Serial1.println(serializeJson(doc, Serial));
}

// Draws the SerialNumber to dsp 0
// Overload for a char array instead of a string
void writeSerialNum(char intext[])
{
  String strtext = intext;
  writeSerialNum(strtext);
}

// Draws a port to the specified display
// Host is responsible for avoiding collisions, bad data will not be drawn but there might be a missing port caused by a collision
void writePorts(uint8_t displayNum, byte activeports)
{
  // 1 = Port
  // 0 = not used
  // 1 = DVI
  // 1 = Parallel
  // 1 = PS/2
  // 1 = RJ45
  // 1 = Serial
  // 1 = SteroRCA
  bool dviActive = activeports & 0b00100000;
  bool parallelActive = activeports & 0b00010000;
  bool ps2Active = activeports & 0b00001000;
  bool rj45Active = activeports & 0b00000100;
  bool serialActive = activeports & 0b00000010;
  bool rcaActive = activeports & 0b00000001;
  // Code needs to be updated to support multiple displays
  // displays[displayNum].fillScreen(GxEPD_WHITE);
  const uint8_t VAL_TO_ADD = 67;

  if (displayNum < 4)
  { // If the display number is on the Conductor node populate the display with the correct artwork
    BlackCanvas.fillScreen(GxEPD_WHITE);
    RedCanvas.fillScreen(GxEPD_WHITE);
    uint8_t longporty = 6;
    uint8_t shortportx = 0;
    if (rj45Active)
    {
      BlackCanvas.drawBitmap(shortportx, 6, epd_bitmap_rj45_Black, epd_bitmap_rj45_size[0], epd_bitmap_rj45_size[1], GxEPD_BLACK);
      longporty += VAL_TO_ADD;
      shortportx += epd_bitmap_rj45_size[0] + 5;
    }
    if (ps2Active)
    {
      BlackCanvas.drawBitmap(shortportx, 6, epd_bitmap_PS2_Black, epd_bitmap_PS2_size[0], epd_bitmap_PS2_size[1], GxEPD_BLACK);
      if (longporty < 30)
      {
        longporty += VAL_TO_ADD;
      }
      shortportx += epd_bitmap_PS2_size[0] + 5;
    }
    if (rcaActive)
    {
      BlackCanvas.drawBitmap(213, 6, epd_bitmap_RCA_Black, epd_bitmap_RCA_size[0], epd_bitmap_RCA_size[1], GxEPD_BLACK);
      RedCanvas.drawBitmap(213, 6, epd_bitmap_RCA_Black, epd_bitmap_RCA_size[0], epd_bitmap_RCA_size[1], GxEPD_RED);
    }
    bool pass = false;
    if (dviActive)
    {
      BlackCanvas.drawBitmap(0, longporty, epd_bitmap_DVI_Black, epd_bitmap_DVI_size[0], epd_bitmap_DVI_size[1], GxEPD_BLACK);
      if (longporty > 30)
      {
        pass = true;
      }
      else
      {
        longporty += VAL_TO_ADD;
      }
    }
    if (serialActive && !pass)
    {
      BlackCanvas.drawBitmap(0, longporty, epd_bitmap_Serial_Black, epd_bitmap_Serial_size[0], epd_bitmap_Serial_size[1], GxEPD_BLACK);
      if (longporty > 30)
      {
        pass = true;
      }
      else
      {
        longporty += VAL_TO_ADD;
      }
    }
    if (parallelActive && !rcaActive && !pass)
    {
      BlackCanvas.drawBitmap(0, longporty, epd_bitmap_Parallel_Black, epd_bitmap_Parallel_size[0], epd_bitmap_Parallel_size[1], GxEPD_BLACK);
    }
  }

  switch (displayNum)
  {
  case 0:
    ArtDSP0.writeScreenBuffer(GxEPD_WHITE);
    ArtDSP0.writeImage(BlackCanvas.getBuffer(), RedCanvas.getBuffer(), 0, 6, 250, 122, true, false, true);
    break;
  case 1:
    ArtDSP1.writeScreenBuffer(GxEPD_WHITE);
    ArtDSP1.writeImage(BlackCanvas.getBuffer(), RedCanvas.getBuffer(), 0, 6, 250, 122, true, false, true);
    break;
  }
}

void drawIndicator(uint8_t indiNumber, bool lit, char lbl[3])
{
  digitalWrite(indicatorLEDMap[indiNumber], lit);
}

// Clear a EPD DSP
void clearEPDDSPS()
{
  digitalWrite(DSP0_CS_PIN, false);
  digitalWrite(DSP1_CS_PIN, false);
  digitalWrite(DSP2_CS_PIN, false);
  digitalWrite(DSP3_CS_PIN, false);
  digitalWrite(DSP4_CS_PIN, false);
  digitalWrite(DSP5_CS_PIN, false);
  IndicatorDSP1.clearScreen(GxEPD_WHITE);
  digitalWrite(DSP0_CS_PIN, true);
  digitalWrite(DSP1_CS_PIN, true);
  digitalWrite(DSP2_CS_PIN, true);
  digitalWrite(DSP3_CS_PIN, true);
  digitalWrite(DSP4_CS_PIN, true);
  digitalWrite(DSP5_CS_PIN, true);

  for (int i = indicatorLEDMap.length(); i >= 0; i--)
  {
    digitalWrite(indicatorLEDMap[i], false);
  }

  Serial1.println("{\"Clear\":\"All\"}");
}

void randomizeArtDSP()
{
}

// Determins what art to draw onto the display or to send that code to the follower device
void processSideArt(uint8_t displayNum, byte artcode)
{
  DynamicJsonDocument doc(50);
  switch (displayNum) {
    case 0:
    case 1:
      // first bit equals art type Battery or Port
      // 0 = Battery
      // 1 = Port
      bool artType = artcode & 0b10000000;
      if (artType)
      {
        writePorts(displayNum, artcode);
      }
      else
      {
        // Battery
        // 0 = Battery
        // 0 = not used
        // 0 = not used
        // 0 = not used
        // 0 = not used
        // 0 = not used
        // 0 = not used
        // 1 = 1 is AA  0 is D
        bool AAorD = artcode & 0b000000001;
        writeBattery(displayNum, AAorD);
      }
      break;
    case 2:
      doc["write"]["ArtDisp2"] = artcode;
      Serial1.println(serializeJson(doc, Serial1));
      break;
    case 3:
      doc["write"]["ArtDisp3"] = artcode;
      Serial1.println(serializeJson(doc, Serial1));
      break;
  }
}

// Sends out our output buffer
void requestEvent()
{
  if (bytesToSend > 0)
  {
    Wire.write(outgoingI2CData, bytesToSend);
  }
  bytesToSend = 0;
}

// Copy the incoming data into our input buffer
void receiveEvent(int numBytes)
{
  bytesReceived = numBytes;
  for (int i = 0; i < numBytes; i++)
  {
    if (i > 10)
    {
      Wire.read();
    }
    else
    {
      incomeingI2CData[i] = Wire.read();
    }
  }
}

// Process the incoming data
void I2CCommandProcessor()
{
  switch (incomeingI2CData[0] >> 4)
  {
  case 0x1: // Set
    switch (incomeingI2CData[0] & 0xF)
    {
    case 0x0: // Set Serial Number
    {
      char serialNum[] = {
          incomeingI2CData[1],
          incomeingI2CData[2],
          incomeingI2CData[3],
          incomeingI2CData[4],
          incomeingI2CData[5],
          incomeingI2CData[6],
          incomeingI2CData[7],
          incomeingI2CData[8]};
      drawSerialNum(serialNum);
    }
    break;
    case 0x1: // Set Indicator
    {
      char indicatorlbl[] = {
          incomeingI2CData[2],
          incomeingI2CData[3],
          incomeingI2CData[4]};
      // Lit               // Label
      drawIndicator(indiNumber, incomeingI2CData[1], indicatorlbl);
    }

      indiNumber++;
      break;
    case 0x2: // Set SideArt
      processSideArt(incomeingI2CData[1]);
      break;
    }
    break;

  case 0x2: // Clear
    switch (incomeingI2CData[0] & 0xF)
    {
    case 0x0: // Clear Serial Number
      clearEPDDSP(0);
      break;
    case 0x1: // Clear Indicator
      for (int i = 0; i < indiNumber; i++)
      {
        clearIndicator(indiNumber);
      }
      indiNumber = 0;
      break;
    case 0x2: // Clear Side Art
      for (int i = 0; i < sizeof(portBatDSPMap); i++)
      {
        for (int j = 0; j < sizeof(portBatDSPMap[i]); j++)
        {
          clearEPDDSP(portBatDSPMap[i][j]);
        }
      }
      break;
    }
    break;
  case 0x3: // Other Functions
    switch (incomeingI2CData[0] & 0xF)
    {
    case 0x0: // Randomize Art DSP
      randomizeArtDSP();
      break;
    }
    break;
  }
}

void setup()
{
  // Set SPI Pins
  SPI.setTX(EPD_MOSI_PIN);
  SPI.setSCK(EPD_SCK_PIN);
  // Set S2S Serial
  Serial1.setTX(S2S_SERIAL_TX);
  Serial1.setRX(S2S_SERIAL_RX);
  Serial1.begin(S2S_SERIAL_SPEED);
  // put your setup code here, to run once:
  Serial.begin(9600);

  setupDisplays();
  Serial.println("Display initialized");

  for (int i = 0; i < sizeof(indicatorLEDMap); i++) // Set up all indicators
  {
    pinMode(indicatorLEDMap[i], OUTPUT);
    digitalWrite(indicatorLEDMap[i], LOW);
  }

  Serial.println("Display test");

  Serial.println("Drawing Ports ");
  drawPorts(0, 0b00011100);
  // delay(1000);

  Serial.println("Display Done");
  pixels->setPixelColor(0, pixels->Color(0xFF, 0xFF, 0xFF));
  pixels->show();

  delay(500);
  pixels->setPixelColor(0, pixels->Color(0, 0, 0));
  pixels->show();
  Serial.println("End Of Line");
}

void loop()
{
  // put your main code here, to run repeatedly:
}