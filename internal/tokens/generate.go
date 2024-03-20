package tokens

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ./../../bin/minimock -g -i TokenOperations -o ./mocks/ -s "_minimock.go"
