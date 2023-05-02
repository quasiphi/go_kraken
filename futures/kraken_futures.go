package futures

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Kraken_Futures struct {
	key    string
	secret string
}

func New(key, secret string) *Kraken_Futures {
	if key == "" || secret == "" {
		log.Print("[WARNING] You have not set api key and secret!")
	}
	return &Kraken_Futures{
		key:    key,
		secret: secret,
	}
}

func (api *Kraken_Futures) Set_Auth_Headers(request *http.Request, PostData, EndpointPath string) error {
	// * Sets Auth headers required for authorization for the endpoint
	// ? Build using: https://docs.futures.kraken.com/#http-api-http-api-introduction-authentication
	nonce := time.Now().UnixMicro()

	conc := PostData + strconv.Itoa(int(nonce)) + EndpointPath
	step2 := sha256.Sum256([]byte(conc))

	data, err := base64.StdEncoding.DecodeString(api.secret)
	if err != nil {
		return err
	}

	h := hmac.New(sha512.New, []byte(data))
	h.Write(step2[:])

	authent := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// * Set Headers
	request.Header.Add("APIKey", api.key)
	request.Header.Add("Nonce", strconv.Itoa(int(nonce)))
	request.Header.Add("Authent", authent)
	return nil
}

func (api *Kraken_Futures) Request(RequestType, Full_URL, API_Path, Post_Body string, is_Private bool) (Response_Body []byte, Error error) {

	if is_Private {
		Full_URL = Full_URL + "?" + Post_Body
	}

	req, err := http.NewRequest(RequestType, Full_URL, nil)

	if err != nil {
		return nil, err
	}

	if is_Private {
		err = api.Set_Auth_Headers(req, Post_Body, API_Path)
		if err != nil {
			return nil, err
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
