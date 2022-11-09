package webhooks

import (
	"context"
	"errors"
	"fmt"
	"nebula/cmd/bot/inject"
	"nebula/pkg/infra/profiles"
	"strings"
	"sync"
	"time"

	"nebula/pkg/logger"

	"github.com/andersfylling/snowflake"
	"github.com/lithammer/shortuuid/v3"
	"github.com/nickname32/discordhook"

	"github.com/ashwanthkumar/slack-go-webhook"
)

var (
	ErrWebhookNotFound = errors.New("webhook does not exist")
	webhooks           = make(map[string]WebhookData)
	webhooksMutex      = sync.RWMutex{}

	webhookQueue      []CheckoutData
	webhookQueueMutex = sync.RWMutex{}

	colors = map[OrderStatus]int{
		ColorStatusFailed: 15888006,
		ColorStatusPlaced: 8893325,
	}

	initialized bool
)

func AddWebhook(webhook WebhookData) string {
	webhooksMutex.Lock()
	defer webhooksMutex.Unlock()

	webhookId := shortuuid.New()
	webhooks[webhookId] = webhook
	return webhookId
}

func EditWebhook(id string, webhook WebhookData) error {
	webhooksMutex.Lock()
	defer webhooksMutex.Unlock()

	_, ok := webhooks[id]
	if !ok {
		return ErrWebhookNotFound
	}

	webhooks[id] = webhook
	return nil
}

func doesWebhookExist(id string) bool {
	webhooksMutex.RLock()
	_, ok := webhooks[id]
	webhooksMutex.RUnlock()
	return ok
}

func DeleteWebhook(webhookId string) error {
	if !doesWebhookExist(webhookId) {
		return ErrWebhookNotFound
	}

	webhooksMutex.RLock()
	defer webhooksMutex.RUnlock()
	delete(webhooks, webhookId)
	return nil
}

func GetWebhookIds() []string {
	webhooksMutex.RLock()
	defer webhooksMutex.RUnlock()

	var webhookIds []string
	for id := range webhooks {
		webhookIds = append(webhookIds, id)
	}
	return webhookIds
}

func GetWebhook(id string) (WebhookData, error) {
	webhooksMutex.RLock()
	defer webhooksMutex.RUnlock()

	if wh, ok := webhooks[id]; ok {
		return wh, nil
	} else {
		return WebhookData{}, ErrWebhookNotFound
	}
}

func Queue(data CheckoutData) {
	if !initialized {
		go Poll()
		initialized = true
	}

	webhookQueueMutex.Lock()
	webhookQueue = append(webhookQueue, data)
	webhookQueueMutex.Unlock()
}

func RemoveUpToWebhooks(url string) string {
	return url[strings.Index(url, "webhooks/")+9:]
}

func checkProfile(profs []Profile, taskId string) bool {
	profile, err := profiles.GetProfileForTask(taskId)
	if err != nil {
		return true
	}
	profileGroupId, err := profiles.GetProfileGroup(profile.Id)
	if err != nil {
		return true
	}

	for _, p := range profs {
		if p.ID == profileGroupId {
			return false
		}
	}

	return true
}

func addField(name WebhookField, data CheckoutData, sensitivity bool) []*discordhook.EmbedField {
	values := map[string]string{}
	inline := true

	switch name {
	case TaskField:
		if sensitivity {
			values = map[string]string{
				"Task Group": fmt.Sprintf("||%s||", data.Task.GroupName),
				"Task": fmt.Sprintf("||%s||", data.Task.ID),
			}
		} else {
			values = map[string]string{
				"Task Group": data.Task.GroupName,
				"Task": data.Task.ID,
			}
		}
	case ProductField:
		if sensitivity {
			values = map[string]string{
				"Product": fmt.Sprintf("||[%s](%s)||", data.Product.Name, data.Product.Url),
				"Price": fmt.Sprintf("||%s||", data.Product.Price),
				"Size": fmt.Sprintf("||%s||", data.Product.Size),
				"Quantity": fmt.Sprintf("||%s||", data.Product.Quantity),
			}
		} else {
			values = map[string]string{
				"Product": fmt.Sprintf("[%s](%s)", data.Product.Name, data.Product.Url),
				"Price": fmt.Sprintf("%s", data.Product.Price),
				"Size": fmt.Sprintf("%s", data.Product.Size),
				"Quantity": fmt.Sprintf("%s", data.Product.Quantity),
			}
		}
	case StoreField:
		if sensitivity {
			values = map[string]string{
				"Store": fmt.Sprintf("||[%s](%s)||", data.Store.Name, data.Store.Url),
			}
		} else {
			values = map[string]string{
				"Store": fmt.Sprintf("[%s](%s)", data.Store.Name, data.Store.Url),
			}
		}
	case ProfileField:
		if sensitivity {
			values = map[string]string{
				"Profile Group": fmt.Sprintf("||%s||", data.Profile.GroupName),
				"Profile": fmt.Sprintf("||%s||", data.Profile.Name),
			}
		} else {
			values = map[string]string{
				"Profile Group": fmt.Sprintf("%s", data.Profile.GroupName),
				"Profile": fmt.Sprintf("%s", data.Profile.Name),
			}
		}
	case ProxyField:
		if sensitivity {
			values = map[string]string{
				"Proxy Group": fmt.Sprintf("||%s||", data.Proxy.GroupName),
				"Proxy": fmt.Sprintf("||%s||", data.Proxy.Name),
			}
		} else {
			values = map[string]string{
				"Proxy Group": fmt.Sprintf("%s", data.Proxy.GroupName),
				"Proxy": fmt.Sprintf("%s", data.Proxy.Name),
			}
		}
	case SizeField:
		if sensitivity {
			values = map[string]string{
				"Size": fmt.Sprintf("||%s||", data.Product.Size),
			}
		} else {
			values = map[string]string{
				"Size": fmt.Sprintf("%s", data.Product.Size),
			}
		}
	case QuantityField:
		if sensitivity {
			values = map[string]string{
				"Quantity": fmt.Sprintf("||%s||", data.Product.Quantity),
			}
		} else {
			values = map[string]string{
				"Quantity": fmt.Sprintf("%s", data.Product.Quantity),
			}
		}
	case OrderField:
		var orderId = data.Order.ID
		if orderId == "" {
			orderId = "None"
		}

		var orderUrl = data.Order.Url
		if orderUrl == "" {
			orderUrl = data.Store.Url
		}

		if sensitivity {
			values = map[string]string{
				"Order": fmt.Sprintf("||[%s](%s)||", orderId, orderUrl),
			}
		} else {
			values = map[string]string{
				"Order": fmt.Sprintf("[%s](%s)", orderId, orderUrl),
			}
		}
	case EmailField:
		inline = true
		if sensitivity {
			values = map[string]string{
				"Email": fmt.Sprintf("||%s||", data.Profile.Email),
			}
		} else {
			values = map[string]string{
				"Email": fmt.Sprintf("%s", data.Profile.Email),
			}
		}
	}

	var fields []*discordhook.EmbedField
	for k, v := range values {
		if v != "" {
			fields = append(fields, &discordhook.EmbedField{
				Name:   k,
				Value:  v,
				Inline: inline,
			})
		}
	}

	 return fields
}

