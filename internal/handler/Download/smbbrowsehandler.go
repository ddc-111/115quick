package Download

import (
	"net/http"

	"115Quick_server/internal/logic/Download"
	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SMBBrowseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SMBBrowseReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := Download.NewSMBBrowseLogic(r.Context(), svcCtx)
		resp, err := l.SMBBrowse(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
