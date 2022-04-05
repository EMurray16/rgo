library(shiny)

# We only want to load the shared library if we haven't already
if (!is.loaded("TestMatrix")) {
  print("loading shared library")
  dyn.load("demo.so")
}
print("shared library is loaded")

#This function parses CSV text into a vector
ParseCSV <- function(csvtext, AsInt=FALSE) {
    x = as.numeric(unlist(strsplit(csvtext, ",")))
    if (AsInt) {
        return(as.integer(x))
    }
    return(x)
}
sexp.ParseGoMatrix <- function(sexpOutput) {
  # The matrix is the second element of the list, while the first is an integer vector c(nrow,ncol)
  outMat = matrix(data=sexpOutput[[2]], nrow=sexpOutput[[1]][1], ncol=sexpOutput[[1]][2])
  return(outMat)
}

# Define UI for application that draws a histogram
ui <- fluidPage(
    headerPanel(strong("Running Go in R")),
    sidebarPanel( h4(strong("What should Go do?")), hr(),
                  textInput("InNum", label="Enter comma-separated numbers", value="0, 1, 2.71828, 3.14159"),
                  numericInput("InNrow", label="Number of matrix rows", value=2, min=1, max=10),
                  actionButton("RunFloat", label="Run TestFloat"),
                  actionButton("RunInt", label="Run TestInt"),
                  actionButton("RunMat", label="Run TestMatrix"),
                  hr(),
                  textInput("Rmessage", label="Enter a message to Gopher", value="Hello Gopher!"),
                  actionButton("RunString", label="Run TestString"),
                  br(), br(),
                  h5("Note that I was (and continue to be) too lazy to implement good error handling for the purposes of this demo.")
    ),
    mainPanel( 
        textOutput("RunMessage1"),
        verbatimTextOutput("YouEntered"),
        textOutput("RunMessage2"),
        verbatimTextOutput("GopherResponded")
    ), 
    title="Running Go in R"
)


# Define server logic required to draw a histogram
server <- function(input, output) {

    observeEvent(input$RunFloat, handlerExpr = {
        #Parse the string into numbers
        doubles = ParseCSV(input$InNum)
        #Render the first messages
        output$RunMessage1 = renderText("You provided the vector:")
        output$YouEntered = renderPrint({
            cat("print(doubles) \n")
            print(doubles)
        })
        
        output$RunMessage2 = renderText("The TestFloat function in Go reads the SEXP vector of doubles
    			you provided into a slice of float64s, doubles each element, and appends the sum. Then
    			it passes the result as a new SEXP back to R, which looks like this:"
        )
        #Call Go and print the result
        doubles2 = .Call("TestFloat", doubles)
        output$GopherResponded = renderPrint({
            cat("doubles2 = .Call('TestFloat', doubles) \n")
            cat("print(doubles2) \n")
            print(doubles2)
        })
    })
    
    #Handle the TestI function
    observeEvent(input$RunInt, handlerExpr = {
        #Parse teh string into Ints
        integers = ParseCSV(input$InNum, AsInt=TRUE)
        #Render the first message
        output$RunMessage1 = renderText("You provided the vector:")
        output$YouEntered = renderPrint({
            cat("print(integers) \n")
            print(integers)
        })
        output$RunMessage2 = renderText("The TestInt function in go reads the SEXP vector of integers
    			you provided into a slice of int64s, doubles each element, and appends the sum. Then it
    			passes the result as a new SEXP back to R, which looks like this:"
        )
        
        #Call Go and print the result
        integers2 = .Call("TestInt", integers)
        output$GopherResponded = renderPrint({
            cat("integers2 = .Call('TestInt', integers) \n")
            cat("print(integers2) \n")
            print(integers2)
        })
    })
    
    #Handle the TestS function
    observeEvent(input$RunString, handlerExpr = {
        #render the first messgae
        output$RunMessage1 = renderText("You said to Gopher:")
        output$YouEntered = renderPrint({
            cat("print(input$Rmessage) \n")
            print(input$Rmessage)
        })
        output$RunMessage2 = renderText("The TestString function reads the string provided in R into a string in Go.
    			Then, it thanks you for the message by acknowledging what you sent, sending the result back 
    			(twice for good measure)."
        )
        
        #Do the things
        GoResponse = .Call("TestString", input$Rmessage)

        output$GopherResponded = renderPrint({
            cat("GoResponse = .Call('TestString', input$Rmessage) \n")
            print(GoResponse)
        })
    })
    
    #Handle the TestMat button
    observeEvent(input$RunMat, handlerExpr = {
        ##Parse the string into numbers
        doubles = ParseCSV(input$InNum)
        #Make the matrix
        doublesMatrix = matrix(nrow=input$InNrow, ncol=length(doubles))
        for (r in 1:nrow(doublesMatrix)) {
            doublesMatrix[r,] = doubles
        }

        #render the first message
        output$RunMessage1 = renderText("Gopher is starting with this matrix:")
        output$YouEntered = renderPrint({
            cat("print(doublesMatrix) \n")
            print(doublesMatrix)
        })
        output$RunMessage2 = renderText("The TestMatrix function reads the matrix provided, a d then adds a column. For each row, 
            the added element is the sum of that row. Then, it adds a row which holds the sum of each column. Once that's done, 
            it reformats the new matrix into a special list, sends it to R, and then R can create a new matrix from that list
            . The whole process looks like this:"
        )
        
        #Pretend to do the work
        size = c(nrow(doublesMatrix),ncol(doublesMatrix))
        GoMatList = .Call("TestMatrix", doublesMatrix, size)
        if (length(GoMatList) != 2) {
          output$GopherResponded = renderPrint({
            cat("oops! Looks like there was an error in the TestMatrix function:\n")
            print(GoMatList)
          })
          return()
        }
        doublesMatrix2 = sexp.ParseGoMatrix(GoMatList)
        
        output$GopherResponded = renderPrint({
            cat("size = c(nrow(doublesMatrix),ncol(doublesMatrix)) \ndoublesMatrix2 = sexp.ParseGoMatrix(.Call('TestMatrix', doublesMatrix, size)) \nprint(doublesMatrix2)\n")
            print(doublesMatrix2)
        })
    })
}

# Run the application 
shinyApp(ui = ui, server = server)
