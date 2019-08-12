
########## first skeleton ##########

this is the standard Go structure
since you've learnt c/c++ it will be similar

```
import (
    xxxx
    yyyy
)


func main( xxx ) {
   read a configuration file or CLI parameter for port XXXX

   i will run a web server bind to port XXXX
      if error 
          print the error and exit the app
      if success
	  print app success running on port XXXX	
   
   create a handler to process all incoming http request
   
}
```

########## lesson learnt from 1st skeleton ##########

from #1 i know at least i will import
net/http -> for web server
fmt -> for printing to console
library about configfile
library about cli

i will know also at least i will have 3 major jobs that can be translated as functions also
1. function for reading configuration file/cli parameter
2. job for running webserver on port XXXX
3. function to handle HTTP request 


########## 2nd skeleton ##########

now you can expand the pseudo codes

```
import (
    fmt
    net/http
    configfile things
    cli things
)


func Configuration() {
   do skeleton here and later on expand
}

func Handler {
   do skeleton here and later on expand
}

func main( xxx ) {
   read a configuration file or CLI parameter for port XXXX

   i will run a web server bind to port XXXX
      if error 
          print the error and exit the app
      if success
	  print app success running on port XXXX	
   
   create a handler to process all incoming http request
   
}
```

########## now what have you learn from 2nd skeleton? ##########
- more functions maybe
- variable type that will be used
- data structs ?
- how to test the logics and functions using both positive and negative condition
- continue the skeleton until you found the final state of your own 
- or you've already visualized all the lego components already and continue to start building the real one with Go (or whatever language) codes
