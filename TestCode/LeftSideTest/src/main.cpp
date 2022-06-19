// #include <Ardiuno.h>
// #include <SPI.h>

#include <GxEPD2_3C.h> 



#define EPD_RESET_PIN 5 // can set to -1 and share with microcontroller Reset!
#define EPD_DC_PIN 4
#define EPD_MOSI_PIN 3
#define EPD_SCK_PIN 2
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


GxEPD2_3C<GxEPD2_290_C90c, GxEPD2_290_C90c::HEIGHT> IndicatorDSP1(GxEPD2_290_C90c(/*CS=*/DSP1_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/DSP1_BUS_PIN)); //Indicator
GxEPD2_3C<GxEPD2_290_C90c, GxEPD2_290_C90c::HEIGHT> IndicatorDSP2(GxEPD2_290_C90c(/*CS=*/DSP2_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/DSP2_BUS_PIN)); //Ind
GxEPD2_3C<GxEPD2_290_C90c, GxEPD2_290_C90c::HEIGHT> ArtDSP1(GxEPD2_290_C90c(/*CS=*/DSP3_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/DSP3_BUS_PIN));
GxEPD2_3C<GxEPD2_290_C90c, GxEPD2_290_C90c::HEIGHT> ArtDSP2(GxEPD2_290_C90c(/*CS=*/DSP4_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/DSP4_BUS_PIN));
GxEPD2_3C<GxEPD2_290_C90c, GxEPD2_290_C90c::HEIGHT> IndicatorDSP3(GxEPD2_290_C90c(/*CS=*/DSP5_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/DSP5_BUS_PIN));
GxEPD2_3C<GxEPD2_290_C90c, GxEPD2_290_C90c::HEIGHT> IndicatorDSP4(GxEPD2_290_C90c(/*CS=*/DSP6_CS_PIN, /*DC=*/EPD_DC_PIN, -1, /*BUSY=*/DSP6_BUS_PIN));

void setup(){
    // MbedSPI SPI(0, EPD_MOSI_PIN, EPD_SCK_PIN);
    SPI.setTX(EPD_MOSI_PIN);
    SPI.setSCK(EPD_SCK_PIN);
    Serial.begin(9600);
    Serial.println("Start");
    while(!Serial){
        Serial.println(".");
        delay(100);
    }
    delay(1000);
    Serial.println("Start");

    pinMode(EPD_RESET_PIN, OUTPUT);
    pinMode(DSP1_LED_PIN, OUTPUT);
    pinMode(DSP2_LED_PIN, OUTPUT);
    pinMode(DSP5_LED_PIN, OUTPUT);
    pinMode(DSP6_LED_PIN, OUTPUT);

    digitalWrite(EPD_RESET_PIN, true);
    digitalWrite(DSP1_LED_PIN, false);
    digitalWrite(DSP2_LED_PIN, false);
    digitalWrite(DSP5_LED_PIN, false);
    digitalWrite(DSP6_LED_PIN, false);

    Serial.println("DSP RESET");
    digitalWrite(EPD_RESET_PIN, true);
    delay(50);
    digitalWrite(EPD_RESET_PIN, false);
    delay(50);
    digitalWrite(EPD_RESET_PIN, true);

    Serial.println("DSP Init");
    IndicatorDSP1.init();
    IndicatorDSP2.init();
    ArtDSP1.init();
    ArtDSP2.init();
    IndicatorDSP3.init();
    IndicatorDSP4.init();

    digitalWrite(DSP1_CS_PIN, true);
    digitalWrite(DSP2_CS_PIN, true);
    digitalWrite(DSP3_CS_PIN, true);
    digitalWrite(DSP4_CS_PIN, true);
    digitalWrite(DSP5_CS_PIN, true);
    digitalWrite(DSP6_CS_PIN, true);

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

    IndicatorDSP1.firstPage();
    IndicatorDSP1.setTextColor(GxEPD_BLACK); 
    IndicatorDSP1.println("Hello World!");
    IndicatorDSP1.setTextColor(GxEPD_RED);
    IndicatorDSP1.println(" Want pi?");
    IndicatorDSP1.displayWindow(0,0,250,122);
}

void loop(){
    Serial.println("Done");
    delay(1000);
}