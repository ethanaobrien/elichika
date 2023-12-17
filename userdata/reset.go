package userdata

import (
	"elichika/enum"
	"elichika/model"
	"elichika/utils"
)

func (session *Session) RemoveUserProgress(table string) {
	count, err := session.Db.Table(table).Where("user_id = ?", session.UserStatus.UserID).Delete()
	utils.CheckErr(err)
	if table == "u_story_event_history" {
		session.AddResource(
			model.Content{
				ContentType:   enum.ContentTypeStoryEventUnlock,
				ContentID:     17001,
				ContentAmount: int64(count),
			})
	}
}

func (session *Session) MarkIsNew(table string, isNew bool) {
	_, err := session.Db.Table(table).Where("user_id = ?", session.UserStatus.UserID).Update(map[string]interface{}{"is_new": isNew})
	utils.CheckErr(err)

}
