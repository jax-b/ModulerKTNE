#include <SPI.h>
#include <GxEPD2_3C.h>
#include <Fonts/FreeMonoBold18pt7b.h>
#include <Fonts/FreeMonoBold24pt7b.h>
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

GxEPD2_290_C90c IndicatorDSP4(/*CS=*/DSP0_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/EPD_BUS_PIN); // Indicator
GxEPD2_290_C90c IndicatorDSP5(/*CS=*/DSP1_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/EPD_BUS_PIN); // Ind
GxEPD2_290_C90c ArtDSP2(/*CS=*/DSP2_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/EPD_BUS_PIN);
GxEPD2_290_C90c ArtDSP3(/*CS=*/DSP3_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/EPD_BUS_PIN);
GxEPD2_290_C90c SerialNumberDSP(/*CS=*/DSP4_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/EPD_BUS_PIN);

GFXcanvas1 BlackCanvas(IndicatorDSP4.WIDTH, IndicatorDSP4.HEIGHT);
GFXcanvas1 RedCanvas(IndicatorDSP4.WIDTH, IndicatorDSP4.HEIGHT);
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
  pinMode(EPD_BUS_PIN, INPUT);

  // Release All Displays
  digitalWrite(DSP0_CS_PIN, true);
  digitalWrite(DSP1_CS_PIN, true);
  digitalWrite(DSP2_CS_PIN, true);
  digitalWrite(DSP3_CS_PIN, true);
  digitalWrite(DSP4_CS_PIN, true);

  // Reset All Displays
  Serial.println("DSP RESET");
  digitalWrite(EPD_RESET_PIN, true);
  delay(100);
  digitalWrite(EPD_RESET_PIN, false);
  delay(100);
  digitalWrite(EPD_RESET_PIN, true);

  // Init and attempt to fill the displays buffers
  Serial.println("DSP INIT");
  IndicatorDSP4.init();
  IndicatorDSP4.writeScreenBuffer(GxEPD_WHITE);
  IndicatorDSP5.init();
  IndicatorDSP5.writeScreenBuffer(GxEPD_WHITE);
  ArtDSP2.init();
  ArtDSP2.writeScreenBuffer(GxEPD_WHITE);
  ArtDSP3.init();
  ArtDSP3.writeScreenBuffer(GxEPD_WHITE);
  SerialNumberDSP.init();
  SerialNumberDSP.writeScreenBuffer(GxEPD_WHITE);

  // Grab all dislays and clear them then release them
  digitalWrite(DSP0_CS_PIN, false);
  digitalWrite(DSP1_CS_PIN, false);
  digitalWrite(DSP2_CS_PIN, false);
  digitalWrite(DSP3_CS_PIN, false);
  digitalWrite(DSP4_CS_PIN, false);
  // IndicatorDSP0.clearScreen(GxEPD_WHITE);
  digitalWrite(DSP0_CS_PIN, true);
  digitalWrite(DSP1_CS_PIN, true);
  digitalWrite(DSP2_CS_PIN, true);
  digitalWrite(DSP3_CS_PIN, true);
  digitalWrite(DSP4_CS_PIN, true);

  delay(10);
}

void drawDisplays()
{
  digitalWrite(DSP0_CS_PIN, false);
  digitalWrite(DSP1_CS_PIN, false);
  digitalWrite(DSP2_CS_PIN, false);
  digitalWrite(DSP3_CS_PIN, false);
  digitalWrite(DSP4_CS_PIN, false);
  IndicatorDSP4.refresh();
  digitalWrite(DSP0_CS_PIN, true);
  digitalWrite(DSP1_CS_PIN, true);
  digitalWrite(DSP2_CS_PIN, true);
  digitalWrite(DSP3_CS_PIN, true);
  digitalWrite(DSP4_CS_PIN, true);
}

