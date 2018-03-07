package kv

type KVClietn interface{
	Put(...interface{})
}