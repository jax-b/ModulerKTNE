function parseTimeToString(strtime) {
    // If we have minutes remaining
    if (strtime.includes("m")) {
        strtime = strtime.split("m");
        if (strtime[0].length == 1) {
            strtime[0]= "0" + strtime[0];
        }
        sec = strtime[1].replace('s', '').split(".");
        if (sec[0].length == 1) {
            sec[0] = "0" + sec[0];
        }
        return strtime[0] + ":" + sec[0];
    } else {
        sec = strtime.split(".");
        if (sec[0].length  == 1) {
            sec[0]= "0" + sec[0];
        }
        return sec[0] + "." + sec[1].substring(0, 2);
    }
}