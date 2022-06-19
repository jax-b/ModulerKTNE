#include <Adafruit_NeoPixel.h>
#include <SPI.h>
#include <GxEPD2_3C.h>
#include <Fonts/FreeMonoBold12pt7b.h>
#include <Adafruit_I2CDevice.h>
#include <Ardiuno.h>

// #include "BatteryPictures.h"
// #include "PortPictures.h"

/// Buffers and data tracking for I2C communication
byte incomeingI2CData[10];
byte outgoingI2CData[10];
uint8_t bytesToSend = 0;
uint8_t bytesReceived = 0;

#define EPD_RESET_PIN 5 // can set to -1 and share with microcontroller Reset!
#define EPD_DC_PIN 3
#define DSP1_CS_PIN 6
#define DSP1_BUS_PIN 7
#define DSP1_LED_PIN 8
#define DSP2_CS_PIN 9
#define DSP2_BUS_PIN 10
#define DSP2_LED_PIN 11
#define DSP3_CS_PIN 12
#define DSP3_BUS_PIN 13
#define DSP4_CS_PIN 14
#define DSP4_BUS_PIN 15
#define DSP5_CS_PIN 16
#define DSP5_BUS_PIN 17
#define DSP5_LED_PIN 18
#define DSP6_CS_PIN 19
#define DSP6_BUS_PIN 20
#define DSP6_LED_PIN 21

GxEPD2_3C<GxEPD2_213c, 50> IndicatorDSP1(GxEPD2_213c(/*CS=*/DSP1_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/DSP1_BUS_PIN)); //Indicator
GxEPD2_3C<GxEPD2_213c, 50> IndicatorDSP2(GxEPD2_213c(/*CS=*/DSP2_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/DSP2_BUS_PIN)); //Ind
GxEPD2_3C<GxEPD2_213c, 50> ArtDSP1(GxEPD2_213c(/*CS=*/DSP3_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/DSP3_BUS_PIN));
GxEPD2_3C<GxEPD2_213c, 50> ArtDSP2(GxEPD2_213c(/*CS=*/DSP4_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/DSP4_BUS_PIN));
GxEPD2_3C<GxEPD2_213c, 50> IndicatorDSP3(GxEPD2_213c(/*CS=*/DSP5_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/DSP5_BUS_PIN));
GxEPD2_3C<GxEPD2_213c, 50> IndicatorDSP4(GxEPD2_213c(/*CS=*/DSP6_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/DSP6_BUS_PIN));


// Tracking Variables
uint8_t indiNumber = 0;
const uint8_t ARTDSPRANGE[] = {1, 5};
const uint8_t MAXDSPNUM = ARTDSPRANGE[1];
uint8_t batDSPDrawCycle = 0;
uint8_t portDSPDrawCycle = 0;
uint8_t portBatDSPMap[2][5] = {
    {1, 2, 0, 0, 0}, // Port DSPS
    {3, 4, 5, 6, 7}  // Bat DSPS
};

uint8_t indicatorLEDMap[MAXINDICATOR] = {
  INDICATOR1PIN, INDICATOR2PIN, INDICATOR3PIN, INDICATOR4PIN, INDICATOR5PIN, INDICATOR6PIN, INDICATOR7PIN, INDICATOR8PIN
};

void setupDisplay(uint8_t DisplayNum)
{
  display1.init();
  display1.setRotation(1);
  display1.setFont(&FreeMonoBold12pt7b);
  display1.setFullWindow();
  display1.setTextColor(GxEPD_BLACK);
}

// Draws a Battery to the specified display
// True is AA false is D
void drawBattery(uint8_t displayNum, bool AAorD)
{
  // Code needs to be updated to support multiple displays
  // displays[displayNum].fillScreen(GxEPD_WHITE);

  if (AAorD)
  {
    display1.firstPage();
    do
    {
      display1.fillScreen(GxEPD_WHITE);
      display1.drawInvertedBitmap(0, 6, epd_bitmap_AA, 250, 122, GxEPD_BLACK);
    } while (display1.nextPage());
  }
  else
  {
    display1.firstPage();
    do
    {
      display1.fillScreen(GxEPD_WHITE);
      display1.drawInvertedBitmap(0, 6, epd_bitmap_D, 250, 122, GxEPD_BLACK);
    } while (display1.nextPage());
  }
}

// Draws the SerialNumber to dsp 0
void drawSerialNum(String text)
{
  // Code needs to be updated to support multiple displays
  // displays[0].fillScreen(GxEPD_WHITE);
  int16_t x1, y1;
  uint16_t w, h;
  display1.getTextBounds(text, 0, 0, &x1, &y1, &w, &h);
  uint16_t x = ((display1.width() - w) / 2) - x1;
  uint16_t y = ((display1.height() - h) / 2) - y1;
  display1.firstPage();
  do
  {
    display1.fillScreen(GxEPD_WHITE);
    display1.setCursor(x, y);
    display1.print(text);
  } while (display1.nextPage());
}

