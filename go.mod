module github.com/edusig/pokedexcli

go 1.21.1

require internal/pokeapi v1.0.0

replace internal/pokeapi => ./internal/pokeapi

require internal/pokecache v1.0.0

replace internal/pokecache => ./internal/pokecache
