package executor

import (
	"gae-backend-analysis/domain"
	"gae-backend-analysis/infrastructure/log"
)

type ConsumeExecutor struct {
	talentEvent domain.TalentEvent
}

func (d *ConsumeExecutor) SetupConsume() {

	d.talentEvent.ConsumeRepo()
	log.GetTextLogger().Info("Get UnRank Cleansing Talent Queue Start")

	d.talentEvent.ConsumeContributors()
	log.GetTextLogger().Info("ALL-----QUEUE----START-----SUCCESSFULLY")
	//在这里全部启动消费者逻辑
}

func NewConsumeExecutor(t domain.TalentEvent) *ConsumeExecutor {
	return &ConsumeExecutor{talentEvent: t}
}
