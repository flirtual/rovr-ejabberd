package main

import (
	"strconv"
	"time"
	"strings"

	"github.com/shelomentsevd/ejabberd-go"
	"github.com/gomodule/redigo/redis"
	rg "github.com/redislabs/redisgraph-go"
)

type Ergauth struct {
}

func (ergauth Ergauth) Auth(user, server, password string) bool {
	conn, _ := redis.Dial("tcp", "redis.example.com:12345", redis.DialPassword("..."))
	defer conn.Close()
	graph := rg.GraphNew("kwerc", conn)

	query := "MATCH (u:user)-[:SESSION]->(s:session {id: '" + password + "'}) WHERE toLower(u.username) = '" + strings.ToLower(user) + "' RETURN exists(s), s.expiry, s.expiryabs"
	result, _ := graph.Query(query)
	result.Next()
	record := result.Record()
	exists := rg.ToString(record.GetByIndex(0))
	expiry, _ := strconv.ParseInt(rg.ToString(record.GetByIndex(1)), 10, 64)
	expiryabs, _ := strconv.ParseInt(rg.ToString(record.GetByIndex(2)), 10, 64)

	now := time.Now().Unix()

	return exists == "true" && now < expiry && now < expiryabs
}

func (ergauth Ergauth) IsUser(user, server string) bool {
	conn, _ := redis.Dial("tcp", "redis.example.com:12345", redis.DialPassword("..."))
	defer conn.Close()
	graph := rg.GraphNew("kwerc", conn)

	query := "MATCH (u:user) WHERE toLower(u.username) = '" + strings.ToLower(user) + "' RETURN exists(u)"
	result, _ := graph.Query(query)
	result.Next()
	record := result.Record()
	exists := rg.ToString(record.GetByIndex(0))

	return exists == "true"
}

// Unsupported
func (ergauth Ergauth) SetPassword(user, server, password string) bool {
	return false
}
func (ergauth Ergauth) TryRegister(user, server, password string) bool {
	return false
}
func (ergauth Ergauth) RemoveUser(user, server string) bool {
	return false
}
func (ergauth Ergauth) RemoveUser3(user, server, password string) bool {
	return false
}

func main() {
	ergauth := Ergauth{}
	external := ejabberd.NewExternal(ergauth)
	external.Start()
}
