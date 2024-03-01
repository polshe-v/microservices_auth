package service

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ./../../bin/minimock -g -i UserService -o ./mocks/ -s "_minimock.go"
