// This will wait for the astilectron namespace to be ready
document.addEventListener('astilectron-ready', function() {
    // This will listen to messages sent by GO
    astilectron.onMessage(function(message) {
        decmsg = DecodeMessage(message);
        console.log(decmsg);
        switchscreen(decmsg.screen)
        updatestrikes(decmsg.strike)
    });
})

function DecodeMessage(message) {
    var decoded = atob(message);
    var json = JSON.parse(decoded);
    return json;
}