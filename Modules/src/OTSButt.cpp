
class OtsButt: public basemodule{
    bool solved = false;
    int time; //placeholder until I do the time later this week.
    int color;
    int strip;
    int batteries;
    char[][] indicators;
    int chosenWord;//index of the chosen word rather than having it be a whole character array. just makes things faster
    //only 4 possible words, placing them into an array then using the seed to choose which word will be displayed
    char[4][8]  words= {{' ',' ','A','b','o','r','t',' '}, {'D','e','t','o','n','a','t','e'}, {' ',' ','H','o','l','d',' ',' '}, {' ',' ','P','r','e','s','s',' '}};//use null character instead of spaces, forget what it is in c++
    void setupModule() override {
        color = seed%5;
        strip = seed+7%4;//adding 7 so when I scale to 10 colors each strip will be different than color. 7 was chosen randomly i just like 7 and will only ever add prime numbers
        chosenWord = seed+11%4;
        if(strip == 0){
            //strip display blue
        }
        else if(strip == 1){
            //strip display white
        }
        else if(strip == 2){
            //strip display yellow
        }
        else{
            //strip display above purple. Can be a set of 7 different colors all of which have the same effect so im simplifying to purple for now.
        }
        if(color == 0){
            //button display blue
        }
        else if(color == 1){
            //button display red
        }
        else if(color == 2){
            //button display yellow
        }
        else if(color == 3){
            //button display red
        }
        else{
            //button display white
        }
        /*
            Important segment! display words[chosenWord] on the button for the word thingy
        */

        //SETUP TIMER I AM WAY TOO TIRED TO READ HOW TO DO THIS :P
    }

    void runModule() override{// else if ladder to get into the given segment then go from there
        if(strip == 0 && chosenWord == 0 || (/*Indicator says car && */ color == 4) || (color == 2)){// part 1 part 3 part 5
            while(!solved){
                relHeldButt();
            }
        }

        else if((batteries > 1 && chosenWord == 1) || (batteries >2 /*FRK in indicator */)){// part 2, part 4 and part 6 since steps are the same
                int timeHeld = 0; // still need to read documentation
                while(!solved){
                    if( true){//Button pressed
                        if(timeHeld < 1){
                            solveModule();
                        }
                        else{
                            failModule;
                        }
                    }
                }

        }
        else{
            while( !solved){
                relHeldButt();
            }
        }
        

        
    }
    void solveModule(){
        solved = true;
        //send solved signal to controller
    }
    void failModule(){
        //send failure signal to controller
    }
    void relHeldButt(){
        //I am assuming the timer is in seconds
        int release = 0;
        if(strip == 0){
            release = 4;
        }
        else if(strip == 2){
            release == 5;
        }
        else{// this is for a white strip or any other color strip so theres no reason to specify white itll fall in here anyway
            release == 1;
        }

        //Now we release when specified...
        while(true){//BUTTON IS HELD
            if(false){//button is released
                if(timer % 10 == release || timer/10 % 10 == release || timer/60 == release){
                    solveModule();
                }
                else{
                    failModule();
                }
            }
        }
    }



}