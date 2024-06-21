## LIBRARY SETUP INSTRUCTIONS:  
Find the library from library branch  
Create a new Go Project, with a go.mod file and a .go file.  

 - Go mod can be initialised with:  
  `go mod init <module_name>`  

In the terminal of that folder, to get the library in your go module use:  
  `go get github.com/pragadeesh-mcw/Go-Mini-Project@library`


Now you can import this library as a package in your .go file.

Example:
> import  
> (  
>  cache "github.com/pragadeesh-mcw/Go-Mini-Project" )   
> func main()   
> {  
>  r := cache.Entry()  
>   r.Run(":8080")  
> }  
