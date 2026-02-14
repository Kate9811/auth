package service

//go:generate cmd /c "if exist mocks rmdir /s /q mocks"
//go:generate cmd /c mkdir mocks
//go:generate ..\..\bin\minimock -i AuthService -o ./mocks/ -s "_minimock.go"
