#!/bin/bash
# TODO сделать запуск тестов, определение тестового покрытия и генерацию html-ки с результатами анализа

go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