// Draws the SerialNumber to dsp 0
// Overload for a char array instead of a string
void drawSerialNum(char intext[])
{
  String strtext = intext;
  drawSerialNum(strtext);
}

// Draws a port to the specified display
// Host is responsible for avoiding collisions, bad data will not be drawn but there might be a missing port caused by a collision
void drawPorts(uint8_t displayNum, byte activeports)
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
  display1.firstPage();
  const uint8_t VAL_TO_ADD = 67;
  do
  {
    display1.fillScreen(GxEPD_WHITE);
    uint8_t longporty = 6;
    uint8_t shortportx = 0;
    if (rj45Active)
    {
      display1.drawInvertedBitmap(shortportx, 6, epd_bitmap_rj45, epd_bitmap_rj45_size[0], epd_bitmap_rj45_size[1], GxEPD_BLACK);
      longporty += VAL_TO_ADD;
      shortportx += epd_bitmap_rj45_size[0] + 5;
    }
    if (ps2Active)
    {
      display1.drawInvertedBitmap(shortportx, 6, epd_bitmap_PS2, epd_bitmap_PS2_size[0], epd_bitmap_PS2_size[1], GxEPD_BLACK);
      if (longporty < 30)
      {
        longporty += VAL_TO_ADD;
      }
      shortportx += epd_bitmap_PS2_size[0] + 5;
    }
    if (rcaActive)
    {
      display1.drawInvertedBitmap(213, 6, epd_bitmap_RCA, epd_bitmap_RCA_size[0], epd_bitmap_RCA_size[1], GxEPD_BLACK);
    }

    if (dviActive)
    {
      display1.drawInvertedBitmap(0, longporty, epd_bitmap_DVI, epd_bitmap_DVI_size[0], epd_bitmap_DVI_size[1], GxEPD_BLACK);
      if (longporty > 30)
      {
        continue;
      }
      else
      {
        longporty += VAL_TO_ADD;
      }
    }
    if (serialActive)
    {
      display1.drawInvertedBitmap(0, longporty, epd_bitmap_Serial, epd_bitmap_Serial_size[0], epd_bitmap_Serial_size[1], GxEPD_BLACK);
      if (longporty > 30)
      {
        continue;
      }
      else
      {
        longporty += VAL_TO_ADD;
      }
    }
    if (parallelActive && !rcaActive)
    {
      display1.drawInvertedBitmap(0, longporty, epd_bitmap_Parallel, epd_bitmap_Parallel_size[0], epd_bitmap_Parallel_size[1], GxEPD_BLACK);
    }
  } while (display1.nextPage());
}

void drawIndicator(uint8_t indiNumber, bool lit, char lbl[3])
{
  digitalWrite(indicatorLEDMap[indiNumber], lit);
}

// Clear a EPD DSP
void clearEPDDSP(uint8_t displayNum)
{
  // Code needs to be updated to support multiple displays
  // displays[0].fillScreen(GxEPD_WHITE);
  display1.firstPage();
  do
  {
    display1.fillScreen(GxEPD_WHITE);
  } while (display1.nextPage());
}

void clearIndicator(uint8_t indiNumber)
{
  digitalWrite(indicatorLEDMap[indiNumber], false);
}

void randomizeArtDSP()
{
  
}

// Determans what art to draw onto the display
void processSideArt(byte artcode)
{
  // first bit equals art type Battery or Port
  // 0 = Battery
  // 1 = Port
  bool artType = artcode & 0b10000000;
  if (artType)
  {
    drawPorts(portBatDSPMap[0][portDSPDrawCycle], artcode);
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
    drawBattery(portBatDSPMap[1][batDSPDrawCycle], AAorD);
    batDSPDrawCycle++;
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
  SPI.setTX(19);
  SPI.setSCK(18);
  // put your setup code here, to run once:
  Serial.begin(115200);
  // while (!Serial)
  ; // wait for serial port to connect. Needed for native USB port only
  delay(10);

  setupDisplay(0);
  Serial.println("Display initialized");

  for (int i = 0; i < sizeof(indicatorLEDMap); i++) // Set up all indicators
  {
    pinMode(indicatorLEDMap[i], OUTPUT);
    digitalWrite(indicatorLEDMap[i], LOW);
  }

  pixels = new Adafruit_NeoPixel(1, 16, NEO_GRB + NEO_KHZ800);
  pixels->begin();
  pixels->setPixelColor(0, pixels->Color(255, 0, 0));
  pixels->show();
  delay(100);
  pixels->setPixelColor(0, pixels->Color(0xAA, 0x0, 0xAA));
  pixels->show();

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