package v2RatelimitLimit

import (
	"errors"
	"net/http"
	"time"

	openapi "github.com/unkeyed/unkey/go/api"
	"github.com/unkeyed/unkey/go/internal/services/keys"
	"github.com/unkeyed/unkey/go/internal/services/ratelimit"
	"github.com/unkeyed/unkey/go/pkg/database"
	"github.com/unkeyed/unkey/go/pkg/fault"
	"github.com/unkeyed/unkey/go/pkg/logging"
	zen "github.com/unkeyed/unkey/go/pkg/zen"
)

type Request = openapi.V2RatelimitLimitRequestBody
type Response = openapi.V2RatelimitLimitResponseBody

type Services struct {
	Logger logging.Logger
	Keys   keys.KeyService
	DB     database.Database

	Ratelimit ratelimit.Service
}

func New(svc Services) zen.Route {
	return zen.NewRoute("POST", "/v2/ratelimit.limit", func(s *zen.Session) error {
		req := new(Request)
		err := s.BindBody(req)
		if err != nil {
			return fault.Wrap(err,
				fault.WithDesc("binding failed", "We're unable to parse the request body as json."),
			)
		}

		limitRequest := ratelimit.RatelimitRequest{
			Identifier: req.Identifier,
			Limit:      req.Limit,
			Duration:   time.Duration(req.Duration) * time.Millisecond,
			Cost:       1,
		}
		if req.Cost != nil {
			limitRequest.Cost = *req.Cost
		}

		namespace, err := svc.DB.FindRatelimitNamespaceByName(s.Context(), s.AuthorizedWorkspaceID(), req.Namespace)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {

				return fault.Wrap(
					err,
					fault.WithTag(fault.NOT_FOUND),
					fault.WithDesc("namespace not found", "This namespace does not exist."),
				)
			}

			return fault.Wrap(
				err,
				fault.WithDesc("unable to load namespace", ""),
			)
		}

		overrides, err := svc.DB.FindRatelimitOverridesByIdentifier(s.Context(), s.AuthorizedWorkspaceID(), namespace.ID, req.Identifier)

		if err != nil && !errors.Is(err, database.ErrNotFound) {
			return fault.Wrap(
				err,
				fault.WithDesc("unable to load overrides", ""),
			)
		}

		usedOverrideID := ""
		for _, override := range overrides {
			usedOverrideID = override.ID
			limitRequest.Limit = int64(override.Limit)
			limitRequest.Duration = override.Duration

			if override.Identifier == req.Identifier {
				// we found an exact match, which takes presedence over wildcard matches
				break
			}
		}
		if usedOverrideID != "" {
			s.AddHeader("X-Unkey-Override-Used", usedOverrideID)
		}

		res, err := svc.Ratelimit.Ratelimit(s.Context(), limitRequest)
		if err != nil {
			return fault.Wrap(err,
				fault.WithDesc("ratelimit failed", "We're unable to ratelimit the request."),
			)
		}

		return s.JSON(http.StatusOK, Response{
			Limit:     req.Limit,
			Remaining: res.Remaining,
			Reset:     res.Reset,
			Success:   res.Success,
		})
	})
}
