#include <SPI.h>
#include <GxEPD2_3C.h>
#include <Fonts/FreeMonoBold18pt7b.h>
#include <ArduinoJson.h>
#include <Adafruit_GFX.h>
#include "BatteryPictures.h"
#include "PortPictures.h"

/// Buffers and data tracking for I2C communication
byte incomeingI2CData[15];
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
#define CONTROLLER_I2C_SCL 1
#define CONTROLLER_I2C_SDA 0
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
#define CONTROLLER_I2C_SCL 1
#define CONTROLLER_I2C_SDA 0
#endif

#define I2C_ADDRESS 0x50
#define S2S_SERIAL_SPEED 115200

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
uint8_t artNumber = 0;
const uint8_t ARTDSPRANGE[] = {1, 5};
const uint8_t MAXDSPNUM = ARTDSPRANGE[1];
uint8_t batDSPDrawCycle = 0;
uint8_t portDSPDrawCycle = 0;
uint8_t indiMap[] = {0, 1, 2, 3, 4, 5};
uint8_t artMap[] = {0, 1, 2, 3};

// Initializes and clears the connected displays
void setupDisplays()
{
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

  // Setup LED Pins
  pinMode(DSP0_LED_PIN, OUTPUT);
  pinMode(DSP1_LED_PIN, OUTPUT);
  pinMode(DSP4_LED_PIN, OUTPUT);
  pinMode(DSP5_LED_PIN, OUTPUT);

  // Release All Displays
  digitalWrite(DSP0_CS_PIN, true);
  digitalWrite(DSP1_CS_PIN, true);
  digitalWrite(DSP2_CS_PIN, true);
  digitalWrite(DSP3_CS_PIN, true);
  digitalWrite(DSP4_CS_PIN, true);
  digitalWrite(DSP5_CS_PIN, true);

  // Reset All Displays
  // Serial.println("DSP RESET");
  digitalWrite(EPD_RESET_PIN, true);
  delay(100);
  digitalWrite(EPD_RESET_PIN, false);
  delay(100);
  digitalWrite(EPD_RESET_PIN, true);

  // Init and attempt to fill the displays buffers
  // Serial.println("DSP INIT");
  IndicatorDSP0.init();
  IndicatorDSP0.writeScreenBuffer(GxEPD_WHITE);
  IndicatorDSP1.init();
  IndicatorDSP1.writeScreenBuffer(GxEPD_WHITE);
  ArtDSP0.init();
  ArtDSP0.writeScreenBuffer(GxEPD_WHITE);
  ArtDSP1.init();
  ArtDSP1.writeScreenBuffer(GxEPD_WHITE);
  IndicatorDSP2.init();
  IndicatorDSP2.writeScreenBuffer(GxEPD_WHITE);
  IndicatorDSP3.init();
  IndicatorDSP3.writeScreenBuffer(GxEPD_WHITE);

  // Grab all dislays and clear them then release them
  digitalWrite(DSP0_CS_PIN, false);
  digitalWrite(DSP1_CS_PIN, false);
  digitalWrite(DSP2_CS_PIN, false);
  digitalWrite(DSP3_CS_PIN, false);
  digitalWrite(DSP4_CS_PIN, false);
  digitalWrite(DSP5_CS_PIN, false);
  // IndicatorDSP0.clearScreen(GxEPD_WHITE);
  digitalWrite(DSP0_CS_PIN, true);
  digitalWrite(DSP1_CS_PIN, true);
  digitalWrite(DSP2_CS_PIN, true);
  digitalWrite(DSP3_CS_PIN, true);
  digitalWrite(DSP4_CS_PIN, true);
  digitalWrite(DSP5_CS_PIN, true);

  delay(10);
}

void drawDisplays()
{
  digitalWrite(DSP0_CS_PIN, false);
  digitalWrite(DSP1_CS_PIN, false);
  digitalWrite(DSP2_CS_PIN, false);
  digitalWrite(DSP3_CS_PIN, false);
  digitalWrite(DSP4_CS_PIN, false);
  digitalWrite(DSP5_CS_PIN, false);
  IndicatorDSP0.refresh();
  digitalWrite(DSP0_CS_PIN, true);
  digitalWrite(DSP1_CS_PIN, true);
  digitalWrite(DSP2_CS_PIN, true);
  digitalWrite(DSP3_CS_PIN, true);
  digitalWrite(DSP4_CS_PIN, true);
  digitalWrite(DSP5_CS_PIN, true);
}

