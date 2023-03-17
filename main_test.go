package feishuRobotGo

import (
	"context"
	"testing"
)

func TestRobot_SendPost(t *testing.T) {
	type args struct {
		ctx     context.Context
		content *RobotContent
	}
	tests := []struct {
		name    string
		robot   *Robot
		args    args
		wantErr bool
	}{
		{
			name: "随便发一条消息",
			robot: &Robot{
				Secret: "xxx",
				Hook:   "https://open.feishu.cn/open-apis/bot/v2/hook/xxx-xx-xx-xx-xxx",
			},
			args: args{
				ctx: context.Background(),
				content: &RobotContent{
					Post: RobotPostContent{
						ZhCn: RobotPostContentGroup{
							Title: "测试标题",
							Content: [][]RobotPostContentGroupContent{
								{
									{
										Tag:    "text",
										Text:   "文本",
										Href:   "",
										UserId: "",
									},
									{
										Tag:  "a",
										Text: "超链接",
										Href: "https://bing.com",
									},
									{
										Tag:    "at",
										UserId: "xxxxxx",
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.robot.SendPost(tt.args.ctx, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("SendPost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
