Fetching a list of essays and counting the top 10 words from all the
essays combined.
A valid word will:
a) Contain at least 3 characters. 
b) Contain only alphabetic characters.
c) Be part of our bank of words (not all the words in the bank are valid according to the
   previous rules)
   The output should be pretty JSON printed to the stdout.

1. Install Go:
   If you haven't installed Go on your machine, download and install it from the official Go Downloads page.

2. Set Up Your Go Workspace:
   Make sure you have a proper Go workspace. The workspace typically includes a src, pkg, and bin directory. You can set up your workspace as follows:
   ```
   export GOPATH=~/go
   export PATH=$PATH:$GOPATH/bin
   ```
   Replace ~/go with the path to your desired workspace.

3. Clone Your Repository:
   Clone your Go project repository from GitHub to your local machine:
   ```
   git clone https://github.com/michaelein/Top10words
   cd your-repo
   ```
4. Install Dependencies :
   If your project has external dependencies, use the go get command to install them:
   ```
   go get -u ./...
   ```
5. Run Your Program:
   Navigate to the directory containing your Go code and run the main program:
   ```
   go run main.go
   ```
