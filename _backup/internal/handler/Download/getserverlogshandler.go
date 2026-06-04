package Download

import (
	"net/http"

	"115Quick_server/internal/logic/Download"
	"115Quick_server/internal/respType"
	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetServerLogsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetServerLogsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := Download.NewGetServerLogsLogic(r.Context(), svcCtx)
		resp, err := l.GetServerLogs(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, respType.SuccessResp{
				Code: 200,
				Data: resp,
			})
		}
	}
}
