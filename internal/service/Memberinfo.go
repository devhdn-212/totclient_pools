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

	"github.com/devhdn-212/totclient_pools/domain"
	"github.com/devhdn-212/totclient_pools/internal/config"
	"github.com/devhdn-212/totclient_pools/internal/connection"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/singleflight"
)

const RedisCompany = "client:company"

type MemberinfoService struct {
	companyRepo domain.CompanyRepository
	balanceAPI  config.BalanceAPI
	httpClient  *http.Client
	sf          singleflight.Group
}

func NewMemberinfoService(
	companyRepo domain.CompanyRepository,
	balanceAPI config.BalanceAPI,
) domain.MemberinfoService {
	return &MemberinfoService{
		companyRepo: companyRepo,
		balanceAPI:  balanceAPI,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

// getCompany looks up a company by idcompany, cached in Redis for 24h (same
// pattern as totclient_api's pasaranService.FindID) — this is called on
// every wallet debit, so it can't hit Postgres every time. Urlapitoto (the
// per-agent wallet API base URL) comes from here since each agent/idcompany
// can point at its own wallet API.
func (d *MemberinfoService) getCompany(ctx context.Context, idcompany string) (domain.Company, error) {
	redisKey := RedisCompany + ":" + strings.ToLower(idcompany)

	cached, found, err := connection.GetRedis(redisKey)
	if err != nil {
		return domain.Company{}, err
	}
	var record domain.Company
	if found {
		if err := json.Unmarshal([]byte(cached), &record); err == nil {
			return record, nil
		}
	}

	result, err, _ := d.sf.Do(redisKey, func() (any, error) {
		return d.fetchCompany(ctx, idcompany, redisKey)
	})
	if err != nil {
		return domain.Company{}, err
	}
	return result.(domain.Company), nil
}

func (d *MemberinfoService) fetchCompany(ctx context.Context, idcompany, redisKey string) (domain.Company, error) {
	company, err := d.companyRepo.FindByID(ctx, idcompany)
	if err != nil {
		return domain.Company{}, err
	}
	go connection.SetRedis(redisKey, company, 24*time.Hour)
	return company, nil
}

// balanceAPIResponse mirrors the upstream member-wallet service's
// /api/public/balance envelope — same shape totclient_api's fetchBalance
// parses (see that repo's Memberinfo.go), duplicated here since these are
// separate Go modules.
type balanceAPIResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Record  struct {
		Username string `json:"username"`
		Balance  string `json:"balance"`
	} `json:"record"`
}

// CheckBalance re-fetches the player's live balance from the wallet API —
// called right before SubmitTransaction so a checkout event sitting in
// Kafka for a while (consumer lag, restart, etc.) doesn't debit against a
// balance that's gone stale since totclient_api's own check.
func (d *MemberinfoService) CheckBalance(ctx context.Context, idcompany, token string) (decimal.Decimal, error) {
	company, err := d.getCompany(ctx, idcompany)
	if err != nil {
		return decimal.Zero, fmt.Errorf("resolve company for balance check: %w", err)
	}

	body, err := json.Marshal(map[string]string{"token": token})
	if err != nil {
		return decimal.Zero, fmt.Errorf("build balance API request body: %w", err)
	}

	url := strings.TrimRight(company.Urlapitoto, "/") + "/api/public/balance"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return decimal.Zero, fmt.Errorf("build balance API request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-KEY", d.balanceAPI.APIKey)

	resp, err := d.httpClient.Do(httpReq)
	if err != nil {
		return decimal.Zero, fmt.Errorf("call balance API: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return decimal.Zero, fmt.Errorf("read balance API response: %w", err)
	}

	var parsed balanceAPIResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return decimal.Zero, fmt.Errorf("decode balance API response: %w", err)
	}
	if resp.StatusCode != http.StatusOK || parsed.Status != http.StatusOK {
		return decimal.Zero, fmt.Errorf("balance API rejected (status=%d message=%q)", parsed.Status, parsed.Message)
	}

	balance, err := decimal.NewFromString(parsed.Record.Balance)
	if err != nil {
		return decimal.Zero, fmt.Errorf("parse balance %q from balance API: %w", parsed.Record.Balance, err)
	}
	return balance, nil
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
// deduction of the player's balance; everything totclient_api validated
// locally before publishing the Kafka event only ever decided whether it
// was worth asking mothership to do this.
//
// Only "insufficient balance" maps to domain.ErrInsufficientBalance — that's
// mothership's own live balance check, more authoritative than whatever
// totclient_api validated against its own (possibly stale-by-now) cached
// balance. Every other rejection (bad API key, malformed request, mothership
// itself down, unrecognized username, ...) is returned as a plain wrapped
// error — the consumer just logs it (see internal/consumer), there's no
// player-facing response path here to shape the error for.
func (d *MemberinfoService) SubmitTransaction(ctx context.Context, req domain.MothershipTransaction) (*domain.MothershipTransactionResult, error) {
	company, err := d.getCompany(ctx, req.Idcompany)
	if err != nil {
		return nil, fmt.Errorf("resolve company for mothership transaction: %w", err)
	}

	body, err := json.Marshal(map[string]interface{}{
		"invoice":       req.Invoice,
		"pasaran":       req.Pasaran,
		"playerinvoice": req.Playerinvoice,
		"username":      req.Username,
		"credit":        req.Credit.InexactFloat64(),
		"debit":         req.Debit.InexactFloat64(),
	})
	if err != nil {
		fmt.Println("Wallet API - GAGAL playerinvoice:", req.Playerinvoice, "build request body:", err)
		return nil, fmt.Errorf("build mothership transaction request body: %w", err)
	}

	url := strings.TrimRight(company.Urlapitoto, "/") + "/api/public/transaction"
	fmt.Println("Wallet API - kirim playerinvoice:", req.Playerinvoice, "url:", url, "payload:", string(body))

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		fmt.Println("Wallet API - GAGAL playerinvoice:", req.Playerinvoice, "build request:", err)
		return nil, fmt.Errorf("build mothership transaction request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-KEY", d.balanceAPI.APIKey)

	resp, err := d.httpClient.Do(httpReq)
	if err != nil {
		fmt.Println("Wallet API - GAGAL playerinvoice:", req.Playerinvoice, "call API:", err)
		return nil, fmt.Errorf("call mothership transaction API: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Wallet API - GAGAL playerinvoice:", req.Playerinvoice, "read response:", err)
		return nil, fmt.Errorf("read mothership transaction response: %w", err)
	}

	var parsed mothershipTransactionResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		fmt.Println("Wallet API - GAGAL playerinvoice:", req.Playerinvoice, "decode response:", err)
		return nil, fmt.Errorf("decode mothership transaction response: %w", err)
	}

	if resp.StatusCode != http.StatusOK || parsed.Status != http.StatusOK {
		fmt.Println("Wallet API - GAGAL playerinvoice:", req.Playerinvoice, "username:", req.Username, "pesan:", parsed.Message)
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
		fmt.Println("Wallet API - GAGAL playerinvoice:", req.Playerinvoice, "parse balance:", parsed.Record.Balance, err)
		return nil, fmt.Errorf("parse balance %q from mothership transaction response: %w", parsed.Record.Balance, err)
	}

	fmt.Println("Wallet API - BERHASIL playerinvoice:", req.Playerinvoice, "username:", req.Username, "balance:", balance.String())
	return &domain.MothershipTransactionResult{Balance: balance, Status: parsed.Record.Status}, nil
}
