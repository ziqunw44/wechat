package response

import "github.com/ziqunw44/wechat/internal/model"

type SearchResponse struct {
	User  model.User  `json:"user"`
	Group model.Group `json:"group"`
}
