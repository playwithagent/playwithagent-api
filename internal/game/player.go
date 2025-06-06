package game


type Player int 

const (
	EmptyPlayer	Player = iota
	PlayerX
	PlayerO
)	

func (p Player) String() string {
	switch p {
	case EmptyPlayer:
		return " "
	case PlayerX:
		return "X"
	case PlayerO:
		return "O"
	default:
		return " "
	}
}