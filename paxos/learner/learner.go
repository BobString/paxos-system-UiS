package learner

import (
  "acceptor"
	"fmt"
	"connector"
)

var(

)

type Pair struct {
	nv int
	val string
}

func receivingMsgs (incLearnerMsgs chan string)
{
	res = strings.Split(mesg,"@")
	
	x = res[1]
	y = res[2]
	p = Pair {res[1], res[2]}

	PairMap [Pair] int
	
	_,ok = PairMap [Pair{x, y}]
}
