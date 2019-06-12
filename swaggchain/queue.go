package swaggchain

import "container/list"

type Queue struct { *list.List }


func NewQueue() *list.List {

	Q := Queue.Init(nil)
	return Q

}

func (q *Queue) Add(element list.Element) {

	q.PushBack(element)

}

func (q *Queue) GetFrontItem() *list.Element {

	return q.Front()

}

func (q *Queue) ProcessQueue(fn func(element *list.Element) (bool, error)) {

	for q.Len() > 0 {

		e := q.GetFrontItem()

		_, err := fn(e)
		if err != nil {
			log.Errorf("an error processing internal queue: %s", err)
		}

		q.Remove(e)

	}

}

type TransactionQueue struct {
	Q *Queue
}

func NewTransactionQueue() *TransactionQueue {

	queue := new(Queue)
	txQueue := &TransactionQueue{Q:queue}
	return txQueue

}

func (txq TransactionQueue) Count() int {
	return txq.Q.Len()
}


func ProcessTransaction(el *list.Element) (bool, error) {

	return false, nil

}

func (txq *TransactionQueue) Process() {

	txq.Q.ProcessQueue(ProcessTransaction)

}


type BlockQueue struct {
	Q *Queue
}

func NewBlockQueue() *BlockQueue {

	queue := new(Queue)
	bQueue := &BlockQueue{Q:queue}
	return bQueue

}

func (bq *BlockQueue) Count() int {
	return bq.Q.Len()
}

func ProcessBlock(el *list.Element) (bool, error) {

	return false, nil

}

func (bq *BlockQueue) Process() {

		bq.Q.ProcessQueue(ProcessBlock)

}