// Draws a Battery to the specified display
// True is AA false is D
void writeBattery(uint8_t displayNum, bool AAorD)
{
  // Code needs to be updated to support multiple displays
  // displays[displayNum].fillScreen(GxEPD_WHITE);
  if (displayNum < 2)
  {
    BlackCanvas.fillScreen(GxEPD_WHITE);
    RedCanvas.fillScreen(GxEPD_WHITE);
    if (AAorD)
    {
      BlackCanvas.drawBitmap(0, 6, epd_bitmap_AA_Black, epd_bitmap_AA_size[0], epd_bitmap_AA_size[1], GxEPD_BLACK);
      RedCanvas.drawBitmap(0, 6, epd_bitmap_AA_Red, epd_bitmap_AA_size[0], epd_bitmap_AA_size[1], GxEPD_BLACK);
    }
    else
    {
      BlackCanvas.drawBitmap(0, 6, epd_bitmap_D_Black, epd_bitmap_D_size[0], epd_bitmap_D_size[1], GxEPD_BLACK);
      RedCanvas.drawBitmap(0, 6, epd_bitmap_D_Red, epd_bitmap_D_size[0], epd_bitmap_D_size[1], GxEPD_BLACK);
    }
  }

  switch (displayNum)
  {
  case 0:
    ArtDSP0.writeScreenBuffer(GxEPD_WHITE);
    ArtDSP0.writeImage(BlackCanvas.getBuffer(), RedCanvas.getBuffer(), 0, 0, 128, 256, true, false, false);
    break;
  case 1:
    ArtDSP1.writeScreenBuffer(GxEPD_WHITE);
    ArtDSP1.writeImage(BlackCanvas.getBuffer(), RedCanvas.getBuffer(), 0, 0, 128, 256, true, false, false);
    break;
  }
}

// Draws the SerialNumber to dsp 0
void writeSerialNum(String text)
{
  DynamicJsonDocument doc(60);
  doc["command"] = "serial";
  doc["serial"] = text;
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
void writePort(uint8_t displayNum, byte activeports)
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

  if (displayNum < 2)
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
      BlackCanvas.drawBitmap(198, 6, epd_bitmap_RCA_Black, epd_bitmap_RCA_size[0], epd_bitmap_RCA_size[1], GxEPD_BLACK);
      RedCanvas.drawBitmap(198, 6, epd_bitmap_RCA_Red, epd_bitmap_RCA_size[0], epd_bitmap_RCA_size[1], GxEPD_BLACK);
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
    switch (displayNum)
    {
    case 0:
      ArtDSP0.writeScreenBuffer(GxEPD_WHITE);
      ArtDSP0.writeImage(BlackCanvas.getBuffer(), RedCanvas.getBuffer(), 0, 0, 128, 256, false, false, false);
      break;
    case 1:
      ArtDSP1.writeScreenBuffer(GxEPD_WHITE);
      ArtDSP1.writeImage(BlackCanvas.getBuffer(), RedCanvas.getBuffer(), 0, 0, 128, 256, false, false, false);
      break;
    }
  }
}

