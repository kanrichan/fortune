package fortune

import (
	"fmt"

	"yaya/core"
)

func sendPicture(botID int64, type_ int64, groupID int64, userID int64, path string) {
	core.SendMsgEX_V2(
		botID,
		type_,
		groupID,
		userID,
		fmt.Sprintf("[pic=%s]", path),
		0,
		false,
		"",
	)
}

func sendMessage(botID int64, type_ int64, groupID int64, userID int64, message string) {
	core.SendMsgEX_V2(
		botID,
		type_,
		groupID,
		userID,
		message,
		0,
		false,
		"",
	)
}
