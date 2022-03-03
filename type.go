package notification

import (
	"time"
)

const (
	AllClients = "all"
)

type json map[string]string
type notification_bool_exp map[string]interface{}
type notification_set_input map[string]interface{}

// AndroidBackgroundLayout allows setting a background image for the notification. This is a JSON object containing the following keys.
// https://documentation.onesignal.com/docs/android-customizations#section-background-images
type AndroidBackgroundLayout struct {
	// Asset file, android resource name, or URL to remote image.
	Image string `json:"image,omitempty"`
	// Title text color ARGB Hex format. Example(Blue): "FF0000FF".
	HeadingsColor string `json:"headings_color,omitempty"`
	// Body text color ARGB Hex format. Example(Red): "FFFF0000"
	ContentsColor string `json:"contents_color,omitempty"`
}

// NotificationButton action button to the notification. The id field is required.
type NotificationButton struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Icon string `json:"icon,omitempty"`
	URL  string `json:"url,omitempty"`
}

type NotificationMetadata struct {
	Color     string            `json:"color"`
	URL       string            `json:"url"`
	ImageURL  string            `json:"image_url"`
	Subtitles map[string]string `json:"subtitles"`

	SmallIcon    string            `json:"small_icon,omitempty"`
	LargeIcon    string            `json:"large_icon,omitempty"`
	GroupID      string            `json:"group_id,omitempty"`
	GroupMessage map[string]string `json:"group_message,omitempty"`
	// iOS 15+ Relevance Score is a score to be set per notification to indicate how it should be displayed when grouped.
	// https://documentation.onesignal.com/docs/ios-relevance-score
	IOSRelevanceScore float32 `json:"ios_relevance_score,omitempty"`
	// iOS 15+ Focus Modes and Interruption Levels indicate the priority and delivery timing of a notification, to ‘interrupt’ the user.
	IOSInterruptionLevel string `json:"ios_interruption_level,omitempty"`
	// Describes whether to set or increase/decrease your app's iOS badge count by the ios_badgeCount specified count.
	// Can specify None, SetTo, or Increase.
	IOSBadgeType string `json:"ios_badge_type,omitempty"`
	// Used with ios_badgeType, describes the value to set or amount to increase/decrease your app's iOS badge count by.
	// You can use a negative number to decrease the badge count when used with an ios_badgeType of Increase.
	IOSBadgeCount int `json:"ios_badge_count,omitempty"`
	// iOS: Category APS payload, use with registerUserNotificationSettings:categories in your Objective-C / Swift code.
	// Example: calendar category which contains actions like accept and decline
	// iOS 10+ This will trigger your UNNotificationContentExtension whose ID matches this category.
	IOSCategory string `json:"ios_category,omitempty"`
	// In iOS you can specify the type of icon to be used in an Action button as being either ['system', 'custom']
	IconType string `json:"icon_type,omitempty"`
	// Sets the background color of the notification circle to the left of the notification text.
	// Only applies to apps targeting Android API level 21+ on Android 5.0+ devices.
	AccentColor string `json:"accent_color"`
	// Android Allowing setting a background image for the notification. This is a JSON object containing the following keys.
	// https://documentation.onesignal.com/docs/android-customizations#section-background-images
	AndroidBackgroundLayout *AndroidBackgroundLayout `json:"android_background_layout,omitempty"`
	// iOS 10+, Android Only one notification with the same id will be shown on the device.
	// Use the same id to update an existing notification instead of showing a new one. Limit of 64 characters.
	CollapseID string `json:"collapse_id,omitempty"`
	// Buttons to add to the notification.
	Buttons []NotificationButton `json:"buttons,omitempty"`
	// Delivery priority through the push server (example GCM/FCM).
	// Pass 10 for high priority or any other integer for normal priority.
	// Defaults to normal priority for Android and high for iOS.
	// For Android 6.0+ devices setting priority to high will wake the device out of doze mode.
	Priority uint `json:"priority,omitempty"`
	// Time To Live - In seconds. The notification will be expired if the device does not come back online within this time.
	// The default is 259,200 seconds (3 days).
	// Max value to set is 2419200 seconds (28 days).
	TTL uint `json:"ttl,omitempty"`
	// Apps with throttling enabled
	// - does not work with timezone or intelligent delivery, throttling limits will take precedence. Set to 0 if you want to use timezone or intelligent delivery.
	// - the parameter value will be used to override the default application throttling value set from the dashboard settings.
	// - parameter value 0 indicates not to apply throttling to the notification.
	// - if the parameter is not passed then the default app throttling value will be applied to the notification.
	// Apps with throttling disabled
	// - this parameter can be used to throttle delivery for the notification even though throttling is not enabled at the application level.
	// https://documentation.onesignal.com/docs/throttling
	ThrottleRatePerMinute uint `json:"throttle_rate_per_minute,omitempty"`
	// When frequency capping is enabled for the app, sending true will apply the frequency capping to the notification.
	// If the parameter is not included, the default behavior is to apply frequency capping if the setting is enabled for the app.
	// Setting the parameter to false will override the frequency capping, meaning that the notification will be sent without considering frequency capping.
	EnableFrequencyCap bool `json:"enable_frequency_cap,omitempty"`
	// iOS Set the value to voip for sending VoIP Notifications
	// This field maps to the APNS header apns-push-type.
	// Note: alert and background are automatically set by OneSignal
	// https://documentation.onesignal.com/docs/voip-notifications
	APNSPushTypeOverride string `json:"apns_push_type_override,omitempty"`
	// If send_after is used, this takes effect after the send_after time has elapsed.
	// Cannot be used if Throttling enabled. Set throttle_rate_per_minute to 0 to disable throttling if enabled to use these features.
	DelayedOption string `json:"delayed_option,omitempty"`
	// Use with delayed_option=timezone.
	DeliveryTimeOfDay string `json:"delivery_time_of_day,omitempty"`
	// Use to target a specific experience in your App Clip, or to target your notification to a specific window in a multi-scene App.
	// https://documentation.onesignal.com/docs/app-clip-support
	TargetContentIdentifier string `json:"target_content_identifier,omitempty"`
	// Use "data" or "message" depending on the type of notification you are sending
	// https://documentation.onesignal.com/docs/data-notifications
	HuaweiMsgType string `json:"huawei_msg_type,omitempty"`
	// The Android Oreo Notification Category to send the notification under.
	AndroidChannelID string `json:"android_channel_id,omitempty"`
	// Use this if you have client side Android Oreo Channels you have already defined in your app with code.
	ExistingAndroidChannelID string `json:"existing_android_channel_id,omitempty"`
	// The Android Oreo Notification Category to send the notification under
	HuaweiChannelID string `json:"huawei_channel_id,omitempty"`
	// Use this if you have client side Android Oreo Channels you have already defined in your app with code.
	HuaweiExistingChannelID string `json:"huawei_existing_channel_id,omitempty"`
}

