// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"time"
)

type AuditStream struct {
	ID             *int               `json:"id,omitempty"`
	Enabled        bool               `json:"enabled,omitempty"`
	StreamType     string             `json:"stream_type,omitempty"`
	StreamDetails  string             `json:"stream_details,omitempty"`
	CreatedAt      time.Time          `json:"created_at,omitempty"`
	UpdatedAt      time.Time          `json:"updated_at,omitempty"`
	PausedAt       *time.Time         `json:"paused_at,omitempty"`
	VendorSpecific *map[string]string `json:"vendor_specific,omitempty"`
}

type AuditLogStreamKey struct {
	KeyID     string `json:"key_id"`
	PublicKey string `json:"public_key"`
}

// CreateAuditStream creates an audit log stream
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/audit-log#create-an-audit-log-streaming-configuration-for-an-enterprise
//
//meta:operation POST /enterprises/{enterprise}/audit-log/streams
func (s *EnterpriseService) CreateAuditStream(ctx context.Context, enterprise string, config *AuditStream) (*AuditStream, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/audit-log/streams", enterprise)

	req, err := s.client.NewRequest("POST", u, config)
	if err != nil {
		return nil, nil, err
	}

	out := new(AuditStream)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListAuditStreams lists all audit log streams
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/audit-log#list-audit-log-stream-configurations-for-an-enterprise
//
//meta:operation GET /enterprises/{enterprise}/audit-log/streams
func (s *EnterpriseService) ListAuditStreams(ctx context.Context, enterprise string) ([]*AuditStream, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/audit-log/streams", enterprise)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var streams []*AuditStream
	resp, err := s.client.Do(ctx, req, &streams)
	if err != nil {
		return nil, resp, err
	}
	return streams, resp, nil
}

// GetAuditStream returns a single audit log stream by ID.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/audit-log#list-one-audit-log-streaming-configuration-via-a-stream-id
//
//meta:operation GET /enterprises/{enterprise}/audit-log/streams/{stream_id}
func (s *EnterpriseService) GetAuditStream(ctx context.Context, enterprise string, streamID int) (*AuditStream, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/audit-log/streams/%d", enterprise, streamID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	out := new(AuditStream)
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// DeleteAuditStream deletes an audit log stream
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/audit-log#delete-an-audit-log-streaming-configuration-for-an-enterprise
//
//meta:operation DELETE /enterprises/{enterprise}/audit-log/streams/{stream_id}
func (s *EnterpriseService) DeleteAuditStream(ctx context.Context, enterprise string, streamID int) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/audit-log/streams/%d", enterprise, streamID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// GetAuditStreamKey retrieves the audit log key for encrypting secrets
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/audit-log#get-the-audit-log-stream-key-for-encrypting-secrets
//
//meta:operation GET /enterprises/{enterprise}/audit-log/stream-key
func (s *EnterpriseService) GetAuditStreamKey(ctx context.Context, enterprise string) (*AuditLogStreamKey, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/audit-log/stream-key", enterprise)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var key AuditLogStreamKey
	resp, err := s.client.Do(ctx, req, &key)
	if err != nil {
		return nil, resp, err
	}
	return &key, resp, nil
}
