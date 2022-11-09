package autosolve

import (
	"encoding/json"
	"errors"
	"time"

	"strings"

	"strconv"

	"net/http"

	"github.com/streadway/amqp"
)

const (
	hostname             = "amqp.autosolve.aycd.io"
	vHost                = "oneclick"
	directExchangePrefix = "exchanges.direct"
	fanoutExchangePrefix = "exchanges.fanout"

	responseQueuePrefix = "queues.response.direct"

	requestTokenRoutePrefix        = "routes.request.token"
	requestTokenCancelRoutePrefix  = "routes.request.token.cancel"
	responseTokenRoutePrefix       = "routes.response.token"
	responseTokenCancelRoutePrefix = "routes.response.token.cancel"

	autoAckQueue   = true
	exclusiveQueue = false
)

var (
	clientId                           string
	account                            Account
	captchaTokenResponseListener       CaptchaTokenResponseListener
	captchaTokenCancelResponseListener CaptchaTokenCancelResponseListener
	errorListener                      ErrorListener
	statusListener                     StatusListener
)

var (
	connection          *amqp.Connection
	directChannel       *amqp.Channel
	fanoutChannel       *amqp.Channel
	reconnectTimeouts         = [7]time.Duration{2, 3, 5, 8, 13, 21, 34}
	reconnectAttempt          = 0
	connectAttemptDelay int64 = 5
	lastConnectAttempt  int64 = 0
	connectionPending   bool  = false
	explicitShutdown          = false
)

var (
	directExchangeName          string
	fanoutExchangeName          string
	responseQueueName           string
	responseTokenRouteKey       string
	responseTokenCancelRouteKey string
	requestTokenRouteKey        string
	requestTokenCancelRouteKey  string
)

func Load(cId string, tokenL CaptchaTokenResponseListener,
	tokenCancelL CaptchaTokenCancelResponseListener,
	statusL StatusListener, errL ErrorListener) error {
	if len(clientId) > 0 {
		return errors.New("autosolve cannot be loaded more than once")
	}
	clientId = cId
	captchaTokenResponseListener = tokenL
	captchaTokenCancelResponseListener = tokenCancelL
	statusListener = statusL
	errorListener = errL
	return nil
}

func Connect(accessToken string, apiKey string) (ConnectResult, error) {
	if connectionPending {
		return ConnectionPending, nil
	}
	_ = Close()
	explicitShutdown = false
	if !isNotEmpty(apiKey) {
		return InvalidApiKey, nil
	}
	if isNotEmpty(accessToken) {
		tokenParts := strings.Split(accessToken, "-")
		if len(tokenParts) > 0 {
			aId, err := strconv.Atoi(tokenParts[0])
			if !checkError(err) {
				account = Account{
					id:           aId,
					rId:          strconv.Itoa(aId),
					accessToken:  accessToken,
					rAccessToken: replaceAllDashes(accessToken),
					apiKey:       apiKey,
					rApiKey:      replaceAllDashes(apiKey),
				}
				return connect(false)
			}
		}
	}
	return InvalidAccessToken, nil
}

func SendTokenRequest(request CaptchaTokenRequest) error {
	if IsClosed() {
		return errors.New("the connection is not available")
	}
	request.ApiKey = account.apiKey
	request.CreatedAt = getCurrentUnixTime()
	var jsonData []byte
	jsonData, err := json.Marshal(request)
	if !checkError(err) {
		err = directChannel.Publish(
			directExchangeName,
			requestTokenRouteKey,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        jsonData,
			})
	}
	return err
}

func SendTokenCancelRequest(request CaptchaTokenCancelRequest) error {
	if IsClosed() {
		return errors.New("the connection is not available")
	}
	request.ApiKey = account.apiKey
	request.CreatedAt = getCurrentUnixTime()
	var jsonData []byte
	jsonData, err := json.Marshal(request)
	if !checkError(err) {
		err = fanoutChannel.Publish(
			fanoutExchangeName,
			requestTokenCancelRouteKey,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        jsonData,
			})
	}
	return err
}

func SendTokenCancelAllRequest() error {
	var request = CaptchaTokenCancelRequest{
		ResponseRequired: true,
	}
	return SendTokenCancelRequest(request)
}

func IsConnected() bool {
	return !IsClosed()
}

func IsClosed() bool {
	return connection == nil || connection.IsClosed()
}

func Close() error {
	var err error = nil
	explicitShutdown = true
	if !IsClosed() {
		err = connection.Close()
	}
	return err
}

func connect(reconnect bool) (ConnectResult, error) {
	if connectionPending {
		return ConnectionPending, nil
	}
	connectionPending = true
	if !reconnect {
		updateStatus(Connecting)
	}
	currentTime := getCurrentUnixTime()
	elapsedTime := currentTime - lastConnectAttempt
	if elapsedTime < connectAttemptDelay {
		sleepTime := time.Duration(connectAttemptDelay - elapsedTime)
		duration := sleepTime * time.Second
		time.Sleep(duration)
	}
	lastConnectAttempt = getCurrentUnixTime()
	var result, err = verifyCredentials()
	if result == Success {
		createKeys()
		err = startConnection()
		if !checkError(err) {
			err = registerConsumers()
			if !checkError(err) {
				updateStatus(Connected)
				reconnectAttempt = 0
				go registerNotifyClose()
				err = nil
			}
		}
		if err != nil {
			if connection != nil {
				defer connection.Close()
			}
			result = ConnectionError
		}
	}
	if err != nil && !reconnect {
		updateStatus(Disconnected)
	}
	connectionPending = false
	return result, err
}