type SendNotificationInput struct {
	AppID       string                `json:"api_id,omitempty"`
	ClientName  string                `json:"client_name,omitempty"`
	TemplateID  string                `json:"template_id,omitempty"`
	Broadcast   bool                  `json:"broadcast"`
	Headings    map[string]string     `json:"headings,omitempty"`
	Contents    map[string]string     `json:"contents,omitempty"`
	SubjectType string                `json:"subject_type,omitempty"`
	SubjectID   string                `json:"subject_id,omitempty"`
	Topics      []string              `json:"topics,omitempty"`
	UserIDs     []string              `json:"user_ids,omitempty"`
	SendAfter   time.Time             `json:"send_after,omitempty"`
	Data        map[string]string     `json:"data,omitempty"`
	Metadata    *NotificationMetadata `json:"metadata,omitempty"`
	Visible     bool                  `json:"visible,omitempty"`
	Save        bool                  `json:"save,omitempty"`
}

type SendResponse struct {
	Success           bool        `json:"success" graphql:"success"`
	RateLimitExceeded bool        `json:"rate_limit_exceeded" graphql:"rate_limit_exceeded"`
	ClientName        string      `json:"client_name,omitempty" graphql:"client_name"`
	RequestID         string      `json:"request_id,omitempty" graphql:"request_id"`
	MessageID         string      `json:"message_id,omitempty" graphql:"message_id"`
	Error             interface{} `json:"error,omitempty" graphql:"error"`
}

type SendNotificationOutput struct {
	Responses    []*SendResponse `json:"responses" graphql:"responses"`
	SuccessCount int             `json:"success_count" graphql:"success_count"`
	FailureCount int             `json:"failure_count" graphql:"failure_count"`
}