void writeIndicator(uint8_t indiNumber, bool lit, char *lbl)
{
  if (indiNumber < 4)
  {
    BlackCanvas.fillScreen(GxEPD_WHITE);
    BlackCanvas.setCursor(-2, 98);
    BlackCanvas.setTextColor(GxEPD_BLACK);
    BlackCanvas.setFont(&FreeMonoBold18pt7b);
    BlackCanvas.setTextSize(4);
    BlackCanvas.print(lbl);
  }
  DynamicJsonDocument doc(60);
  switch (indiNumber)
  {
  case 0:
    IndicatorDSP0.writeScreenBuffer(GxEPD_WHITE);
    IndicatorDSP0.writeImage(BlackCanvas.getBuffer(), NULL, 0, 0, 128, 256, false, false, false);
    digitalWrite(DSP0_LED_PIN, lit);
    break;
  case 1:
    IndicatorDSP1.writeScreenBuffer(GxEPD_WHITE);
    IndicatorDSP1.writeImage(BlackCanvas.getBuffer(), NULL, 0, 0, 128, 256, false, false, false);
    digitalWrite(DSP1_LED_PIN, lit);
    break;
  case 2:
    IndicatorDSP2.writeScreenBuffer(GxEPD_WHITE);
    IndicatorDSP2.writeImage(BlackCanvas.getBuffer(), NULL, 0, 0, 128, 256, false, false, false);
    digitalWrite(DSP4_LED_PIN, lit);
    break;
  case 3:
    IndicatorDSP3.writeScreenBuffer(GxEPD_WHITE);
    IndicatorDSP3.writeImage(BlackCanvas.getBuffer(), NULL, 0, 0, 128, 256, false, false, false);
    digitalWrite(DSP5_LED_PIN, lit);
    break;
  case 4:
    doc["command"] = "indicator";
    doc["display"] = 4;
    doc["label"] = lbl;
    doc["lit"] = lit;
    Serial1.println(serializeJson(doc, Serial1));
    break;
  case 5:
    doc["command"] = "indicator";
    doc["display"] = 5;
    doc["label"] = lbl;
    doc["lit"] = lit;
    Serial1.println(serializeJson(doc, Serial1));
    break;
  }
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

  digitalWrite(DSP0_LED_PIN, false);
  digitalWrite(DSP1_LED_PIN, false);
  digitalWrite(DSP4_LED_PIN, false);
  digitalWrite(DSP5_LED_PIN, false);

  Serial1.println("{\"Clear\":\"All\"}");
}

void randomizeDSPS()
{
  for (int i = 6; i > 0; i--) {
    indiMap[i-1] = 10;
  } 
  for (int i = 6; i > 0; i--) {
    bool Stashed = false;
    while (!Stashed) {
      int r = random(4);
      if (indiMap[r] == 10) {
        indiMap[r] = i-1;
        Stashed = true;
      }
    }
  }
  for (int i = 4; i > 0; i--) {
    artMap[i-1] = 10;
  }
  for (int i = 4; i > 0; i--) {
    bool Stashed = false;
    while (!Stashed) {
      int r = random(4);
      if (artMap[r] == 10) {
        artMap[r] = i-1;
        Stashed = true;
      }
    }
  }
}

// Determins what art to draw onto the display or to send that code to the follower device
void processSideArt(uint8_t displayNum, byte artcode)
{
  DynamicJsonDocument doc(60);
  switch (displayNum)
  {
  case 0:
  case 1:
    // first bit equals art type Battery or Port
    // 0 = Battery
    // 1 = Port
    {
      bool artType = artcode & 0b10000000;
      if (artType)
      {
        writePort(displayNum, artcode);
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
    }
    break;
  case 2:
    doc["command"] = "art";
    doc["display"] = 2;
    doc["artcode"] = artcode;
    Serial1.println(serializeJson(doc, Serial1));
    break;
  case 3:
    doc["command"] = "art";
    doc["display"] = 3;
    doc["artcode"] = artcode;
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
        uint32_t command1 = 0x05000000 & (incomeingI2CData[1]);                                                                     // Construct command
        uint32_t data1 = (incomeingI2CData[2]) << 24 & (incomeingI2CData[3] << 16) & (incomeingI2CData[4] << 8) & (incomeingI2CData[5]); // Construct data
        uint32_t data2 = (incomeingI2CData[6]) << 24 & (incomeingI2CData[7] << 16) & (incomeingI2CData[8]) & (incomeingI2CData[9]); // Construct data
        rp2040.fifo.push(command1);                                                                                                 // Send Command
        rp2040.fifo.push(data1);                                                                                                    // Send Data
        rp2040.fifo.push(data2);                                                                                                    // Send Data
      }
      break;
    case 0x1: // Set Indicator
      {
        uint32_t command1 = 0x03000000 & (indiMap[indiNumber] << 16); // Construct Command
        uint32_t indicatorValue = (incomeingI2CData[1] << 24) & (incomeingI2CData[2] << 16) & (incomeingI2CData[3] << 8) & incomeingI2CData[4]; // Construct Indicator Value
        rp2040.fifo.push(command1);                                                                                                 // Send Command
        rp2040.fifo.push(indicatorValue);                                                                                           // Send Indicator Value
        indiNumber++;
      }
      break;
    case 0x2: // Set SideArt
      {
        uint32_t command1 = 0x04000000 & (artMap[artNumber] << 16) & (incomeingI2CData[1] << 8) ; // Construct command
        rp2040.fifo.push(command1); // Send Command
        artNumber++;
      }
      break;
    }
    break;

  case 0x0: // Other Functions
    switch (incomeingI2CData[0] & 0xF)
    {
    case 0x0: // Randomize Art DSP
      randomizeDSPS();
      break;
    case 0x1:
      rp2040.fifo.push(0x02000000); //Draw
      break;
    case 0x2:
      rp2040.fifo.push(0x01000000); //Clear
      artNumber = 0;
      indiNumber = 0;
      break;
    }
  }
}

