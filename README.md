# simple ascii
this is an simple project, an program made with go and raylib that enables drawing ascii without the problems of an text editor like a line break, characters positioning, etc, it has an infinite canvas where you can draw a character typing with the keyboard. 

to save the art you must press **ctrl + s**, this will create a new folder if not exists called "save" with an art{num}.txt, the result art is only the area that you have draw

## Commands
- Ctrl + S - saves the art
- Left Click + drag - selects an area
- Ctrl + X - cut selection
- Ctrl + C - copy selection
- Ctrl + v - pastes content from clipboard to the canvas
- Ctrl + h - shows the help commands
- typing keys to the canvas
- you can drag the camera arround the canvas with space + left button
- space + left button - drag camera"

## Some arts i made

```
             .                                                 
            /|\                                                
          /  |  \                                              
        /    |    \                                            
      /     |       \                                          
    /       |         \              An platonic form emerges..
  /        |            \                       An "Octaedran"!
/__        |           ___\                                    
\  --___   |     __----   /                                    
  \     --__ __--       /                                      
    \       |         /                                        
      \     |       /                                          
        \   |     /                                            
          \  |  /                                              
            \|/                                                
             '                                                 
                .  .                                           
                \ \ \                                          
                 . \.\                                         
          .      / /\ ;                                        
         / \    / / / |                                        
          \ -..\   /  .                                        
           '    .    /                                         
           \        /                                          
           '       -                                           
           /  ---/                                             
         ./     /                                              
        /      /                                               
       /      /                                                
     --      /                                                 
            .                                                  
```

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

```
    /\                               
   /  \       .                      
   |  |                              
   |  |    .      .                  
   |  |                              
   | .|        .   It's a sword      
   | .|                 or an dagger?
   |..|                              
   |..|                       ..     
   | .|                              
   | .|                              
__ | .| __                   .       
  \ ---/                             
   ||||                              
   ||||                     .        
   ||||                              
  /.  .\                             
  | . .|                             
   `---                              
```


```
         |                                        
    \ ...|... /                                   
     \ ..| . /                                    
   ...\ .|../ . .                                 
   ....\.-./ .         a star*                    
-------(@@@)-------                                 
     .. -_- . .                                   
    .. / | \ .                                    
      /  |  \                                     
     /   |   \                                    
         |                                        
                .---------------.                 
                |\               \                
                | \               \               
                |  \               \              
                |   \               \      an cube
                |    \               \            
                |     . ______________.           
                |     |               |           
                \     |               |           
                 \    |               |           
                  \   |               |           
                   \  |               |           
                    \ |               |           
                     \|_______________|           
```

# installing 
install the version of your operating system from the releases page of the repo

## compiling 
after installed the lastest version of **go**

in order to build to windows you will need mingw
in order to build to macOs you will need osxcross

run the command **build.sh** this will try to build to all operating systems

```
chmod +x ./build.sh
./build.sh
```

run (linux)
```
./build/simple-ascii-linux
```

## Improvements
some improvements can be made in the future such as enabling custom fonts, and custom color, multiple layers of ascii
select area or region and cut, and copy to another