// Draws a Battery to the specified display
// True is AA false is D
void writeBattery(uint8_t displayNum, bool AAorD)
{
  // Code needs to be updated to support multiple displays
  // displays[displayNum].fillScreen(GxEPD_WHITE);
  if (displayNum == 2 || displayNum == 3)
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
  case 2:
    ArtDSP2.writeScreenBuffer(GxEPD_WHITE);
    ArtDSP2.writeImage(BlackCanvas.getBuffer(), RedCanvas.getBuffer(), 0, 0, 128, 256, true, false, false);
    break;
  case 3:
    ArtDSP3.writeScreenBuffer(GxEPD_WHITE);
    ArtDSP3.writeImage(BlackCanvas.getBuffer(), RedCanvas.getBuffer(), 0, 0, 128, 256, true, false, false);
    break;
  }
}

// Draws the SerialNumber to dsp 0
void writeSerialNum(String inText)
{
  inText.trim();
  BlackCanvas.fillScreen(GxEPD_WHITE);
  RedCanvas.fillScreen(GxEPD_WHITE);
  BlackCanvas.setTextColor(GxEPD_BLACK);
  RedCanvas.fillRect(0, 7, 250, 35, GxEPD_BLACK);
  RedCanvas.fillRect(0, 94, 250, 35, GxEPD_BLACK);
  BlackCanvas.setTextSize(1);
  BlackCanvas.setFont(&FreeMonoBold24pt7b);
  int16_t tbx, tby;
  uint16_t tbw, tbh;
  BlackCanvas.getTextBounds(inText, 0, 0, &tbx, &tby, &tbw, &tbh);
  uint16_t x = ((BlackCanvas.width() - 45 - tbw) / 2) - tbx;
  uint16_t y = ((BlackCanvas.height() + 6 - tbh) / 2) - tby;
  BlackCanvas.setCursor(x, y);
  BlackCanvas.print(inText);
  SerialNumberDSP.writeScreenBuffer(GxEPD_WHITE);
  SerialNumberDSP.writeImage(BlackCanvas.getBuffer(), RedCanvas.getBuffer(), 0, 0, 128, 256, false, false, false);
}

// // Draws the SerialNumber to dsp 0
// // Overload for a char array instead of a string
// void writeSerialNum(string)
// {
//   char strtext[] = intext;
//   writeSerialNum(strtext);
// }

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

  if (displayNum == 2 || displayNum == 3)
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
    case 2:
      ArtDSP2.writeScreenBuffer(GxEPD_WHITE);
      ArtDSP2.writeImage(BlackCanvas.getBuffer(), RedCanvas.getBuffer(), 0, 0, 122, 256, false, false, true);
      break;
    case 3:
      ArtDSP3.writeScreenBuffer(GxEPD_WHITE);
      ArtDSP3.writeImage(BlackCanvas.getBuffer(), RedCanvas.getBuffer(), 0, 0, 122, 256, false, false, true);
      break;
    }
  }
}

