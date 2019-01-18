package service

import (
	"github.com/giantswarm/auto-oncall/flag/service/oncall"
)

type Service struct {
	Oncall oncall.Oncall
}
