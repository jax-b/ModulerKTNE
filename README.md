# ModulerKTNE
A real life keep talking and nobody explodes
The goal of this project it to make a modular [Keep Talking and Nobody Explodes game](https://keeptalkinggame.com/).
I also want this project to use the original game manual!
My eventual aim is to replicate every module.

Lets start out with my overall plan!
Each module should be self contained, meaning that that module knows weather or not it is solved.
Each module should be able to be configured from the central controller to replicate game scenarios.
The central controller should have a easy to use web interface.
Each module should look like it does in the game.
The game should know if it is being tampered with.

Coms between the central controller and the modules will likely use I2C with a master interrupt line.
