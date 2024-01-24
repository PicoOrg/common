package eip

import (
	"github.com/picoorg/common/context"
)

type EIP interface {
	Refresh(ctx context.Context, tables []RefreshTable) (config *gatewayConfig, err error)
}
