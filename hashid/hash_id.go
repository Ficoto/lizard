package hashid

import (
	"fmt"

	"github.com/speps/go-hashids"
)

const (
	defaultAlphabet   string = "abcdefghijklmnopqrstuvwxyz1234567890"
	minAlphabetLength int    = 6
)

// Option configures a Hash
type Option func(h *Hash) error

// WithSalt sets the underlying salt.
func WithSalt(salt string) Option {
	return func(h *Hash) error {
		err := checkSalt(salt)
		if err != nil {
			return err
		}
		if h.d.Salt != "" {
			return fmt.Errorf("can not change the salt when having a value")
		}
		h.d.Salt = salt
		return nil
	}
}

func checkSalt(salt string) error {
	if salt == "" {
		return fmt.Errorf("salt can not be a zero string")
	}
	return nil
}

// WithAlphabet sets the underlying alphabet.
func WithAlphabet(alphabet string) Option {
	return func(h *Hash) error {
		if alphabet == "" {
			return fmt.Errorf("alphabet can not be a zero string")
		}
		if h.d.Alphabet != "" {
			return fmt.Errorf("can not change the alphabet when having a value")
		}
		h.d.Alphabet = alphabet
		return nil
	}
}

// Hash represents a mapping.
type Hash struct {
	d *hashids.HashIDData
	h *hashids.HashID
}

// New generates a Hash.
func New(oo ...Option) (*Hash, error) {
	var err error
	h := &Hash{
		d: &hashids.HashIDData{},
	}
	h.d.MinLength = minAlphabetLength
	h.d.Alphabet = defaultAlphabet
	for _, o := range oo {
		if err := o(h); err != nil {
			return nil, err
		}
	}
	h.h, err = hashids.NewWithData(h.d)
	if err != nil {
		return nil, err
	}
	err = checkSalt(h.d.Salt)
	if err != nil {
		return nil, err
	}
	return h, nil
}

// Encode encodes integers to strings.
func (h *Hash) Encode(ii ...int64) (map[int64]string, error) {
	ret := make(map[int64]string)
	for _, i := range ii {
		s, err := h.h.EncodeInt64([]int64{i})
		if err != nil {
			return nil, err
		}
		ret[i] = s
	}
	return ret, nil
}

// Decode decodes strings to integers.
func (h *Hash) Decode(ss ...string) (map[string]int64, error) {
	ret := make(map[string]int64)
	for _, s := range ss {
		ii, err := h.h.DecodeInt64WithError(s)
		if err != nil {
			return nil, err
		}
		ret[s] = ii[0]
	}
	return ret, nil
}
