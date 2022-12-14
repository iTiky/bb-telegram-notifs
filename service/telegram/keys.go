package telegram

// Telegram commands and their formats.
const (
	cmdRepos = "repos" // list registered repos with inline keyboard to subscribe/unsubscribe

	cmdSetBBEmail       = "set_bb_email" // set BitBucket email
	cmdSetBBEmailFormat = `/set\_bb\_email <email>`
)

// Telegram callback data.
const (
	callbackDataSubscribeAll          = "subscribe_all"           // {data}/{userID}/{repoID}
	callbackDataSubscribeReviewerOnly = "subscribe_reviewer_only" // {data}/{userID}/{repoID}
	callbackDataUnsubscribe           = "unsubscribe"             // {data}/{userID}/{repoID}
)
