package task

import (
	"testing"
)


func TestCreateInboxQueueMessage(t *testing.T) {
	var testData SQSPushLambdaPayload
	
	testData = SQSPushLambdaPayload{
		PushMessage:  PushMessage{SkipInbox: false, GUID: "", PlayerID: 0, MessageTemplateID: "", Badge: false, BgImage: "", CallToAction: "", CallToActionEvent: "", CallToActionPoints: "", CompletedAt: "", CompleteOnCallToAction: "", Component: "", ComponentParams: "", CreatedAt: "", CreatedDate: "", Goto: "", Image: "", MemberID: "", NextMessage: "", PushMessage: "", S3ProgramImageExtension: "", S3ProgramImageName: "", Title: "", ViewPoints: "", Body: "", TemplateBody: ""},
		BonusMessage: Bonus{Guid: "", AwardedPoints: "", CreatedAt: "", UpdatedAt: "", TemplateId: "", ExpiresAt: "", ExpiresIn: "", AwardedAt: "", RewardTitle: "", CompletedAt: "", EventName: "", EventCount: "", EventCounter: "", Type: ""},
		PushMessages: []Push{},
	}

	//testData.BonusMessage.Guid = "Ador"

	CreateInboxQueueMessage(testData)

	if (testData.PushMessage.SkipInbox == Data["skip_inbox"]){
		t.Logf("passed")
	}else{
		t.Errorf("failed")
	}
}


// func TestAdd(t *testing.T) {

	
// 	a :=  5
// 	b := 7
	
// 	if (Add(a,b) != 12){
// 		t.Errorf("Add(5,2) failed, we expected")
// 	}else{
// 		t.Logf("passed")
// 	}
// }