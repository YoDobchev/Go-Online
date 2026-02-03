package gogame

// gameId -> game
var GameInstances = make(map[string]*Game)

// player username -> game
var PlayerToGame = make(map[string]*Game)
