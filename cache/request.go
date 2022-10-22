package cache

import (
	"ti-ticket/DAO"
)

type Request struct {
	Proposer       string
	Privilege_code int
	Next           *Request
}

var (
	queue_head *Request  = nil
	queue_end  **Request = &queue_head
)

func push_request(proposer string, code int) {
	(*queue_end) = &Request{
		Proposer:       proposer,
		Privilege_code: code,
	}

	queue_end = &(*queue_end).Next
}

func pop_request() *Request {
	rp := queue_head
	queue_head = (*rp).Next
	return rp
}

func ReceiveRequest(proposer string, code int) {
	push_request(proposer, code)
}

func ListRequest() *[]map[string]string {
	var requests []map[string]string
	for it := queue_head; it != nil; it = (*it).Next {
		str := ""
		for mask, priv := range DAO.Privilege_table() {
			if (mask & (*it).Privilege_code) != 0 {
				str += "," + priv
			}
		}
		requests = append(requests, map[string]string{
			"proposer":           (*it).Proposer,
			"required_privilege": str[1:],
		})
	}
	return &requests
}

func HandleRequest(op bool) (string, error) {
	request := *pop_request()
	if !op {
		return "Reject", nil
	}
	pri_str, err := DAO.GrantUser(request.Proposer, request.Privilege_code)

	return pri_str, err
}