func sendDiscordWebhook(wh WebhookData, data CheckoutData) {
	if checkProfile(wh.Profiles, data.Task.ID) {
		return
	}
	if !wh.Declines && data.Status != OrderStatusPlaced {
		return
	}

	now := time.Now()

	color := colors[ColorStatusFailed]
	statusMessage := OrderStatusFailed
	if data.Status == OrderStatusFailed {
		color = colors[ColorStatusPlaced]
		statusMessage = OrderStatusFailed
	}

	urlParts := strings.Split(RemoveUpToWebhooks(wh.Url), "/")
	wa, err := discordhook.NewWebhookAPI(snowflake.ParseSnowflakeString(urlParts[0]), urlParts[1], true, nil)
	if err != nil {
		logger.Error("error creating discord webhook: ", err)
	}

	var fields []*discordhook.EmbedField
	for _, name := range wh.Fields {
		f := addField(name, data, wh.Sensitivity)
		if len(f) != 0 {
			for _, a := range f {
				fields = append(fields, a)
			}
		}
	}

	embed := discordhook.Embed{
		Title:       string(statusMessage),
		Timestamp:   nil,
		Color:       color,
		Footer:      &discordhook.EmbedFooter{
			Text: fmt.Sprintf("Nebulabots (%s)", inject.CommitHash),
			IconURL: "https://images-ext-2.discordapp.net/external/ouj9wWdCNDBTjuD_DHaeU9iBU7sdkmhRSw2XdfwCA-4/https/nebulabots.s3.amazonaws.com/nebula-logo.png",
		},
		Thumbnail:   &discordhook.EmbedThumbnail{
			URL:      data.Product.Image,
		},
		Fields:      fields,
	}

	for _, f := range wh.Fields {
		if f == "Timestamp" {
			embed.Timestamp = &now
			break
		}
	}

	_, err = wa.Execute(
		context.TODO(), // TODO: Add a channel context
		&discordhook.WebhookExecuteParams{Embeds: []*discordhook.Embed{ &embed }},
		nil,
		"",
	)
}

func sendSlackWebhook(wh WebhookData, data CheckoutData) {
	statusMessage := "Order failed"
	if data.Status == "placed" {
		statusMessage = "Order placed"
	}

	webhookUrl := wh.Url
	attachment1 := slack.Attachment{}
	attachment1.AddField(slack.Field{Title: "Status", Value: statusMessage}).
		AddField(slack.Field{Title: "Product", Value: data.Product.Name, Short: false}).
		AddField(slack.Field{Title: "Price", Value: data.Product.Price, Short: true}).
		AddField(slack.Field{Title: "Quantity", Value: data.Product.Quantity, Short: true}).
		AddField(slack.Field{Title: "Size", Value: data.Product.Size, Short: true})

	payload := slack.Payload{
		Text:        "",
		Username:    "",
		Channel:     "",
		Attachments: []slack.Attachment{attachment1},
	}
	err := slack.Send(webhookUrl, "", payload)
	if len(err) > 0 {
		fmt.Printf("error: %s\n", err)
	}
}

func sendOutWebhook(wh WebhookData, data CheckoutData) {
	if wh.Type == "slack" {
		sendSlackWebhook(wh, data)
	} else if wh.Type == "discord" {
		sendDiscordWebhook(wh, data)
	}
}

func SendToWebhook(id string, data CheckoutData) {
	wh := webhooks[id]
	sendOutWebhook(wh, data)
}

func Poll() {
	for range time.Tick(time.Second * 1) {
		if len(webhookQueue) > 0 {
			webhookQueueMutex.Lock()
			data := webhookQueue[0]
			for _, wh := range webhooks {
				if wh.Active {
					sendOutWebhook(wh, data)
				}
			}
			webhookQueue = webhookQueue[1:]
			webhookQueueMutex.Unlock()
		}
	}
}
