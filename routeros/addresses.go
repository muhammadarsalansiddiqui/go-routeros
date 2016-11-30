package routeros

import (
	"strconv"
)

type Address struct {
	id        string
	Address   string
	Interface string
	Dynamic   bool
	Invalid   bool
	Disabled  bool
}

func (s *Session) DescribeAddresses() ([]Address, error) {
	r, err := s.Request(Request{Sentence{
		Command:    "ip/address/print",
		Attributes: map[string]string{}}})
	if err != nil {
		return nil, err
	}
	a := []Address{}
	for _, s := range r.Sentences {
		if "re" == s.Command {
			a = append(a, Address{
				id:        s.Attributes[".id"],
				Address:   s.Attributes["address"],
				Interface: s.Attributes["interface"],
				Dynamic:   parseBool(s.Attributes["dynamic"]),
				Invalid:   parseBool(s.Attributes["invalid"]),
				Disabled:  parseBool(s.Attributes["disabled"])})
		}
	}
	return a, err
}

func (s *Session) AddAddress(a Address) error {
	_, err := s.Request(Request{Sentence{
		Command: "ip/address/add",
		Attributes: map[string]string{
			"address":   a.Address,
			"interface": a.Interface}}})
	return err
}

func (s *Session) RemoveAddress(a Address) error {
	return s.withAddressIndex(a, func(pos int) error {
        _, err := s.Request(Request{Sentence{
            Command: "ip/address/remove",
            Attributes: map[string]string{
                "numbers": strconv.Itoa(pos)},
            Query: map[string]string{}}})
        return err
    })
}

func (s *Session) withAddressIndex(a Address, action func(int) error) error {
    addresses, err := s.DescribeAddresses()
    if err != nil { return err }
    for i, address := range addresses {
        if address.id == a.id || address.Address == a.Address {
            return action(i)
        }
    }
    return nil
}
