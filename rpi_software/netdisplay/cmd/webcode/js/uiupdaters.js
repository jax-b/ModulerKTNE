function switchscreen(screen) {
    homescreen = document.querySelector("#WaitingScreen");
    boomscreen = document.querySelector("#GameBoom");
    winscreen = document.querySelector("#GameWin");
    countscreen = document.querySelector("#GameRunClock");
    switch (screen) {
        case "home":
            homescreen.hidden = false;
            boomscreen.hidden = true;
            winscreen.hidden = true;
            countscreen.hidden = true;
            break;
        case "boom":
            homescreen.hidden = true;
            boomscreen.hidden = false;
            winscreen.hidden = true;
            countscreen.hidden = true;
            break;
        case "win":
            homescreen.hidden = true;
            boomscreen.hidden = true;
            winscreen.hidden = false;
            countscreen.hidden = true;
            break;
        case "gametime":
            homescreen.hidden = true;
            boomscreen.hidden = true;
            winscreen.hidden = true;
            countscreen.hidden = false;
            break;
        default:
            console.log("Unknown screen: " + screen)
    }
}
function updatestrikes(strikes) {
    strike1 = document.querySelector("#strike1");
    strike2 = document.querySelector("#strike2");
    strikeText = document.querySelector("#strikeText");
    switch (strikes) {
        case 0:
            strike1.hidden = false;
            strike2.hidden = false;
            strikeText.hidden = true;
            strike1.classList.remove("on");
            strike2.classList.remove("on");
            strike1.classList.add("off");
            strike2.classList.add("off");
            break;
        case 1:
            strike1.hidden = false;
            strike2.hidden = false;
            strikeText.hidden = true;
            strike1.classList.remove("off");
            strike2.classList.remove("on");
            strike1.classList.add("on");
            strike2.classList.add("off");
            break;
        case 2:
            strike1.hidden = false;
            strike2.hidden = false;
            strikeText.hidden = true;
            strike1.classList.remove("off");
            strike2.classList.remove("off");
            strike1.classList.add("on");
            strike2.classList.add("on");
            break;
        case 3: 
            strike1.hidden = false;
            strike2.hidden = false;
            strikeText.hidden = true;
            strike1.classList.remove("on");
            strike2.classList.remove("off");
            strike1.classList.add("off");
            strike2.classList.add("on");
            break;
        default:
            strike1.hidden = true;
            strike2.hidden = true;
            strikeText.hidden = false;
            strikeText.innerHTML = strikes;
    }
}
function updateClock(time) {
    display1 = document.getElementById('display-1');
    display2 = document.getElementById('display-2');
    displayDot = document.getElementById('colendot');
    display3 = document.getElementById('display-3');
    display4 = document.getElementById('display-4');

    baseclass = "display-container display-size-12";
    dotclass = "display-container display-size-12-indicator"

    display1.classList = baseclass + " display-no-" + time.charAt(0);
    display2.classList = baseclass + " display-no-" + time.charAt(1);
    display3.classList = baseclass + " display-no-" + time.charAt(3);
    display4.classList = baseclass + " display-no-" + time.charAt(4);
    if (time.charAt(2) == ':') {
        displayDot.classList = dotclass + " display-colon";
    } else if (time.charAt(2) == '.') {
        displayDot.classList = dotclass + " display-dot";
    } else {
        displayDot.classList = dotclass + " display-off";
    }
    
}