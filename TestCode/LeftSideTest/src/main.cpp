// #include <Ardiuno.h>
// #include <SPI.h>
#include <Adafruit_GFX.h>
#include <GxEPD2_3C.h> 
#include "bitmaps/Bitmaps104x212.h"
#include "bitmaps/Bitmaps3c104x212.h"
#define EPD_RESET_PIN 5
#define EPD_DC_PIN 4
#define EPD_MOSI_PIN 3
#define EPD_SCK_PIN 2
#define EPD_BUS_PIN 18
#define DSP1_CS_PIN 6
#define DSP1_LED_PIN 7
#define DSP2_CS_PIN 8
#define DSP2_LED_PIN 9
#define DSP3_CS_PIN 10
#define DSP4_CS_PIN 11
#define DSP5_CS_PIN 12
#define DSP5_LED_PIN 13
#define DSP6_CS_PIN 14
#define DSP6_LED_PIN 15

#define DSP_Master true

GxEPD2_3C<GxEPD2_290_C90c, 1> IndicatorDSP1(GxEPD2_290_C90c(/*CS=*/DSP1_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/EPD_BUS_PIN)); //Indicator
GxEPD2_3C<GxEPD2_290_C90c, 1> IndicatorDSP2(GxEPD2_290_C90c(/*CS=*/DSP2_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/EPD_BUS_PIN)); //Ind
GxEPD2_3C<GxEPD2_290_C90c, 1> ArtDSP1(GxEPD2_290_C90c(/*CS=*/DSP3_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/EPD_BUS_PIN));
GxEPD2_3C<GxEPD2_290_C90c, 1> ArtDSP2(GxEPD2_290_C90c(/*CS=*/DSP4_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/EPD_BUS_PIN));
GxEPD2_3C<GxEPD2_290_C90c, 1> IndicatorDSP3(GxEPD2_290_C90c(/*CS=*/DSP5_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/EPD_BUS_PIN));
#ifdef DSP_Master
GxEPD2_3C<GxEPD2_290_C90c, 1> IndicatorDSP4(GxEPD2_290_C90c(/*CS=*/DSP6_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/EPD_BUS_PIN));
#endif

GFXcanvas1 canvas(250, 122); 

#define MAX_HEIGHT_3C(EPD) (120)
#define pgm_spisettings SPISettings(4000000, MSBFIRST, SPI_MODE0)
void dspinit() {
    pinMode(EPD_RESET_PIN, OUTPUT);
    pinMode(EPD_DC_PIN, OUTPUT);
    pinMode(DSP1_CS_PIN, OUTPUT);
    pinMode(DSP2_CS_PIN, OUTPUT);
    pinMode(DSP3_CS_PIN, OUTPUT);
    pinMode(DSP4_CS_PIN, OUTPUT);
    pinMode(DSP5_CS_PIN, OUTPUT);
    pinMode(DSP6_CS_PIN, OUTPUT);
    pinMode(EPD_BUS_PIN, INPUT);

    digitalWrite(DSP1_CS_PIN, false);
    digitalWrite(DSP2_CS_PIN, false);
    digitalWrite(DSP3_CS_PIN, false);
    digitalWrite(DSP4_CS_PIN, false);
    digitalWrite(DSP5_CS_PIN, false);
    digitalWrite(DSP6_CS_PIN, false);
    

    Serial.println("DSP RESET");
    digitalWrite(EPD_RESET_PIN, true);
    delay(100);
    digitalWrite(EPD_RESET_PIN, false);
    delay(100);
    digitalWrite(EPD_RESET_PIN, true);

    Serial.println("DSP INIT");

    IndicatorDSP1.init();
    IndicatorDSP1.writeScreenBuffer(0xFF);
    IndicatorDSP2.writeScreenBuffer(0xFF);
    ArtDSP1.writeScreenBuffer(0xFF);
    ArtDSP2.writeScreenBuffer(0xFF);
    IndicatorDSP3.writeScreenBuffer(0xFF);
    #ifdef DSP_Master
    IndicatorDSP4.writeScreenBuffer(0xFF);
    #endif
    // Release Displays
    digitalWrite(DSP1_CS_PIN, true);
    digitalWrite(DSP2_CS_PIN, true);
    digitalWrite(DSP3_CS_PIN, true);
    digitalWrite(DSP4_CS_PIN, true);
    digitalWrite(DSP5_CS_PIN, true);
    digitalWrite(DSP6_CS_PIN, true);
    delay(10);
}

void setup(){
    // MbedSPI SPI(0, EPD_MOSI_PIN, EPD_SCK_PIN);
    SPI.setTX(EPD_MOSI_PIN);
    SPI.setSCK(EPD_SCK_PIN);
    Serial.begin(9600);
    Serial.println("Start");
    for (int i = 0; i < 10; i++) {
        Serial.println(".");
        delay(100);
    }
    delay(1000);
    Serial.println("Start");

    
    pinMode(DSP6_LED_PIN, OUTPUT);
    pinMode(DSP1_LED_PIN, OUTPUT);
    pinMode(DSP2_LED_PIN, OUTPUT);
    pinMode(DSP5_LED_PIN, OUTPUT);
    pinMode(DSP6_LED_PIN, OUTPUT);



    

    Serial.println("DSP Init");
    dspinit();

    
    Serial.println("LED Light Test");
    digitalWrite(DSP1_LED_PIN, true);
    delay(500);
    digitalWrite(DSP1_LED_PIN, false);
    delay(100);
    digitalWrite(DSP2_LED_PIN, true);
    delay(500);
    digitalWrite(DSP2_LED_PIN, false);
    delay(100);
    digitalWrite(DSP5_LED_PIN, true);
    delay(500);
    digitalWrite(DSP5_LED_PIN, false);
    delay(100);
    digitalWrite(DSP6_LED_PIN, true);
    delay(500);
    digitalWrite(DSP6_LED_PIN, false);
    delay(100);

    Serial.println("DSP Draw");

    IndicatorDSP1.writeImage(WS_Bitmap104x212, 0, 0, 104, 212, 0,0,1);
    ArtDSP1.writeImage(WS_Bitmap3c104x212_black, WS_Bitmap3c104x212_red, 0, 0, 104, 212, 0,0,1);


    digitalWrite(DSP1_CS_PIN, false);
    digitalWrite(DSP2_CS_PIN, false);
    digitalWrite(DSP3_CS_PIN, false);
    digitalWrite(DSP4_CS_PIN, false);
    digitalWrite(DSP5_CS_PIN, false); 
    digitalWrite(DSP6_CS_PIN, false);
    IndicatorDSP1.refresh(false);
    digitalWrite(DSP1_CS_PIN, true);
    digitalWrite(DSP2_CS_PIN, true);
    digitalWrite(DSP3_CS_PIN, true);
    digitalWrite(DSP4_CS_PIN, true);
    digitalWrite(DSP5_CS_PIN, true);
    digitalWrite(DSP6_CS_PIN, true);
}

void loop(){
    Serial.println("Done");
    delay(1000);
}