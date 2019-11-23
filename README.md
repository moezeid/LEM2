# LEM2





### Compiling and Running On Linux

   This is targeted towards users who wants to compile and run 
   this project that may not have sudo access on their system
    to download the Golang compiler
    
First:
    
   Navigate to your home directory
      
      ~$ cd ~/
     
   Download the go project
      
      ~$ curl -O https://storage.googleapis.com/golang/go1.11.2.linux-amd64.tar.gz
   
   Extract the package and set an environment variable
    
      ~$ tar -xvf go1.11.2.linux-amd64.tar.gz
      ~$ GO=~/go/bin/go
         
   If you already have the code in a zipfile:
     
        $ unzip LEM2.zip      
         
         
   If you don't already have the code:
   Now download the MLEM2 code
     
     ~$ git clone https://github.com/moezeid/LEM2.git 
    
  
   Go into the directory 
    
     ~$ cd LEM2
     
   Build and run the code:
     
     ~/LEM2  $GO build run/main.go
     ~/LEM2 ./main
      
    
   An output file will be produced in the same directory containing the results.
    
