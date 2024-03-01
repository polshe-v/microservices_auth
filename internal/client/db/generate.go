package db

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ./../../../bin/minimock -g -i Transactor -o ./mocks/ -s "_minimock.go"
//go:generate ./../../../bin/minimock -g -i github.com/jackc/pgx/v5.Tx -o ./mocks/ -s "_minimock.go"
