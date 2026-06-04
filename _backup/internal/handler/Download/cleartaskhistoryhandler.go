package Download

import (
	"net/http"

	"115Quick_server/internal/logic/Download"
	"115Quick_server/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ClearTaskHistoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := Download.NewClearTaskHistoryLogic(r.Context(), svcCtx)
		resp, err := l.ClearTaskHistory()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
