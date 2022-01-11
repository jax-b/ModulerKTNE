#include "OTS_Button.h"

void OTS_Button::OTS_Button() {
    modID = "bttn"
}

void OTS_Button::setupModule() override
{
    // Setup button input
    pinMode(BUTTON_PIN, INPUT_PULLUP);
    // Setup the NeoPixel strip
    _pixels = new Adafruit_NeoPixel(4, NEOPIXEL_PIN, NEO_GRB + NEO_KHZ800);
    _pixels->begin();
    for (uint8_t i = 0; i < 4; i++)
    {
        _pixels->setPixelColor(i, 0x0);
    }
    _pixels->show();
    // Set up the Eink Display

    _display.init();
    _display.setRotation(1);
    _display.setFont(&FreeMonoBold12pt7b);
    _display.setFullWindow();
    _display.setTextColor(GxEPD_BLACK);
}

uint16_t OTS_Button::btnDebounce()
{
    // Get the reading from the button
    int reading = !digitalRead(BUTTON_PIN);

    if (reading != lastBTNState)
    { // reset the debouncing timer if it has changed
        timeLastBtn = millis();
    }

    uint16_t tpress = millis() - lastBTNState; // calculate the press time in ms
    if (tpress > 5)
    { // Debounce delay to prevent flickering must be > 5ms
        // if the button state has changed:
        if (reading != lastBTNState)
        {
            lastBTNState = reading; // update the button state
            if (lastBTNState == true)
            { // if the button is pressed return the press time
                _pixels->setPixelColor(0, stripColorHex);
                _pixels->setPixelColor(1, stripColorHex);
                return tpress
            }
            else
            {
                _pixels->setPixelColor(0, 0x0);
                _pixels->setPixelColor(1, 0x0);
            }
        }
    }
    return 0;
}

void OTS_Button::setSeed(uint16_t inSeed) override
{
    seed = inSeed;
    // SETUP SEED
    buttonColor = seed % 5;
    stripColor = seed + 7 % 4; // adding 7 so when I scale to 10 colors each strip will be different than color. 7 was chosen randomly i just like 7 and will only ever add prime numbers
    chosenWord = seed + 11 % 4;
    switch (stripColor)
    {
    case 0: // strip display blue
        stripColorHex = 0x0000FF;
        break;
    case 1: // strip display white
        stripColorHex = 0xFFFFFF;
        break;
    case 2: // strip display yellow
        stripColorHex = 0xFFFF00;
        break;
    default: // strip display above purple. Can be a set of 7 different colors all of which have the same effect so im simplifying to purple for now.
        stripColorHex = 0xFF00FF;
        break;
    }
    switch (buttonColor)
    {
    case 0: // button display blue
        _pixels->setPixelColor(2, 0x0000FF);
        _pixels->setPixelColor(3, 0x0000FF);
        break;
    case 1: // button display red
    case 3: // Cooper had it twice so...
        _pixels->setPixelColor(2, 0xFF0000);
        _pixels->setPixelColor(3, 0xFF0000);
        break;
    case 2: // button display yellow
        _pixels->setPixelColor(2, 0xFFFF00);
        _pixels->setPixelColor(3, 0xFFFF00);
        break;
    default:
        _pixels->setPixelColor(2, 0xFF00FF);
        _pixels->setPixelColor(3, 0xFF00FF);
        break;
    }
    _pixels->show();
    // Display the word on the screen
    int16_t x1, y1;
    uint16_t w, h;
    _display.getTextBounds(possibleButtonWords[chosenWord], 0, 0, &x1, &y1, &w, &h);
    uint16_t x = ((_display.width() - w) / 2) - x1;
    uint16_t y = ((_display.height() - h) / 2) - y1;
    _display.firstPage();
    do
    {
        _display.fillScreen(GxEPD_WHITE);
        _display.setCursor(x, y);
        _display.print(possibleButtonWords[chosenWord]);

    } while (_display.nextPage());
}

void OTS_Button::tickModule(uint16_t currentGameTime) override
{
    uint16_t timeBTNPressed = btnDebounce();
    _pixels->show();
    // We need to wait if the master controller has not cleared out the failure flag
    if (failureTriggered)
    {
        return;
    }

    // Instructions 2, 4 and 6: end results are the same
    if ((batteries > 1 && chosenWord == 1) || (batteries > 2 && this->checkIndicator('FRK')) || (buttonColor == 1 && chosenWord == 4)) // See 1/3 in set seed line 74
    {
        if (timeBTNPressed >= 1 && timeBTNPressed < 300)
        { // if the button is pressed for less than 1 second
            successTriggered = true;
        }
        else if (timeBTNPressed == 0)
        {
            return;
        }
        else
        {
            failureTriggered = true;
        }
    }
    else
    {
        if (timeBTNPressed >= 1)
        {
            relHeldButt(currentGameTime);
        }
    }
}

void OTS_Button::relHeldButton(uint16_t currentGameTime)
{   
    uint16_t hundrethsTime = currentGameTime/ 10;
    // I am assuming the timer is in seconds
    int releaseDigit;
    if (stripColor == 0)
    {
        releaseDigit = 4;
    }
    else if (stripColor == 2)
    {
        releaseDigit = 5;
    }
    else
    { // this is for a white strip or any other color strip so theres no reason to specify white it'll fall in here anyway
        releaseDigit = 1;
    }
    // min:sec.tenhundosec 00:00.00
    //               00:00.0_                          00:00._0                               00:0_.00                                 00:_0.00                                  0_:00.00                                  _0:00.00
    if (hundrethsTime % 10 == release || hundrethsTime / 10 % 10 == release || hundrethsTime / 100 % 10 == release || hundrethsTime /1000 % 10 % 6 == release || hundrethsTime /100 / 60 % 10 == release || time /100 / 60 / 10 % 10 == release)
    {
        successTriggered = true;
    }
    else
    {
        failureTriggered = false;
    }
}
