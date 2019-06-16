package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	mrand "math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

// LoginService are they valid credentials to login
func LoginService(loginReq string) (GeneratedCode, error) {
	l := Login{}
	err := json.Unmarshal([]byte(loginReq), &l)
	if err != nil {
		return GeneratedCode{}, err
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/identify/%s/%s", os.Getenv("IDENTITY_URI"), l.Email, l.Plate), nil)
	if err != nil {
		return GeneratedCode{}, err
	}
	key := os.Getenv("IDENTITY_AUTH_KEY")
	req.Header.Set("X-Authorization", key)
	resp, err := client.Do(req)
	if err != nil {
		return GeneratedCode{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GeneratedCode{}, err
	}

	ident := Identity{}
	err = json.Unmarshal(body, &ident)
	if err != nil {
		return GeneratedCode{}, err
	}

	code, err := ident.generateCode()
	if err != nil {
		return code, err
	}

	return GeneratedCode{}, nil
}

func (i Identity) generateCode() (GeneratedCode, error) {
	data, err := i.generateData()
	if err != nil {
		return GeneratedCode{}, err
	}

	code := generateCode()
	block, err := aes.NewCipher([]byte(code))
	if err != nil {
		return GeneratedCode{}, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return GeneratedCode{}, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return GeneratedCode{}, err
	}
	cipherText := gcm.Seal(nonce, nonce, []byte(data), nil)

	id := makeId()
	err = StoreData(cipherText, id)
	if err != nil {
		return GeneratedCode{}, err
	}

	return GeneratedCode{
		ID: id,
		Code: code,
	}, nil
}

func (i Identity) generateData() ([]byte, error) {
	loginIdent := LoginIdent{
		ID: i.Ident.ID,
	}

	j, err := json.Marshal(loginIdent)
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}

func compressLZW(testStr string) []int {
	code := 256
	dictionary := make(map[string]int)
	for i := 0; i < 256; i++ {
		dictionary[string(i)] = i
	}

	currChar := ""
	result := make([]int, 0)
	for _, c := range []byte(testStr) {
		phrase := currChar + string(c)
		if _, isTrue := dictionary[phrase]; isTrue {
			currChar = phrase
		} else {
			result = append(result, dictionary[currChar])
			dictionary[phrase] = code
			code++
			currChar = string(c)
		}
	}
	if currChar != "" {
		result = append(result, dictionary[currChar])
	}
	return result
}

func makeId() int {
	t := time.Now()
	t = t.Add(5 * time.Minute)
	ints := compressLZW(t.String())

	ret := 0
	for _, i := range ints {
		ret += i
	}

	return ret
}

func generateCode() string {
	t := time.Now()
	s := os.Getenv("ENCRYPTION_KEY")

	k := t.String() + s
	i := compressLZW(k)

	r := 0
	for _, j := range i {
		r += j
	}
	r *= mrand.Int()
	rs := strconv.Itoa(r)
	ret := ""

	if len(rs) >= 20 {
		if rs[0:1] == "-" {
			ret = fmt.Sprintf("9%s-%s-%s-%s-%s", rs[1:4], rs[4:8], rs[8:12], rs[12:16], rs[16:20])
		} else {
			ret = fmt.Sprintf("%s-%s-%s-%s-%s", rs[0:4], rs[4:8], rs[8:12], rs[12:16], rs[16:20])
		}
	} else {
		return generateCode()
	}

	return ret
}
