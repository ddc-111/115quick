package Download

import (
	"net/http"

	"115Quick_server/internal/logic/Download"
	"115Quick_server/internal/respType"
	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 重命名
func StartReNameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StartReNameReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := Download.NewStartReNameLogic(r.Context(), svcCtx)
		resp, err := l.StartReName(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, respType.SuccessResp{
				Code: http.StatusOK,
				Data: resp,
			})
		}
	}
}