void writeIndicator(uint8_t indiNumber, bool lit, char lbl[3])
{
  if (indiNumber == 4 || indiNumber == 5)
  {
    BlackCanvas.fillScreen(GxEPD_WHITE);
    BlackCanvas.setCursor(-2, 98);
    BlackCanvas.setTextColor(GxEPD_BLACK);
    BlackCanvas.setFont(&FreeMonoBold18pt7b);
    BlackCanvas.setTextSize(4);
    BlackCanvas.print(lbl);
  }
  switch (indiNumber)
  {
  case 4:
    IndicatorDSP4.writeScreenBuffer(GxEPD_WHITE);
    IndicatorDSP4.writeImage(BlackCanvas.getBuffer(), NULL, 0, 0, 122, 256, false, false, false);
    digitalWrite(DSP0_LED_PIN, lit);
    break;
  case 5:
    IndicatorDSP5.writeScreenBuffer(GxEPD_WHITE);
    IndicatorDSP5.writeImage(BlackCanvas.getBuffer(), NULL, 0, 0, 122, 256, false, false, false);
    digitalWrite(DSP1_LED_PIN, lit);
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
  IndicatorDSP4.clearScreen(GxEPD_WHITE);
  digitalWrite(DSP0_CS_PIN, true);
  digitalWrite(DSP1_CS_PIN, true);
  digitalWrite(DSP2_CS_PIN, true);
  digitalWrite(DSP3_CS_PIN, true);
  digitalWrite(DSP4_CS_PIN, true);
}

// Determins what art to draw onto the display or to send that code to the follower device
void processSideArt(uint8_t displayNum, byte artcode)
{
  switch (displayNum)
  {
  case 2:
  case 3:
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
  }
}


enum command_type
{
  cmd_startup,
  cmd_clear,
  cmd_draw,
  cmd_indicator,
  cmd_art,
  cmd_serial,
  cmd_bad
};

command_type commandHash(String const &inString)
{
  if (inString == "startup")
    return cmd_startup;
  if (inString == "clear")
    return cmd_clear;
  if (inString == "draw")
    return cmd_draw;
  if (inString == "serial")
    return cmd_serial;
  if (inString == "art")
    return cmd_art;
  if (inString == "indicator")
    return cmd_indicator;
  return cmd_bad;
}


void setup()
{
  Serial1.setTX(S2S_SERIAL_TX);
  Serial1.setRX(S2S_SERIAL_RX);
  Serial1.begin(S2S_SERIAL_SPEED);
}

void setup1()
{
  // Set SPI Pins
  SPI.setTX(EPD_MOSI_PIN);
  SPI.setSCK(EPD_SCK_PIN);
  // Set S2S Serial

  BlackCanvas.setRotation(1);
  RedCanvas.setRotation(1);
}

DynamicJsonDocument doc(1024);
bool newMessage = false;
void loop()
{
  while (Serial1.available())
  {
    String SerialBuffer = Serial1.readStringUntil('\n');
    deserializeJson(doc, SerialBuffer);
    newMessage = true;
  }
  if (newMessage)
  {
    newMessage = false;
    String command = doc["command"];
    uint8_t display = doc["display"];
    switch (commandHash(command))
    {
    case cmd_startup:
      rp2040.fifo.push(0x00000000); // Startup
      break;
    case cmd_clear:
      rp2040.fifo.push(0x01000000); // Clear
      break;
    case cmd_draw:
      rp2040.fifo.push(0x02000000); // Draw
      break;
    case cmd_indicator:
      {
        uint32_t command1 = 0x03000000 & (display << 16); // Construct Command
        bool lit = doc["lit"];
        char lbl[3] = {doc["lbl"][0], doc["lbl"][1], doc["lbl"][2]};
        uint32_t indicatorValue = (lit << 24) & (lbl[0] << 16) & (lbl[1] << 8) & lbl[2]; // Construct Indicator Value
        rp2040.fifo.push(command1);                                                      // Send Command
        rp2040.fifo.push(indicatorValue);                                                // Send Indicator Value
      }
      break;
    case cmd_art:
      {
        uint8_t artValue = doc["artcode"];                                  // Get Artcode
        uint32_t command1 = 0x04000000 & (display << 16) & (artValue << 8); // Construct command
        rp2040.fifo.push(command1);                                         // Send Command
      }
      break;
    case cmd_serial:
      {
        String serialValue = doc["serial"];                                                                                                 // Get Serial
        uint32_t command1 = 0x05000000 & (serialValue.charAt(0));                                                                           // Construct command
        uint32_t data1 = (serialValue.charAt(1)) << 24 & (serialValue.charAt(2) << 16) & (serialValue.charAt(3)) & (serialValue.charAt(4)); // Construct data
        uint32_t data2 = (serialValue.charAt(5)) << 24 & (serialValue.charAt(6) << 16) & (serialValue.charAt(7)) & (serialValue.charAt(8)); // Construct data
        rp2040.fifo.push(command1);                                                                                                         // Send Command
        rp2040.fifo.push(data1);                                                                                                            // Send Data
        rp2040.fifo.push(data2);                                                                                                            // Send Data
      }
      break;
    }
  }
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