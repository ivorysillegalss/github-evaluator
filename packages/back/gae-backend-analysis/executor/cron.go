package executor

import (
	"gae-backend-analysis/domain"
)

type CronExecutor struct {
	TalentCron domain.TalentCron
}

// SetupCron 启动定时任务
func (d *CronExecutor) SetupCron() {

	go d.TalentCron.AnalyseTalent()
}

func NewCronExecutor(t domain.TalentCron) *CronExecutor {
	return &CronExecutor{TalentCron: t}
}
