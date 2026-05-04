# erlc-go 

> simple go wrapper for erlc 

![code quality](https://img.shields.io/badge/code%20quality-clean-brightgreen)
![status](https://img.shields.io/badge/status-archived-lightgrey)
![license](https://img.shields.io/badge/license-apache--2.0-blue)

## install
go get github.com/joyclens/erlc-go

## usage
client, _ := erlc.NewClient("key")
defer client.Close()

server, _ := client.GetServer(ctx)
println(server.Name)

## features
- full api v1
- retries + rate limit
- caching
- context support
- no deps

## methods
client.GetServer(ctx)
client.GetPlayers(ctx)
client.GetBans(ctx)
client.GetCommandLogs(ctx)
