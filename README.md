# simple ascii
this is an simple project, an program made with pygame that enables drawing ascii without the problems of an text editor like a line break, characters positioning, etc, it has an infinite canvas where you can draw a character typing with the keyboard. 

to save the art you must press **ctrl + s**, this will create a new folder if not exists called "save" with an art{num}.txt, the result art is only the area that you have draw

## Commands

argument commands:
-- load - loads an art in form of text

keyboard commands:
- Ctrl + S - saves the art
- Ctrl + X - clears the screen

## Some arts i made

```
         .-----------                     
        /  \         \                    
       /    \         \                   
      /   |  \         \                  
     /  __._  \---------\ some box.txt    
    / ._  |   /         /   what's inside?
   /  | \ |  /     -.../                  
  /   \_.   /     /.../                   
 /___      /     /. ./                    
/ \ /     /     /.../                     
\  \ -.  /      \../                      
 \  \./ /       \;/                       
  \    /       \ /                        
   \  /         /                         
    \/________\/                          
```


## Requirements

- the python version used in this project is the current last version at the time **Python 3.13.7**
- The only requirement is **pygame==2.6.1**


## Improvements

some improvements can be made in the future such as an proper executable that runs on mac, linux and windows.
enabling custom fonts, and custom color, multiple layers of ascii

## update about making an executable
python's pyinstaller is unable to make a cross platform build from one machine, my conclusion is that python is an bad language for such things as generate an executable.

maybe i remade this project in another language.