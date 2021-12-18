// This will wait for the astilectron namespace to be ready
document.addEventListener('astilectron-ready', function() {
    // This will listen to messages sent by GO
    astilectron.onMessage(function(message) {
        decmsg = DecodeMessage(message);
        console.log(decmsg);
        switchScreen(decmsg.screen)
        updateStrikes(decmsg.strike)
        updateClock(parseTimeToString(decmsg.time));
    });
})

function DecodeMessage(message) {
    var decoded = atob(message);
    var json = JSON.parse(decoded);
    return json;
}