func registerNotifyClose() {
	closeCh := make(chan *amqp.Error)
	connection.NotifyClose(closeCh)
	connErr := <-closeCh
	if explicitShutdown {
		updateStatus(Disconnected)
	} else if connErr != nil {
		errorListener(connErr)
		go reconnect()
	}
}

func reconnect() {
	for IsClosed() && !explicitShutdown {
		updateStatus(Reconnecting)
		timeout := getNextReconnectTimeout()
		duration := timeout * time.Second
		time.Sleep(duration)
		_, err := connect(true)
		if checkError(err) {
			errorListener(err)
		}
	}
}

func createKeys() {
	directExchangeName = createKeyWithAccountId(directExchangePrefix)
	fanoutExchangeName = createKeyWithAccountId(fanoutExchangePrefix)

	responseQueueName = createKeyWithAccountIdAndApiKey(responseQueuePrefix)
	responseTokenRouteKey = createKeyWithAccountIdAndApiKey(responseTokenRoutePrefix)
	responseTokenCancelRouteKey = createKeyWithAccountIdAndApiKey(responseTokenCancelRoutePrefix)

	requestTokenRouteKey = createKeyWithAccessToken(requestTokenRoutePrefix)
	requestTokenCancelRouteKey = createKeyWithAccessToken(requestTokenCancelRoutePrefix)
}

func startConnection() error {
	conn, err := amqp.Dial("amqp://" + strconv.Itoa(account.id) + ":" + account.accessToken + "@" + hostname + ":5672/" + vHost + "?heartbeat=10")
	if checkError(err) {
		return err
	}
	connection = conn
	directCh, err := conn.Channel()
	if !checkError(err) {
		fanoutCh, err := conn.Channel()
		if !checkError(err) {
			directChannel = directCh
			fanoutChannel = fanoutCh
			return nil
		}
	}
	return err
}

func registerConsumers() error {
	err := bindQueues()
	if !checkError(err) {
		messages, err := directChannel.Consume(
			responseQueueName,
			"",
			autoAckQueue,
			exclusiveQueue,
			false,
			false,
			nil,
		)
		if !checkError(err) {
			go func() {
				for delivery := range messages {
					switch delivery.RoutingKey {
					case responseTokenRouteKey:
						processTokenMessage(&delivery)
					case responseTokenCancelRouteKey:
						processTokenCancelMessage(&delivery)
					}
				}
			}()
		}
	}
	return err
}

func processTokenMessage(message *amqp.Delivery) {
	var tokenResponse CaptchaTokenResponse
	var err = json.Unmarshal(message.Body, &tokenResponse)
	if !checkError(err) {
		captchaTokenResponseListener(tokenResponse)
	} else {
		errorListener(err)
	}
}

func processTokenCancelMessage(message *amqp.Delivery) {
	var tokenCancelResponse CaptchaTokenCancelResponse
	var err = json.Unmarshal(message.Body, &tokenCancelResponse)
	if !checkError(err) {
		captchaTokenCancelResponseListener(tokenCancelResponse)
	} else {
		errorListener(err)
	}
}

func bindQueues() error {
	var err error
	err = directChannel.QueueBind(
		responseQueueName,
		responseTokenRouteKey,
		directExchangeName,
		false,
		nil,
	)
	if !checkError(err) {
		err = directChannel.QueueBind(
			responseQueueName,
			responseTokenCancelRouteKey,
			directExchangeName,
			false,
			nil,
		)
	}
	return err
}

func verifyCredentials() (ConnectResult, error) {
	url := "https://dash.autosolve.aycd.io/rest/" + account.accessToken + "/verify/" + account.apiKey + "?clientId=" + clientId
	response, err := http.Get(url)
	if checkError(err) {
		return VerificationError, err
	}
	defer response.Body.Close()
	switch response.StatusCode {
	case http.StatusOK:
		return Success, nil
	case http.StatusBadRequest:
		return InvalidClientId, nil
	case http.StatusUnauthorized:
		return InvalidCredentials, nil
	case http.StatusTooManyRequests:
		return TooManyRequests, nil
	default:
		return UnknownError, nil
	}
}

func updateStatus(status Status) {
	statusListener(status)
}

func replaceAllDashes(key string) string {
	return strings.Replace(key, "-", "", -1)
}

func checkError(err error) bool {
	return err != nil
}

func isNotEmpty(str string) bool {
	return len(str) > 0
}

func createKeyWithAccessToken(prefix string) string {
	return prefix + "." + account.rAccessToken
}

func createKeyWithAccountId(prefix string) string {
	return prefix + "." + account.rId
}

func createKeyWithAccountIdAndApiKey(prefix string) string {
	return prefix + "." + account.rId + "." + account.rApiKey
}

func getCurrentUnixTime() int64 {
	return time.Now().Unix()
}

func getNextReconnectTimeout() time.Duration {
	var index int
	if reconnectAttempt >= len(reconnectTimeouts) {
		index = len(reconnectTimeouts) - 1
	} else {
		index = reconnectAttempt
		reconnectAttempt += 1
	}
	return reconnectTimeouts[index]
}
