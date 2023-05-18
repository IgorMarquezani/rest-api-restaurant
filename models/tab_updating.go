package models

import (
	"errors"
	"sort"
)

type UpdatingTab struct {
	RoomId   int     `json:"room"`
	Number   int     `json:"number"`
	PayValue float64 `json:"pay_value"`
	Maded    string  `json:"time_maded"`
	Table    int     `json:"table"`

	Requests []UpdatingRequest `json:"requests"`
}

func (ut UpdatingTab) ToNormalTab() Tab {
	var requests []Request

	for _, updatingRequest := range ut.Requests {
		request := Request{
			TabRoom:     updatingRequest.TabRoom,
			TabNumber:   updatingRequest.TabNumber,
			ProductName: updatingRequest.ProductName,
			Quantity:    updatingRequest.Quantity,
		}

		requests = append(requests, request)
	}

	return Tab{
		RoomId:   ut.RoomId,
		Number:   ut.Number,
		PayValue: ut.PayValue,
		Maded:    ut.Maded,
		Table:    ut.Table,
		Requests: requests,
	}
}

func (ut UpdatingTab) ToNormalRequest(i int) (Request, error) {
	if i < 0 || i > ut.Len()-1 {
		return Request{}, errors.New("Index out of range")
	}

	return Request{
		TabRoom:         ut.RoomId,
		TabNumber:       ut.Number,
		ProductName:     ut.Requests[i].ProductName,
		ProductListRoom: ut.RoomId,
		Quantity:        ut.Requests[i].Quantity,
	}, nil
}

func (ut *UpdatingTab) Len() int {
	return len(ut.Requests)
}

func (ut *UpdatingTab) Less(i, j int) bool {
	return ut.Requests[i].ProductName < ut.Requests[j].ProductName
}

func (ut *UpdatingTab) Swap(i, j int) {
	request := ut.Requests[i]
	ut.Requests[i] = ut.Requests[j]
	ut.Requests[j] = request
}

func (ut *UpdatingTab) GroupRequests() {
	var (
		total   uint
		grouped []UpdatingRequest
	)

	sort.Sort(ut)

	for i := 0; i < ut.Len(); i++ {
		request := ut.Requests[i]

		if request.Operation == "" {
			request.Operation = "inserting"
		}

		for ; i < ut.Len()-1 && ut.Requests[i].ProductName == ut.Requests[i+1].ProductName; i++ {
			if ut.Requests[i+1].Operation == "deleting" {
				request.Quantity -= ut.Requests[i+1].Quantity
				continue
			}

			request.Quantity += ut.Requests[i+1].Quantity
		}

		grouped = append(grouped, request)
		total++
	}

	ut.Requests = grouped[0:total]
}
