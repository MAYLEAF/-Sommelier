# Sommelier

Sommelier는 TCP 부하테스트 툴로 json 포맷의 데이터를 주고 받습니다.

Sommelier is TCP LoadTesting Tool sending data with json format.

TCP LoadTesting Tool with golang


## Install
### Source
You need go installed and GOBIN in your PATH.

$ go get -u github.com/MAYLEAF/Sommelier

## Usage
Usage: Sommelier [global flags]

global flags:

   -request string ex: request.json
      json file saved for request message
    
   -value string ex: value.csv
      csv file saved for thread value
