#include "OTS_Button.h"

//#define DEBUG_MODE true

OTS_Button::OTS_Button()
{
    const char realModID[] = "bttn";
    for (uint8_t i = 0; i < 4; i++)
    {
        modID[i] = realModID[i];
    }
}

void OTS_Button::clearModule()
{
    OTS_Button::_pixels->setPixelColor(1, 0);
    OTS_Button::_pixels->setPixelColor(0, 0);
    OTS_Button::_pixels->setPixelColor(2, 0);
    OTS_Button::_pixels->setPixelColor(3, 0);
    OTS_Button::_pixels->show();
    OTS_Button::_display.firstPage();
    do
    {
        OTS_Button::_display.fillScreen(GxEPD_WHITE);
    } while (OTS_Button::_display.nextPage());;
    OTS_Button::failureTriggered = false;
    OTS_Button::successTriggered = false;
    OTS_Button::numBatteries = 0;
    OTS_Button::numStrike = 0;
    OTS_Button::activePorts = 0;
    OTS_Button::chosenWord = 0;
    OTS_Button::stripColor = 0x0;
    OTS_Button::stripColorHex = 0x0;
    OTS_Button::buttonColor = 0x0;
}

void OTS_Button::setupModule()
{
    // Setup button input
    pinMode(BUTTON_PIN, INPUT_PULLUP);
    // Setup the NeoPixel strip
    OTS_Button::_pixels = new Adafruit_NeoPixel(4, NEOPIXEL_PIN, NEO_GRB + NEO_KHZ800);
    OTS_Button::_pixels->begin();
    for (uint8_t i = 0; i < 4; i++)
    {
        OTS_Button::_pixels->setPixelColor(i, 0x0);
    }
    OTS_Button::_pixels->show();
    // Set up the Eink Display
    OTS_Button::_display.init();
    OTS_Button::_display.setRotation(1);
    OTS_Button::_display.setFont(&FreeMonoBold12pt7b);
    OTS_Button::_display.setFullWindow();
    OTS_Button::_display.setTextColor(GxEPD_BLACK);
}

uint16_t OTS_Button::btnDebounce()
{
    // Get the reading from the button
    int reading = !digitalRead(BUTTON_PIN);

    uint16_t tpress = millis() - OTS_Button::timeLastBtn; // calculate the press time in ms
    if (tpress > 10)
    { // Debounce delay to prevent flickering must be > ?ms
        // if the button state has changed:
        if (reading != OTS_Button::lastBTNState)
        {
            OTS_Button::lastBTNState = reading; // update the button state
            if (reading == true)
            { // if the button is pressed return the press time
                OTS_Button::_pixels->setPixelColor(0, OTS_Button::stripColorHex);
                OTS_Button::_pixels->setPixelColor(1, OTS_Button::stripColorHex);
                OTS_Button::timeLastBtn = millis();
            }
            else
            {
                OTS_Button::_pixels->setPixelColor(0, 0x0);
                OTS_Button::_pixels->setPixelColor(1, 0x0);
                OTS_Button::timeLastBtn = millis();
                return tpress;
            }
        }
    } 
    return 0;
}

