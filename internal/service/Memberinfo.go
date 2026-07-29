package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/dto"
	"github.com/devhdn-212/totclient_api/internal/config"
	"github.com/devhdn-212/totclient_api/internal/connection"
	"github.com/shopspring/decimal"

	"go.uber.org/zap"
)

const (
	RedisMemberInfo = "client:memberinfo"
)

type MemberinfoService struct {
	companyRepo domain.CompanyRepository
	pasaranRepo domain.PasaranRepository
	balanceAPI  config.BalanceAPI
	httpClient  *http.Client
}

func NewMemberinfoService(
	companyRepo domain.CompanyRepository,
	pasaranRepo domain.PasaranRepository,
	balanceAPI config.BalanceAPI,
) domain.MemberinfoService {
	return &MemberinfoService{
		companyRepo: companyRepo,
		pasaranRepo: pasaranRepo,
		balanceAPI:  balanceAPI,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

// balanceAPIResponse mirrors the upstream member-wallet service's envelope —
// balance comes back as a decimal string ("1000000.00"), not a JSON number,
// so it's parsed with decimal.NewFromString rather than a numeric type.
type balanceAPIResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Record  struct {
		Username string `json:"username"`
		Balance  string `json:"balance"`
	} `json:"record"`
}

// fetchBalance resolves a launch token into username + live balance by
// calling the upstream member-wallet service directly (net/http, matching
// the only other outbound HTTP call in this codebase — util.SendTelegramAlert
// — rather than pulling in a client library used nowhere else here).
//
// Only an explicit rejection from the upstream (non-200) maps to
// domain.ErrInvalidToken — a permanent, non-retryable answer the frontend
// treats as a dead end. Network/decode failures return a plain error instead,
// so a transient wallet-service outage surfaces as the generic "Aplikasi
// bermasalah" internal error (retryable), not as "token invalid".
func (d *MemberinfoService) fetchBalance(ctx context.Context, token string) (*domain.Memberinfo, error) {
	body, err := json.Marshal(map[string]string{"token": token})
	if err != nil {
		return nil, fmt.Errorf("build balance API request body: %w", err)
	}

	url := strings.TrimRight(d.balanceAPI.URL, "/") + "/api/public/balance"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build balance API request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-KEY", d.balanceAPI.APIKey)

	resp, err := d.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("call balance API: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read balance API response: %w", err)
	}

	var parsed balanceAPIResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("decode balance API response: %w", err)
	}

	if resp.StatusCode != http.StatusOK || parsed.Status != http.StatusOK {
		connection.Log.Warn("Token is invalid",
			zap.String("token", token),
			zap.Int("http_status", resp.StatusCode),
			zap.Int("body_status", parsed.Status),
			zap.String("message", parsed.Message),
		)
		return nil, domain.ErrInvalidToken
	}

	balance, err := decimal.NewFromString(parsed.Record.Balance)
	if err != nil {
		return nil, fmt.Errorf("parse balance %q from balance API: %w", parsed.Record.Balance, err)
	}

	return &domain.Memberinfo{
		Username: parsed.Record.Username,
		Balance:  balance,
		Token:    token,
	}, nil
}

// mothershipTransactionResponse mirrors the upstream member-wallet
// service's /api/public/transaction envelope.
type mothershipTransactionResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Record  struct {
		Invoice  string `json:"invoice"`
		Username string `json:"username"`
		Balance  string `json:"balance"`
		Status   string `json:"status"`
	} `json:"record"`
}

// SubmitTransaction reports one checkout chunk's ledger entry (a debit) to
// the upstream member-wallet service — this is the actual, authoritative
// deduction of the player's balance; everything checkout does locally
// (limittotal/limitglobal, the chunkTotal-vs-cached-balance check) only ever
// decides whether it's even worth asking mothership to do this.
//
// Only "insufficient balance" maps to domain.ErrInsufficientBalance — that's
// mothership's own live balance check, more authoritative than checkout's
// own (the cached balance from CheckToken could be stale by the time this
// runs). Every other rejection (bad API key, malformed request, mothership
// itself down, unrecognized username, ...) is this service's own bug or a
// technical failure, not a business rejection — it's returned as a plain
// wrapped error so it gets the generic "Aplikasi bermasalah" treatment
// (hidden from the player, alerted to Telegram/console) instead of being
// misreported as something the player did wrong.
func (d *MemberinfoService) SubmitTransaction(ctx context.Context, req domain.MothershipTransaction) (*domain.MothershipTransactionResult, error) {
	body, err := json.Marshal(map[string]interface{}{
		"invoice":       req.Invoice,
		"pasaran":       req.Pasaran,
		"playerinvoice": req.Playerinvoice,
		"username":      req.Username,
		"credit":        req.Credit.InexactFloat64(),
		"debit":         req.Debit.InexactFloat64(),
	})
	if err != nil {
		return nil, fmt.Errorf("build mothership transaction request body: %w", err)
	}

	url := strings.TrimRight(d.balanceAPI.URL, "/") + "/api/public/transaction"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build mothership transaction request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-KEY", d.balanceAPI.APIKey)

	resp, err := d.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("call mothership transaction API: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read mothership transaction response: %w", err)
	}

	var parsed mothershipTransactionResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("decode mothership transaction response: %w", err)
	}

	if resp.StatusCode != http.StatusOK || parsed.Status != http.StatusOK {
		if parsed.Message == "insufficient balance" {
			return nil, domain.ErrInsufficientBalance
		}
		return nil, fmt.Errorf(
			"mothership transaction rejected (playerinvoice=%s username=%s): %s",
			req.Playerinvoice, req.Username, parsed.Message,
		)
	}

	balance, err := decimal.NewFromString(parsed.Record.Balance)
	if err != nil {
		return nil, fmt.Errorf("parse balance %q from mothership transaction response: %w", parsed.Record.Balance, err)
	}

	return &domain.MothershipTransactionResult{Balance: balance, Status: parsed.Record.Status}, nil
}

func (d *MemberinfoService) CheckToken(ctx context.Context, req dto.MemberinfoResponse) (*domain.Memberinfo, error) {
	// Agent and market are validated before the token itself: they identify
	// which operator/pasaran the launch belongs to, and market is looked up
	// scoped to the agent below, so a bad agent has to be caught first or it
	// would be misreported as a bad market instead.
	company, err := d.companyRepo.FindByID(ctx, req.Agen)
	if err != nil {
		return nil, err
	}
	if company.IDcompany == "" {
		msg := fmt.Sprintf(
			"Agent tidak ditemukan di database: token=%s agen=%s market=%s",
			req.Token, req.Agen, req.Market,
		)
		connection.Log.Error(msg)
		return nil, domain.ErrInvalidAgent
	}

	pasaran, err := d.pasaranRepo.FindByID(ctx, req.Agen, strings.ToUpper(req.Market))
	if err != nil {
		return nil, err
	}
	if pasaran.IDcomppasaran == "" {
		msg := fmt.Sprintf(
			"Market tidak ditemukan di database: token=%s agen=%s market=%s",
			req.Token, req.Agen, req.Market,
		)
		connection.Log.Error(msg)
		return nil, domain.ErrInvalidMarket
	}

	member, err := d.fetchBalance(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	member.AgentCode = req.Agen
	return member, nil
}