void setup()
{
  // Seed random
  randomSeed(analogRead(26));

  // Start I2C Communications 
  Wire.setSCL(CONTROLLER_I2C_SCL);
  Wire.setSDA(CONTROLLER_I2C_SDA);
  Wire.begin(I2C_ADDRESS);
  Wire.onReceive(receiveEvent);
  Wire.onRequest(requestEvent);
}

void setup1() {
  // Set SPI Pins
  SPI.setTX(EPD_MOSI_PIN);
  SPI.setSCK(EPD_SCK_PIN);
 // Set S2S Serial
  Serial1.setTX(S2S_SERIAL_TX);
  Serial1.setRX(S2S_SERIAL_RX);
  Serial1.begin(S2S_SERIAL_SPEED);

  BlackCanvas.setRotation(1);
  RedCanvas.setRotation(1);

  // Tell remote to startup
  Serial1.println("{\"command\":\"startup\"}");
  // Startup Own Displays
  setupDisplays();
}

void loop()
{
  I2CCommandProcessor();
}

void loop1()
{
  if (rp2040.fifo.available())
  {
    uint32_t data = rp2040.fifo.pop();
    uint8_t command = (data & 0xFF000000) >> 24;
    uint8_t displayNum = (data & 0x00FF0000) >> 16;
    uint8_t artcode = (data & 0x0000FF00) >> 8;
    uint8_t sparebyte = (data & 0x000000FF);
    switch (command)
    {
    case 0:
      setupDisplays();
      break;
    case 1:
      clearEPDDSPS();
      break;
    case 2:
      drawDisplays();
      break;
    case 3:
      if (rp2040.fifo.available())
      {
        uint32_t data = rp2040.fifo.pop();
        bool lit = (data & 0xFF000000) >> 24;
        char letter1 = (data & 0x00FF0000) >> 16;
        char letter2 = (data & 0x0000FF00) >> 8;
        char letter3 = data & 0x000000FF;
        char lbl[3] = {letter1, letter2, letter3};
        writeIndicator(displayNum, lit, lbl);
      }
      break;
    case 4:
      processSideArt(displayNum, artcode);
      break;
    case 5:
      if (rp2040.fifo.available())
      {
        uint32_t data = rp2040.fifo.pop();
        uint32_t data1 = rp2040.fifo.pop();
        char serialnum[10] = {
            sparebyte,
            (data & 0xFF000000) >> 24,
            (data & 0x00FF0000) >> 16,
            (data & 0x0000FF00) >> 8,
            (data & 0x000000FF),
            (data1 & 0xFF000000) >> 24,
            (data1 & 0x00FF0000) >> 16,
            (data1 & 0x0000FF00) >> 8,
            (data1 & 0x000000FF),
            0x0};
        String serial = String(*serialnum);
        writeSerialNum(serial);
      }
      break;
    }
  }
}