package repository

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ./../../bin/minimock -g -i UserRepository -o ./mocks/ -s "_minimock.go"
//go:generate ./../../bin/minimock -g -i KeyRepository -o ./mocks/ -s "_minimock.go"
//go:generate ./../../bin/minimock -g -i AccessRepository -o ./mocks/ -s "_minimock.go"
//go:generate ./../../bin/minimock -g -i LogRepository -o ./mocks/ -s "_minimock.go"
