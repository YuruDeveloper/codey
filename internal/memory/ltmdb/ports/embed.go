package ports

import "github.com/YuruDeveloper/codey/internal/memory/ltmdb/types"



type AppEmbed interface {
	Embed(input string) (types.Vector, error)
}