void OTS_Button::setSeed(uint16_t inSeed)
{
    OTS_Button::seed = inSeed;
    if (seed == 0)
    {
        return;
    }
    // SETUP SEED
    OTS_Button::buttonColor = seed % 5;
    OTS_Button::stripColor = seed + 7 % 4; // adding 7 so when I scale to 10 colors each strip will be different than color. 7 was chosen randomly i just like 7 and will only ever add prime numbers
    OTS_Button::chosenWord = seed + 11 % 4;
    switch (stripColor)
    {
    case 0: // strip display blue
        OTS_Button::stripColorHex = 0x0000FF;
        break;
    case 1: // strip display white
        OTS_Button::stripColorHex = 0xFFFFFF;
        break;
    case 2: // strip display yellow
        OTS_Button::stripColorHex = 0xFFFF00;
        break;
    default: // strip display above purple. Can be a set of 7 different colors all of which have the same effect so im simplifying to purple for now.
        OTS_Button::stripColorHex = 0xFF00FF;
        break;
    }
    switch (OTS_Button::buttonColor)
    {
    case 0: // button display blue
        OTS_Button::_pixels->setPixelColor(2, 0x0000FF);
        OTS_Button::_pixels->setPixelColor(3, 0x0000FF);
        break;
    case 1: // button display red
    case 3: // Cooper had it twice so...
        OTS_Button::_pixels->setPixelColor(2, 0xFF0000);
        OTS_Button::_pixels->setPixelColor(3, 0xFF0000);
        break;
    case 2: // button display yellow
        OTS_Button::_pixels->setPixelColor(2, 0xFFFF00);
        OTS_Button::_pixels->setPixelColor(3, 0xFFFF00);
        break;
    default:
        // button display purple
        OTS_Button::_pixels->setPixelColor(2, 0xFF00FF);
        OTS_Button::_pixels->setPixelColor(3, 0xFF00FF);
        break;
    }
    OTS_Button::_pixels->setPixelColor(1, 0);
    OTS_Button::_pixels->setPixelColor(0, 0);
    OTS_Button::_pixels->show();

    if (OTS_Button::chosenWord >= 4)
    { // Prevent out of bounds
        OTS_Button::chosenWord = 3;
    }

    // Display the word on the screen
    const char possibleButtonWords[4][8] = {
        {'A', 'b', 'o', 'r', 't', '\0', '\0', '\0'},
        {'D', 'e', 't', 'o', 'n', 'a', 't', 'e'},
        {'H', 'o', 'l', 'd', '\0', '\0', '\0', '\0'},
        {'P', 'r', 'e', 's', 's', '\0', '\0', '\0'}};
    String strChosenWord = "";
    for (uint8_t i = 0; i < 8; i++)
    {
        strChosenWord += (char)possibleButtonWords[OTS_Button::chosenWord][i];
        if (possibleButtonWords[OTS_Button::chosenWord][i] == '\0')
        {
            break;
        }
    }

#ifdef DEBUG_MODE
    Serial.print("Chosen word num: ");
    Serial.print(OTS_Button::chosenWord);
    Serial.print(":");
    Serial.println(strChosenWord);
#endif

    // Center Text
    int16_t x1, y1;
    uint16_t w, h;
    OTS_Button::_display.getTextBounds(strChosenWord, 0, 0, &x1, &y1, &w, &h);
    uint16_t x = ((OTS_Button::_display.width() - w) / 2) - x1;
    uint16_t y = ((OTS_Button::_display.height() - h) / 2) - y1;
    OTS_Button::_display.firstPage();
    do
    {
        OTS_Button::_display.fillScreen(GxEPD_WHITE);
        OTS_Button::_display.setCursor(x, y);
        OTS_Button::_display.print(strChosenWord);

    } while (OTS_Button::_display.nextPage());
}

void OTS_Button::tickModule(uint16_t currentGameTime)
{
    uint16_t timeBTNPressed = btnDebounce();
    OTS_Button::_pixels->show();
    // We need to wait if the master controller has not cleared out the failure flag
    if (OTS_Button::failureTriggered)
    {
        return;
    }

    const char step4indiLabel[3] = {'F', 'R', 'K'};
    // Instructions 2, 4 and 6: end results are the same
    if ((OTS_Button::numBatteries > 1 && OTS_Button::chosenWord == 1) || (OTS_Button::numBatteries > 2 && this->checkIndicator(step4indiLabel)) || (OTS_Button::buttonColor == 1 && OTS_Button::chosenWord == 4)) // See 1/3 in set seed line 74
    {
        if (timeBTNPressed >= 1 && timeBTNPressed < 300)
        { // if the button is pressed for less than 1 second
#ifdef DEBUG_MODE
            Serial.println("Immediately Button Pressed Successfully");
#endif
            OTS_Button::successTriggered = true;
        }
        else if (timeBTNPressed == 0) // Button Not Pressed
        {
            return;
        }
        else
        {
#ifdef DEBUG_MODE
            Serial.println("Immediately Button Pressed Failed");
#endif
            OTS_Button::failureTriggered = true;
        }
    }
    else
    {
        if (timeBTNPressed >= 1)
        {
            this->relHeldButton(currentGameTime);
        }
    }
}

void OTS_Button::relHeldButton(uint16_t currentGameTime)
{
    // I am assuming the timer is in seconds
    char releaseDigit;
    switch (stripColor)
    {
    case 0:
        releaseDigit = '4';
        break;
    case 2:
        releaseDigit = '5';
        break;
    default: // this is for a white strip or any other color strip so theres no reason to specify white it'll fall in here anyway
        releaseDigit = '1';
        break;
    }

#ifdef DEBUG_MODE
    Serial.print("Release digit: ");
    Serial.println(releaseDigit);
#endif
    uint16_t hundrethsTime = currentGameTime / 10;
    String ctime = String(hundrethsTime / 100 / 60 / 10 % 10) + String(hundrethsTime / 100 / 60 % 10) + ":";
    ctime += String(hundrethsTime / 1000 % 10) + String(hundrethsTime / 100 % 10) + ".";
    ctime += String(hundrethsTime / 10 % 10) + String(hundrethsTime % 10);
    if ( ctime.charAt(ctime.indexOf(releaseDigit)) == releaseDigit)
    {
#ifdef DEBUG_MODE
        Serial.println("Held Button Successfully");
#endif
        OTS_Button::successTriggered = true;
    }
    else
    {
#ifdef DEBUG_MODE
        Serial.println("Held Button Failed");
#endif
        OTS_Button::failureTriggered = true;
    }
}
